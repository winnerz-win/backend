package jmap

import (
	"jtools/cc"
	"sync/atomic"
)

type KeyLocker[T any] interface {
	Lock(key T) bool
	Unlock(key T)
	Action(key T, f func()) bool

	Keys() []T
	RemoveKey(key T) bool
}

type async_key_lock[T any] struct {
	data MAP[T, *int32]
}

func NewKeyLock[T any]() KeyLocker[T] {
	return &async_key_lock[T]{
		data: New[T, *int32](),
	}
}

func (my *async_key_lock[T]) Lock(key T) bool {
	lock_value := int32(1)
	if pre_lock, loaded := my.data.LoadOrStore(key, &lock_value); loaded {
		if !atomic.CompareAndSwapInt32(pre_lock, 0, 1) {
			return false
		}
	}
	return true
}

func (my *async_key_lock[T]) Unlock(key T) {
	un_lock_value := int32(0)
	my.data.Store(key, &un_lock_value)
}

func (my *async_key_lock[T]) Action(key T, f func()) bool {
	if my.Lock(key) {
		defer func() {
			if e := recover(); e != nil {
				cc.Red(e)
			}
			my.Unlock(key)
		}()

		f()

		return true
	}
	return false
}

func (my async_key_lock[T]) Keys() []T {
	keys := []T{}
	my.data.Range(func(key T, value *int32) (isContinue bool) {
		keys = append(keys, key)
		return true
	})
	return keys
}

func (my *async_key_lock[T]) RemoveKey(key T) bool {
	_, loaded := my.data.LoadAndDelete(key)
	return loaded
}
