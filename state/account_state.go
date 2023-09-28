package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/account"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/crypto"
)

func (s *Session) FindAccount(id types.IdType) (*account.Account, error) {
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
	if _, err := Codec.Unmarshal(data, out); err != nil {
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

	return s.FindAccount(types.NewIdType(data))
}

func (s *Session) CreateAccount(in *account.Account) error {
	err := s.create(true, func(id types.IdType) error {
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

func (s *Session) FindAccountMetaData(id types.IdType) (*account.AccountMetaDataObject, error) {
	key := getObjectKeyByIndex(&account.AccountMetaDataObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &account.AccountMetaDataObject{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindAccountMetaDataByName(name name.AccountName) (*account.AccountMetaDataObject, error) {
	key := getObjectKeyByIndex(&account.AccountMetaDataObject{Name: name}, "byName")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindAccountMetaData(types.NewIdType(data))
}

func (s *Session) CreateAccountMetaData(in *account.AccountMetaDataObject) error {
	err := s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ModifyAccountMetaData(in *account.AccountMetaDataObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	return nil
}

func (s *Session) FindCodeObject(id types.IdType) (*account.CodeObject, error) {
	key := getObjectKeyByIndex(&account.CodeObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &account.CodeObject{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) FindCodeObjectByCodeHash(codeHash crypto.Sha256, vmType uint8, vmVersion uint8) (*account.CodeObject, error) {
	key := getObjectKeyByIndex(&account.CodeObject{CodeHash: codeHash, VmType: vmType, VmVersion: vmVersion}, "byCodeHash")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindCodeObject(types.NewIdType(data))
}

func (s *Session) CreateCodeObject(in *account.CodeObject) error {
	err := s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ModifyCodeObject(in *account.CodeObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	return nil
}

func (s *Session) RemoveCodeObject(in *account.CodeObject) error {
	return s.remove(in)
}
