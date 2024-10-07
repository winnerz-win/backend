package doc

import (
	"jtools/jnet/chttp"
	"net/http"
	"strings"
)

// URL : http
func (my *Doc) URL(v string, params ...pairArray) *Doc {
	my.url = strings.TrimSpace(v)
	if len(params) > 0 {
		my.url += getParams(params[0]...)
	}
	return my
}

type pairArray []string

func Pair(pair ...string) pairArray {
	pa := pairArray{}
	for _, v := range pair {
		pa = append(pa, v)
	}
	return pa
}

func getParams(s ...string) string {
	if len(s)%2 != 0 {
		return ""
	}
	if len(s) == 0 {
		return ""
	}

	msg := "?"
	for i := 0; i < len(s); i += 2 {
		key := s[i]
		val := s[i+1]
		msg = msg + key + "=" + val

		if i+2 < len(s) {
			msg += "&"
		}
	}

	return msg
}

//GetParam :
// func (my *Doc) GetParam(pair ...string) *Doc {
// 	if len(pair)%2 != 0 {
// 		return my
// 	}

// 	if my.getparam == "" {
// 		if my.url != "" {
// 			my.getparam = my.url + "?"
// 		} else {
// 			my.getparam = "address/api?"
// 		}
// 	}

// 	for i := 0; i < len(pair); i += 2 {
// 		key := pair[i]
// 		val := pair[i+1]
// 		my.getparam = my.getparam + key + "=" + val

// 		if i+2 < len(pair) {
// 			my.getparam += "&"
// 		}
// 	}
// 	return my
// }

// URLS : [path , :args,user, :args,member, ...]
func (my *Doc) URLS(kvs ...string) *Doc {
	path := kvs[0]
	args := kvs[1:]

	my.urls.do = true
	my.urls.path = path
	my.urls.kvs = map[string][]string{}
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		val := args[i+1]
		if list, do := my.urls.kvs[key]; do {
			list = append(list, val)
			my.urls.kvs[key] = list
		} else {
			my.urls.kvs[key] = []string{val}
		}
	} //for
	return my
}

// WS : websocket
func (my *Doc) WS(v string) *Doc {
	my.url = webSocketTag + strings.TrimSpace(v)
	return my
}

// Method :
func (my *Doc) Method(v string) *Doc {
	my.method = strings.TrimSpace(v)
	return my
}

// POST :
func (my *Doc) POST() *Doc {
	my.Method("POST")
	return my
}

// GET :
func (my *Doc) GET() *Doc {
	my.Method("GET")
	return my
}

// Header :
func (my *Doc) Header(kvlist ...keyValue) *Doc {
	for _, kv := range kvlist {
		my.headers = append(my.headers, kv)
	}
	my.headers.KeyAlign()
	return my
}

// Param :
func (my *Doc) Param(kvlist ...keyValue) *Doc {
	for _, kv := range kvlist {
		my.params = append(my.params, kv)
	}
	my.params.KeyAlign()
	return my
}

// JParam :
func (my *Doc) JParam(v interface{}, tag ...string) *Doc {
	ev := etcValue{
		void: v,
		tags: map[string]string{},
	}

	for i := 0; i < len(tag); i += 2 {
		key := tag[i]
		val := tag[i+1]
		ev.tags[key] = val
	} //for

	my.jparam = &ev

	return my
}

// JResultOK :
func (my *Doc) JResultOK(v interface{}, tag ...string) *Doc {
	ev := etcValue{
		void: v,
		tags: map[string]string{},
	}

	for i := 0; i < len(tag); i += 2 {
		key := tag[i]
		val := tag[i+1]
		ev.tags[key] = val
	} //for

	my.jResultOK = &ev

	return my
}

func (my *Doc) JAckOK(v interface{}, tag ...string) *Doc {
	ev := etcValue{
		void: v,
		tags: map[string]string{},
	}

	for i := 0; i < len(tag); i += 2 {
		key := tag[i]
		val := tag[i+1]
		ev.tags[key] = val
	} //for

	my.jAckOK = &ev
	return my
}
func (my *Doc) JAckNone() *Doc {
	my.JAckOK(struct{}{})
	return my
}

type IAckError interface {
	Code() int
	Desc() string
}

type AckErrorPair struct {
	err IAckError
	tag string
}

