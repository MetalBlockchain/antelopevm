package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/inconshreveable/log15"
)

type Notified struct {
	Receiver      core.AccountName
	ActionOrdinal int
}

type ApplyContext struct {
	Control                    *Controller
	TrxContext                 *TransactionContext
	RecurseDepth               uint32
	FirstReceiverActionOrdinal int
	ActionOrdinal              int
	Trace                      *core.ActionTrace
	Act                        *core.Action
	Receiver                   core.AccountName
	ContextFree                bool

	Privileged         bool
	UsedContestFreeApi bool
	Notified           []Notified
	InlineActions      []int
	CfaInlineActions   []int

	State         state.State
	Authorization *AuthorizationManager
	//AccountRamDeltas   []types.AccountDelta
}

func NewApplyContext(trxContext *TransactionContext, actionOrdinal int, recurseDepth uint32) (*ApplyContext, error) {
	trace, err := trxContext.GetActionTrace(actionOrdinal)

	if err != nil {
		return nil, err
	}

	applyContext := &ApplyContext{
		Control:                    trxContext.Control,
		TrxContext:                 trxContext,
		RecurseDepth:               recurseDepth,
		FirstReceiverActionOrdinal: actionOrdinal,
		ActionOrdinal:              actionOrdinal,

		Trace:              trace,
		Act:                &trace.Action,
		Receiver:           trace.Receiver,
		ContextFree:        trace.ContextFree,
		Privileged:         false,
		UsedContestFreeApi: false,

		State:         trxContext.State,
		Authorization: &trxContext.AuthorizationManager,
	}

	return applyContext, nil
}

func (a *ApplyContext) Exec() error {
	a.Notified = append(a.Notified, Notified{a.Receiver, a.ActionOrdinal})

	if err := a.execOne(); err != nil {
		return err
	}

	a.IncrementActionId()

	for i := 1; i < len(a.Notified); i++ {
		a.ActionOrdinal = a.Notified[i].ActionOrdinal
		a.Receiver = a.Notified[i].Receiver

		if err := a.execOne(); err != nil {
			return err
		}

		a.IncrementActionId()
	}

	if len(a.CfaInlineActions) > 0 || len(a.InlineActions) > 0 {
		if a.RecurseDepth >= uint32(GetDefaultConfig().MaxInlineActionDepth) {
			return fmt.Errorf("inline action recursion depth reached")
		}
	}

	// Execute context free inlines
	for _, ordinal := range a.CfaInlineActions {
		a.TrxContext.ExecuteAction(ordinal, a.RecurseDepth+1)
	}

	// Execute non-context free inlines
	for _, ordinal := range a.InlineActions {
		a.TrxContext.ExecuteAction(ordinal, a.RecurseDepth+1)
	}

	return nil
}

func (a *ApplyContext) execOne() error {
	start := core.Now()
	receiverAccount, _ := a.State.GetAccountByName(a.Receiver)
	a.Privileged = receiverAccount.Privileged

	if !a.ContextFree {
		native := a.Control.FindApplyHandler(a.Receiver, a.Act.Account, a.Act.Name)

		if native != nil {
			if err := native(a); err != nil {
				log15.Error("failed to exec native", "error", err)
				return err
			}
		}
	}

	trace, _ := a.TrxContext.GetActionTrace(a.ActionOrdinal)
	recvSequence := a.NextRecvSequence(receiverAccount)
	receipt := core.ActionReceipt{
		Receiver:       a.Receiver,
		ActDigest:      *crypto.Hash256(a.Act),
		GlobalSequence: 1,
		RecvSequence:   recvSequence,
		AuthSequence:   make(core.AuthSequenceSet, 0),
	}

	var firstReceiverAccount *core.Account
	var err error

	if a.Act.Account == receiverAccount.Name {
		firstReceiverAccount = receiverAccount
	} else {
		firstReceiverAccount, err = a.State.GetAccountByName(a.Act.Account)

		if err != nil {
			return err
		}
	}

	receipt.CodeSequence = core.Vuint32(firstReceiverAccount.CodeSequence)
	receipt.AbiSequence = core.Vuint32(firstReceiverAccount.AbiSequence)

	for _, k := range a.Act.Authorization {
		receipt.AuthSequence.Set(k.Actor, a.NextAuthSequence(k.Actor))
	}

	trace.Receipt = receipt
	a.TrxContext.ExecutedActionReceiptDigests = append(a.TrxContext.ExecutedActionReceiptDigests, receipt.Digest())

	a.FinalizeTrace(trace, start)

	return nil
}

func (a *ApplyContext) FinalizeTrace(trace *core.ActionTrace, start core.TimePoint) {
	trace.Elapsed = uint64(core.Now() - start)
}

func (a *ApplyContext) RequireAuthorization(account int64) error {
	for _, v := range a.Act.Authorization {
		if v.Actor == core.AccountName(account) {
			return nil
		}
	}

	return fmt.Errorf("missing authority of %s", core.NameToString(uint64(account)))
}

func (a *ApplyContext) IncrementActionId() {
	a.TrxContext.ActionId += 1
}

func (a *ApplyContext) NextRecvSequence(account *core.Account) uint64 {
	err := a.State.UpdateAccount(account, func(ra *core.Account) {
		ra.RecvSequence += 1
	})

	if err != nil {
		panic(err)
	}

	return account.RecvSequence
}

func (a *ApplyContext) NextAuthSequence(accountName core.AccountName) uint64 {
	account, err := a.State.GetAccountByName(accountName)

	if err != nil {
		panic(err)
	}

	err = a.State.UpdateAccount(account, func(ra *core.Account) {
		ra.AuthSequence += 1
	})

	if err != nil {
		panic(err)
	}

	return account.AuthSequence
}
