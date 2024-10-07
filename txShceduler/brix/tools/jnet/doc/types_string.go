package doc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const _QT = `"`

func wrap_field_name(v string) string {
	omit_text := ""
	if strings.Contains(v, ",omitempty") {
		v = strings.Replace(v, ",omitempty", "", 1)
		omit_text = "(omitempty)"
	}
	return "<cc_bold>" + _QT + v + _QT + "</cc_bold>" + omit_text
}
func wrap_type_name(v string) string {
	return "<cc_green>" + v + "</cc_green>"
}

func parse_field_inline(add func(s string), tap string, t reflect.Type, ev *etcValue) {
	space := tap
	if t.Kind() == reflect.Struct {
		name := t.String()
		add(space + "< " + name + " >\n")
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tagName, ok := field.Tag.Lookup("json")
			if tagName == ",inline" {
				space_inline := space + tap
				parse_field_inline(add, space_inline, field.Type, ev)
				continue
			}

			tag := ev.getTag(tagName)
			if ok {
				if tagName == "-" || tagName == "" {
					continue
				}
				add(space + strings.TrimSpace(wrap_field_name(tagName)) + " : " + wrap_type_name(field.Type.String()) + tag + "\n")
			} else {
				add(space + strings.TrimSpace(wrap_field_name(field.Name)) + " : " + wrap_type_name(field.Type.String()) + tag + "\n")
			}
		} //for
		add(space + "----------\n")
	} else {
		name := t.Name()
		add(space + fmt.Sprintf("%v\n", name))
	}
}

