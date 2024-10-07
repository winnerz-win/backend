package chttp

type cTextAlignTapper struct {
	keys     []string
	vals     []interface{}
	taps     []string
	isTapped bool
	space    string
}

func (my cTextAlignTapper) Count() int {
	return len(my.keys)
}

/////////////////////////////////////////////////////////////////////

type TextAlignTapperPair interface {
	SetKeyPairs(sl ...interface{}) TextAlignTapperPair
	ForPair(space string, f func(i int, item, tap string, val interface{}))
	Count() int
}

func NewTextAlignTapperPair() TextAlignTapperPair {
	return &cTextAlignTapper{}
}

func (my *cTextAlignTapper) SetKeyPairs(sl ...interface{}) TextAlignTapperPair {
	my.isTapped = false
	for i := 0; i < len(sl); i += 2 {
		my.keys = append(my.keys, Cat(sl[i]))
		my.vals = append(my.vals, sl[i+1])
	}
	return my
}

func (my *cTextAlignTapper) ForPair(space string, f func(i int, item, tap string, val interface{})) {
	if !my.isTapped || my.space != space {
		my._set_align(space)
	}
	for i := 0; i < len(my.keys); i++ {
		f(i, my.keys[i], my.taps[i], my.vals[i])
	}
}

/////////////////////////////////////////////////////////////////////

type TextAlignTapperKey interface {
	SetKey(sl ...interface{}) TextAlignTapperKey
	For(space string, f func(i int, item, tap string))
	Count() int
}

func NewTextAlignTapperKey() TextAlignTapperKey {
	return &cTextAlignTapper{}
}
func (my *cTextAlignTapper) SetKey(sl ...interface{}) TextAlignTapperKey {
	my.isTapped = false
	for i := 0; i < len(sl); i++ {
		my.keys = append(my.keys, Cat(sl[i]))
	}

	return my
}

func (my *cTextAlignTapper) For(space string, f func(i int, item, tap string)) {
	if !my.isTapped || my.space != space {
		my._set_align(space)
	}
	for i := 0; i < len(my.keys); i++ {
		f(i, my.keys[i], my.taps[i])
	}
}

func (my *cTextAlignTapper) _set_align(space string) {
	my.taps = []string{}

	max := 0
	for _, key := range my.keys {
		size := len(key)
		if max < size {
			max = size
		}
	} //for

	for _, key := range my.keys {
		cnt := max - len(key)
		n := ""
		for cnt > 0 {
			n += " "
			cnt--
		}
		my.taps = append(my.taps, n+space)
	}
	my.space = space
	my.isTapped = true
}
