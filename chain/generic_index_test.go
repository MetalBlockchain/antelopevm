package chain

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/state"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

type Idx64Record struct {
	ssn  uint64
	name uint64
}

func TestIdx64General(t *testing.T) {
	receiver := name.StringToName("eosio.token")
	table := name.StringToName("myindextable")
	records := []Idx64Record{
		{ssn: 265, name: uint64(name.StringToName("alice"))},
		{ssn: 781, name: uint64(name.StringToName("bob"))},
		{ssn: 234, name: uint64(name.StringToName("charlie"))},
		{ssn: 650, name: uint64(name.StringToName("allyson"))},
		{ssn: 540, name: uint64(name.StringToName("bob"))},
		{ssn: 976, name: uint64(name.StringToName("emily"))},
		{ssn: 110, name: uint64(name.StringToName("joe"))},
	}
	applyContext, err := setupEnvironment(t)
	assert.NoError(t, err)

	for _, record := range records {
		_, err := applyContext.Idx64.Store(receiver, table, receiver, record.ssn, record.name)

		if err != nil {
			panic(err)
		}
	}

	// find_primary
	var sec uint64
	itr := applyContext.Idx64.FindPrimary(receiver, receiver, table, &sec, 999)
	assert.True(t, itr < 0 && sec == 0, "idx64_general - db_idx64_find_primary")
	itr = applyContext.Idx64.FindPrimary(receiver, receiver, table, &sec, 110)
	assert.True(t, itr >= 0 && sec == uint64(name.StringToName("joe")), "idx64_general - db_idx64_find_primary")
	var primNext uint64
	itrNext, err := applyContext.Idx64.NextSecondary(itr, &primNext)
	assert.NoError(t, err)
	assert.True(t, itrNext < 0 && primNext == 0, "idx64_general - db_idx64_find_primary")

	// iterate forward starting with charlie
	sec = 0
	itr = applyContext.Idx64.FindPrimary(receiver, receiver, table, &sec, 234)
	assert.True(t, itr >= 0 && sec == uint64(name.StringToName("charlie")), "idx64_general - db_idx64_find_primary")

	primNext = 0
	itrNext, err = applyContext.Idx64.NextSecondary(itr, &primNext)
	assert.NoError(t, err)
	assert.True(t, itrNext >= 0 && primNext == 976, "idx64_general - db_idx64_find_primary")
	var secNext uint64
	itrNextExpected := applyContext.Idx64.FindPrimary(receiver, receiver, table, &secNext, primNext)
	assert.True(t, itrNext == itrNextExpected && secNext == uint64(name.StringToName("emily")))

	itrNext, err = applyContext.Idx64.NextSecondary(itrNext, &primNext)
	assert.NoError(t, err)
	assert.True(t, itrNext >= 0 && primNext == 110)
	itrNextExpected = applyContext.Idx64.FindPrimary(receiver, receiver, table, &secNext, primNext)
	assert.True(t, itrNext == itrNextExpected && secNext == uint64(name.StringToName("joe")))

	itrNext, err = applyContext.Idx64.NextSecondary(itrNext, &primNext)
	assert.NoError(t, err)
	assert.True(t, itrNext < 0 && primNext == 110)

	// iterate backward staring with second bob
	sec = 0
	itr = applyContext.Idx64.FindPrimary(receiver, receiver, table, &sec, 781)
	assert.True(t, itr >= 0 && sec == uint64(name.StringToName("bob")))

	var primPrev uint64
	itrPrev, err := applyContext.Idx64.PreviousSecondary(itr, &primPrev)
	assert.NoError(t, err)
	assert.True(t, itrPrev >= 0 && primPrev == 540)

	var secPrev uint64
	itrPrevExpected := applyContext.Idx64.FindPrimary(receiver, receiver, table, &secPrev, primPrev)
	assert.True(t, itrPrev == itrPrevExpected && secPrev == uint64(name.StringToName("bob")))

	itrPrev, err = applyContext.Idx64.PreviousSecondary(itrPrev, &primPrev)
	assert.NoError(t, err)
	assert.True(t, itrPrev >= 0 && primPrev == 650)
	itrPrevExpected = applyContext.Idx64.FindPrimary(receiver, receiver, table, &secPrev, primPrev)
	assert.True(t, itrPrev == itrPrevExpected && secPrev == uint64(name.StringToName("allyson")))

	itrPrev, err = applyContext.Idx64.PreviousSecondary(itrPrev, &primPrev)
	assert.NoError(t, err)
	assert.True(t, itrPrev >= 0 && primPrev == 265)
	itrPrevExpected = applyContext.Idx64.FindPrimary(receiver, receiver, table, &secPrev, primPrev)
	assert.True(t, itrPrev == itrPrevExpected && secPrev == uint64(name.StringToName("alice")))

	itrPrev, err = applyContext.Idx64.PreviousSecondary(itrPrev, &primPrev)
	assert.NoError(t, err)
	assert.True(t, itrPrev < 0 && primPrev == 265)

	// find_secondary
	var prim uint64
	sec = uint64(name.StringToName("bob"))
	itr = applyContext.Idx64.FindSecondary(receiver, receiver, table, &sec, &prim)
	assert.True(t, itr >= 0 && prim == 540)

	sec = uint64(name.StringToName("emily"))
	itr = applyContext.Idx64.FindSecondary(receiver, receiver, table, &sec, &prim)
	assert.True(t, itr >= 0 && prim == 976)

	sec = uint64(name.StringToName("frank"))
	itr = applyContext.Idx64.FindSecondary(receiver, receiver, table, &sec, &prim)
	assert.True(t, itr < 0 && prim == 976)

	// update and remove
	oneMoreBob := uint64(name.StringToName("bob"))
	ssn := uint64(421)
	itr, err = applyContext.Idx64.Store(receiver, table, receiver, ssn, oneMoreBob)
	assert.NoError(t, err)
	newName := uint64(name.StringToName("billy"))
	err = applyContext.Idx64.Update(itr, receiver, newName)
	assert.NoError(t, err)
	sec = 0
	secItr := applyContext.Idx64.FindPrimary(receiver, receiver, table, &sec, ssn)
	assert.True(t, secItr == itr && sec == newName)
	err = applyContext.Idx64.Remove(itr)
	assert.NoError(t, err)
	itrf := applyContext.Idx64.FindPrimary(receiver, receiver, table, &sec, ssn)
	assert.True(t, itrf < 0)
}

