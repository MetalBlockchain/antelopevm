package state

import (
	"encoding/binary"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/MetalBlockchain/metalgo/database"
)

var (
	permissionIncrementKey          = []byte("Permission__id")
	permissionIdKey                 = []byte("Permission__id__")
	permissionNameKey               = []byte("Permission__byName__")
	permissionParentKey             = []byte("Permission__byParent__")
	permissionOwnerKey              = []byte("Permission__byOwner__")
	permissionLinkIncrementKey      = []byte("Permission__id")
	permissionLinkIdKey             = []byte("PermissionLink__id__")
	permissionLinkActionNameKey     = []byte("PermissionLink__byActionName__")
	permissionLinkPermissionNameKey = []byte("PermissionLink__byPermissionName__")
	separator                       = []byte("__")
)

var _ PermissionState = &permissionState{}

type PermissionState interface {
	GetPermission([]byte) (*core.Permission, error)
	GetPermissionByOwner(core.AccountName, core.PermissionName) (*core.Permission, error)
	UpdatePermission(*core.Permission) error
	PutPermission(*core.Permission) error
	RemovePermission(*core.Permission) error

	GetPermissionLink([]byte) (*core.PermissionLink, error)
	GetPermissionLinkByActionName(core.AccountName, core.AccountName, core.ActionName) (*core.PermissionLink, error)
	GetPermissionLinksByPermissionName(core.AccountName, core.PermissionName) database.Iterator
	UpdatePermissionLink(*core.PermissionLink, func(*core.PermissionLink)) error
	PutPermissionLink(*core.PermissionLink) error
	RemovePermissionLink(*core.PermissionLink) error
}

type permissionState struct {
	accCache cache.Cacher
	db       database.Database
}

func NewPermissionState(db database.Database) PermissionState {
	return &permissionState{
		accCache: &cache.LRU{Size: accountCacheSize},
		db:       db,
	}
}

