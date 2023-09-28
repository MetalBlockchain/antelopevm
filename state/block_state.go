package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
)

const (
	// maximum block capacity of the cache
	blockCacheSize = 8192
)

var (
	lastAcceptedKey = []byte("lastAccepted")
)

func (s *Session) FindBlock(id types.IdType) (*Block, error) {
	key := getObjectKeyByIndex(&Block{Index: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &Block{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindBlockByIndex(index uint64) (*Block, error) {
	return s.FindBlock(types.IdType(index))
}

func (s *Session) FindBlockByHash(hash block.BlockHash) (*Block, error) {
	key := getObjectKeyByIndex(&Block{Hash: hash}, "byHash")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindBlock(types.NewIdType(data))
}

func (s *Session) CreateBlock(in *Block) error {
	in.Index = types.IdType(in.Header.BlockNum())

	return s.create(false, nil, in)
}

// GetLastAccepted returns last accepted block ID
func (s *Session) GetLastAccepted() (ids.ID, error) {
	// check if we already have lastAccepted ID in state memory
	if s.state.lastAccepted != ids.Empty {
		return s.state.lastAccepted, nil
	}

	// get lastAccepted bytes from database with the fixed lastAcceptedKey
	lastAcceptedBytes, err := s.transaction.Get(lastAcceptedKey)

	if err != nil {
		return ids.ID{}, err
	}

	// parse bytes to ID
	if value, err := lastAcceptedBytes.ValueCopy(nil); err == nil {
		if lastAccepted, err := ids.ToID(value); err == nil {
			s.state.lastAccepted = lastAccepted

			return lastAccepted, nil
		}
	}

	return ids.ID{}, err
}

// SetLastAccepted persists lastAccepted ID into both cache and database
func (s *Session) SetLastAccepted(lastAccepted ids.ID) error {
	if s.state.lastAccepted == lastAccepted {
		return nil
	}

	if err := s.transaction.Set(lastAcceptedKey, lastAccepted[:]); err != nil {
		return err
	}

	s.state.lastAccepted = lastAccepted

	return nil
}

func (s *Session) AcceptBlock(block *Block) error {
	block.SetStatus(choices.Accepted) // Change state of this block
	blkID := block.ID()

	// Persist data
	if err := s.CreateBlock(block); err != nil {
		return err
	}

	// Set last accepted ID to this block ID
	if err := s.SetLastAccepted(blkID); err != nil {
		return err
	}

	// Commit changes to database
	return nil
}
