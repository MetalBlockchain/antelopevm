package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/global"
)

func (s *Session) FindGlobalPropertyObject(id core.IdType) (*global.GlobalPropertyObject, error) {
	key := getObjectKeyByIndex(&global.GlobalPropertyObject{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &global.GlobalPropertyObject{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Session) CreateGlobalPropertyObject(in *global.GlobalPropertyObject) error {
	err := s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ModifyGlobalPropertyObject(in *global.GlobalPropertyObject, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	return nil
}
