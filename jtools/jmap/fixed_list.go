package jmap

type fixed_list[T any] struct {
	list   []T
	is_set []bool
}

func (my fixed_list[T]) Size() int {
	return len(my.list)
}
func (my fixed_list[T]) List() []T {
	sl := []T{}
	for i, v := range my.list {
		if my.is_set[i] {
			sl = append(sl, v)
		}
	} //for
	return sl
}

func (my fixed_list[T]) Count() int {
	cnt := 0
	for _, is_set := range my.is_set {
		if is_set {
			cnt++
		}
	} //for
	return cnt
}

func (my *fixed_list[T]) Set(index int, data T) {
	if index < 0 || index >= len(my.list) {
		return
	}
	my.list[index] = data
	my.is_set[index] = true
}

func (my fixed_list[T]) Do() (T, bool) {
	for i, is_set := range my.is_set {
		if is_set {
			return my.list[i], true
		}
	} //for
	var empty T
	return empty, false
}

func MakeFixedList[T any](count int) fixed_list[T] {
	my := fixed_list[T]{
		list:   make([]T, count),
		is_set: make([]bool, count),
	}
	return my
}
