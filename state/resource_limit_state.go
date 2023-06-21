package state

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/core/resource"
)

func (s *Session) FindResourceLimits(id core.IdType) (*resource.ResourceLimits, error) {
	if obj, found := s.resourceLimitsCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&resource.ResourceLimits{ID: id}, "id")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	out := &resource.ResourceLimits{}

	if _, err := out.UnmarshalMsg(data); err != nil {
		return nil, err
	}

	s.resourceLimitsCache.Put(id, out)

	return out, nil
}

func (s *Session) FindResourceLimitsByOwner(pending bool, name name.AccountName) (*resource.ResourceLimits, error) {
	key := getObjectKeyByIndex(&resource.ResourceLimits{Pending: pending, Owner: name}, "byOwner")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindResourceLimits(core.NewIdType(data))
}

func (s *Session) CreateResourceLimits(in *resource.ResourceLimits) error {
	err := s.create(true, func(id core.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.resourceLimitsCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyResourceLimits(in *resource.ResourceLimits, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.resourceLimitsCache.Put(in.ID, in)

	return nil
}
