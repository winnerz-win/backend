package spaceup

import (
	"fmt"
	"jtools/dbg"
	"sort"
)

func Align(sl []string, tap_space int) []string {

	getMax := func(tap_space int, sl ...string) int {
		max := 0
		for _, v := range sl {
			if max < len(v) {
				max = len(v)
			}
		}
		if tap_space < 0 {
			tap_space = 0
		}
		max += tap_space
		return max
	}

	max := getMax(tap_space, sl...)

	rl := []string{}
	for _, v := range sl {
		tap := getTap(max, v)
		rl = append(rl, v+tap)
	} //for
	return rl
}

func WithFixed(text string, max_size int, is_force_cut ...bool) string {
	if max_size < 0 {
		max_size = 0
	}
	t_size := len(text)

	if t_size >= max_size {
		if !dbg.IsTrue(is_force_cut) {
			return text
		} else {
			return text[:max_size]
		}
	}

	gap := max_size - t_size
	return text + getSpace(gap)
}

func _toSl(names ...any) []string {
	sl := []string{}
	for _, v := range names {
		sl = append(sl, fmt.Sprintf("%v", v))
	}
	return sl
}

type spaceUpItem struct {
	space     []string
	tap_space int
}

func New(tap_space ...int) *spaceUpItem {
	my := &spaceUpItem{}
	if len(tap_space) > 0 {
		my.tap_space = tap_space[0]
		if my.tap_space < 0 {
			my.tap_space = 0
		}
	}
	return my
}

func (my spaceUpItem) List() []string { return my.space }

func (my spaceUpItem) Text(i int, text string) string {
	if i >= len(my.space) || i < 0 {
		return text
	}
	return text + my.space[i]
}

func (my *spaceUpItem) Add(names ...any) *spaceUpItem {
	return my.Append(my.tap_space, names...)
}

func (my *spaceUpItem) Append(tap_space int, names ...any) *spaceUpItem {
	sl := _toSl(names...)

	max := getMax2(tap_space, my.space, sl)

	rl := []string{}

	is_first := false
	loop := len(my.space)
	if loop == 0 {
		is_first = true
	}

	if len(my.space) < len(sl) {
		loop = len(sl)
	}

	first_text := ""
	for i := 0; i < loop; i++ {
		a := ""
		if i < len(my.space) {
			a = my.space[i]
		} else {
			if !is_first {
				a = getTap(len(first_text), "")
			}
		}
		if i == 0 {
			first_text = a
		}
		b := ""
		if i < len(sl) {
			b = sl[i]
		}

		text := a + b
		tap := getTap(max, text)
		rl = append(rl, text+tap)
	} //for

	my.space = rl

	return my
}

func getMax2(tap_space int, pre, next []string) int {
	max := 0

	is_first := false
	loop := len(pre)
	if loop == 0 {
		is_first = true
	}

	if len(pre) < len(next) {
		loop = len(next)
	}

	first_text := ""

	for i := 0; i < loop; i++ {
		a := ""
		if i < len(pre) {
			a = pre[i]
		} else {
			if !is_first {
				a = getTap(len(first_text), "")
			}
		}
		if i == 0 {
			first_text = a
		}

		b := ""
		if i < len(next) {
			b = next[i]
		}

		text := a + b
		if max < len(text) {
			max = len(text)
		}
	} //for
	if tap_space < 0 {
		tap_space = 0
	}
	max += tap_space
	return max
}

func getTap(cnt int, name string) string {
	cnt = cnt - len(name)
	return getSpace(cnt)
}

func getSpace(cnt int) string {
	space := ""
	for i := 0; i < cnt; i++ {
		space += " "
	}
	return space
}

// ///////////////////////////////////////////////////////////////
type _table struct {
	tap_space int
	texts     [][]string
	max       []int

	line int
}

func Table(tap_space int, names ...any) *_table {
	my := &_table{
		tap_space: tap_space,
		line:      1,
	}
	sl := _toSl(names...)
	my.texts = append(my.texts, sl)

	my.calc_max(sl...)
	return my
}

func (my *_table) calc_max(names ...string) {
	for i, v := range names {
		if i == len(my.max) {
			my.max = append(my.max, len(v))
			continue
		}

		if my.max[i] < len(v) {
			my.max[i] = len(v)
		}
	} //for
}

func (my *_table) Add(names ...any) *_table {
	sl := _toSl(names...)

	if len(names) < len(my.texts[my.line-1]) {
		cnt := len(my.texts[my.line-1]) - len(names)
		for i := 0; i < cnt; i++ {
			sl = append(sl, "")
		}
	} else if len(names) > len(my.texts[my.line-1]) {
		cnt := len(names) - len(my.texts[my.line-1])
		loop := my.line - 1
		for i := loop; i >= 0; i-- {
			for j := 0; j < cnt; j++ {
				my.texts[i] = append(my.texts[i], "")
			}
		}
	}

	my.texts = append(my.texts, sl)
	my.line++
	my.calc_max(sl...)

	return my
}

func (my *_table) List(is_sorts ...bool) []string {
	rl := []string{}
	cells := len(my.texts[0])

	isSort := false
	if len(is_sorts) > 0 && is_sorts[0] {
		isSort = true
	}

	al := []string{}
	for i, list := range my.texts {
		text := ""
		for c := 0; c < cells; c++ {
			v := list[c]
			tap := getTap(my.max[c]+my.tap_space, v)
			text += v + tap
		}
		if i == 0 {
			al = append(al, text)
		} else {
			rl = append(rl, text)
		}
	} //for

	if isSort {
		sort.Slice(rl, func(i, j int) bool {
			return rl[i] < rl[j]
		})
	}

	al = append(al, rl...)
	return al
}
