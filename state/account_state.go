package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/account"
	"github.com/MetalBlockchain/antelopevm/core/name"
)

func (s *Session) FindAccount(id core.IdType) (*account.Account, error) {
	if obj, found := s.accountCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&account.Account{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &account.Account{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	s.accountCache.Put(id, out)

	return out, nil
}

func (s *Session) FindAccountByName(name name.AccountName) (*account.Account, error) {
	key := getObjectKeyByIndex(&account.Account{Name: name}, "byName")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindAccount(core.NewIdType(data))
}

func (s *Session) CreateAccount(in *account.Account) error {
	err := s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.accountCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyAccount(in *account.Account, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.accountCache.Put(in.ID, in)

	return nil
}