// String :
func (my Doc) String() string {

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

		} else if !json.Valid(b) {
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

	//Comment
	add(fmt.Sprintf("<cc_purple>Comment</cc_purple> : <cc_bold>%v</cc_bold>\n", my.comment))
	//add(fmt.Sprintf(hColor("Comment", "blue")+" : %v\n", my.comment))

	//URLS
	if my.urls.do {
		for _, url := range my.urls.ToArray() {
			add(fmt.Sprintf("<cc_purple>URL</cc_purple>        : <cc_bold>%v</cc_bold>\n", url))
		}
	} else {
		//URL
		if !strings.HasPrefix(my.url, webSocketTag) {
			add(fmt.Sprintf("<cc_purple>URL</cc_purple>        : <cc_bold>%v</cc_bold>\n", my.url))
		} else {
			my.url = strings.ReplaceAll(my.url, webSocketTag, "")
			add(fmt.Sprintf("<cc_purple>WS</cc_purple>         : <cc_bold>%v</cc_bold>\n", my.url))
		}
	}

	//Method
	add(fmt.Sprintf("<cc_purple>Method</cc_purple>  : <cc_bold>%v</cc_bold>\n", my.method))

	//Headers
	if len(my.headers) < 1 {
		//add(fmt.Sprintf("<cc_purple>Headers</cc_purple>  : %v\n", "{ }"))
	} else {
		hspace := "       "
		add("<cc_purple>Headers</cc_purple>  : {\n")
		for _, v := range my.headers {
			add("<cc_bold>" + v.String(hspace) + "</cc_bold>")
		}
		add("}\n")
	}

	//get-param
	// if my.getparam != "" {
	// 	add("GET-Example : " + my.getparam + "\n")
	// }

	//Params
	if len(my.params) < 1 {
		//add(fmt.Sprintf("Params   : %v\n", "{ }"))
	} else {
		hspace := "       "
		add("<cc_purple>Params</cc_purple>   : {\n")
		for _, v := range my.params {
			add("<cc_bold>" + v.String(hspace) + "</cc_bold>")
		}
		add("}\n")
	}

	fEtcValuePtr := func(kind string, ev *etcValue, spaces ...string) {
		space := ""
		if len(spaces) > 0 {
			space = spaces[0]
		}
		parseStruct := func(rt reflect.Type, isArray bool, sp string) {
			space := "    "
			space += sp

			if !isArray {
				add(kind + " {\n")
			} else {
				add(kind + " [ //array \n")
				add(sp + " {\n")
			}

			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				tagName, ok := field.Tag.Lookup("json")
				if tagName == ",inline" {
					parse_field_inline(add, space, field.Type, ev)
					continue
				}
				tag := ev.getTag(tagName)
				if ok {
					if tagName == "-" || tagName == "" {
						continue
					}
					add(space + strings.TrimSpace(wrap_field_name(tagName)) + " : " + wrap_type_name(field.Type.String()) + tag + "\n")
				} else {
					add(space + strings.TrimSpace(wrap_field_name(field.Name)) + " : " + wrap_type_name(field.Type.String()) + tag + "\n")
				}
			}

			if !isArray {
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
			preKind := kind
			_ = preKind
			elem := rt.Elem()
			kind := elem.Kind()
			if kind == reflect.Struct {
				//reflect.MakeSlice(reflect.SliceOf(elem))
				rv := reflect.New(elem)
				parseStruct(rv.Elem().Type(), true, space)
			} else {
				add(preKind)
				add(space + fmt.Sprintf("%v\n", toETCString(ev.void)))
			}
		} else {
			add(kind)
			add(space + fmt.Sprintf("%v\n", toETCString(ev.void)))
			//add("")
		}
	}

	//Json-Params
	if my.jparam != nil {
		fEtcValuePtr("<cc_purple>Param</cc_purple>", my.jparam)
	}

	add("<cc_blue>--- Response ---</cc_blue>\n")
	for _, key := range my.resultOrder {
		val := my.results[key]

		keyString := getStatusString(key)
		if len(val) < 1 {
			add("" + keyString + " { nil }\n")
		} else if len(val) == 1 {
			add("" + keyString + " { " + val[0].String("", true) + " }\n")
		} else {
			hspace := "       "
			add(keyString + " : {\n")
			for _, v := range val {
				add(v.String(hspace))
			}
			add("}\n")
		}
	} //for
	if my.jResultOK != nil {
		fEtcValuePtr("StatusOK", my.jResultOK)
	}
	if my.jAckOK != nil {
		if _, do := my.jAckOK.void.(struct{}); do {
			add(`{
	"<cc_bold>success</cc_bold>" : <cc_blue>true</cc_blue>,
	"<cc_bold>data</cc_bold>" : <cc_bold>{}</cc_bold>,
 }`)
		} else {
			ackPrefix := `{
	"<cc_bold>success</cc_bold>" : <cc_blue>true</cc_blue>,
	"<cc_bold>data</cc_bold>" : `
			fEtcValuePtr(ackPrefix, my.jAckOK, "    ")
			add("\n }\n")
		}

	}

	if len(my.resultERRR) > 0 {
		add("<cc_blue>--- Fail Response ( 응답이 200 이 아닐경우 )---</cc_blue>\n")
		for _, ed := range my.resultERRR {
			//add(`{ "code":` + v.SInt() + `, "message":string } --- ` + v.ToString() + "\n")
			if ed.tag != "" {
				add(ed.err.Desc() + " (" + ed.tag + ")\n")
			} else {
				add(ed.err.Desc() + "\n")
			}

		} //for

	}

	if len(my.jAckError) > 0 {
		add("<cc_blue>--- 실패 응답 ---</cc_blue>\n")
		add(`{
	"<cc_bold>success</cc_bold>" : <cc_red>false</cc_red>,
	"<cc_bold>data</cc_bold>" : `)
		for i, err := range my.jAckError {
			if i == 0 {
				add(err.String() + "\n")
			} else {
				add("              " + err.String() + "\n")
			}
		} //for
		add(" }\n")
	}

	if len(my.etcs) > 0 {
		add("<cc_blue>--< 구조체 필드 타입 >-------------------------------------------------------------</cc_blue>\n")
	}
	for index, ev := range my.etcs {
		rt := reflect.TypeOf(ev.void)
		if rt.Kind() == reflect.Struct {
			space := "    "
			add("<cc_bold>" + rt.Name() + "</cc_bold>" + " {\n")
			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				tagName, ok := field.Tag.Lookup("json")
				if tagName == ",inline" {
					parse_field_inline(add, space, field.Type, &ev)
					continue
				}

				tag := ev.getTag(tagName)
				if ok {
					if tagName == "-" || tagName == "" {
						continue
					}
					add(space + strings.TrimSpace(wrap_field_name(tagName)) + " : " + wrap_type_name(field.Type.String()) + tag + "\n")
				} else {
					add(space + strings.TrimSpace(wrap_field_name(field.Name)) + " : " + wrap_type_name(field.Type.String()) + tag + "\n")
				}
			}
			add("}\n")
		} else {
			add(fmt.Sprintf("%v\n", toETCString(ev.void)))
		}
		if len(my.etcTag[index]) > 0 {
			for _, comment := range my.etcTag[index] {
				add("<cc_bold>☞</cc_bold>\n " + comment + "\n")
			}
		}
		if index < len(my.etcs)-1 {
			add("---------------------------------------------------------------\n")
		}
	} //for

	return msg

}
