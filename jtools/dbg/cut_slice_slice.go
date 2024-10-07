package dbg

import "fmt"

type CutMap[T any] struct {
	data map[string]int
}

func (my *CutMap[T]) _set(t_key T, index int) {
	key := fmt.Sprintf("%v", t_key)
	my.data[key] = index
}

//Do : key가 존재하는지 단순 검사
func (my CutMap[T]) Do(t_key T) (int, bool) {
	key := fmt.Sprintf("%v", t_key)
	index, do := my.data[key]
	return index, do
}

//Do : key가 존재하는지 검사후 key삭제
func (my *CutMap[T]) DoRmv(t_key T) (int, bool) {
	key := fmt.Sprintf("%v", t_key)
	index, do := my.data[key]
	if do {
		delete(my.data, key)
	}
	return index, do
}

//Foreach : f(인덱스,키)
func (my CutMap[T]) Foreach(f func(i int, key string)) {
	for key, index := range my.data {
		f(index, key)
	}
}

func CutSliceMap[T any](dst []T, cutCnt int) ([][]T, CutMap[T]) {
	cm := CutMap[T]{
		data: map[string]int{},
	}
	for i, v := range dst {
		cm._set(v, i)
	} //for

	re := [][]T{}
	for {
		if len(dst) <= cutCnt {
			if len(dst) == 0 {
				return re, cm
			}
			re = append(re, dst)
			return re, cm
		}
		elem := dst[:cutCnt]
		re = append(re, elem)
		dst = dst[cutCnt:]
	} //for
}

func CutSliceSlice[T any](dst []T, cutCnt int) [][]T {
	re := [][]T{}
	for {
		if len(dst) <= cutCnt {
			if len(dst) == 0 {
				return re
			}
			re = append(re, dst)
			return re
		}
		elem := dst[:cutCnt]
		re = append(re, elem)
		dst = dst[cutCnt:]
	} //for
}

func CutSliceSlice2[T any](dst []T, cutCnt int) ([][]T, [][]int) {
	re := [][]T{}
	pos := [][]int{}
	last := 0
	for {
		if len(dst) <= cutCnt {
			if len(dst) == 0 {
				return re, pos
			}
			re = append(re, dst)

			p := make([]int, len(dst))
			for i := 0; i < len(dst); i++ {
				p[i] = last
				last++
			} //for
			pos = append(pos, p)

			return re, pos
		}
		elem := dst[:cutCnt]
		re = append(re, elem)

		p := make([]int, cutCnt)
		for i := 0; i < cutCnt; i++ {
			p[i] = last
			last++
		} //for
		pos = append(pos, p)

		dst = dst[cutCnt:]
	} //for

}

func CloneSlice[T any](src []T) []T {
	clone := make([]T, len(src))
	copy(clone, src)
	return clone
}
