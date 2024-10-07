package jmap

import "sync"

type MAP[TKey, TValue any] interface {
	Delete(key TKey)
	Load(key TKey) (value TValue, ok bool)

	/*
		1. key값이 있으면 삭제후 삭제된 value를 반환 loaded = true |
		2. key값이 없으면 default-value 반환 loaded = false
	*/
	LoadAndDelete(key TKey) (value TValue, loaded bool)

	/*
		1. key값이 있으면 해당 value를 반환 loaded = true |
		2. key값이 없으면 value를 Store한후 value를 반환 loaded = false
	*/
	LoadOrStore(key TKey, value TValue) (actual TValue, loaded bool)

	/*
		맵을 순회 : callback 반환값이 false 이면 종료.
	*/
	Range(f func(key TKey, value TValue) bool)

	Store(key TKey, value TValue)
}

type syncMap[TKey, TValue any] struct {
	data sync.Map
}

func New[TKey, TValue any]() MAP[TKey, TValue] {
	return &syncMap[TKey, TValue]{}
}

func (my *syncMap[TKey, TValue]) Delete(key TKey) {
	my.data.Delete(key)
}
func (my *syncMap[TKey, TValue]) Load(key TKey) (value TValue, ok bool) {
	if v, ok := my.data.Load(key); ok {
		return v.(TValue), ok
	}
	return
}
func (my *syncMap[TKey, TValue]) LoadAndDelete(key TKey) (value TValue, loaded bool) {
	if val, loaded := my.data.LoadAndDelete(key); loaded {
		return val.(TValue), loaded
	}
	return
}
func (my *syncMap[TKey, TValue]) LoadOrStore(key TKey, value TValue) (actual TValue, loaded bool) {
	val, loaded := my.data.LoadOrStore(key, value)
	return val.(TValue), loaded
}
func (my *syncMap[TKey, TValue]) Range(f func(key TKey, value TValue) bool) {
	my.data.Range(
		func(key, value any) bool {
			return f(key.(TKey), value.(TValue))
		},
	)
}
func (my *syncMap[TKey, TValue]) Store(key TKey, value TValue) {
	my.data.Store(key, value)
}
