package state

import (
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
)

var (
	_              State = &s{}
	initializedKey       = []byte("initialized")
)

type State interface {
	BlockState
	AccountState
	PermissionState

	Commit() error
	Close() error
	IsInitialized() (bool, error)
	SetInitialized() error
}

type s struct {
	BlockState
	AccountState
	PermissionState
	baseDB *versiondb.Database
}

func NewState(vm VM, db database.Database) State {
	baseDB := versiondb.New(db)

	return &s{
		BlockState:      NewBlockState(vm, baseDB),
		AccountState:    NewAccountState(baseDB),
		PermissionState: NewPermissionState(baseDB),
		baseDB:          baseDB,
	}
}

func (s *s) Commit() error {
	return s.baseDB.Commit()
}

func (s *s) Close() error {
	return s.baseDB.Close()
}

func (s *s) IsInitialized() (bool, error) {
	return s.baseDB.Has(initializedKey)
}

func (s *s) SetInitialized() error {
	return s.baseDB.Put(initializedKey, []byte{1})
}
