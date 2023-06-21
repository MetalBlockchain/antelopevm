package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func (s *Session) FindPermissionLink(id core.IdType) (*core.PermissionLink, error) {
	key := getObjectKeyByIndex(&core.PermissionLink{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &core.PermissionLink{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindPermissionLinkByActionName(account name.AccountName, code name.AccountName, messageType name.ActionName) (*core.PermissionLink, error) {
	key := getObjectKeyByIndex(&core.PermissionLink{Account: account, Code: code, MessageType: messageType}, "byActionName")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindPermissionLink(core.NewIdType(data))
}

func (s *Session) CreatePermissionLink(in *core.PermissionLink) error {
	return s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)
}

func (s *Session) ModifyPermissionLink(in *core.PermissionLink, modifyFunc func()) error {
	return s.modify(in, modifyFunc)
}

func (s *Session) RemovePermissionLink(in *core.PermissionLink) error {
	return s.remove(in)
}
