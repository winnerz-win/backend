package jmap

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

type ConstraintsKeySet interface {
	constraints.Ordered | constraints.Complex
}

type IKeySet[T ConstraintsKeySet] interface {
	SetExist(key T) bool
	SetFirst(key T) bool
	GetList(sort_opt ...func(i, j T) bool) []T
	Do(key T) bool
	Clear()
}

////////////////////////////////////////////////////////////////////////////////

type KeySet[T ConstraintsKeySet] map[T]struct{}

// SetExist : key값이 중복되지면 true
func (my KeySet[T]) SetExist(key T) bool {
	_, do := my[key]
	my[key] = struct{}{}
	return do
}

// SetFirst : key값이 중복되지 않으면 true
func (my KeySet[T]) SetFirst(key T) bool {
	return !my.SetExist(key)
}

func (my KeySet[T]) SetList(keys ...T) {
	for _, key := range keys {
		my[key] = struct{}{}
	}
}

func (my KeySet[T]) GetList(sort_opt ...func(i, j T) bool) []T {
	list := make([]T, len(my))
	i := 0
	for key, _ := range my {
		list[i] = key
		i++
	}

	if len(sort_opt) > 0 {
		if sort_opt[0] != nil {
			sort.Slice(list, func(i, j int) bool {
				return sort_opt[0](list[i], list[j])
			})
		}
	}

	return list
}

func (my KeySet[T]) Do(key T) bool {
	_, do := my[key]
	return do
}

func (my *KeySet[T]) Clear() {
	*my = KeySet[T]{}
}

////////////////////////////////////////////////////////////////////////////////

type syncKeySet[T ConstraintsKeySet] struct {
	set KeySet[T]
	mu  sync.RWMutex
}

func SyncKeySet[T ConstraintsKeySet]() IKeySet[T] {
	set := KeySet[T]{}
	my := &syncKeySet[T]{set: set}
	return my
}

// SetExist : key값이 중복되지면 true
func (my *syncKeySet[T]) SetExist(key T) bool {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.set.SetExist(key)
}

// SetFirst : key값이 중복되지 않으면 true
func (my *syncKeySet[T]) SetFirst(key T) bool {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.set.SetFirst(key)
}
func (my *syncKeySet[T]) GetList(sort_opt ...func(i, j T) bool) []T {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.set.GetList(sort_opt...)
}
func (my *syncKeySet[T]) Do(key T) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.set.Do(key)
}
func (my *syncKeySet[T]) Clear() {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.set.Clear()
}
