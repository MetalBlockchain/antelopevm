package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/authority"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

func (s *Session) FindPermissionLink(id types.IdType) (*authority.PermissionLink, error) {
	key := getObjectKeyByIndex(&authority.PermissionLink{ID: id}, "id")
	item, err := s.transaction.Get(key)
	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	out := &authority.PermissionLink{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindPermissionLinkByActionName(account name.AccountName, code name.AccountName, messageType name.ActionName) (*authority.PermissionLink, error) {
	key := getObjectKeyByIndex(&authority.PermissionLink{Account: account, Code: code, MessageType: messageType}, "byActionName")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindPermissionLink(types.NewIdType(data))
}

func (s *Session) CreatePermissionLink(in *authority.PermissionLink) error {
	return s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)
}

func (s *Session) ModifyPermissionLink(in *authority.PermissionLink, modifyFunc func()) error {
	return s.modify(in, modifyFunc)
}

func (s *Session) RemovePermissionLink(in *authority.PermissionLink) error {
	return s.remove(in)
}
