package state

import (
	"encoding/binary"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/MetalBlockchain/metalgo/database"
)

const (
	// maximum block capacity of the cache
	accountCacheSize = 8192
)

var (
	accountIncrementKey = []byte("Account__id")
	accountIdKey        = []byte("Account__id__")
	accountNameKey      = []byte("Account__byName__")
)

var _ AccountState = &accountState{}

type AccountState interface {
	GetAccountByName(types.AccountName) (*Account, error)
	PutAccount(*Account) error
	UpdateAccount(*Account, func(*Account)) error
}

type accountState struct {
	accCache cache.Cacher
	db       database.Database
}

func NewAccountState(db database.Database) AccountState {
	return &accountState{
		accCache: &cache.LRU{Size: accountCacheSize},
		db:       db,
	}
}

func (s *accountState) GetAccount(id []byte) (*Account, error) {
	key := append(accountIdKey, id...)
	wrappedBytes, err := s.db.Get(key)

	if err != nil {
		return nil, err
	}

	account := &Account{}

	if _, err := Codec.Unmarshal(wrappedBytes, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountState) GetAccountByName(name types.AccountName) (*Account, error) {
	nameBytes, _ := name.Pack()
	byNameKey := append(accountNameKey, nameBytes...)
	wrappedBytes, err := s.db.Get(byNameKey)

	if err != nil {
		return nil, err
	}

	return s.GetAccount(wrappedBytes)
}

func (s *accountState) UpdateAccount(account *Account, updateFunc func(*Account)) error {
	if _, err := s.GetAccount(account.ID.ToBytes()); err != nil {
		return err
	}

	// Perform updates
	oldIndexKeys := getAccountIndexKeys(*account)
	updateFunc(account)
	newIndexKeys := getAccountIndexKeys(*account)
	batch := s.db.NewBatch()

	for _, key := range oldIndexKeys {
		if err := batch.Delete(key); err != nil {
			return err
		}
	}

	for _, key := range newIndexKeys {
		if err := batch.Put(key, account.ID.ToBytes()); err != nil {
			return err
		}
	}

	wrappedBytes, err := Codec.Marshal(CodecVersion, &account)

	if err != nil {
		return err
	}

	key := append(accountIdKey, account.ID.ToBytes()...)

	if err = batch.Put(key, wrappedBytes); err != nil {
		return err
	}

	return batch.Write()
}

func (s *accountState) PutAccount(account *Account) error {
	wrappedBytes, err := Codec.Marshal(CodecVersion, &account)

	if err != nil {
		return err
	}

	nameBytes, _ := account.Name.Pack()

	id, err := s.GenerateAccountId()

	if err != nil {
		return err
	}

	account.ID = types.IdType(id)
	batch := s.db.NewBatch()
	key := append(accountIdKey, account.ID.ToBytes()...)
	byNameKey := append(accountNameKey, nameBytes...)

	batch.Put(key, wrappedBytes)
	batch.Put(byNameKey, account.ID.ToBytes())

	return batch.Write()
}

func (s *accountState) GenerateAccountId() (uint64, error) {
	var id uint64

	if value, err := s.db.Get(accountIncrementKey); err == nil {
		id = binary.BigEndian.Uint64(value) + 1
	}

	newValue := make([]byte, 8)
	binary.BigEndian.PutUint64(newValue, uint64(id))

	if err := s.db.Put(accountIncrementKey, newValue); err != nil {
		return 0, err
	}

	return id, nil
}

func getAccountIndexKeys(account Account) map[string][]byte {
	keys := make(map[string][]byte)
	nameBytes, _ := account.Name.Pack()
	byName := append(accountNameKey, nameBytes...)

	keys["byName"] = byName

	return keys
}
