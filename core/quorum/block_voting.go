package quorum

import (
	"errors"
	"math/big"
	"sync"

	"gopkg.in/fatih/set.v0"

	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// QuorumBlockVoting is a type of BlockMaker that uses a smart contract
// to determine the canonical chain. Parties that are allowed to
// vote send vote transactions to the voting contract. Based on
// these transactions the parent block is selected where the next
// block will be build on top of.
type QuorumBlockVoting struct {
	bc       *core.BlockChain
	cc       *params.ChainConfig
	txpool   *core.TxPool
	synced   bool
	mux      *event.TypeMux
	db       ethdb.Database
	am       *accounts.Manager
	gasPrice *big.Int

	voteSession  *BlockVotingSession
	callContract *BlockVotingCaller

	ks        *keystore.KeyStore
	voteAcct  accounts.Account
	makerAcct accounts.Account
	coinbase  common.Address

	pStateMu sync.Mutex
	pState   *pendingState
}

// Vote is posted to the event mux when the QuorumBlockVoting instance
// is ordered to send a new vote transaction. Hash is the hash for the
// given number depth.
type Vote struct {
	Hash   common.Hash
	Number *big.Int
	TxHash chan common.Hash
	Err    chan error
}

// CreateBlock is posted to the event mux when the QuorumBlockVoting instance
// is ordered to create a new block. Either the hash of the created
// block is returned is hash or an error.
type CreateBlock struct {
	Hash chan common.Hash
	Err  chan error
}

// NewQuorumBlockVoting creates a new QuorumBlockVoting instance.
// blockMakerKey and/or voteKey can be nil in case this node doesn't create blocks or vote.
// Note, don't forget to call Start.
func NewQuorumBlockVoting(bc *core.BlockChain, chainConfig *params.ChainConfig, txpool *core.TxPool, mux *event.TypeMux, db ethdb.Database, accountMgr *accounts.Manager, isSynchronised bool) *QuorumBlockVoting {
	bv := &QuorumBlockVoting{
		bc:     bc,
		cc:     chainConfig,
		txpool: txpool,
		mux:    mux,
		db:     db,
		am:     accountMgr,
		synced: isSynchronised,
	}

	return bv
}

func (bv *QuorumBlockVoting) resetPendingState(parent *types.Block) {
	publicState, err := bv.bc.State()
	if err != nil {
		panic(fmt.Sprintf("State error", "err", err))
	}

	ps := &pendingState{
		parent:      parent,
		publicState: publicState,
		//privateState:  privateState,
		header: bv.makeHeader(parent),
		gp:     new(core.GasPool),
		//ownedAccounts: accountAddressesSet(bv.am.Accounts()),
	}

	ps.gp.AddGas(ps.header.GasLimit)

	pending, err := bv.txpool.Pending()
	if err != nil {
		panic(err)
	}
	txs := NewTransactionsByPriorityAndNonce(pending)

	lowGasTxs, failedTxs := ps.applyTransactions(txs, bv.mux, bv.bc, bv.cc)
	bv.txpool.RemoveBatch(lowGasTxs)
	bv.txpool.RemoveBatch(failedTxs)

	bv.pStateMu.Lock()
	bv.pState = ps
	bv.pStateMu.Unlock()
}

func (bv *QuorumBlockVoting) makeHeader(parent *types.Block) *types.Header {
	tstart := time.Now()
	tstamp := tstart.Unix()
	if parent.Time().Cmp(new(big.Int).SetInt64(tstamp)) >= 0 {
		tstamp = parent.Time().Int64() + 1
	}
	// this will ensure we're not going off too far in the future
	if now := time.Now().Unix(); tstamp > now+4 {
		wait := time.Duration(tstamp-now) * time.Second
		log.Info("We are too far in the future. Waiting for", wait)
		time.Sleep(wait)
	}

	num := parent.Number()
	header := &types.Header{
		Number:     num.Add(num, common.Big1),
		ParentHash: parent.Hash(),
		Difficulty: ethash.CalcDifficulty(bv.cc, uint64(tstamp), parent.Header()),
		GasLimit:   core.CalcGasLimit(parent),
		GasUsed:    new(big.Int),
		Time:       big.NewInt(tstamp),
	}

	header.Coinbase = bv.makerAcct.Address

	return header
}

// Start runs the event loop.
func (bv *QuorumBlockVoting) Start(client *rpc.Client, strat BlockMakerStrategy, ks *keystore.KeyStore, voteAcct, makerAcct accounts.Account) error {
	bv.ks = ks
	bv.voteAcct = voteAcct
	bv.makerAcct = makerAcct

	ethClient := ethclient.NewClient(client)
	callContract, err := NewBlockVotingCaller(params.QuorumVotingContractAddr, ethClient)
	if err != nil {
		return err
	}
	bv.callContract = callContract

	if voteAcct != (accounts.Account{}) {
		contract, err := NewBlockVoting(params.QuorumVotingContractAddr, ethClient)
		if err != nil {
			return err
		}

		//auth := bind.NewKeyedTransactor(voteAcct)
		addr := voteAcct.Address
		auth := bind.TransactOpts{
			From: addr,
			Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if address != addr {
					return nil, errors.New("not authorized to sign this account")
				}
				signature, err := bv.ks.SignHash(voteAcct, signer.Hash(tx).Bytes())
				if err != nil {
					return nil, err
				}
				return tx.WithSignature(signer, signature)
			},
		}
		bv.voteSession = &BlockVotingSession{
			Contract: contract,
			CallOpts: bind.CallOpts{
				Pending: true,
			},
			TransactOpts: bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
			},
		}
	}

	bv.run(strat)

	log.EmitCheckpoint(log.BlockVotingStarted)

	return nil
}

