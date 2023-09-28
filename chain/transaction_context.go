package chain

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/chain/fc"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/wasm/api"
)

var _ api.TransactionContext = &TransactionContext{}

type TransactionContext struct {
	Control                      *Controller
	PackedTrx                    *transaction.PackedTransaction
	ID                           transaction.TransactionIdType
	ApplyContextFree             bool
	Trace                        *transaction.TransactionTrace
	ActionId                     types.IdType
	ExecutedActionReceiptDigests []crypto.Sha256
	Session                      *state.Session
	AuthorizationManager         AuthorizationManager
	Deadline                     time.TimePoint
	BilledCpuTimeUs              int64
	ExplicitBilledCpuTime        bool

	Published time.TimePoint

	isInitialized bool
	start         time.TimePoint
	pseudoStart   time.TimePoint
	billedTime    time.Microseconds
	//billingTimerDurationLimit core.Microseconds
	objectiveDurationLimit    time.Microseconds
	deadline                  time.TimePoint
	deadlineExceptionCode     int64
	billingTimerExceptionCode int64
	isInput                   bool
}

func NewTransactionContext(control *Controller, s *state.Session, t *transaction.PackedTransaction, trxId transaction.TransactionIdType, block *state.Block) *TransactionContext {
	tc := TransactionContext{
		Control:               control,
		PackedTrx:             t,
		Session:               s,
		AuthorizationManager:  *NewAuthorizationManager(control, s),
		Deadline:              time.MaxTimePoint(),
		BilledCpuTimeUs:       0,
		ExplicitBilledCpuTime: false,

		isInitialized:             false,
		start:                     time.Now(),
		pseudoStart:               time.Now(),
		billedTime:                time.Microseconds(0),
		deadline:                  time.MaxTimePoint(),
		deadlineExceptionCode:     BlockCpuUsageExceededException{}.Code(),
		billingTimerExceptionCode: BlockCpuUsageExceededException{}.Code(),
	}

	tc.Trace = &transaction.TransactionTrace{
		Hash:         trxId,
		BlockNum:     uint64(block.Header.BlockNum()),
		BlockTime:    block.Header.Timestamp.ToTimePoint(),
		ActionTraces: make([]transaction.ActionTrace, 0),
	}

	return &tc
}

func (t *TransactionContext) InitForImplicitTransaction(initialNetUsage uint64) error {
	transaction, err := t.PackedTrx.GetTransaction()

	if err != nil {
		return err
	}

	if len(transaction.TransactionExtensions) > 0 {
		return fmt.Errorf("no transaction extensions supported yet for implicit transactions")
	}

	t.Published = t.Control.PendingBlockTime()

	return t.Init(initialNetUsage)
}

func (t *TransactionContext) InitForInputTransaction(packedTrxUnprunableSize uint64, packedTrxPrunableSize uint64) error {
	transaction, err := t.PackedTrx.GetTransaction()

	if err != nil {
		return err
	}

	if len(transaction.TransactionExtensions) > 0 {
		return fmt.Errorf("no transaction extensions supported yet for input transactions")
	}

	cfg, err := t.Session.FindGlobalPropertyObject(0)

	if err != nil {
		return err
	}

	discountedSizeForPrunedData := packedTrxPrunableSize

	if cfg.Configuration.ContextFreeDiscountNetUsageDen > 0 && cfg.Configuration.ContextFreeDiscountNetUsageNum < cfg.Configuration.ContextFreeDiscountNetUsageDen {
		discountedSizeForPrunedData *= uint64(cfg.Configuration.ContextFreeDiscountNetUsageNum)
		discountedSizeForPrunedData = (discountedSizeForPrunedData + uint64(cfg.Configuration.ContextFreeDiscountNetUsageDen) - 1) / uint64(cfg.Configuration.ContextFreeDiscountNetUsageDen) // rounds up
	}

	initialNetUsage := uint64(cfg.Configuration.BasePerTransactionNetUsage) + packedTrxUnprunableSize + discountedSizeForPrunedData
	t.Published = t.Control.PendingBlockTime()
	t.isInput = true

	if err := t.Init(initialNetUsage); err != nil {
		return err
	}

	id, _ := t.PackedTrx.ID()

	return t.RecordTransaction(*id, transaction.Expiration)
}

func (t *TransactionContext) Init(initialNetUsage uint64) error {
	if t.isInitialized {
		return fmt.Errorf("cannot initialize twice")
	}

	transaction, err := t.PackedTrx.GetTransaction()

	if err != nil {
		return err
	}

	if transaction.DelaySec > 0 {
		return fmt.Errorf("deferred transactions are deprecated")
	}

	t.objectiveDurationLimit = time.Microseconds(config.MaxBlockCpuUsage)
	t.deadline = t.start + time.TimePoint(t.objectiveDurationLimit)

	// Possibly lower objective_duration_limit to the maximum cpu usage a transaction is allowed to be billed
	if config.MaxTransactionCpuUsage <= uint32(t.objectiveDurationLimit.Count()) {
		t.objectiveDurationLimit = time.Microseconds(config.MaxTransactionCpuUsage)
		t.billingTimerExceptionCode = TxCpuUsageExceededException{}.Code()
		t.deadline = t.start + time.TimePoint(t.objectiveDurationLimit)
	}

	// Possibly lower objective_duration_limit to optional limit set in transaction header
	if transaction.MaxCpuUsageMS > 0 {
		trxSpecifiedCpuUsageLimit := time.Milliseconds(int64(transaction.MaxCpuUsageMS))

		if trxSpecifiedCpuUsageLimit <= t.objectiveDurationLimit {
			t.objectiveDurationLimit = trxSpecifiedCpuUsageLimit
			t.billingTimerExceptionCode = TxCpuUsageExceededException{}.Code()
			t.deadline = t.start + time.TimePoint(t.objectiveDurationLimit)
		}
	}

	// Check if deadline is limited by caller-set deadline (only change deadline if billed_cpu_time_us is not set)
	if t.Deadline < t.deadline {
		t.deadline = t.Deadline
		t.deadlineExceptionCode = DeadlineException{}.Code()
	} else {
		t.deadlineExceptionCode = t.billingTimerExceptionCode
	}

	if err := t.CheckTime(); err != nil {
		return err
	}

	t.isInitialized = true

	return nil
}

