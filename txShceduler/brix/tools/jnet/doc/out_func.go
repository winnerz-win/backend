package doc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func etcValuePtrString(kind string, ev *etcValue, spaces ...string) string {
	msg := ""
	add := func(s string) {
		msg += " " + s
	}

	toETCString := func(v interface{}) string {
		if v == nil {
			return ""
		}
		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			return fmt.Sprintln("dbg.ViewJSON :", err, "\n")

		} else if json.Valid(b) == false {
			return fmt.Sprintln("dbg.ViewJSON is not Valid :", string(b), "\n")

		}
		space := []byte(" ")
		_ = space
		hbuf := []byte{}
		hbuf = append(hbuf, space...)
		for i := 0; i < len(b); i++ {
			if b[i] == '\n' {
				hbuf = append(hbuf, space...)
			} else {
				hbuf = append(hbuf, b[i])
			}
		}
		return string(hbuf) + "\n"
	}

	fInline := func(tap string, t reflect.Type, ev *etcValue) {
		space := tap
		if t.Kind() == reflect.Struct {
			add(space + "< " + t.String() + " >\n")
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				tagName, ok := field.Tag.Lookup("json")
				tag := ev.getTag(tagName)
				if ok {
					if tagName == "-" || tagName == "" {
						continue
					}
					add(space + strings.TrimSpace(tagName) + " : " + field.Type.String() + tag + "\n")
				} else {
					add(space + strings.TrimSpace(field.Name) + " : " + field.Type.String() + tag + "\n")
				}
			} //for
			add(space + "----------\n")
		} else {
			add(space + fmt.Sprintf("%v\n", t.Name()))
		}
	}

	space := ""
	if len(spaces) > 0 {
		space = spaces[0]
	}
	parseStruct := func(rt reflect.Type, isArray bool, sp string) {
		space := "    "
		space += sp

		//add(rt.Name() + " {\n")
		if isArray == false {
			add(kind + " {\n")
		} else {
			add(kind + " [ //array \n")
			add(sp + " {\n")
		}

		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			tagName, ok := field.Tag.Lookup("json")
			if tagName == ",inline" {
				fInline(space, field.Type, ev)
				continue
			}
			tag := ev.getTag(tagName)
			if ok {
				if tagName == "-" || tagName == "" {
					continue
				}
				add(space + strings.TrimSpace(tagName) + " : " + field.Type.String() + tag + "\n")
			} else {
				add(space + strings.TrimSpace(field.Name) + " : " + field.Type.String() + tag + "\n")
			}
		}

		if isArray == false {
			add(sp + "}\n")
		} else {
			add(sp + " },\n" + sp + "  。。。。\n")
			add(sp + "]\n")
		}
	}

	rt := reflect.TypeOf(ev.void)
	if rt.Kind() == reflect.Struct {
		parseStruct(rt, false, space)

	} else if rt.Kind() == reflect.Slice {
		elem := rt.Elem()
		kind := elem.Kind()
		if kind == reflect.Struct {
			//reflect.MakeSlice(reflect.SliceOf(elem))
			rv := reflect.New(elem)
			parseStruct(rv.Elem().Type(), true, space)
		} else {
			add(space + fmt.Sprintf("%v\n", toETCString(ev.void)))
		}
	} else {
		add(space + fmt.Sprintf("%v\n", toETCString(ev.void)))
	}
	return msg
}

func StructString(name string, v interface{}, tag ...string) string {
	ev := etcValue{
		void: v,
		tags: map[string]string{},
	}
	for i := 0; i < len(tag); i += 2 {
		key := tag[i]
		val := tag[i+1]
		ev.tags[key] = val
	} //for
	return etcValuePtrString(name, &ev)
}