func (my AckErrorPair) String() string {
	msg := `{ <cc_bold>"error_code"</cc_bold>:<cc_green>` +
		chttp.Cat(my.err.Code()) +
		`</cc_green>, <cc_bold>"error_message"</cc_bold>:<cc_green>"` +
		my.err.Desc() + `"</cc_green> }`

	if my.tag != "" {
		msg += " // " + my.tag
	}
	return msg
}

func (my *Doc) JAckError(err IAckError, tag ...string) *Doc {
	pair := AckErrorPair{}
	pair.err = err
	if len(tag) > 0 {
		pair.tag = tag[0]
	}
	my.jAckError = append(my.jAckError, pair)
	return my
}

// Result :
func (my *Doc) Result(code int, kvlist ...keyValue) *Doc {
	vals := keyValueList{}
	for _, kv := range kvlist {
		vals = append(vals, kv)
	}
	vals.KeyAlign()
	my.results[code] = vals

	isOrder := true
	for _, v := range my.resultOrder {
		if v == code {
			isOrder = false
		}
	} //for
	if isOrder {
		my.resultOrder = append(my.resultOrder, code)
	}
	return my
}

// ResultOK :
func (my *Doc) ResultOK(kvlist ...keyValue) *Doc {
	my.Result(http.StatusOK, kvlist...)
	return my
}

// ResultERRR :
func (my *Doc) ResultERRR(etype IError, tag ...string) *Doc {
	ed := errrData{
		err: etype,
		tag: "",
	}
	if len(tag) > 0 {
		ed.tag = tag[0]
	}
	my.resultERRR = append(my.resultERRR, ed)
	return my
}

// ResultBadParameter :
func (my *Doc) ResultBadParameter(tag ...string) *Doc {
	my.Result(http.StatusBadRequest, KV("message", `"Bad-parameter"`, tag...))
	return my
}

// ResultBadRequest :
func (my *Doc) ResultBadRequest(tag ...string) *Doc {
	my.Result(http.StatusBadRequest, KVMessageString(tag...))
	return my
}

// ResultUnauthorized :
func (my *Doc) ResultUnauthorized(tag ...string) *Doc {
	my.Result(http.StatusUnauthorized, KVMessageString(tag...))
	return my
}

// ResultServerError :
func (my *Doc) ResultServerError(tag ...string) *Doc {
	my.Result(http.StatusInternalServerError, KVMessageString(tag...))
	return my
}

// ResultNotFound :
func (my *Doc) ResultNotFound(tag ...string) *Doc {
	my.Result(http.StatusNotFound, KVMessageString(tag...))
	return my
}

// ResultConflict :
func (my *Doc) ResultConflict(tag ...string) *Doc {
	my.Result(http.StatusConflict, KVMessageString(tag...))
	return my
}

// Etc :
func (my *Doc) Etc(v interface{}, tag ...string) *Doc {
	ev := etcValue{
		void: v,
		tags: map[string]string{},
	}
	my.etcs = append(my.etcs, ev)

	taglist := []string{}
	for _, c := range tag {
		taglist = append(taglist, strings.TrimSpace(c))
	}
	index := len(my.etcs) - 1
	my.etcTag[index] = taglist
	return my
}

// ETC :
func (my *Doc) ETC(ev etcValue, tag ...string) *Doc {
	my.etcs = append(my.etcs, ev)

	taglist := []string{}
	for _, c := range tag {
		taglist = append(taglist, strings.TrimSpace(c))
	}
	index := len(my.etcs) - 1
	my.etcTag[index] = taglist
	return my
}

// ETCVAL :
func (my *Doc) ETCVAL(void interface{}, pair ...string) *Doc {
	ev := EV(void, pair...)
	my.etcs = append(my.etcs, ev)

	taglist := []string{}
	index := len(my.etcs) - 1
	my.etcTag[index] = taglist
	return my
}

// Set : Global.Set()
func (my *Doc) Set(color ...Color) {
	// defaultColor := Black

	// if len(color) > 0 &&
	// 	strings.HasPrefix(string(color[0]), "#") &&
	// 	len(color[0]) == len(defaultColor) {
	// 	defaultColor = color[0]
	// }

	// my.Set2(NSize, Normal, defaultColor)
}

// Set2 :Global.Set2()
func (my *Doc) Set2(size string, w Weight, c Color) {
	// dd := DocumentData{
	// 	Size:   size,
	// 	Weight: w.String(),
	// 	Color:  c.String(),
	// 	Text:   my.HTML(),
	// }
	// doclist = append(doclist, dd)
}
