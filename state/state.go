package state

import (
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
)

var (
	_ State = &state{}
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

type state struct {
	BlockState
	AccountState
	PermissionState
	baseDB *versiondb.Database
}

func NewState(db database.Database) State {
	baseDB := versiondb.New(db)

	return &state{
		BlockState:      NewBlockState(baseDB),
		AccountState:    NewAccountState(baseDB),
		PermissionState: NewPermissionState(baseDB),
		baseDB:          baseDB,
	}
}

func (s *state) Commit() error {
	return s.baseDB.Commit()
}

func (s *state) Close() error {
	return s.baseDB.Close()
}

func (s *state) IsInitialized() (bool, error) {
	return true, nil
}

func (s *state) SetInitialized() error {
	return nil
}
