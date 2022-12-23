package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/inconshreveable/log15"
)

type TransactionContext struct {
	Control                      *Controller
	Trx                          *core.SignedTransaction
	ID                           core.TransactionIdType
	ApplyContextFree             bool
	Trace                        *core.TransactionTrace
	ActionId                     core.IdType
	ExecutedActionReceiptDigests []crypto.Sha256
	State                        state.State
	AuthorizationManager         AuthorizationManager

	isInitialized bool
}

func NewTransactionContext(control *Controller, s state.State, t *core.SignedTransaction, trxId core.TransactionIdType) *TransactionContext {
	tc := TransactionContext{
		Control:              control,
		Trx:                  t,
		State:                s,
		AuthorizationManager: *NewAuthorizationManager(control, s),

		isInitialized: false,
	}

	tc.Trace = &core.TransactionTrace{
		Id:           t.ID(),
		BlockNum:     0,
		BlockTime:    0,
		ActionTraces: make([]core.ActionTrace, 0),
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
			log15.Error("failed to execute action", "index", i, "error", err)
			return err
		}
	}

	return nil
}

func (t *TransactionContext) ScheduleAction(action core.Action, receiver core.AccountName, contextFree bool, creatorActionOrdinal int) int {
	newActionOrdinal := len(t.Trace.ActionTraces) + 1
	actionTrace := core.NewActionTrace(t.Trace, action, receiver, contextFree, core.Vuint32(newActionOrdinal), core.Vuint32(creatorActionOrdinal))
	t.Trace.ActionTraces = append(t.Trace.ActionTraces, *actionTrace)

	return newActionOrdinal
}

func (t *TransactionContext) ExecuteAction(actionOrdinal int, recurseDepth uint32) error {
	applyContext, err := NewApplyContext(t, actionOrdinal, recurseDepth)

	if err != nil {
		return err
	}

	if err := applyContext.Exec(); err != nil {
		log15.Error("failed to exec apply context", "error", err)
		return err
	}

	return nil
}

func (t *TransactionContext) GetActionTrace(actionOrdinal int) (*core.ActionTrace, error) {
	if actionOrdinal < 0 || actionOrdinal > len(t.Trace.ActionTraces) {
		return nil, fmt.Errorf("action_ordinal %v is outside allowed range [1,%v]", actionOrdinal, len(t.Trace.ActionTraces))
	}

	return &t.Trace.ActionTraces[actionOrdinal-1], nil
}

func (t *TransactionContext) Finalize() error {
	if err := t.State.PutTransaction(t.Trace); err != nil {
		return err
	}

	if err := t.State.Commit(); err != nil {
		return err
	}

	return nil
}
