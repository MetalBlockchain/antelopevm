package state

import "github.com/dgraph-io/badger/v3"

type Iterator[T any] struct {
	iterator   *badger.Iterator
	opts       badger.IteratorOptions
	lookupFunc func([]byte) (*T, error)
}

func newIterator[T any](session *Session, prefix []byte, lookupFunc func([]byte) (*T, error)) *Iterator[T] {
	opts := badger.DefaultIteratorOptions

	if len(prefix) > 0 {
		opts.Prefix = prefix
	}

	iterator := session.transaction.NewIterator(opts)

	return &Iterator[T]{
		iterator:   iterator,
		opts:       opts,
		lookupFunc: lookupFunc,
	}
}

func newReverseIterator[T any](session *Session, prefix []byte, lookupFunc func([]byte) (*T, error)) *Iterator[T] {
	opts := badger.DefaultIteratorOptions
	opts.Reverse = true

	if len(prefix) > 0 {
		opts.Prefix = prefix
	}

	iterator := session.transaction.NewIterator(opts)

	return &Iterator[T]{
		iterator:   iterator,
		opts:       opts,
		lookupFunc: lookupFunc,
	}
}

func (i *Iterator[T]) Rewind() {
	i.iterator.Rewind()
}

func (i *Iterator[T]) Seek(key []byte) {
	i.iterator.Seek(key)
}

func (i *Iterator[T]) ValidForPrefix(prefix []byte) bool {
	return i.iterator.ValidForPrefix(prefix)
}

func (i *Iterator[T]) Valid() bool {
	if len(i.opts.Prefix) > 0 {
		return i.iterator.ValidForPrefix(i.opts.Prefix)
	}

	return i.iterator.Valid()
}

func (i *Iterator[T]) Next() {
	i.iterator.Next()
}

func (i *Iterator[T]) Item() (*T, error) {
	item := i.iterator.Item()
	value, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return i.lookupFunc(value)
}

func (i *Iterator[T]) Close() {
	i.iterator.Close()
}
