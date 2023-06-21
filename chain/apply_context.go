package chain

import (
	"context"
	"fmt"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/MetalBlockchain/antelopevm/utils"
	"github.com/MetalBlockchain/antelopevm/wasm"
	wasmApi "github.com/MetalBlockchain/antelopevm/wasm/api"
	"github.com/dgraph-io/badger/v3"
)

var (
	_ wasmApi.ApplyContext = &applyContext{}
	_ ApplyContext         = &applyContext{}

	errDatabaseAccessViolation = fmt.Errorf("db access violation")
)

type Notified struct {
	Receiver      name.AccountName
	ActionOrdinal int
}

type ApplyContext interface {
	GetSession() *state.Session
	GetAuthorizationManager() *AuthorizationManager

	GetAction() core.Action
	RequireAuthorization(name.AccountName) error
}

type applyContext struct {
	Control                    *Controller
	TrxContext                 *TransactionContext
	RecurseDepth               uint32
	FirstReceiverActionOrdinal int
	ActionOrdinal              int
	Trace                      *core.ActionTrace
	Act                        *core.Action
	Receiver                   name.AccountName
	ContextFree                bool
	ConsoleOutput              string

	Privileged         bool
	UsedContestFreeApi bool
	Notified           []Notified
	InlineActions      []int
	CfaInlineActions   []int
	AccountRamDeltas   map[name.AccountName]int64
	ActionReturnValue  []byte

	KeyValueCache *IteratorCache

	Session       *state.Session
	Authorization *AuthorizationManager
	//AccountRamDeltas   []types.AccountDelta

	Idx64         *Idx64
	Idx128        *Idx128
	Idx256        *Idx256
	IdxDouble     *IdxDouble
	IdxLongDouble *IdxLongDouble
}

func NewApplyContext(trxContext *TransactionContext, actionOrdinal int, recurseDepth uint32) (*applyContext, error) {
	trace, err := trxContext.GetActionTrace(actionOrdinal)

	if err != nil {
		return nil, err
	}

	applyContext := &applyContext{
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
		AccountRamDeltas:   make(map[name.Name]int64),
		ActionReturnValue:  make([]byte, 0),

		KeyValueCache: NewIteratorCache(),
		Session:       trxContext.Session,
		Authorization: &trxContext.AuthorizationManager,
	}

	applyContext.Idx64 = &Idx64{Context: applyContext}
	applyContext.Idx128 = &Idx128{Context: applyContext}
	applyContext.Idx256 = &Idx256{Context: applyContext}
	applyContext.IdxDouble = &IdxDouble{Context: applyContext}
	applyContext.IdxLongDouble = &IdxLongDouble{Context: applyContext}

	return applyContext, nil
}

func (a *applyContext) GetSession() *state.Session {
	return a.Session
}

func (a *applyContext) GetAuthorizationManager() *AuthorizationManager {
	return a.Authorization
}

func (a *applyContext) GetIdx64() *Idx64 {
	return a.Idx64
}

