package doc

import (
	"strings"
)

// Objecter :
type Objecter struct {
	*Doc

	urlPath    string
	docTitle   string
	docVersion string
	doclist    DocStringList

	isCustomURL   bool
	isCustomLock  bool
	customAddress string
}

// NewObjecter : http://address:65530/doc/path , cpp[0]="address" , cpp[1]:uplock(true/false)
func NewObjecter(path, titleString, versionString string, cpp ...interface{}) Object {
	obj := &Objecter{
		urlPath:      path,
		docTitle:     titleString,
		docVersion:   versionString,
		isCustomURL:  false,
		isCustomLock: false,
	}
	if len(cpp) >= 2 {
		addr, do := cpp[0].(string)
		if do == true {
			obj.isCustomURL = true
			obj.customAddress = addr
		}
		isLock, do := cpp[1].(bool)
		if do == true {
			obj.isCustomLock = isLock
		}

	}
	return obj
}
func (my *Objecter) ApplyItems(list ItemList) {
	for _, item := range list {
		my.Doc = item.Doc
		my.Apply(item.color)
	}
}

// Message :
func (my *Objecter) Message(v string, color ...Color) {
	headLine := "<--- Message-TAG --->"
	ss := strings.Split(v, "\n")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			headLine = "<--- " + s + " --->"
			break
		}
	} //for

	htmlb := []byte{}
	htmlb = append(htmlb, []byte("&nbsp;&nbsp;")...)
	buf := []byte(v)
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			htmlb = append(htmlb, []byte("<br>&nbsp;")...)
		} else if buf[i] == ' ' {
			htmlb = append(htmlb, []byte("&nbsp;")...)
		} else if buf[i] == '\t' {
			htmlb = append(htmlb, []byte("&nbsp;&nbsp;&nbsp;&nbsp;")...)
		} else {
			htmlb = append(htmlb, buf[i])
		}
	} //for

	defaultColor := Black
	if len(color) > 0 &&
		strings.HasPrefix(string(color[0]), "#") &&
		len(color[0]) == len(defaultColor) {
		defaultColor = color[0]
	}
	dd := DocumentData{
		Href:     headLine,
		Size:     NSize,
		Weight:   Normal.String(),
		Color:    defaultColor.String(),
		Text:     string(htmlb),
		onlyText: true,
	}
	my.doclist = append(my.doclist, dd)
	my.Doc = nil
}

// Comment :
func (my *Objecter) Comment(v string) Object {
	my.Doc = newDoc(v)
	return my
}

// URL :
func (my *Objecter) URL(v string, params ...pairArray) Object {
	my.Doc.URL(v, params...)
	return my
}

// URLS :
func (my *Objecter) URLS(kvs ...string) Object {
	my.Doc.URLS(kvs...)
	return my
}

// WS :
func (my *Objecter) WS(v string) Object {
	my.Doc.WS(v)
	return my
}

// Method :
func (my *Objecter) Method(v string) Object {
	my.Doc.Method(v)
	return my
}

// POST :
func (my *Objecter) POST() Object {
	my.Doc.POST()
	return my
}

// GET :
func (my *Objecter) GET() Object {
	my.Doc.GET()
	return my
}

// Header :
func (my *Objecter) Header(kvlist ...keyValue) Object {
	my.Doc.Header(kvlist...)
	return my
}

// Param :
func (my *Objecter) Param(kvlist ...keyValue) Object {
	my.Doc.Param(kvlist...)
	return my
}

// JParam :
func (my *Objecter) JParam(v interface{}, tag ...string) Object {
	my.Doc.JParam(v, tag...)
	return my
}

// func (my *Objecter) GetParam(pair ...string) Object {
// 	my.Doc.GetParam(pair...)
// 	return my
// }

// JResultOK :
func (my *Objecter) JResultOK(v interface{}, tag ...string) Object {
	my.Doc.JResultOK(v, tag...)
	return my
}

