package state

import (
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/resource"
	"github.com/MetalBlockchain/antelopevm/chain/types"
)

func (s *Session) FindResourceUsage(id types.IdType) (*resource.ResourceUsage, error) {
	if obj, found := s.resourceUsageCache.Get(id); found {
		return obj, nil
	}

	key := getObjectKeyByIndex(&resource.ResourceUsage{ID: id}, "id")
	item, err := s.transaction.Get(key)
	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	out := &resource.ResourceUsage{}
	if _, err := Codec.Unmarshal(data, out); err != nil {
		return nil, err
	}

	s.resourceUsageCache.Put(id, out)

	return out, nil
}

func (s *Session) FindResourceUsageByOwner(name name.AccountName) (*resource.ResourceUsage, error) {
	key := getObjectKeyByIndex(&resource.ResourceUsage{Owner: name}, "byOwner")
	item, err := s.transaction.Get(key)

	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return s.FindResourceUsage(types.NewIdType(data))
}

func (s *Session) CreateResourceUsage(in *resource.ResourceUsage) error {
	err := s.create(true, func(id types.IdType) error {
		in.ID = id
		return nil
	}, in)

	if err != nil {
		return err
	}

	s.resourceUsageCache.Put(in.ID, in)

	return nil
}

func (s *Session) ModifyResourceUsage(in *resource.ResourceUsage, modifyFunc func()) error {
	if err := s.modify(in, modifyFunc); err != nil {
		return err
	}

	s.resourceUsageCache.Put(in.ID, in)

	return nil
}
