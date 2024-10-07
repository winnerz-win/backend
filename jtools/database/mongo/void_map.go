package mongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"jtools/cc"
	"jtools/dbg"
	"jtools/jmath"
	"jtools/unix"
	"strings"
)

var (
	Null = errors.New("null")
)

type MAPLIST []MAP

func (my VOIDS) MAPLIST() MAPLIST {
	list := MAPLIST{}
	for _, v := range my {
		list = append(list, MAP(v))
	}
	return list
}

func (my MAPLIST) String() string { return toString(my) }
func (my MAPLIST) First() MAP {
	if len(my) > 0 {
		return my[0]
	}
	return MAP{}
}
func (my MAPLIST) Last() MAP {
	if len(my) > 0 {
		return my[len(my)-1]
	}
	return MAP{}
}

/////////////////////////////////////////////////////////////////////////////////////////////

func (my VOID) MAP() MAP {
	return MAP(my)
}

type MAP map[string]interface{}

func (my MAP) String() string { return toString(my) }

// IdSelector : { "_id": my.Value("_id") }
func (my MAP) IDSelector() Bson {
	return Bson{
		"_id": my.Value("_id"),
	}
}

func (my MAP) Error() error {
	if my == nil {
		return Null
	}
	return nil
}

// TraceOnKeyFromList : <-- TraceKeyOne
func (my MAP) TraceOneKeyFromList(keys ...string) string {
	for _, key := range keys {
		if _, do := my[key]; do {
			return key
		}
	}
	return ""
}

func (my MAP) Clone() MAP {
	return MakeMap(my)
}

func MakeMap(p interface{}) MAP {
	switch v := p.(type) {

	case MAP:
		clone := MAP{}
		for key, val := range v {
			clone[key] = val
		}
		return clone

	case map[string]interface{}:
		return MAP(v)

	case string:
		data := MAP{}
		if err := json.Unmarshal([]byte(v), &data); err == nil {
			return data
		}

	case []byte:
		data := MAP{}
		if err := json.Unmarshal(v, &data); err == nil {
			return data
		}

	case interface{}: // must last position( p type == v type )
		if b, err := json.Marshal(v); err == nil {
			data := MAP{}
			if err := json.Unmarshal(b, &data); err == nil {
				return data
			}
		}

	} //switch
	return nil
}

func (my MAP) Valid() bool                  { return len(my) > 0 }
func (my MAP) Value(key string) interface{} { return my[key] }
func (my MAP) Float64(key string) float64   { return jmath.Float64(my[key]) }

func (my MAP) Text(key string) string {
	if v, do := my[key]; do {
		return _Cat(v)
	}
	return ""
}
func (my MAP) Text_Int64(key string) int64 {
	return jmath.Int64(my.Text(key))
}
func (my MAP) Text_UnixTime(key string) unix.Time {
	return unix.Time(my.Text_Int64(key))
}
func (my MAP) Text_Bool(key string) bool {
	return dbg.IsTrue(my.Text(key))
}

func (my MAP) Parse(key string, p interface{}) error {
	if val, do := my[key]; !do {
		return Null
	} else {
		return parseStruct(val, p)
	}
}
func (my MAP) ParseSelf(p interface{}) error {
	if !my.Valid() {
		return Null
	}
	return parseStruct(my, p)
}

func (my MAP) Get(key string, f func(val interface{})) bool {
	defer func() {
		if e := recover(); e != nil {
			viewErrStack(e)
		}
	}()
	if val, do := my[key]; !do {
		return false
	} else {
		if f != nil {
			f(val)
		}
	}
	return true
}

type SelectorGet map[string]func(val interface{})

func (my MAP) SelectGet(sg SelectorGet) int {
	cnt := 0
	for key, f := range sg {
		if my.Get(key, f) {
			cnt++
		}
	}
	return cnt
}

func (my MAP) DotGet(keysdot string, f func(val interface{})) bool {
	keys := strings.Split(keysdot, ".")

	var re interface{}
	for _, key := range keys {
		if val, do := my[key]; !do {
			return false
		} else {
			re = val
			my = MakeMap(val)
		}
	}
	if f != nil {
		f(re)
	}
	return true
}
func (my MAP) DotValue(keysdot string) interface{} {
	var re interface{}
	my.DotGet(keysdot, func(val interface{}) { re = val })
	return re
}
func (my MAP) DotText(keysdot string) string {
	if v := my.DotValue(keysdot); v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

func (my MAP) MAP(key string, f func(m MAP)) error {
	m := MAP{}

	if do := my.Get(
		key,
		func(val interface{}) {
			m = MakeMap(val)
		},
	); !do {
		return Null
	}

	if err := m.Error(); err != nil {
		return err
	}

	if f != nil {
		f(m)
	}

	return nil
}
func (my MAP) MAPValue(key string) MAP {
	r := MAP{}
	my.MAP(key, func(m MAP) {
		r = m
	})
	return r
}

func (my MAP) MapError(key string, f func(m MAP) error) error {
	var mErr error
	my.Get(
		key,
		func(val interface{}) {
			mErr = f(MakeMap(val))
		},
	)
	return mErr
}

type SelectorMap map[string]func(m MAP) error

func (my MAP) SelectMAP(sm SelectorMap, isDebug ...bool) error {
	for key, f := range sm {
		var err error
		my.Get(
			key,
			func(val interface{}) {
				if isTrue(isDebug) {
					cc.PurpleItalic("SelectMAP[", key, "]")
				}
				err = f(MakeMap(val))
			},
		)
		if err != nil {
			return err
		}
	} //for

	return nil

}

func (my MAP) DotMap(keysdot string, f func(val MAP)) bool {
	keys := strings.Split(keysdot, ".")
	for _, key := range keys {
		if val, do := my[key]; !do {
			return false
		} else {
			my = MakeMap(val)
		}
	}

	if f != nil {
		f(my)
	}

	return true
}

/////////////////////////////////////////////////////////////////////////////////

type LIST []interface{}

func (my LIST) String() string { return toString(my) }
func (my LIST) Error() error {
	if my == nil {
		return Null
	}
	return nil
}

func MakeList(p interface{}) LIST {
	switch v := p.(type) {
	case []interface{}:
		return v
	case string:
		data := []interface{}{}
		if json.Unmarshal([]byte(v), &data) == nil {
			return data
		}

	case interface{}:
		if b, err := json.Marshal(v); err == nil {
			data := []interface{}{}
			if json.Unmarshal(b, &data) == nil {
				return data
			}
		}
	} //switch

	return nil
}

func (my MAP) List(key string, f func(list LIST)) error {
	list := LIST{}
	if do := my.Get(key, func(val interface{}) {
		list = MakeList(val)
	}); !do {
		return Null
	}

	if list.Error() == nil {
		f(list)
	}
	return nil
}

func (my MAP) ListValue(key string) LIST {
	re := LIST{}
	my.List(key, func(list LIST) {
		re = list
	})
	return re
}

func (my LIST) MapList(f func(ml MAPLIST)) error {
	ml := MAPLIST{}
	if err := parseStruct(my, &ml); err != nil {
		viewErrStack(err)
		return err
	}
	if f != nil {
		f(ml)
	}
	return nil
}

func (my LIST) MapListValue() MAPLIST {
	re := MAPLIST{}
	my.MapList(func(ml MAPLIST) {
		re = ml
	})
	return re
}