func (my *Objecter) JAckOK(v interface{}, tag ...string) Object {
	my.Doc.JAckOK(v, tag...)
	return my
}
func (my *Objecter) JAckNone() Object {
	my.Doc.JAckNone()
	return my
}
func (my *Objecter) JAckError(err IAckError, tag ...string) Object {
	my.Doc.JAckError(err, tag...)
	return my
}

// Result :
func (my *Objecter) Result(code int, kvlist ...keyValue) Object {
	my.Doc.Result(code, kvlist...)
	return my
}

// ResultOK :
func (my *Objecter) ResultOK(kvlist ...keyValue) Object {
	my.Doc.ResultOK(kvlist...)
	return my
}

// ResultERRR :
func (my *Objecter) ResultERRR(etype IError, tag ...string) Object {
	my.Doc.ResultERRR(etype, tag...)
	return my
}

// ResultBadParameter :
func (my *Objecter) ResultBadParameter(tag ...string) Object {
	my.Doc.ResultBadParameter(tag...)
	return my
}

// ResultBadRequest :
func (my *Objecter) ResultBadRequest(tag ...string) Object {
	my.Doc.ResultBadRequest(tag...)
	return my
}

// ResultUnauthorized :
func (my *Objecter) ResultUnauthorized(tag ...string) Object {
	my.Doc.ResultUnauthorized(tag...)
	return my
}

// ResultServerError :
func (my *Objecter) ResultServerError(tag ...string) Object {
	my.Doc.ResultServerError(tag...)
	return my
}

// ResultNotFound :
func (my *Objecter) ResultNotFound(tag ...string) Object {
	my.Doc.ResultNotFound(tag...)
	return my
}

// ResultConflict :
func (my *Objecter) ResultConflict(tag ...string) Object {
	my.Doc.ResultConflict(tag...)
	return my
}

// Etc :
func (my *Objecter) Etc(v interface{}, tag ...string) Object {
	my.Doc.Etc(v, tag...)
	return my
}

// Tag :
func (my *Objecter) Tag(tag ...string) Object {
	my.Doc.Etc("[etc]", tag...)
	return my
}

// ETC :
func (my *Objecter) ETC(ev etcValue, tag ...string) Object {
	my.Doc.ETC(ev, tag...)
	return my
}

// ETCVAL :
func (my *Objecter) ETCVAL(void interface{}, pair ...string) Object {
	my.Doc.ETCVAL(void, pair...)
	return my
}

// Apply :
func (my *Objecter) Apply(color ...Color) {
	defaultColor := Black

	if len(color) > 0 &&
		strings.HasPrefix(string(color[0]), "#") &&
		len(color[0]) == len(defaultColor) {
		defaultColor = color[0]
	}

	my.Apply2(NSize, Normal, defaultColor)
}

// Apply2 :
func (my *Objecter) Apply2(size string, w Weight, c Color) {
	if my.Doc == nil {
		return
	}
	dd := DocumentData{
		Href:   my.comment,
		Size:   size,
		Weight: w.String(),
		Color:  c.String(),
		Text:   my.HTML(),
	}
	my.doclist = append(my.doclist, dd)
	my.Doc = nil
}

// Set :
func (my *Objecter) Set(color ...Color) {
	my.Apply(color...)
}

// Set2 :
func (my *Objecter) Set2(size string, w Weight, c Color) {
	my.Apply2(size, w, c)
}

// Update :
func (my *Objecter) Update(isLocal ...bool) {
	if my.isCustomURL == false {
		update(my.urlPath, my.docTitle, my.docVersion, my.doclist)
	} else {
		if my.isCustomLock == false {
			updateCustom(my.customAddress, my.urlPath, my.docTitle, my.docVersion, my.doclist)
		}
	}
	my.doclist = my.doclist[:0]
}

func (my *Objecter) Bytes() []byte {
	return HTMLBytes(my.urlPath, my.docTitle, my.docVersion, my.doclist)
}