func (a *applyContext) Exec() error {
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
		if a.RecurseDepth >= uint32(config.MaxInlineActionDepth) {
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

func (a *applyContext) execOne() error {
	start := core.Now()
	receiverAccount, err := a.Session.FindAccountByName(a.Receiver)

	if err != nil {
		return fmt.Errorf("could not find receiver account: %v", err)
	}

	a.Privileged = receiverAccount.Privileged

	if !a.ContextFree {
		native := a.Control.FindApplyHandler(a.Receiver, a.Act.Account, a.Act.Name)

		if native != nil {
			if err := native(a); err != nil {
				return err
			}
		}

		if receiverAccount.Code.Size() > 0 && !(a.Act.Account == config.SystemAccountName && a.Act.Name == name.StringToName("setcode") && a.Receiver == config.SystemAccountName) {
			// Check contract blacklist
			if err := a.Control.CheckContractList(receiverAccount.Name); err != nil {
				return err
			}

			a.TrxContext.PauseBillingTimer()
			module := wasm.NewWasmExecutionContext(context.Background(), a.Control, a, a.Authorization, a.GetMutableResourceLimitsManager(), a.Idx64, a.Idx128, a.Idx256, a.IdxDouble, a.IdxLongDouble)

			if err := module.Initialize(); err != nil {
				return err
			}

			if err := module.Exec(receiverAccount.Code); err != nil {
				return err
			}

			a.TrxContext.ResumeBillingTimer()
		}
	}

	trace, _ := a.TrxContext.GetActionTrace(a.ActionOrdinal)
	recvSequence, err := a.NextRecvSequence(receiverAccount.Name)

	if err != nil {
		return err
	}

	receipt := core.ActionReceipt{
		Receiver:       a.Receiver,
		ActDigest:      *crypto.Hash256(a.Act),
		GlobalSequence: 1,
		RecvSequence:   *recvSequence,
		AuthSequence:   core.NewAuthSequenceSet(),
	}

	var firstReceiverAccount *account.Account

	if a.Act.Account == receiverAccount.Name {
		firstReceiverAccount = receiverAccount
	} else {
		firstReceiverAccount, err = a.Session.FindAccountByName(a.Act.Account)

		if err != nil {
			return fmt.Errorf("could not find account by name %s", a.Act.Account)
		} else if firstReceiverAccount == nil {
			return fmt.Errorf("could not find account by name %s", a.Act.Account)
		}
	}

	receipt.CodeSequence = core.Vuint32(firstReceiverAccount.CodeSequence)
	receipt.AbiSequence = core.Vuint32(firstReceiverAccount.AbiSequence)

	for _, k := range a.Act.Authorization {
		if authSequence, err := a.NextAuthSequence(k.Actor); err == nil {
			receipt.AuthSequence.Set(k.Actor, authSequence)
		} else {
			return err
		}
	}

	trace.Receipt = receipt
	a.TrxContext.ExecutedActionReceiptDigests = append(a.TrxContext.ExecutedActionReceiptDigests, receipt.Digest())

	a.FinalizeTrace(trace, start)

	return nil
}

func (a *applyContext) ExecuteInline(action core.Action) error {
	if _, err := a.Session.FindAccountByName(action.Account); err != nil {
		return fmt.Errorf("inline action's code account %s does not exist", action.Account)
	}

	sendToSelf := (action.Account == a.Receiver)
	inheritParentAuthorizations := (sendToSelf && a.Receiver == a.Act.Account)
	inheritedAuthorizations := authority.NewPermissionLevelSet(len(action.Authorization))

	for _, auth := range action.Authorization {
		if _, err := a.Session.FindAccountByName(auth.Actor); err != nil {
			return fmt.Errorf("inline action's authorizing actor %s does not exist", auth.Actor)
		}

		if _, err := a.Authorization.GetPermission(auth); err != nil {
			return fmt.Errorf("inline action's authorizations include a non-existent permission: %v", auth)
		}

		if inheritParentAuthorizations {
			for _, parentAuth := range a.Act.Authorization {
				if auth.Actor == parentAuth.Actor && auth.Permission == parentAuth.Permission {
					inheritedAuthorizations.Insert(auth)
				}
			}
		}
	}

	// TODO: Add inline action size
	if !a.Privileged {
		auth := authority.PermissionLevel{Actor: a.Receiver, Permission: config.EosioCodeName}

		if err := a.Authorization.CheckAuthorization([]*core.Action{&action}, ecc.NewPublicKeySet(0), []authority.PermissionLevel{auth}, false, inheritedAuthorizations); err != nil {
			return fmt.Errorf("authorization failure with inline action")
		}
	}

	if ordinal, err := a.ScheduleActionByAction(action, action.Account, false); err == nil {
		a.InlineActions = append(a.InlineActions, *ordinal)
	} else {
		return err
	}

	return nil
}

func (a *applyContext) FinalizeTrace(trace *core.ActionTrace, start core.TimePoint) {
	trace.Elapsed = uint64(core.Now() - start)
	trace.AccountRamDeltas = make([]core.RamDelta, len(trace.AccountRamDeltas))

	for account, delta := range a.AccountRamDeltas {
		trace.AccountRamDeltas = append(trace.AccountRamDeltas, core.RamDelta{
			Account: account,
			Delta:   delta,
		})
	}
}

func (a *applyContext) ScheduleAction(ordinalOfActionToSchedule int, receiver name.AccountName, contextFree bool) (int, error) {
	scheduledActionOrdinal, err := a.TrxContext.ScheduleActionFromOrdinal(ordinalOfActionToSchedule, receiver, contextFree, a.ActionOrdinal)

	if err != nil {
		return -1, err
	}

	actionTrace, err := a.TrxContext.GetActionTrace(a.ActionOrdinal)

	if err != nil {
		return -1, err
	}

	a.Act = &actionTrace.Action

	return scheduledActionOrdinal, nil
}

func (a *applyContext) ScheduleActionByAction(actionToSchedule core.Action, receiver name.AccountName, contextFree bool) (*int, error) {
	scheduledActionOrdinal := a.TrxContext.ScheduleAction(actionToSchedule, receiver, contextFree, a.ActionOrdinal)
	actionTrace, err := a.TrxContext.GetActionTrace(a.ActionOrdinal)

	if err != nil {
		return nil, err
	}

	a.Act = &actionTrace.Action

	return &scheduledActionOrdinal, nil
}

func (a *applyContext) RequireAuthorization(account name.AccountName) error {
	for _, v := range a.Act.Authorization {
		if v.Actor == account {
			return nil
		}
	}

	return fmt.Errorf("missing authority of %s", name.NameToString(uint64(account)))
}

func (a *applyContext) RequireAuthorizationWithPermission(account name.AccountName, permission name.PermissionName) error {
	for _, v := range a.Act.Authorization {
		if v.Actor == account {
			if v.Permission == permission {
				return nil
			}
		}
	}

	return fmt.Errorf("missing authority of %s", name.NameToString(uint64(account)))
}

func (a *applyContext) HasAuthorization(account name.AccountName) bool {
	for _, v := range a.Act.Authorization {
		if v.Actor == account {
			return true
		}
	}

	return false
}

func (a *applyContext) FindAccount(account name.AccountName) (*account.Account, error) {
	return a.Session.FindAccountByName(account)
}

func (a *applyContext) IsAccount(account name.AccountName) bool {
	if account, err := a.Session.FindAccountByName(account); account != nil && err == nil {
		return true
	}

	return false
}

func (a *applyContext) HasRecipient(account name.AccountName) bool {
	for _, notified := range a.Notified {
		if notified.Receiver == account {
			return true
		}
	}

	return false
}

func (a *applyContext) RequireRecipient(recipient name.AccountName) error {
	if !a.HasRecipient(recipient) {
		if ordinal, err := a.ScheduleAction(a.ActionOrdinal, recipient, false); err == nil {
			a.Notified = append(a.Notified, Notified{
				Receiver:      recipient,
				ActionOrdinal: ordinal,
			})
		} else {
			return err
		}
	}

	return nil
}

func (a *applyContext) GetSender() (*name.ActionName, error) {
	trace, err := a.TrxContext.GetActionTrace(a.ActionOrdinal)

	if err != nil {
		return nil, err
	}

	if trace.CreatorActionOrdinal > 0 {
		if creatorTrace, err := a.TrxContext.GetActionTrace(int(trace.CreatorActionOrdinal)); err == nil {
			return &creatorTrace.Receiver, nil
		}
	}

	return nil, nil
}

func (a *applyContext) IncrementActionId() {
	a.TrxContext.ActionId += 1
}

func (a *applyContext) NextRecvSequence(accountName name.AccountName) (*uint64, error) {
	account, err := a.Session.FindAccountByName(accountName)

	if err != nil {
		return nil, fmt.Errorf("could not find account %s", accountName)
	}

	err = a.Session.ModifyAccount(account, func() {
		account.RecvSequence = account.RecvSequence + 1
	})

	if err != nil {
		return nil, err
	}

	return &account.RecvSequence, nil
}

func (a *applyContext) NextAuthSequence(accountName name.AccountName) (uint64, error) {
	account, err := a.Session.FindAccountByName(accountName)

	if err != nil {
		return 0, fmt.Errorf("could not find account %s", accountName)
	}

	err = a.Session.ModifyAccount(account, func() {
		account.AuthSequence += 1
	})

	if err != nil {
		return 0, err
	}

	return account.AuthSequence, nil
}

func (a *applyContext) GetAction() core.Action {
	return *a.Act
}

func (a *applyContext) GetReceiver() name.AccountName { return a.Receiver }

func (a *applyContext) FindI64(code name.AccountName, scope name.ScopeName, tableName name.TableName, primaryKey uint64) int {
	table, _ := a.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if table == nil {
		return -1
	}

	endIterator := a.KeyValueCache.cacheTable(table)
	keyValue, _ := a.Session.FindKeyValueByScopePrimary(table.ID, primaryKey)

	if keyValue == nil {
		return endIterator
	}

	return a.KeyValueCache.add(keyValue)
}

func (a *applyContext) StoreI64(code name.AccountName, scope name.ScopeName, tableName name.TableName, payer name.AccountName, primaryKey uint64, buffer []byte) (int, error) {
	table, err := a.Session.FindOrCreateTable(code, scope, tableName, payer)

	if err != nil {
		return 0, fmt.Errorf("failed to find or create table: %v", err)
	}

	keyValue := &core.KeyValue{
		TableID:    table.ID,
		PrimaryKey: primaryKey,
		Payer:      payer,
		Value:      buffer,
	}
	err = a.Session.CreateKeyValue(keyValue)

	if err != nil {
		return 0, fmt.Errorf("failed to insert kv object: %v", err)
	}

	if err := a.Session.ModifyTable(table, func() {
		table.Count = table.Count + 1
	}); err != nil {
		return 0, fmt.Errorf("failed to update table: %v", err)
	}

	billableSize := uint64(len(buffer)) + config.GetBillableSize("key_value_object")
	a.UpdateDatabaseUsage(payer, int64(billableSize))

	a.KeyValueCache.cacheTable(table)
	iterator := a.KeyValueCache.add(keyValue)

	return iterator, nil
}

func (a *applyContext) GetI64(iterator int, buffer []byte, bufferSize int) (int, error) {
	obj := (a.KeyValueCache.get(iterator)).(*core.KeyValue)
	size := len(obj.Value)

	if bufferSize == 0 {
		return size, nil
	}

	copySize := utils.MinInt(bufferSize, size)
	copy(buffer[0:copySize], obj.Value[0:copySize])

	return copySize, nil
}

func (a *applyContext) UpdateI64(iterator int, payer name.AccountName, buffer []byte, bufferSize int) error {
	obj := (a.KeyValueCache.get(iterator)).(*core.KeyValue)
	table := a.KeyValueCache.tableCache[obj.TableID]

	if table.table.Code != a.Receiver {
		return fmt.Errorf("db access violation")
	}

	overhead := config.GetBillableSize("key_value_object")
	oldSize := int64(obj.Value.Size()) + int64(overhead)
	newSize := int64(bufferSize) + int64(overhead)

	if payer.Empty() {
		payer = obj.Payer
	}

	if obj.Payer != payer {
		// refund the existing payer
		if err := a.UpdateDatabaseUsage(obj.Payer, -oldSize); err != nil {
			return err
		}
		// charge the new payer
		if err := a.UpdateDatabaseUsage(payer, newSize); err != nil {
			return err
		}
	} else if oldSize != newSize {
		// charge/refund the existing payer the difference
		if err := a.UpdateDatabaseUsage(obj.Payer, newSize-oldSize); err != nil {
			return err
		}
	}

	if err := a.Session.ModifyKeyValue(obj, func() {
		obj.Value = buffer
		obj.Payer = payer
	}); err != nil {
		return fmt.Errorf("failed to update table: %v", err)
	}

	return nil
}

func (a *applyContext) RemoveI64(iterator int) error {
	obj := (a.KeyValueCache.get(iterator)).(*core.KeyValue)
	table, err := a.KeyValueCache.getTable(obj.TableID)

	if err != nil {
		return err
	} else if table.Code != a.Receiver {
		return errDatabaseAccessViolation
	}

	ramDelta := int64(uint64(len(obj.Value))+config.GetBillableSize("key_value_object")) * -1
	a.UpdateDatabaseUsage(obj.Payer, ramDelta)

	if err := a.Session.ModifyTable(table, func() {
		table.Count = table.Count - 1
	}); err != nil {
		return fmt.Errorf("failed to update table: %v", err)
	}

	if err := a.Session.RemoveKeyValue(obj); err != nil {
		return err
	}

	if table.Count == 0 {
		if err := a.RemoveTable(table); err != nil {
			return err
		}
	}

	a.KeyValueCache.remove(iterator)

	return nil
}

func (a *applyContext) NextI64(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 {
		return -1, nil // cannot increment past end iterator of table
	}

	obj := (a.KeyValueCache.get(iterator)).(*core.KeyValue) // check for iterator != -1 happens in this call
	nextKv, err := a.Session.FindNextKeyValue(obj)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return a.KeyValueCache.getEndIteratorByTableId(obj.TableID)
		}

		return -1, err
	}

	if nextKv.TableID != obj.TableID || nextKv.PrimaryKey < obj.PrimaryKey {
		return a.KeyValueCache.getEndIteratorByTableId(obj.TableID)
	}

	*primaryKey = nextKv.PrimaryKey

	return a.KeyValueCache.add(nextKv), nil
}