func TestIdx64Lowerbound(t *testing.T) {
	receiver := name.StringToName("eosio.token")
	table := name.StringToName("myindextable")
	records := []Idx64Record{
		{ssn: 265, name: uint64(name.StringToName("alice"))},
		{ssn: 781, name: uint64(name.StringToName("bob"))},
		{ssn: 234, name: uint64(name.StringToName("charlie"))},
		{ssn: 650, name: uint64(name.StringToName("allyson"))},
		{ssn: 540, name: uint64(name.StringToName("bob"))},
		{ssn: 976, name: uint64(name.StringToName("emily"))},
		{ssn: 110, name: uint64(name.StringToName("joe"))},
	}
	applyContext, err := setupEnvironment(t)
	assert.NoError(t, err)

	for _, record := range records {
		_, err := applyContext.Idx64.Store(receiver, table, receiver, record.ssn, record.name)

		if err != nil {
			panic(err)
		}
	}

	lbPrim := uint64(0)
	lbSec := uint64(name.StringToName("alice"))
	ssn := uint64(265)
	lb := applyContext.Idx64.LowerboundSecondary(receiver, receiver, table, &lbSec, &lbPrim)
	assert.True(t, lbPrim == ssn && lbSec == uint64(name.StringToName("alice")))
	assert.True(t, lb == applyContext.Idx64.FindPrimary(receiver, receiver, table, &lbSec, ssn))

	lbSec = uint64(name.StringToName("billy"))
	lbPrim = 0
	ssn = uint64(540)
	lb = applyContext.Idx64.LowerboundSecondary(receiver, receiver, table, &lbSec, &lbPrim)
	assert.True(t, lbPrim == ssn && lbSec == uint64(name.StringToName("bob")), "lbPrim %v == ssn %v && lbSec %v == bob", lbPrim, ssn, lbSec)
	assert.True(t, lb == applyContext.Idx64.FindPrimary(receiver, receiver, table, &lbSec, ssn))

	lbSec = uint64(name.StringToName("joe"))
	lbPrim = 0
	ssn = uint64(110)
	lb = applyContext.Idx64.LowerboundSecondary(receiver, receiver, table, &lbSec, &lbPrim)
	assert.True(t, lbPrim == ssn && lbSec == uint64(name.StringToName("joe")), "lbPrim %v == ssn %v && lbSec %v == joe", lbPrim, ssn, lbSec)
	assert.True(t, lb == applyContext.Idx64.FindPrimary(receiver, receiver, table, &lbSec, ssn))

	lbSec = uint64(name.StringToName("kevin"))
	lbPrim = 0
	lb = applyContext.Idx64.LowerboundSecondary(receiver, receiver, table, &lbSec, &lbPrim)
	assert.True(t, lbPrim == 0 && lbSec == uint64(name.StringToName("kevin")), "lbPrim %v == 0 %v && lbSec %v == kevin", lbPrim, ssn, lbSec)
	assert.True(t, lb < 0)
}

