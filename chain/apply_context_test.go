package chain

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestFindI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:    trxContext,
		Control:       controller,
		Session:       session,
		KeyValueCache: NewIteratorCache(),
	}
	// Fetching a KV from a non existing table should return -1
	result := applyContext.FindI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 1)
	assert.Equal(t, -1, result)
	_, err = applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	// Stored object should return 0 iterator
	result = applyContext.FindI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 1)
	assert.Equal(t, 0, result)
	// Non-existing KV should return end iterator -2
	result = applyContext.FindI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 2)
	assert.Equal(t, -2, result)
}

func TestNextI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:    trxContext,
		Control:       controller,
		Session:       session,
		KeyValueCache: NewIteratorCache(),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)
	iterator2, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 2, []byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, iterator2, 1)

	// Next I64
	var primaryKey uint64
	result, err := applyContext.NextI64(iterator1, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, primaryKey, uint64(2))
	assert.Equal(t, result, iterator2)

	primaryKey = 0
	result, err = applyContext.NextI64(iterator2, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, primaryKey, uint64(0), "Second next_i64 primary key should be 0")
	assert.Equal(t, -2, result)
}

func TestPreviousI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:    trxContext,
		Control:       controller,
		Session:       session,
		KeyValueCache: NewIteratorCache(),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)
	iterator2, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 2, []byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, iterator2, 1)
	iterator3, err := applyContext.StoreI64(name.StringToName("eosio"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, iterator3, 2)

	// Previous I64 on the first row should return -1
	var primaryKey uint64
	result, err := applyContext.PreviousI64(iterator1, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), primaryKey)
	assert.Equal(t, -1, result)

	// Previous I64 on the first row should return 0 and primary key 1
	primaryKey = 0
	result, err = applyContext.PreviousI64(iterator2, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), primaryKey, "Second previous_i64 primary key should be 0")
	assert.Equal(t, 0, result)

	// Finding previous I64 by end iterator should return last row
	primaryKey = 0
	result, err = applyContext.PreviousI64(-2, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), primaryKey, "Second previous_i64 primary key should be 2")
	assert.Equal(t, iterator2, result, "Returned iterator should equal 1")

	// Finding previous I64 of table with one KV should return -1
	primaryKey = 0
	result, err = applyContext.PreviousI64(iterator3, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), primaryKey, "Primary key should be 0")
	assert.Equal(t, -1, result, "Returned iterator should equal -1")

	// Finding previous I64 by end iterator of table with one KV should return only key in there
	primaryKey = 0
	result, err = applyContext.PreviousI64(-3, &primaryKey)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), primaryKey, "Primary key should be 1")
	assert.Equal(t, iterator3, result, "Returned iterator should equal 2")
}

func TestGetI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:    trxContext,
		Control:       controller,
		Session:       session,
		KeyValueCache: NewIteratorCache(),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)

	// Get I64 should return the actual size of the buffer if provided buffer size is 0
	result, err := applyContext.GetI64(iterator1, make([]byte, 0), 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)

	// Get I64 should return buffer size of 1 if buffer size of 1 is sent
	result, err = applyContext.GetI64(iterator1, make([]byte, 1), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// Get I64 should return buffer size of 2
	result, err = applyContext.GetI64(iterator1, make([]byte, 2), 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)

	// Get I64 should return buffer size of 2 if KV buffer size is less than provider buffer size
	result, err = applyContext.GetI64(iterator1, make([]byte, 3), 3)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)
}

func TestRemoveI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:       trxContext,
		Control:          controller,
		Session:          session,
		KeyValueCache:    NewIteratorCache(),
		Receiver:         name.StringToName("eosio.token"),
		AccountRamDeltas: make(map[name.Name]int64),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)

	// Remove I64 should not return an error
	err = applyContext.RemoveI64(iterator1)
	assert.NoError(t, err)

	iterator1, err = applyContext.StoreI64(name.StringToName("eosio.token2"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)

	// Remove I64 should throw error if receiver doesn't match code
	err = applyContext.RemoveI64(iterator1)
	assert.ErrorIs(t, errDatabaseAccessViolation, err)
}

func TestLowerboundI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:       trxContext,
		Control:          controller,
		Session:          session,
		KeyValueCache:    NewIteratorCache(),
		Receiver:         name.StringToName("eosio.token"),
		AccountRamDeltas: make(map[name.Name]int64),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)
	iterator2, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 2, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator2, 1)
	iterator3, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 3, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator3, 2)

	// Lowerbound I64 should return iterator
	res, err := applyContext.LowerboundI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 2)
	assert.NoError(t, err)
	assert.Equal(t, iterator2, res)

	// Lowerbound I64 on non existing ID should return end iterator
	res, err = applyContext.LowerboundI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 4)
	assert.NoError(t, err)
	assert.Equal(t, -2, res)
}

func TestUpperboundI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:       trxContext,
		Control:          controller,
		Session:          session,
		KeyValueCache:    NewIteratorCache(),
		Receiver:         name.StringToName("eosio.token"),
		AccountRamDeltas: make(map[name.Name]int64),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)
	iterator2, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 2, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator2, 1)
	iterator3, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 4, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator3, 2)

	// Upperbound I64 should return iterator
	res, err := applyContext.UpperboundI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 3)
	assert.NoError(t, err)
	assert.Equal(t, iterator2, res)

	// Upperbound I64 on non existing ID should return end iterator
	res, err = applyContext.UpperboundI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), 0)
	assert.NoError(t, err)
	assert.Equal(t, -2, res)
}

func TestEndI64(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	assert.NoError(t, err)
	state := state.NewState(nil, db)
	controller := &Controller{
		State: state,
	}
	session := state.CreateSession(true)
	trxContext := &TransactionContext{
		Session: session,
		Control: controller,
	}
	applyContext := &applyContext{
		TrxContext:       trxContext,
		Control:          controller,
		Session:          session,
		KeyValueCache:    NewIteratorCache(),
		Receiver:         name.StringToName("eosio.token"),
		AccountRamDeltas: make(map[name.Name]int64),
	}
	iterator1, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 1, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator1, 0)
	iterator2, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 2, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator2, 1)
	iterator3, err := applyContext.StoreI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"), name.StringToName("glenn"), 3, []byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, iterator3, 2)

	// End I64 should return iterator
	res, err := applyContext.EndI64(name.StringToName("eosio.token"), name.StringToName("eosio"), name.StringToName("stat"))
	assert.NoError(t, err)
	assert.Equal(t, -2, res)

	// End I64 on non existing ID should return -1
	res, err = applyContext.EndI64(name.StringToName("eosio.token2"), name.StringToName("eosio"), name.StringToName("stat"))
	assert.NoError(t, err)
	assert.Equal(t, -1, res)
}