func (a *applyContext) PreviousI64(iterator int, primaryKey *uint64) (int, error) {
	if iterator < -1 { // is end iterator
		table, err := a.KeyValueCache.findTableByEndIterator(iterator)

		if err != nil {
			return 0, err
		}

		obj, err := a.Session.UpperboundKeyValueByScope(table)

		if err != nil {
			if err == badger.ErrKeyNotFound {
				return -1, nil
			}

			return 0, err
		} else if obj.TableID != table.ID { // Empty table
			return -1, nil
		}

		*primaryKey = obj.PrimaryKey

		return a.KeyValueCache.add(obj), nil
	}

	obj := (a.KeyValueCache.get(iterator)).(*core.KeyValue)
	nextKv, err := a.Session.FindPreviousKeyValue(obj)

	if err != nil {
		return -1, err
	} else if nextKv.TableID != obj.TableID || nextKv.PrimaryKey > obj.PrimaryKey {
		return -1, nil
	}

	*primaryKey = nextKv.PrimaryKey

	return a.KeyValueCache.add(nextKv), nil
}

func (a *applyContext) LowerboundI64(code name.AccountName, scope name.ScopeName, table name.TableName, id uint64) (int, error) {
	tab, err := a.Session.FindTableByCodeScopeTable(code, scope, table)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return -1, nil
		}

		return -1, err
	}

	endIterator := a.KeyValueCache.cacheTable(tab)
	obj, err := a.Session.LowerboundKeyValueByScopePrimary(tab, id)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return endIterator, nil
		}

		return -1, err
	} else if obj.TableID != tab.ID {
		return endIterator, nil
	}

	return a.KeyValueCache.add(obj), nil
}

