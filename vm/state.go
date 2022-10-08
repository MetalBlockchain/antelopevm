package vm

import (
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/prefixdb"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
)

var (
	singletonStatePrefix = []byte("singleton")
	blockStatePrefix     = []byte("block")

	_ State = &state{}
)

type State interface {
	avax.SingletonState
	BlockState

	Commit() error
	Close() error
}

type state struct {
	avax.SingletonState
	BlockState

	baseDB *versiondb.Database
}

func NewState(db database.Database, vm *VM) State {
	baseDB := versiondb.New(db)
	blockDB := prefixdb.New(blockStatePrefix, baseDB)
	singletonDB := prefixdb.New(singletonStatePrefix, baseDB)

	return &state{
		BlockState:     NewBlockState(blockDB, vm),
		SingletonState: avax.NewSingletonState(singletonDB),
		baseDB:         baseDB,
	}
}

func (s *state) Commit() error {
	return s.baseDB.Commit()
}

func (s *state) Close() error {
	return s.baseDB.Close()
}
