package jset

import "sync"

//iset :
type iset struct {
	name string
	m    map[interface{}]interface{}
	mu   sync.RWMutex
}

//SET :
type SET interface {
	String() string
	Count() int
	IsEmpty() bool
	IsExist(key interface{}) bool
	Set(key interface{})
	TrySet(key interface{}) bool
	Remove(key interface{})
	Clear()

	SetValue(key, value interface{})
	TrySetValue(key, value interface{}) bool
	GetValue(key interface{}) (interface{}, bool)
	GetValue1(key interface{}) interface{}
	PopValueAll() []interface{}
	ForValue(next func(k, v interface{}, rmv func()))
}

//New :
func New(name ...string) SET {
	instance := &iset{
		name: "jset.SET",
		m:    map[interface{}]interface{}{},
	}
	if len(name) > 0 {
		instance.name = name[0]
	}
	return instance
}

//String
func (my *iset) String() string {
	return my.name
}

//Count :
func (my *iset) Count() int {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return len(my.m)
}

//IsEmpty :
func (my *iset) IsEmpty() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return len(my.m) == 0
}

func (my *iset) IsExist(key interface{}) bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	_, do := my.m[key]
	return do
}

func (my *iset) Set(key interface{}) {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.m[key] = struct{}{}
}

func (my *iset) TrySet(key interface{}) bool {
	defer my.mu.Unlock()
	my.mu.Lock()
	if v, do := my.m[key]; do {
		_ = v
		return false
	}
	my.m[key] = struct{}{}
	return true
}

func (my *iset) remove(key interface{}) {
	delete(my.m, key)
}

func (my *iset) Remove(key interface{}) {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.remove(key)
}

func (my *iset) clear() {
	my.m = map[interface{}]interface{}{}
}
func (my *iset) Clear() {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.clear()
}

func (my *iset) SetValue(key, value interface{}) {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.m[key] = value
}

func (my *iset) TrySetValue(key, value interface{}) bool {
	defer my.mu.Unlock()
	my.mu.Lock()
	if v, do := my.m[key]; do {
		_ = v
		return false
	}
	my.m[key] = value
	return true
}

func (my *iset) GetValue(key interface{}) (interface{}, bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	val, do := my.m[key]
	return val, do
}

func (my *iset) GetValue1(key interface{}) interface{} {
	val, do := my.GetValue(key)
	if do == false {
		return nil
	}
	return val
}

func (my *iset) PopValueAll() []interface{} {
	defer my.mu.Unlock()
	my.mu.Lock()

	list := []interface{}{}
	for _, val := range my.m {
		list = append(list, val)
	}
	my.clear()

	return list
}

//ForValue : key,val , removeCallback func을 함께 반환
func (my *iset) ForValue(next func(k, v interface{}, rmv func())) {
	defer my.mu.Unlock()
	my.mu.Lock()
	for key, val := range my.m {
		rmv := func() {
			delete(my.m, key)
		}
		next(key, val, rmv)
	}
}