func (bv *QuorumBlockVoting) run(strat BlockMakerStrategy) {
	if bv.makerAcct != (accounts.Account{}) {
		log.Info("Node configured for block creation", "address", bv.makerAcct.Address.Hex())
	}
	if bv.voteAcct != (accounts.Account{}) {
		log.Info("Node configured for block voting", "address", bv.voteAcct.Address.Hex())
	}

	sub := bv.mux.Subscribe(downloader.StartEvent{},
		downloader.DoneEvent{},
		downloader.FailedEvent{},
		core.ChainHeadEvent{},
		core.TxPreEvent{},
		Vote{},
		CreateBlock{})

	bv.resetPendingState(bv.bc.CurrentBlock())

	go func() {
		defer sub.Unsubscribe()

		strat.Start()

		for {
			select {
			case event, ok := <-sub.Chan():
				if !ok {
					return
				}

				switch e := event.Data.(type) {
				case downloader.StartEvent: // begin synchronising, stop block creation and/or voting
					bv.synced = false
					strat.Pause()
				case downloader.DoneEvent, downloader.FailedEvent: // caught up, or got an error, start block createion and/or voting
					bv.synced = true
					strat.Resume()
				case core.ChainHeadEvent: // got a new header, reset pending state
					bv.resetPendingState(e.Block)
					if bv.synced {
						number := new(big.Int)
						number.Add(e.Block.Number(), common.Big1)
						if tx, err := bv.vote(number, e.Block.Hash()); err == nil {
							log.Debug("Voted for for tx", "hash", e.Block.Hash().Hex(), "height", number, "tx", tx.Hex())
						} else {
							log.Error("Unable to vote", "err", err)
						}
					}
				case core.TxPreEvent: // tx entered pool, apply to pending state
					bv.applyTransaction(e.Tx)
				case Vote:
					if bv.synced {
						txHash, err := bv.vote(e.Number, e.Hash)
						if err == nil && e.TxHash != nil {
							e.TxHash <- txHash
						} else if err != nil && e.Err != nil {
							e.Err <- err
						} else if err != nil {
							log.Error("Unable to vote", "err", err)
						}
					} else {
						e.Err <- fmt.Errorf("Node not synced")
					}
				case CreateBlock:
					block, err := bv.createBlock()
					if err == nil && e.Hash != nil {
						e.Hash <- block.Hash()
					} else if err != nil && e.Err != nil {
						e.Err <- err
					} else if err != nil {
						log.Error("Unable to create block", "err", err)
					}

					if err != nil {
						bv.pStateMu.Lock()
						cBlock := bv.pState.parent
						bv.pStateMu.Unlock()
						num := new(big.Int).Add(cBlock.Number(), common.Big1)
						_, err := bv.vote(num, cBlock.Hash())
						if err != nil {
							log.Error("Unable to vote", "err", err)
							bv.resetPendingState(bv.bc.CurrentBlock())
						}
					}
				}
			}
		}
	}()
}

func (bv *QuorumBlockVoting) applyTransaction(tx *types.Transaction) {
	acc, err := types.Sender(types.HomesteadSigner{}, tx)
	if err != nil {
		panic(fmt.Errorf("invalid transaction", "err", err))
	}
	txs := map[common.Address]types.Transactions{acc: types.Transactions{tx}}
	txset := NewTransactionsByPriorityAndNonce(txs)

	bv.pStateMu.Lock()
	bv.pState.applyTransactions(txset, bv.mux, bv.bc, bv.cc)
	bv.pStateMu.Unlock()
}

