package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/state"
)

type Notified struct {
	Receiver      types.AccountName
	ActionOrdinal int
}

type ApplyContext struct {
	Control                    *Controller
	TrxContext                 *TransactionContext
	RecurseDepth               uint32
	FirstReceiverActionOrdinal int
	ActionOrdinal              int
	Trace                      *types.ActionTrace
	Act                        *types.Action
	Receiver                   types.AccountName
	ContextFree                bool

	Privileged         bool
	UsedContestFreeApi bool
	Notified           []Notified
	InlineActions      []int
	CfaInlineActions   []int
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
	start := types.Now()
	receiverAccount, _ := a.Control.State.GetAccountByName(a.Receiver)
	a.Privileged = receiverAccount.Privileged

	if !a.ContextFree {
		native := a.Control.FindApplyHandler(a.Receiver, a.Act.Account, a.Act.Name)

		if native != nil {
			if err := native(a); err != nil {
				return err
			}
		}
	}

	trace, _ := a.TrxContext.GetActionTrace(a.ActionOrdinal)
	recvSequence, err := a.NextRecvSequence(receiverAccount)

	if err != nil {
		return err
	}

	receipt := types.ActionReceipt{
		Receiver:       a.Receiver,
		ActDigest:      *crypto.Hash256(a.Act),
		GlobalSequence: 1,
		RecvSequence:   recvSequence,
		AuthSequence:   make(map[types.AccountName]uint64),
	}

	var firstReceiverAccount *state.Account

	if a.Act.Account == receiverAccount.Name {
		firstReceiverAccount = receiverAccount
	} else {
		firstReceiverAccount, err = a.Control.State.GetAccountByName(a.Act.Account)

		if err != nil {
			return err
		}
	}

	receipt.CodeSequence = types.Vuint32(firstReceiverAccount.CodeSequence)
	receipt.AbiSequence = types.Vuint32(firstReceiverAccount.AbiSequence)

	for _, k := range a.Act.Authorization {
		receipt.AuthSequence[k.Actor], err = a.NextAuthSequence(k.Actor)

		if err != nil {
			return err
		}
	}

	trace.Receipt = receipt

	a.FinalizeTrace(trace, start)

	return nil
}

func (a *ApplyContext) FinalizeTrace(trace *types.ActionTrace, start types.TimePoint) {
	trace.Elapsed = uint64(types.Now() - start)
}

func (a *ApplyContext) RequireAuthorization(account int64) error {
	for _, v := range a.Act.Authorization {
		if v.Actor == types.AccountName(account) {
			return nil
		}
	}

	return fmt.Errorf("missing authority of %s", types.S(uint64(account)))
}

func (a *ApplyContext) IncrementActionId() {
	a.TrxContext.ActionId += 1
}

func (a *ApplyContext) NextRecvSequence(account *state.Account) (uint64, error) {
	err := a.Control.State.UpdateAccount(account, func(ra *state.Account) {
		ra.RecvSequence += 1
	})

	if err != nil {
		return 0, err
	}

	return account.RecvSequence, nil
}

func (a *ApplyContext) NextAuthSequence(accountName types.AccountName) (uint64, error) {
	account, err := a.Control.State.GetAccountByName(accountName)

	if err != nil {
		return 0, err
	}

	err = a.Control.State.UpdateAccount(account, func(ra *state.Account) {
		ra.AuthSequence += 1
	})

	if err != nil {
		return 0, err
	}

	return account.AuthSequence, nil
}
