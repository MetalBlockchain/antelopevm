package chain

import (
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type TransactionContext struct {
	Control          *Controller
	Trx              *types.SignedTransaction
	ID               types.TransactionIdType
	ApplyContextFree bool
	Trace            *types.TransactionTrace

	isInitialized bool
}

func NewTransactionContext(control *Controller, t *types.SignedTransaction, trxId types.TransactionIdType) *TransactionContext {
	tc := TransactionContext{
		Control: control,
		Trx:     t,

		isInitialized: false,
	}

	tc.Trace = &types.TransactionTrace{
		ID:              trxId,
		BlockNum:        0,                     // TODO: Fix
		BlockTime:       0,                     // TODO: Fix
		ProducerBlockId: crypto.NewSha256Nil(), // TODO: Fix
		ActionTraces:    []types.ActionTrace{},
	}

	return &tc
}

func (t *TransactionContext) Init() {
	if t.Trx.DelaySec > 0 {
		panic("Deferred transactions are deprecated")
	}

	t.isInitialized = true
}

func (t *TransactionContext) Exec() {
	if !t.isInitialized {
		panic("must first initialize")
	}

	if t.ApplyContextFree {
		for _, act := range t.Trx.ContextFreeActions {
			t.Trace.ActionTraces = append(t.Trace.ActionTraces, types.ActionTrace{})
			t.DispatchAction(&t.Trace.ActionTraces[len(t.Trace.ActionTraces)-1], act, act.Account, true, 0)
		}
	}

	for _, act := range t.Trx.Actions {
		t.Trace.ActionTraces = append(t.Trace.ActionTraces, types.ActionTrace{})
		t.DispatchAction(&t.Trace.ActionTraces[len(t.Trace.ActionTraces)-1], act, act.Account, false, 0)
	}
}

func (t *TransactionContext) DispatchAction(trace *types.ActionTrace, action *types.Action, receiver types.AccountName, contextFree bool, recurseDepth uint32) {
	applyContext := NewApplyContext(t, action, recurseDepth)
	applyContext.ContextFree = contextFree
	applyContext.Receiver = receiver
	applyContext.Exec(trace)
}

func (t *TransactionContext) Finalize() {
	err := t.Control.State.Commit()

	if err != nil {
		panic(err)
	}
}