func (s *permissionState) GetPermission(id []byte) (*core.Permission, error) {
	key := append(permissionIdKey, id...)
	wrappedBytes, err := s.db.Get(key)

	if err != nil {
		return nil, err
	}

	permission := &core.Permission{}

	if _, err := Codec.Unmarshal(wrappedBytes, permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionState) GetPermissionByOwner(owner core.AccountName, name core.PermissionName) (*core.Permission, error) {
	nameBytes, _ := name.Pack()
	ownerBytes, _ := owner.Pack()
	byOwnerKey := append(permissionOwnerKey, ownerBytes...)
	byOwnerKey = append(byOwnerKey, separator...)
	byOwnerKey = append(byOwnerKey, nameBytes...)
	wrappedBytes, err := s.db.Get(byOwnerKey)

	if err != nil {
		return nil, err
	}

	return s.GetPermission(wrappedBytes)
}

func (s *permissionState) UpdatePermission(perm *core.Permission) error {
	if _, err := s.GetPermission(perm.ID.ToBytes()); err != nil {
		return err
	}

	wrappedBytes, err := Codec.Marshal(CodecVersion, &perm)

	if err != nil {
		return err
	}

	key := append(permissionIdKey, perm.ID.ToBytes()...)

	if err = s.db.Put(key, wrappedBytes); err != nil {
		return err
	}

	return nil
}

func (s *permissionState) PutPermission(perm *core.Permission) error {
	wrappedBytes, err := Codec.Marshal(CodecVersion, &perm)

	if err != nil {
		return err
	}

	nameBytes, _ := perm.Name.Pack()
	ownerBytes, _ := perm.Owner.Pack()
	id, err := s.GeneratePermissionId()

	if err != nil {
		return err
	}

	perm.ID = core.IdType(id)
	batch := s.db.NewBatch()

	key := append(permissionIdKey, perm.ID.ToBytes()...)
	byNameKey := append(permissionNameKey, nameBytes...)
	byNameKey = append(byNameKey, separator...)
	byNameKey = append(byNameKey, perm.ID.ToBytes()...)

	byParentKey := append(permissionParentKey, perm.Parent.ToBytes()...)
	byParentKey = append(byParentKey, separator...)
	byParentKey = append(byParentKey, perm.ID.ToBytes()...)

	byOwnerKey := append(permissionOwnerKey, ownerBytes...)
	byOwnerKey = append(byOwnerKey, separator...)
	byOwnerKey = append(byOwnerKey, nameBytes...)

	batch.Put(key, wrappedBytes)
	batch.Put(byNameKey, perm.ID.ToBytes())
	batch.Put(byParentKey, perm.ID.ToBytes())
	batch.Put(byOwnerKey, perm.ID.ToBytes())

	return batch.Write()
}

func (s *permissionState) RemovePermission(perm *core.Permission) error {
	if _, err := s.GetPermission(perm.ID.ToBytes()); err != nil {
		return err
	}

	// Perform updates
	indexKeys := getPermissionIndexKeys(*perm)
	batch := s.db.NewBatch()

	for _, key := range indexKeys {
		if err := batch.Delete(key); err != nil {
			return err
		}
	}

	key := append(permissionIdKey, perm.ID.ToBytes()...)

	if err := batch.Delete(key); err != nil {
		return err
	}

	return batch.Write()
}

func (s *permissionState) GeneratePermissionId() (uint64, error) {
	var id uint64

	if value, err := s.db.Get(permissionIncrementKey); err == nil {
		id = binary.BigEndian.Uint64(value) + 1
	}

	newValue := make([]byte, 8)
	binary.BigEndian.PutUint64(newValue, uint64(id))

	if err := s.db.Put(permissionIncrementKey, newValue); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *permissionState) GetPermissionLink(id []byte) (*core.PermissionLink, error) {
	key := append(permissionLinkIdKey, id...)
	wrappedBytes, err := s.db.Get(key)

	if err != nil {
		return nil, err
	}

	link := &core.PermissionLink{}

	if _, err := Codec.Unmarshal(wrappedBytes, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *permissionState) GetPermissionLinkByActionName(account core.AccountName, code core.AccountName, messageType core.ActionName) (*core.PermissionLink, error) {
	accountBytes, _ := account.Pack()
	codeBytes, _ := code.Pack()
	messageTypeBytes, _ := messageType.Pack()
	byActionName := append(permissionLinkActionNameKey, accountBytes...)
	byActionName = append(byActionName, separator...)
	byActionName = append(byActionName, codeBytes...)
	byActionName = append(byActionName, separator...)
	byActionName = append(byActionName, messageTypeBytes...)
	wrappedBytes, err := s.db.Get(byActionName)

	if err != nil {
		return nil, err
	}

	return s.GetPermissionLink(wrappedBytes)
}

func (s *permissionState) GetPermissionLinksByPermissionName(account core.AccountName, permission core.PermissionName) database.Iterator {
	accountBytes, _ := account.Pack()
	requiredPermissionBytes, _ := permission.Pack()
	byPermissionName := append(permissionLinkPermissionNameKey, accountBytes...)
	byPermissionName = append(byPermissionName, separator...)
	byPermissionName = append(byPermissionName, requiredPermissionBytes...)
	byPermissionName = append(byPermissionName, separator...)
	iterator := s.db.NewIteratorWithPrefix(byPermissionName)

	return iterator
}

func (s *permissionState) UpdatePermissionLink(link *core.PermissionLink, updateFunc func(*core.PermissionLink)) error {
	if _, err := s.GetPermissionLink(link.ID.ToBytes()); err != nil {
		return err
	}

	// Perform updates
	oldIndexKeys := getPermissionLinkIndexKeys(*link)
	updateFunc(link)
	newIndexKeys := getPermissionLinkIndexKeys(*link)
	batch := s.db.NewBatch()

	for _, key := range oldIndexKeys {
		if err := batch.Delete(key); err != nil {
			return err
		}
	}

	for _, key := range newIndexKeys {
		if err := batch.Put(key, link.ID.ToBytes()); err != nil {
			return err
		}
	}

	wrappedBytes, err := Codec.Marshal(CodecVersion, &link)

	if err != nil {
		return err
	}

	key := append(permissionLinkIdKey, link.ID.ToBytes()...)

	if err = batch.Put(key, wrappedBytes); err != nil {
		return err
	}

	return batch.Write()
}

func (s *permissionState) RemovePermissionLink(link *core.PermissionLink) error {
	if _, err := s.GetPermissionLink(link.ID.ToBytes()); err != nil {
		return err
	}

	// Perform updates
	indexKeys := getPermissionLinkIndexKeys(*link)
	batch := s.db.NewBatch()

	for _, key := range indexKeys {
		if err := batch.Delete(key); err != nil {
			return err
		}
	}

	key := append(permissionLinkIdKey, link.ID.ToBytes()...)

	if err := batch.Delete(key); err != nil {
		return err
	}

	return batch.Write()
}

func (s *permissionState) PutPermissionLink(link *core.PermissionLink) error {
	id, err := s.GeneratePermissionLinkId()

	if err != nil {
		return err
	}

	link.ID = core.IdType(id)

	wrappedBytes, err := Codec.Marshal(CodecVersion, &link)

	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	key := append(permissionLinkIdKey, link.ID.ToBytes()...)
	indexKeys := getPermissionLinkIndexKeys(*link)

	batch.Put(key, wrappedBytes)

	for _, indexKey := range indexKeys {
		batch.Put(indexKey, link.ID.ToBytes())
		batch.Put(indexKey, link.ID.ToBytes())
	}

	return batch.Write()
}

func (s *permissionState) GeneratePermissionLinkId() (uint64, error) {
	var id uint64

	if value, err := s.db.Get(permissionLinkIncrementKey); err == nil {
		id = binary.BigEndian.Uint64(value) + 1
	}

	newValue := make([]byte, 8)
	binary.BigEndian.PutUint64(newValue, uint64(id))

	if err := s.db.Put(permissionLinkIncrementKey, newValue); err != nil {
		return 0, err
	}

	return id, nil
}

func getPermissionIndexKeys(perm core.Permission) map[string][]byte {
	keys := make(map[string][]byte)
	nameBytes, _ := perm.Name.Pack()
	ownerBytes, _ := perm.Owner.Pack()
	byNameKey := append(permissionNameKey, nameBytes...)
	byNameKey = append(byNameKey, separator...)
	byNameKey = append(byNameKey, perm.ID.ToBytes()...)

	byParentKey := append(permissionParentKey, perm.Parent.ToBytes()...)
	byParentKey = append(byParentKey, separator...)
	byParentKey = append(byParentKey, perm.ID.ToBytes()...)

	byOwnerKey := append(permissionOwnerKey, ownerBytes...)
	byOwnerKey = append(byOwnerKey, separator...)
	byOwnerKey = append(byOwnerKey, nameBytes...)

	keys["byName"] = byNameKey
	keys["byParent"] = byParentKey
	keys["byOwner"] = byOwnerKey

	return keys
}

func getPermissionLinkIndexKeys(link core.PermissionLink) map[string][]byte {
	keys := make(map[string][]byte)
	accountBytes, _ := link.Account.Pack()
	codeBytes, _ := link.Code.Pack()
	messageTypeBytes, _ := link.MessageType.Pack()
	requiredPermissionBytes, _ := link.RequiredPermission.Pack()
	byActionName := append(permissionLinkActionNameKey, accountBytes...)
	byActionName = append(byActionName, separator...)
	byActionName = append(byActionName, codeBytes...)
	byActionName = append(byActionName, separator...)
	byActionName = append(byActionName, messageTypeBytes...)

	byPermissionName := append(permissionLinkPermissionNameKey, accountBytes...)
	byPermissionName = append(byPermissionName, separator...)
	byPermissionName = append(byPermissionName, requiredPermissionBytes...)
	byPermissionName = append(byPermissionName, separator...)
	byPermissionName = append(byPermissionName, link.ID.ToBytes()...)

	keys["byActionName"] = byActionName
	keys["byPermissionName"] = byPermissionName

	return keys
}