func (a *applyContext) UpperboundI64(code name.AccountName, scope name.ScopeName, table name.TableName, id uint64) (int, error) {
	tab, err := a.Session.FindTableByCodeScopeTable(code, scope, table)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return -1, nil
		}

		return -1, err
	}

	endIterator := a.KeyValueCache.cacheTable(tab)
	obj, err := a.Session.UpperboundKeyValueByScopePrimary(tab, id)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return endIterator, nil
		}

		return -1, err
	} else if obj.TableID != tab.ID {
		return endIterator, nil
	}

	return a.KeyValueCache.add(obj), nil
}

func (a *applyContext) EndI64(code name.AccountName, scope name.ScopeName, table name.TableName) (int, error) {
	tab, err := a.Session.FindTableByCodeScopeTable(code, scope, table)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return -1, nil
		}

		return -1, err
	}

	return a.KeyValueCache.cacheTable(tab), nil
}

func (a *applyContext) RemoveTable(table *core.Table) error {
	ramDelta := int64(config.GetBillableSize("table_id_object")) * -1

	if err := a.UpdateDatabaseUsage(table.Payer, ramDelta); err != nil {
		return err
	}

	return a.Session.RemoveTable(table)
}

func (a *applyContext) FindOrCreateTable(code name.AccountName, scope name.ScopeName, tableName name.TableName, payer name.AccountName) (*core.Table, error) {
	table, err := a.Session.FindTableByCodeScopeTable(code, scope, tableName)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			ramDelta := int64(config.GetBillableSize("table_id_object"))

			if err := a.UpdateDatabaseUsage(payer, ramDelta); err != nil {
				return nil, err
			}

			table = &core.Table{
				Code:  code,
				Scope: scope,
				Table: tableName,
				Payer: payer,
			}

			err = a.Session.CreateTable(table)

			if err != nil {
				return nil, err
			}

			return table, nil
		}

		return nil, err
	}

	return table, nil
}

