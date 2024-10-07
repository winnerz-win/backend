package mongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"jtools/cc"
	"reflect"
	"runtime/debug"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (my IDS) SelectorMongoID() Bson { return Bson{"_id": my.ID} }
func (my IDS) View() {
	switch v := my.ID.(type) {
	case primitive.ObjectID:
		cc.Yellow("aa", v)
	case string:
		cc.Green("bb:", v)
	default:
		cc.Red("cc", v)
	}
}

type IndexingDBFunc interface {
	IndexingDB()
}

func StartIndexingDB(list ...IndexingDBFunc) {
	defer cc.YellowItalic("indexingDB ------ END")
	cc.YellowItalic("indexingDB ------ START")

	total := len(list)
	for i, v := range list {
		v.IndexingDB()

		cc.YellowItalic("indexingDB[", reflect.TypeOf(v).Name(), "] (", i+1, "/", total, ")")
	}
}

// PrefixDot :
func PrefixDot(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return prefix
	}

	if !strings.HasPrefix(prefix, ".") {
		prefix = prefix + "."
	}
	return prefix
}

func toString(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func _Cat(a ...interface{}) string {
	sl := []string{}
	for _, v := range a {
		sl = append(sl, fmt.Sprintf("%v", v))
	}
	return strings.Join(sl, "")
}

func _Error(a ...interface{}) error {
	if len(a) == 0 {
		return errors.New("dbg.Error")
	}
	sl := []string{}
	for i := 0; i < len(a); i++ {
		sl = append(sl, "%v")
	}
	sf := strings.Join(sl, " ")

	return fmt.Errorf(sf, a...)

}

func parseStruct(src, dst interface{}) error {
	if src == nil {
		return _Error("src is nil")
	}
	if dst == nil {
		return _Error("dst is nil")
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, dst)
	case string:
		return json.Unmarshal([]byte(v), dst)
	}

	b, err := json.Marshal(src)
	if err != nil {
		return _Error("[Marshal]", err)
	}

	return json.Unmarshal(b, dst)
}

func _Stack(skip_line_ ...int) string {
	skip_line := 5
	if len(skip_line_) > 0 && skip_line_[0] > skip_line {
		skip_line = skip_line_[0]
	}
	sl := strings.Split(string(debug.Stack()), "\n")
	if len(sl) >= skip_line {
		/*
			[ 0 ] goroutine 6 [running]:
			[ 1 ] runtime/debug.Stack()
			[ 2 ] 	C:/Program Files/Go/src/runtime/debug/stack.go:24 +0x7a
			[ 3 ] jtools/dbg._print_stack()
			[ 4 ] 	d:/work/go/src/brix_pkg/jtools/dbg/error.go:32 +0x2e
		*/
		sl = sl[skip_line:]
	}
	// for i, v := range sl {
	// 	Println("[", i, "]", v)
	// }

	return strings.Join(sl, "\n")
}

func viewErrStack(e interface{}) {
	defer cc.Red("----------------------------------------------------")
	cc.RedItalic("[error]", e)
	cc.RedItalic(_Stack(7))

}

func isTrue(p interface{}) bool {
	switch v := p.(type) {
	case bool:
		return v
	case []bool:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		return strings.ToLower(strings.TrimSpace(v)) == "true"
	case []interface{}:
		if len(v) > 0 {
			return isTrue(v[0])
		}
	}
	return false
}
