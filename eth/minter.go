package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"sync"
	"time"
)

type work struct {
	config     *params.ChainConfig
	signer     types.Signer
	state      *state.StateDB
	Block      *types.Block
	header     *types.Header
	tcount     int
	failedTxes types.Transactions
	txs        []*types.Transaction
	receipts   []*types.Receipt
}

type minter struct {
	config   *params.ChainConfig
	engine   consensus.Engine
	mu       sync.Mutex
	mux      *event.TypeMux
	eth      miner.Backend
	chain    *core.BlockChain
	chainDb  ethdb.Database
	coinbase common.Address
}

func newMinter(config *params.ChainConfig, engine consensus.Engine, mux *event.TypeMux, eth miner.Backend) *minter {
	return &minter{
		config:   config,
		engine:   engine,
		mux:      mux,
		eth:      eth,
		chain:    eth.BlockChain(),
		chainDb:  eth.ChainDb(),
		coinbase: common.HexToAddress("0x0000000000000000000000000000000000badbad"),
	}
}

// Assumes mu is held.
func (minter *minter) createWork() *work {
	parent := minter.chain.CurrentBlock()
	parentNumber := parent.Number()
	blockNum := big.NewInt(0).Add(parentNumber, common.Big1)

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     blockNum,
		Difficulty: ethash.CalcDifficulty(minter.config, blockNum.Uint64(), parent.Header()),
		GasLimit:   core.CalcGasLimit(parent),
		GasUsed:    new(big.Int),
		Coinbase:   minter.coinbase,
		Time:       blockNum,
	}

	state, err := minter.chain.StateAt(parent.Root())
	if err != nil {
		panic(fmt.Sprint("failed to get parent state: ", err))
	}

	return &work{
		config: minter.config,
		signer: types.NewEIP155Signer(minter.config.ChainId),
		state:  state,
		header: header,
	}
}

func (minter *minter) mintNewBlock(transactions *types.TransactionsByPriceAndNonce) {
	minter.mu.Lock()
	defer minter.mu.Unlock()

	work := minter.createWork()

	header := work.header

	if err := minter.engine.Prepare(minter.chain, header); err != nil {
		log.Error("Failed to prepare header for mining", "err", err)
		return
	}

	log.Info("minting", "txes", transactions)
	//committedTxes, publicReceipts, logs := work.commitTransactions(minter.mux, transactions, minter.chain)
	work.commitTransactions(minter.mux, transactions, minter.chain)

	minter.eth.TxPool().RemoveBatch(work.failedTxes)

	log.Info("committed")

	//// commit state root after all state transitions.
	//ethash.AccumulateRewards(work.state, header, nil)
	//header.Root = work.state.IntermediateRoot(minter.chain.Config().IsEIP158(work.header.Number))
	//
	//header.Bloom = types.CreateBloom(work.receipts)
	//
	//// update block hash since it is now available, but was not when the
	//// receipt/log of individual transactions were created:
	//headerHash := header.Hash()
	//for _, l := range logs {
	//	l.BlockHash = headerHash
	//}
	//
	//block := types.NewBlock(header, committedTxes, nil, work.receipts)
	var err error
	block := work.Block
	if block, err = minter.engine.Finalize(minter.chain, header, work.state, work.txs, []*types.Header{}, work.receipts); err != nil {
		log.Error("Failed to finalize block for sealing", "err", err)
		return
	}

	log.Info("Generated next block", "block num", block.Number())

	if _, err := work.state.CommitTo(minter.chainDb, minter.chain.Config().IsEIP158(block.Number())); err != nil {
		panic(fmt.Sprint("error committing public state: ", err))
	}

	_, err = minter.chain.InsertChain([]*types.Block{block})

	if err != nil {
		panic(fmt.Sprintf("failed to extend chain: %s", err.Error()))
	}

	elapsed := time.Since(time.Unix(0, header.Time.Int64()))
	log.Info("💎  Minted block", "num", block.Number(), "hash", fmt.Sprintf("%x", block.Hash().Bytes()[:4]), "elapsed", elapsed)
}

func (env *work) commitTransactions(mux *event.TypeMux, txes *types.TransactionsByPriceAndNonce, bc *core.BlockChain) {
	gp := new(core.GasPool).AddGas(env.header.GasLimit)

	var coalescedLogs []*types.Log

	for {
		tx := txes.Peek()
		if tx == nil {
			break
		}
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(env.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !env.config.IsEIP155(env.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", env.config.EIP155Block)

			txes.Pop()
			continue
		}
		// Start executing the transaction
		env.state.Prepare(tx.Hash(), common.Hash{}, env.tcount)

		err, logs := env.commitTransaction(tx, bc, gp)
		switch err {
		case core.ErrGasLimitReached:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Trace("Gas limit exceeded for current block", "sender", from)
			txes.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			env.tcount++
			txes.Shift()

		default:
			// Pop the current failed transaction without shifting in the next from the account
			log.Trace("Transaction failed, will be removed", "hash", tx.Hash(), "err", err)
			env.failedTxes = append(env.failedTxes, tx)
			txes.Pop()
		}
	}

	//return committedTxes, publicReceipts, logs
	if len(coalescedLogs) > 0 || env.tcount > 0 {
		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		go func(logs []*types.Log, tcount int) {
			if len(logs) > 0 {
				mux.Post(core.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(core.PendingStateEvent{})
			}
		}(cpy, env.tcount)
	}
}

func (env *work) commitTransaction(tx *types.Transaction, bc *core.BlockChain, gp *core.GasPool) (error, []*types.Log) {
	snap := env.state.Snapshot()

	var author *common.Address
	receipt, _, err := core.ApplyTransaction(env.config, bc, author, gp, env.state, env.header, tx, env.header.GasUsed, vm.Config{})
	if err != nil {
		env.state.RevertToSnapshot(snap)
		return err, nil
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return nil, receipt.Logs
}