func (a *applyContext) UpdateDatabaseUsage(payer name.AccountName, delta int64) error {
	if delta > 0 {
		if a.Receiver != payer {
			return fmt.Errorf("cannot charge RAM to other accounts during notify")
		}
	}

	a.AddRamUsage(payer, delta)

	return nil
}

func (a *applyContext) AddRamUsage(account name.AccountName, delta int64) {
	a.AccountRamDeltas[account] += delta
}

func (a *applyContext) CheckAuthorization(actions []core.Action, providedKeys ecc.PublicKeySet) error {
	//a.Authorization.CheckAuthorization()
	return nil
}

func (a *applyContext) ConsoleAppend(value string) {
	a.ConsoleOutput += value
}

func (a *applyContext) SetActionReturnValue(value []byte) {
	a.ActionReturnValue = value
}

func (a *applyContext) GetPackedTransaction() *core.PackedTransaction {
	return a.TrxContext.PackedTrx
}

func (a *applyContext) IsContextPrivileged() bool {
	return a.Privileged
}

func (a *applyContext) IsPrivileged(name name.AccountName) (bool, error) {
	account, err := a.Session.FindAccountByName(name)

	if err != nil {
		return false, err
	}

	return account.Privileged, nil
}

func (a *applyContext) SetPrivileged(name name.AccountName, privileged bool) error {
	account, err := a.Session.FindAccountByName(name)

	if err != nil {
		return err
	}

	return a.Session.ModifyAccount(account, func() {
		account.Privileged = privileged
	})
}

func (a *applyContext) GetMutableResourceLimitsManager() *ResourceLimitsManager {
	return a.Control.GetResourceLimitsManager(a.Session)
}