func (bv *QuorumBlockVoting) Pending() (*types.Block /**state.StateDB,*/, *state.StateDB) {
	bv.pStateMu.Lock()
	defer bv.pStateMu.Unlock()
	return types.NewBlock(bv.pState.header, bv.pState.txs, nil, bv.pState.receipts), bv.pState.publicState.Copy() // , bv.pState.privateState.Copy()
}

func (bv *QuorumBlockVoting) createBlock() (*types.Block, error) {
	if bv.makerAcct == (accounts.Account{}) {
		return nil, fmt.Errorf("Node not configured for block creation")
	}

	log.Info("createBlock")

	ch, err := bv.canonHash(bv.pState.header.Number.Uint64())
	log.Info("canonicalHash", "ch", ch.Hex(), "err", err, "bv.pState.header", bv.pState.header.Hash().Hex(), "header number", bv.pState.header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	if ch != bv.pState.parent.Hash() {
		return nil, fmt.Errorf("invalid canonical hash: expected %v, received %v", ch.Hex(), bv.pState.header.Hash().Hex())
	}

	bv.pStateMu.Lock()
	defer bv.pStateMu.Unlock()

	state := bv.pState.publicState // shortcut
	header := bv.pState.header
	receipts := bv.pState.receipts

	ethash.AccumulateRewards(state, header, nil)

	header.Root = state.IntermediateRoot(false)

	// Quorum blocks contain a signature of the header in the Extra field.
	// This signature is verified during block import and ensures that the
	// block is created by a party that is allowed to create blocks.
	//signature, err := crypto.Sign(header.QuorumHash().Bytes(), bv.makerAcct)
	signature, err := bv.ks.SignHash(bv.makerAcct, header.QuorumHash().Bytes())
	if err != nil {
		return nil, err
	}
	header.Extra = signature

	// update block hash in receipts and logs now it is available
	for _, r := range receipts {
		for _, l := range r.Logs {
			l.BlockHash = header.Hash()
		}
	}

	header.Bloom = types.CreateBloom(receipts)

	block := types.NewBlock(header, bv.pState.txs, nil, receipts)
	if _, err := bv.bc.InsertChain(types.Blocks{block}); err != nil {
		return nil, err
	}

	bv.mux.Post(core.NewMinedBlockEvent{Block: block})

	return block, nil
}

func (bv *QuorumBlockVoting) vote(height *big.Int, hash common.Hash) (common.Hash, error) {
	if bv.voteSession == nil {
		return common.Hash{}, fmt.Errorf("Node is not configured for voting")
	}
	cv, err := bv.callContract.CanVote(nil, bv.voteSession.TransactOpts.From)
	log.Info("can vote?", "cv", cv)
	if err != nil {
		return common.Hash{}, err
	}
	if !cv {
		return common.Hash{}, fmt.Errorf("not allowed to vote", "node id", bv.voteSession.TransactOpts.From.Hex())
	}

	log.Info("vote on block", "block", hash.Hex(), "height", height)

	nonce := bv.txpool.State().GetNonce(bv.voteSession.TransactOpts.From)
	bv.voteSession.TransactOpts.Nonce = new(big.Int).SetUint64(nonce)
	defer func() { bv.voteSession.TransactOpts.Nonce = nil }()

	tx, err := bv.voteSession.Vote(height, hash)
	log.Info("vote result", "tx", tx, "err", err)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

// CanonHash returns the canonical block hash on the given height.
func (bv *QuorumBlockVoting) canonHash(height uint64) (common.Hash, error) {
	opts := &bind.CallOpts{Pending: true}
	return bv.callContract.GetCanonHash(opts, new(big.Int).SetUint64(height))
}

// isVoter returns an indication if the given address is allowed
// to vote.
func (bv *QuorumBlockVoting) isVoter(addr common.Address) (bool, error) {
	return bv.callContract.IsVoter(nil, addr)
}

// isBlockMaker returns an indication if the given address is allowed
// to make blocks
func (bv *QuorumBlockVoting) isBlockMaker(addr common.Address) (bool, error) {
	return bv.callContract.IsBlockMaker(nil, addr)
}

func accountAddressesSet(accounts []accounts.Account) *set.Set {
	accountSet := set.New()
	for _, account := range accounts {
		accountSet.Add(account.Address)
	}
	return accountSet
}
