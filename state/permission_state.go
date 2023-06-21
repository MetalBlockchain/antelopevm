package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func (s *Session) FindPermission(id core.IdType) (*core.Permission, error) {
	key := getObjectKeyByIndex(&core.Permission{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &core.Permission{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindPermissionByOwner(owner name.AccountName, name name.PermissionName) (*core.Permission, error) {
	key := getObjectKeyByIndex(&core.Permission{Owner: owner, Name: name}, "byOwner")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindPermission(core.NewIdType(data))
}

func (s *Session) FindPermissionsByOwner(owner name.AccountName) *Iterator[core.Permission] {
	key := getPartialKey("byOwner", &core.Permission{}, owner)

	return newIterator(s, key, func(b []byte) (*core.Permission, error) {
		return s.FindPermission(core.NewIdType(b))
	})
}

func (s *Session) CreatePermission(in *core.Permission) error {
	return s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)
}

func (s *Session) ModifyPermission(in *core.Permission, modifyFunc func()) error {
	return s.modify(in, modifyFunc)
}

func (s *Session) RemovePermission(in *core.Permission) error {
	return s.remove(in)
}
