package doc

import (
	"fmt"
	"strings"
)

const (
	Struct     = "struct{}"
	ClassArray = "Array(struct)"
	String     = "string"
	Bool       = "bool"
	Int        = "int"
	Int64      = "int64"

	Strings = "[]string"
	Bools   = "[]bool"
	Ints    = "[]int"
	Int64s  = "[]int64"
	Bytes   = "[]bytes"

	webSocketTag = "{ws}"
	notKey       = "-D-"
)

type cURL struct {
	do   bool
	path string
	kvs  map[string][]string // [key]string
}

func (my cURL) ToArray() []string {
	urls := []string{}
	isFirst := true
	for key, list := range my.kvs {
		if isFirst {
			for _, str := range list {
				path := strings.ReplaceAll(my.path, key, str)
				urls = append(urls, path)
			} //for
		} else {
			clones := urls
			urls = []string{}
			for _, str := range list {
				for _, clone := range clones {
					path := strings.ReplaceAll(clone, key, str)
					urls = append(urls, path)
				}
			} //for
		}
	} //for

	return urls
}

//Doc :
type Doc struct {
	comment string
	url     string
	urls    cURL
	method  string
	headers keyValueList
	params  keyValueList
	jparam  *etcValue // json-type params
	//getparam string    //get - param

	results   map[int]keyValueList
	jResultOK *etcValue // json-type params
	jAckOK    *etcValue // ack-format
	jAckError []AckErrorPair

	resultOrder []int
	resultERRR  []errrData
	//etcs        []interface{}
	etcs   []etcValue
	etcTag map[int][]string
}

type errrData struct {
	err IError
	tag string
}

type etcValue struct {
	void interface{}
	tags map[string]string
}

func (my etcValue) getTag(jsonKey string) string {
	jsonKey = strings.TrimSpace(jsonKey)
	if tag, do := my.tags[jsonKey]; do == true {
		return "   // " + tag
	}
	return ""
}

//EV :
func EV(void interface{}, pair ...string) etcValue {
	tmap := map[string]string{}
	for i := 0; i < len(pair); i += 2 {
		key := strings.TrimSpace(pair[i])
		val := strings.TrimSpace(pair[i+1])
		tmap[key] = val
	}
	return etcValue{
		void: void,
		tags: tmap,
	}
}

func newDoc(comment string) *Doc {
	return &Doc{
		comment: strings.TrimSpace(comment),
		results: map[int]keyValueList{},
		etcTag:  map[int][]string{},
	}
}

type Item struct {
	*Doc
	color Color
}
type ItemList []Item

//Add : doc.cItem
func (my *ItemList) Add(item Item) {
	*my = append(*my, item)
}

func NewItem(comment string, c ...Color) Item {
	item := Item{
		Doc:   newDoc(comment),
		color: Black,
	}
	if len(c) > 0 {
		item.color = c[0]
	}
	return item
}

//KV :
type keyValue struct {
	Key    interface{}
	_space string
	Val    interface{}
	Tag    string
}

//KVList :
type keyValueList []keyValue

func (my keyValueList) KeyAlign() {
	max := 0
	for _, data := range my {
		strKey, do := data.Key.(string)
		if do == false {
			continue
		}
		if max < len(strKey) {
			max = len(strKey)
		}
	} //for

	for index, data := range my {
		strKey, do := data.Key.(string)
		if do == false {
			continue
		}
		if len(strKey) < max {
			space := ""
			spaceCnt := max - len(strKey) + 1
			for i := 0; i < spaceCnt; i++ {
				space += " "
			}
			my[index]._space = space
		}
	} //for

}

func (my keyValue) getKey() string {
	return fmt.Sprintf(`"%v"%v`, my.Key, my._space)
}
func (my keyValue) getVal() string {
	return fmt.Sprintf("%v", my.Val)
}

func (my keyValue) String(space string, isNotNewLine ...bool) string {
	newLine := ",\n"
	if len(isNotNewLine) > 0 && isNotNewLine[0] == true {
		newLine = ""
	}
	my.Tag = strings.TrimSpace(my.Tag)

	if my.Key == notKey {
		if my.Tag != "" {
			return space + my.getVal() + `	// ` + my.Tag + newLine
		}
		return space + my.getVal() + newLine
	}

	if my.Tag != "" {
		return space + my.getKey() + " : " + my.getVal() + `	// ` + my.Tag + newLine
	}
	return space + my.getKey() + " : " + my.getVal() + newLine
}

//HTML :
func (my Doc) HTML() string {
	return stringLineToTag(my.String())
}

func stringLineToTag(s string) string {
	newLine := []byte("<br>")
	space1 := []byte("&nbsp;")
	space2 := []byte("&nbsp;&nbsp;&nbsp;&nbsp;")
	b := []byte(s)
	hbuf := []byte{}
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' {
			hbuf = append(hbuf, newLine...)
		} else if b[i] == ' ' {
			hbuf = append(hbuf, space1...)
		} else if b[i] == '\t' {
			hbuf = append(hbuf, space2...)
		} else {
			hbuf = append(hbuf, b[i])
		}
	}

	return string(hbuf)
}

//View :
func (my *Doc) View() *Doc {
	fmt.Println(my.String())
	return my
}