func (t *TransactionContext) Exec() error {
	if !t.isInitialized {
		return fmt.Errorf("must first initialize")
	}

	transaction, err := t.PackedTrx.GetTransaction()

	if err != nil {
		return err
	}

	if t.ApplyContextFree {
		for _, act := range transaction.ContextFreeActions {
			t.ScheduleAction(*act, act.Account, true, 0)
		}
	}

	if transaction.DelaySec == 0 {
		for _, act := range transaction.Actions {
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

func (t *TransactionContext) ScheduleActionFromOrdinal(actionOrdinal int, receiver name.AccountName, contextFree bool, creatorActionOrdinal int) (int, error) {
	newActionOrdinal := len(t.Trace.ActionTraces) + 1
	trace, err := t.GetActionTrace(actionOrdinal)

	if err != nil {
		return 0, err
	}

	t.Trace.ActionTraces = append(t.Trace.ActionTraces, *transaction.NewActionTrace(t.Trace, trace.Action, receiver, contextFree, fc.UnsignedInt(newActionOrdinal), fc.UnsignedInt(creatorActionOrdinal)))

	return newActionOrdinal, nil
}

func (t *TransactionContext) ScheduleAction(action transaction.Action, receiver name.AccountName, contextFree bool, creatorActionOrdinal int) int {
	newActionOrdinal := len(t.Trace.ActionTraces) + 1
	actionTrace := transaction.NewActionTrace(t.Trace, action, receiver, contextFree, fc.UnsignedInt(newActionOrdinal), fc.UnsignedInt(creatorActionOrdinal))
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

func (t *TransactionContext) GetActionTrace(actionOrdinal int) (*transaction.ActionTrace, error) {
	if actionOrdinal < 0 || actionOrdinal > len(t.Trace.ActionTraces) {
		return nil, fmt.Errorf("action_ordinal %v is outside allowed range [1,%v]", actionOrdinal, len(t.Trace.ActionTraces))
	}

	return &t.Trace.ActionTraces[actionOrdinal-1], nil
}

func (t *TransactionContext) Finalize() error {
	now := time.Now()
	t.Trace.Elapsed = time.Microseconds(now - t.start)
	t.UpdateBilledCpuTime(now)

	return nil
}

func (t *TransactionContext) Commit() error {
	// TODO: Handle
	/* if err := t.State.PutTransaction(t.Trace); err != nil {
		return err
	} */

	if err := t.Session.Commit(); err != nil {
		return err
	}

	return nil
}

func (t *TransactionContext) CheckTime() error {
	now := time.Now()

	if now < t.deadline {
		return nil
	}

	duration := time.Microseconds(now - t.pseudoStart).Count()

	if t.deadlineExceptionCode == (DeadlineException{}).Code() {
		return fmt.Errorf("deadline exceeded %dus", duration)
	} else if t.deadlineExceptionCode == (BlockCpuUsageExceededException{}).Code() {
		return fmt.Errorf("not enough time left in block to complete executing transaction %dus", duration)
	} else if t.deadlineExceptionCode == (TxCpuUsageExceededException{}).Code() {
		return fmt.Errorf("transaction was executing for too long %dus", duration)
	} else if t.deadlineExceptionCode == (LeewayDeadlineException{}).Code() {
		return fmt.Errorf("the transaction was unable to complete by deadline, but it is possible it could have succeeded if it were allowed to run to completion")
	}

	return nil
}

func (t *TransactionContext) PauseBillingTimer() {
	if t.pseudoStart == 0 {
		return
	}

	now := time.Now()
	t.billedTime = time.Microseconds(now - t.pseudoStart)
	t.deadlineExceptionCode = DeadlineException{}.Code()
	t.pseudoStart = time.TimePoint(0)
}

func (t *TransactionContext) ResumeBillingTimer() {
	if t.pseudoStart != 0 {
		return
	}

	now := time.Now()
	t.pseudoStart = now - time.TimePoint(t.billedTime)
}

func (t *TransactionContext) UpdateBilledCpuTime(now time.TimePoint) {
	billed := now - t.pseudoStart

	if billed < time.TimePoint(config.MinTransactionCpuUsage) {
		t.BilledCpuTimeUs = int64(config.MinTransactionCpuUsage)
	} else {
		t.BilledCpuTimeUs = int64(billed)
	}
}

func (t *TransactionContext) RecordTransaction(id transaction.TransactionIdType, expire time.TimePointSec) error {
	return t.Session.CreateTransactionObject(&transaction.TransactionObject{
		TrxId:      id,
		Expiration: expire,
	})
}

func (t *TransactionContext) GetPublicationTime() time.TimePoint {
	return t.Published
}
