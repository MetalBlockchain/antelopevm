package chain

import (
	"errors"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type ApplyContext struct {
	Control            *Controller
	TrxContext         *TransactionContext
	Act                *types.Action
	Receiver           types.AccountName
	UsedAuthorizations []bool
	RecurseDepth       uint32
	Privileged         bool
	ContextFree        bool
	UsedContestFreeApi bool
	Notified           []types.AccountName
	InlineActions      []types.Action
	CfaInlineActions   []types.Action
	AccountRamDeltas   []types.AccountDelta
}

func NewApplyContext(trxContext *TransactionContext, act *types.Action, recurseDepth uint32) *ApplyContext {

	applyContext := &ApplyContext{
		Control:            trxContext.Control,
		TrxContext:         trxContext,
		Act:                act,
		Receiver:           act.Account,
		UsedAuthorizations: make([]bool, len(act.Authorization)),
		RecurseDepth:       recurseDepth,

		Privileged:         false,
		ContextFree:        false,
		UsedContestFreeApi: false,
	}

	applyContext.Notified = []types.AccountName{}
	applyContext.InlineActions = []types.Action{}
	applyContext.CfaInlineActions = []types.Action{}

	return applyContext
}

func (a *ApplyContext) Exec(trace *types.ActionTrace) {
	a.Notified = append(a.Notified, a.Receiver)
	a.execOne(trace)

	for k, r := range a.Notified {
		if k == 0 {
			continue
		}

		a.Receiver = r

		t := types.ActionTrace{}
		trace.InlineTraces = append(trace.InlineTraces, t)
		a.execOne(&trace.InlineTraces[len(trace.InlineTraces)-1])
	}

	if len(a.CfaInlineActions) > 0 || len(a.InlineActions) > 0 {
		if a.RecurseDepth >= uint32(GetDefaultConfig().MaxInlineActionDepth) {
			panic("inline action recursion depth reached")
		}
	}

	// Execute context free inlines
	for _, inlineAction := range a.CfaInlineActions {
		trace.InlineTraces = append(trace.InlineTraces, types.ActionTrace{})
		a.TrxContext.DispatchAction(&trace.InlineTraces[len(trace.InlineTraces)-1], &inlineAction, inlineAction.Account, true, a.RecurseDepth+1)
	}

	// Execute non-context free inlines
	for _, inlineAction := range a.InlineActions {
		trace.InlineTraces = append(trace.InlineTraces, types.ActionTrace{})
		a.TrxContext.DispatchAction(&trace.InlineTraces[len(trace.InlineTraces)-1], &inlineAction, inlineAction.Account, false, a.RecurseDepth+1)
	}
}

func (a *ApplyContext) execOne(trace *types.ActionTrace) {
	//start := types.Now()

	r := types.ActionReceipt{}
	r.Receiver = a.Receiver
	r.ActDigest = *crypto.Hash256(a.Act)

	trace.TrxId = a.TrxContext.ID
	trace.BlockNum = 0 // TODO: Fix
	trace.BlockTime = 0
	trace.ProducerBlockId = crypto.NewSha256Nil()
	trace.Act = *a.Act
	trace.ContextFree = a.ContextFree
	account, _ := a.Control.State.GetAccountByName(a.Receiver)
	a.Privileged = account.Privileged
	native := a.Control.FindApplyHandler(a.Receiver, a.Act.Account, a.Act.Name)

	if native != nil {
		native(a)
	}
}

func (a *ApplyContext) RequireAuthorization(account int64) error {
	for k, v := range a.Act.Authorization {
		if v.Actor == types.AccountName(account) {
			a.UsedAuthorizations[k] = true

			return nil
		}
	}

	return errors.New("missing authority of " + types.S(uint64(account)))
}
