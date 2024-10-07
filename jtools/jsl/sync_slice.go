package jsl

import (
	"fmt"
	"sort"
	"sync"
)

type SLICE[T any] interface {
	Append(a ...T)
	Sort(f func(i, j T) bool)
	Remove(idx int) bool
	Clear() []T
	////////////////////////////////////////////////////
	String() string
	Count() int
	List() []T
	Foreach(f func(idx int, item T) bool)
}

type syncSlice[T any] struct {
	mu   sync.RWMutex
	list []T
}

func New[T any]() SLICE[T] {
	return &syncSlice[T]{}
}

func Get[T any](list []T) SLICE[T] {
	return &syncSlice[T]{
		list: list,
	}
}

func (my *syncSlice[T]) Append(a ...T) {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.list = append(my.list, a...)
}

func (my *syncSlice[T]) Sort(f func(i, j T) bool) {
	my.mu.Lock()
	defer my.mu.Unlock()

	sort.Slice(my.list, func(i, j int) bool {
		return f(my.list[i], my.list[j])
	})
}

func (my *syncSlice[T]) Remove(idx int) bool {
	my.mu.Lock()
	defer my.mu.Unlock()
	if idx < 0 || idx >= len(my.list) {
		return false
	}

	my.list = append(my.list[:idx], my.list[idx+1:]...)

	return true
}
func (my *syncSlice[T]) Clear() []T {
	my.mu.Lock()
	defer my.mu.Unlock()
	if len(my.list) > 0 {
		sl := my.list
		my.list = []T{}
		return sl
	}
	return []T{}
}

////////////////////////////////////////////////////

func (my *syncSlice[T]) String() string {
	my.mu.RLock()
	defer my.mu.RUnlock()
	return fmt.Sprint(my.list)
}

func (my *syncSlice[T]) Count() int {
	my.mu.RLock()
	defer my.mu.RUnlock()
	return len(my.list)
}

func (my *syncSlice[T]) List() []T {
	my.mu.RLock()
	defer my.mu.RUnlock()
	sl := []T{}
	sl = append(sl, my.list...)
	return sl
}

func (my *syncSlice[T]) Foreach(f func(idx int, item T) bool) {
	my.mu.RLock()
	defer my.mu.RUnlock()
	for i, v := range my.list {
		if f(i, v) {
			break
		}
	} //for
}
