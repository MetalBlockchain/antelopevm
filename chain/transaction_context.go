package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/chain/types"
)

type TransactionContext struct {
	Control          *Controller
	Trx              *types.SignedTransaction
	ID               types.TransactionIdType
	ApplyContextFree bool
	Trace            *types.TransactionTrace
	ActionId         types.IdType

	isInitialized bool
}

func NewTransactionContext(control *Controller, t *types.SignedTransaction, trxId types.TransactionIdType) *TransactionContext {
	tc := TransactionContext{
		Control: control,
		Trx:     t,

		isInitialized: false,
	}

	tc.Trace = &types.TransactionTrace{
		Id:           t.ID(),
		BlockNum:     0,
		BlockTime:    0,
		ActionTraces: make([]types.ActionTrace, 0),
	}

	return &tc
}

func (t *TransactionContext) Init() error {
	if t.Trx.DelaySec > 0 {
		return fmt.Errorf("Deferred transactions are deprecated")
	}

	t.isInitialized = true

	return nil
}

func (t *TransactionContext) Exec() error {
	if !t.isInitialized {
		return fmt.Errorf("must first initialize")
	}

	if t.ApplyContextFree {
		for _, act := range t.Trx.ContextFreeActions {
			t.ScheduleAction(*act, act.Account, true, 0)
		}
	}

	if t.Trx.DelaySec == 0 {
		for _, act := range t.Trx.Actions {
			t.ScheduleAction(*act, act.Account, false, 0)
		}
	}

	for i := 1; i <= len(t.Trace.ActionTraces); i++ {
		if err := t.ExecuteAction(i, 0); err != nil {
			return err
		}
	}

	return nil
}

func (t *TransactionContext) ScheduleAction(action types.Action, receiver types.AccountName, contextFree bool, creatorActionOrdinal int) int {
	newActionOrdinal := len(t.Trace.ActionTraces) + 1
	actionTrace := types.NewActionTrace(t.Trace, action, receiver, contextFree, newActionOrdinal, creatorActionOrdinal)
	t.Trace.ActionTraces = append(t.Trace.ActionTraces, *actionTrace)

	return newActionOrdinal
}

func (t *TransactionContext) ExecuteAction(actionOrdinal int, recurseDepth uint32) error {
	applyContext, err := NewApplyContext(t, actionOrdinal, recurseDepth)

	if err != nil {
		return err
	}

	if err := applyContext.Exec(); err != nil {
		return err
	}

	return nil
}

func (t *TransactionContext) GetActionTrace(actionOrdinal int) (*types.ActionTrace, error) {
	if actionOrdinal < 0 || actionOrdinal > len(t.Trace.ActionTraces) {
		return nil, fmt.Errorf("action_ordinal %v is outside allowed range [1,%v]", actionOrdinal, len(t.Trace.ActionTraces))
	}

	return &t.Trace.ActionTraces[actionOrdinal-1], nil
}

func (t *TransactionContext) Finalize() error {
	if err := t.Control.State.Commit(); err != nil {
		return err
	}

	return nil
}
