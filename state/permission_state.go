package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

func (s *Session) FindPermission(id types.IdType) (*authority.Permission, error) {
	key := getObjectKeyByIndex(&authority.Permission{ID: id}, "id")
	item, err := s.transaction.Get(key)
	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	out := &authority.Permission{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindPermissionByOwner(owner name.AccountName, name name.PermissionName) (*authority.Permission, error) {
	key := getObjectKeyByIndex(&authority.Permission{Owner: owner, Name: name}, "byOwner")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindPermission(types.NewIdType(data))
}

func (s *Session) FindPermissionsByOwner(owner name.AccountName) *Iterator[authority.Permission] {
	key := getPartialKey("byOwner", &authority.Permission{}, owner)

	return newIterator(s, key, func(b []byte) (*authority.Permission, error) {
		return s.FindPermission(types.NewIdType(b))
	})
}

func (s *Session) CreatePermission(in *authority.Permission) error {
	return s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)
}

func (s *Session) ModifyPermission(in *authority.Permission, modifyFunc func()) error {
	return s.modify(in, modifyFunc)
}

func (s *Session) RemovePermission(in *authority.Permission) error {
	return s.remove(in)
}
