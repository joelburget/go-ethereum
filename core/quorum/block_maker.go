package quorum

import (
	"time"

	"container/heap"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// vvv

type TransactionsByPriorityAndNonce struct {
	txs   map[common.Address]types.Transactions
	heads TxByPriority
}

// TxByPriority implements both sort and the heap interface, making it useful
// for all at once sorting as well as individual adding and removing elements.
//
// It will prioritise transaction to the voting contract.
type TxByPriority types.Transactions

func (s TxByPriority) Len() int { return len(s) }
func (s TxByPriority) Less(i, j int) bool {
	var (
		iRecipient = s[i].To()
		jRecipient = s[j].To()
	)

	// in case iReceipt is towards the voting contract and jRecipient is not towards the voting contract
	// iReceipt is "smaller".
	return iRecipient != nil && *iRecipient == params.QuorumVotingContractAddr && (jRecipient == nil || *jRecipient != params.QuorumVotingContractAddr)
}

func (s TxByPriority) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s *TxByPriority) Push(x interface{}) {
	*s = append(*s, x.(*types.Transaction))
}
func (s *TxByPriority) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// NewTransactionsByPriorityAndNonce creates a transaction set that can retrieve
// vote tx sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// it after providing it to the constructor.
func NewTransactionsByPriorityAndNonce(txs map[common.Address]types.Transactions) *TransactionsByPriorityAndNonce {
	heads := make(TxByPriority, 0, len(txs))
	for acc, accTxs := range txs {
		heads = append(heads, accTxs[0])
		txs[acc] = accTxs[1:]
	}
	heap.Init(&heads)

	return &TransactionsByPriorityAndNonce{
		txs:   txs,
		heads: heads,
	}
}

func (t *TransactionsByPriorityAndNonce) Peek() *types.Transaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0]
}

func (t *TransactionsByPriorityAndNonce) Shift() {
	acc, err := types.Sender(types.HomesteadSigner{}, t.heads[0])
	if err != nil {
		panic(fmt.Errorf("invalid transaction: %v", err))
	}
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

func (t *TransactionsByPriorityAndNonce) Pop() {
	heap.Pop(&t.heads)
}

// ^^^

// BlockMaker defines the interface that block makers must provide.
type BlockMaker interface {
	// Pending returns the pending block and pending state.
	Pending() (*types.Block, *state.StateDB)
}

type pendingState struct {
	publicState, privateState *state.StateDB
	tcount                    int // tx count in cycle
	gp                        *core.GasPool
	//ownedAccounts             *set.Set
	txs       types.Transactions // set of transactions
	lowGasTxs types.Transactions
	failedTxs types.Transactions
	parent    *types.Block

	header   *types.Header
	receipts types.Receipts
	logs     []*types.Log

	createdAt time.Time
}

func (ps *pendingState) applyTransaction(tx *types.Transaction, bc *core.BlockChain, cc *params.ChainConfig) (error, []*types.Log) {
	publicSnaphot, privateSnapshot := ps.publicState.Snapshot(), ps.privateState.Snapshot()

	// this is a bit of a hack to force jit for the miners
	config := vm.Config{} // XXX
	if !(config.EnableJit && config.ForceJit) {
		config.EnableJit = false
	}
	config.ForceJit = false // disable forcing jit

	var author *common.Address
	publicReceipt, _, err := core.ApplyTransaction(cc, bc, author, ps.gp, ps.publicState, ps.header, tx, ps.header.GasUsed, config)
	if err != nil {
		ps.publicState.RevertToSnapshot(publicSnaphot)
		ps.privateState.RevertToSnapshot(privateSnapshot)

		return err, nil
	}
	ps.txs = append(ps.txs, tx)
	ps.receipts = append(ps.receipts, publicReceipt)

	return nil, publicReceipt.Logs
}

func (ps *pendingState) applyTransactions(txs *TransactionsByPriorityAndNonce, mux *event.TypeMux, bc *core.BlockChain, cc *params.ChainConfig) (types.Transactions, types.Transactions) {
	var (
		lowGasTxs types.Transactions
		failedTxs types.Transactions
	)

	var coalescedLogs []*types.Log
	for {
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}

		from, err := types.Sender(types.HomesteadSigner{}, tx)
		if err != nil {
			log.Error(fmt.Sprintf("invalid transaction: %v", err))
		}

		// Start executing the transaction
		ps.publicState.Prepare(tx.Hash(), common.Hash{}, 0)

		err, logs := ps.applyTransaction(tx, bc, cc)
		switch {
		case err == vm.ErrOutOfGas:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Info("Gas limit reached for (%x) in this block. Continue to try smaller txs\n", from[:4])
			txs.Pop()
		case err != nil:
			// Pop the current failed transaction without shifting in the next from the account
			log.Info("Transaction (%x) failed, will be removed: %v\n", tx.Hash().Bytes()[:4], err)
			failedTxs = append(failedTxs, tx)
			txs.Pop()
		default:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			ps.tcount++
			log.EmitCheckpoint(log.TxAccepted, tx.Hash().Hex())
			txs.Shift()
		}
	}
	if len(coalescedLogs) > 0 || ps.tcount > 0 {
		go func(logs []*types.Log, tcount int) {
			if len(logs) > 0 {
				mux.Post(core.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(core.PendingStateEvent{})
			}
		}(coalescedLogs, ps.tcount)
	}

	return lowGasTxs, failedTxs
}