func TestIdx64Upperbound(t *testing.T) {
	receiver := name.StringToName("eosio.token")
	table := name.StringToName("myindextable")
	records := []Idx64Record{
		{ssn: 265, name: uint64(name.StringToName("alice"))},
		{ssn: 781, name: uint64(name.StringToName("bob"))},
		{ssn: 234, name: uint64(name.StringToName("charlie"))},
		{ssn: 650, name: uint64(name.StringToName("allyson"))},
		{ssn: 540, name: uint64(name.StringToName("bob"))},
		{ssn: 976, name: uint64(name.StringToName("emily"))},
		{ssn: 110, name: uint64(name.StringToName("joe"))},
	}
	applyContext, err := setupEnvironment(t)
	assert.NoError(t, err)

	for _, record := range records {
		_, err := applyContext.Idx64.Store(receiver, table, receiver, record.ssn, record.name)

		if err != nil {
			panic(err)
		}
	}

	ubPrim := uint64(0)
	ubSec := uint64(name.StringToName("alice"))
	allysonSsn := uint64(650)
	ub := applyContext.Idx64.UpperboundSecondary(receiver, receiver, table, &ubSec, &ubPrim)
	assert.True(t, ubPrim == allysonSsn && ubSec == uint64(name.StringToName("allyson")), "ubPrim %v == %v && ubSec %v == allyson", ubPrim, allysonSsn, name.NameToString(ubSec))
	assert.True(t, ub == applyContext.Idx64.FindPrimary(receiver, receiver, table, &ubSec, allysonSsn))

	ubPrim = uint64(0)
	ubSec = uint64(name.StringToName("billy"))
	bobSsn := uint64(540)
	ub = applyContext.Idx64.UpperboundSecondary(receiver, receiver, table, &ubSec, &ubPrim)
	assert.True(t, ubPrim == bobSsn && ubSec == uint64(name.StringToName("bob")), "ubPrim %v == %v && ubSec %v == bob", ubPrim, bobSsn, name.NameToString(ubSec))
	assert.True(t, ub == applyContext.Idx64.FindPrimary(receiver, receiver, table, &ubSec, bobSsn))

	ubPrim = uint64(0)
	ubSec = uint64(name.StringToName("joe"))
	ub = applyContext.Idx64.UpperboundSecondary(receiver, receiver, table, &ubSec, &ubPrim)
	assert.True(t, ubPrim == 0 && ubSec == uint64(name.StringToName("joe")), "ubPrim %v == 0 && ubSec %v == joe", ubPrim, name.NameToString(ubSec))
	assert.True(t, ub < 0)

	ubPrim = uint64(0)
	ubSec = uint64(name.StringToName("kevin"))
	ub = applyContext.Idx64.UpperboundSecondary(receiver, receiver, table, &ubSec, &ubPrim)
	assert.True(t, ubPrim == 0 && ubSec == uint64(name.StringToName("kevin")), "ubPrim %v == 0 && ubSec %v == kevin", ubPrim, name.NameToString(ubSec))
	assert.True(t, ub < 0)
}

func setupEnvironment(t *testing.T) (*applyContext, error) {
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	if err != nil {
		return nil, err
	}
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
	applyContext.Idx64 = &Idx64{Context: applyContext}

	return applyContext, nil
}
