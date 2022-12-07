package state

import (
	"encoding/binary"
	"fmt"

	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/prefixdb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/choices"
)

const (
	// maximum block capacity of the cache
	blockCacheSize = 8192
)

var (
	blockIdKey      = []byte("Block__id__")
	blockIndexKey   = []byte("Block__byIndex__")
	lastAcceptedKey = []byte("Block_lastAccepted")
)

var _ BlockState = &blockState{}

// BlockState defines methods to manage state with Blocks and LastAcceptedIDs.
type BlockState interface {
	GetBlock(blkID ids.ID) (*Block, error)
	GetBlockByIndex(uint64) (*Block, error)
	PutBlock(blk *Block) error
	GetLastAccepted() (ids.ID, error)
	SetLastAccepted(ids.ID) error
}

// blockState implements BlocksState interface with database and cache.
type blockState struct {
	// cache to store blocks
	blkCache cache.Cacher
	// block database
	db           database.Database
	blockDB      database.Database
	blockIndexDB database.Database
	lastAccepted ids.ID
	vm           VM
}

// blkWrapper wraps the actual blk bytes and status to persist them together
type blkWrapper struct {
	Blk    []byte         `serialize:"true"`
	Status choices.Status `serialize:"true"`
}

// NewBlockState returns BlockState with a new cache and given db
func NewBlockState(vm VM, db database.Database) BlockState {
	return &blockState{
		vm:           vm,
		blkCache:     &cache.LRU{Size: blockCacheSize},
		blockDB:      prefixdb.New(blockIdKey, db),
		blockIndexDB: prefixdb.New(blockIndexKey, db),
		db:           db,
	}
}

// GetBlock gets Block from either cache or database
func (s *blockState) GetBlock(blkID ids.ID) (*Block, error) {
	// Check if cache has this blkID
	if blkIntf, cached := s.blkCache.Get(blkID); cached {
		// there is a key but value is nil, so return an error
		if blkIntf == nil {
			return nil, database.ErrNotFound
		}
		// We found it return the block in cache
		return blkIntf.(*Block), nil
	}

	// get block bytes from db with the blkID key
	wrappedBytes, err := s.blockDB.Get(blkID[:])
	if err != nil {
		// we could not find it in the db, let's cache this blkID with nil value
		// so next time we try to fetch the same key we can return error
		// without hitting the database
		if err == database.ErrNotFound {
			s.blkCache.Put(blkID, nil)
		}
		// could not find the block, return error
		return nil, err
	}

	// first decode/unmarshal the block wrapper so we can have status and block bytes
	blkw := blkWrapper{}
	if _, err := Codec.Unmarshal(wrappedBytes, &blkw); err != nil {
		return nil, err
	}

	// now decode/unmarshal the actual block bytes to block
	blk := &Block{}
	if _, err := Codec.Unmarshal(blkw.Blk, blk); err != nil {
		return nil, err
	}

	// initialize block with block bytes, status and vm
	blk.Initialize(s.vm, blkw.Status)

	// put block into cache
	s.blkCache.Put(blkID, blk)

	return blk, nil
}

func (s *blockState) GetBlockByIndex(index uint64) (*Block, error) {
	indexValue := make([]byte, 8)
	binary.BigEndian.PutUint64(indexValue, uint64(index))

	// get block bytes from db with the blkID key
	wrappedBytes, err := s.blockIndexDB.Get(indexValue)

	if err != nil {
		return nil, err
	}

	blkID, err := ids.ToID(wrappedBytes)

	if err != nil {
		return nil, err
	}

	return s.GetBlock(blkID)
}

// PutBlock puts block into both database and cache
func (s *blockState) PutBlock(blk *Block) error {
	// create block wrapper with block bytes and status
	blkw := blkWrapper{
		Blk:    blk.Bytes(),
		Status: blk.Status(),
	}

	// encode block wrapper to its byte representation
	wrappedBytes, err := Codec.Marshal(CodecVersion, &blkw)
	if err != nil {
		return err
	}

	blkID := blk.ID()
	// put actual block to cache, so we can directly fetch it from cache
	s.blkCache.Put(blkID, blk)

	// put wrapped block bytes into database
	indexValue := make([]byte, 8)
	binary.BigEndian.PutUint64(indexValue, uint64(blk.Index))

	if err := s.blockDB.Put(blkID[:], wrappedBytes); err != nil {
		return err
	}

	if err := s.blockIndexDB.Put(indexValue, blkID[:]); err != nil {
		return err
	}

	return nil
}

// DeleteBlock deletes block from both cache and database
func (s *blockState) DeleteBlock(blkID ids.ID) error {
	return fmt.Errorf("Who in their right mind deletes blocks?")
}

// GetLastAccepted returns last accepted block ID
func (s *blockState) GetLastAccepted() (ids.ID, error) {
	// check if we already have lastAccepted ID in state memory
	if s.lastAccepted != ids.Empty {
		return s.lastAccepted, nil
	}

	// get lastAccepted bytes from database with the fixed lastAcceptedKey
	lastAcceptedBytes, err := s.db.Get(lastAcceptedKey)
	if err != nil {
		return ids.ID{}, err
	}
	// parse bytes to ID
	lastAccepted, err := ids.ToID(lastAcceptedBytes)
	if err != nil {
		return ids.ID{}, err
	}
	// put lastAccepted ID into memory
	s.lastAccepted = lastAccepted
	return lastAccepted, nil
}

// SetLastAccepted persists lastAccepted ID into both cache and database
func (s *blockState) SetLastAccepted(lastAccepted ids.ID) error {
	if s.lastAccepted == lastAccepted {
		return nil
	}

	s.lastAccepted = lastAccepted

	return s.db.Put(lastAcceptedKey, lastAccepted[:])
}

func (s *blockState) AcceptBlock(block *Block) error {
	block.SetStatus(choices.Accepted) // Change state of this block
	blkID := block.ID()

	// Persist data
	if err := s.PutBlock(block); err != nil {
		return err
	}

	// Set last accepted ID to this block ID
	if err := s.SetLastAccepted(blkID); err != nil {
		return err
	}

	// Commit changes to database
	return nil
}
