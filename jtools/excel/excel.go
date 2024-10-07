package excel

import (
	"jtools/cc"
	"jtools/dbg"
	"sync"

	"github.com/tealeg/xlsx"
)

type Excel struct {
	*xlsx.File
	name string
	sh   map[string]*Sheet
	mu   sync.RWMutex
}

// Excel :
func New(name string) *Excel {
	excel := Excel{
		File: xlsx.NewFile(),
		name: name,
		sh:   map[string]*Sheet{},
	}
	return &excel
}
func (my *Excel) Save() {
	my.File.Save(my.name + ".xlsx")
}
func (my *Excel) Name(onlyName ...bool) string {
	n := my.name

	if !dbg.IsTrue(onlyName) {
		n += ".xlsx"
	}
	return n
}

type Opener struct {
	*xlsx.File
}

func (my Opener) Valid() bool { return my.File != nil }

// Open :
func Open(fullname string) Opener {
	f, err := xlsx.OpenFile(fullname)
	if err != nil {
		cc.Red("excel.OpenFile :", err)
		return Opener{
			File: nil,
		}
	}
	cc.Yellow("excel.OpenFile :", fullname)
	return Opener{
		File: f,
	}
}
