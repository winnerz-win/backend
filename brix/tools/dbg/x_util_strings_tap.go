package dbg

type xTextAlignTap struct {
	items   []string
	taps    []string
	isAlign bool
}

func NewTextAlignTap() TextAlignTapper {
	return &xTextAlignTap{}
}

type TextAlignTapper interface {
	Count() int
	AddText(text string)
	AddTexts(texts ...string)
	SetAlign(space string)
	GetItems() []string
	GetTaps() []string
	GetPair(i int) (string, string)
	For(callback func(i int, item, tap string))
}

func (my xTextAlignTap) Count() int { return len(my.items) }

func (my *xTextAlignTap) AddText(text string) {
	my.items = append(my.items, text)
}
func (my *xTextAlignTap) AddTexts(texts ...string) {
	my.items = append(my.items, texts...)
}

func (my *xTextAlignTap) SetAlign(space string) {
	my.taps = []string{}
	my.taps = append(my.taps, my.items...)

	max := 0
	for _, name := range my.taps {
		size := len(name)
		if max < size {
			max = size
		}
	}

	for i := 0; i < len(my.taps); i++ {
		cnt := max - len(my.taps[i])
		n := ""
		for i := 0; i < cnt; i++ {
			n += " "
		}
		n += space
		my.taps[i] = n
	}

	my.isAlign = true

}

func (my xTextAlignTap) GetItems() []string {
	return my.items
}

func (my xTextAlignTap) GetTaps() []string {
	return my.taps
}

func (my *xTextAlignTap) GetPair(i int) (string, string) {
	if !my.isAlign {
		my.SetAlign("")
	}
	if i < 0 || i >= len(my.items) {
		i = 0
	}
	return my.items[i], my.taps[i]
}

func (my *xTextAlignTap) For(callback func(i int, item, tap string)) {
	if !my.isAlign {
		my.SetAlign("")
	}
	for i, item := range my.items {
		callback(i, item, my.taps[i])
	}
}
