package chttp

import (
	"bufio"
	"mime/multipart"
	"net/http"
	"strconv"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

//MultipartForm :
type mpForm struct {
	form *multipart.Form
}

//MultipartForm :
func MultipartForm(req *http.Request) mpForm {
	BindingDummy(req)
	return mpForm{
		form: req.MultipartForm,
	}
}

//Value :
func (my mpForm) Value(name string, index int) string {
	if index < 0 {
		dbg.Red("mpForm.Value.index err :", index)
		return ""
	}
	mf := my.form
	if array, do := mf.Value[name]; do == false || len(array) == 0 {
		dbg.Red("mpForm.Value is not :", name)
		return ""
	} else if index >= len(array) {
		dbg.Red("mpForm.Value index over :", index)
		return ""
	} else {
		return array[index]
	}
}

func (my mpForm) Int64(name string, index int) int64 {
	str := my.Value(name, index)
	var vErr error
	v := jmath.NewBigDecimal(str, &vErr)
	if vErr != nil {
		dbg.Red("mpForm.Int64 not format.")
		return 0
	}
	return int64(v.ToBigInteger().Uint64())
}

//Float64 :
func (my mpForm) Float64(name string, index int) float64 {
	str := my.Value(name, index)
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		dbg.Red("mpForm.Int64 not format.")
		return 0
	}
	return f64
}

//Bool :
func (my mpForm) Bool(name string, index int) bool {
	str := my.Value(name, index)
	str = dbg.TrimToLower(str)
	switch str {
	case "true":
		return true
	case "false":
		return false
	}
	dbg.Red("mpForm.Bool not format.")
	return false
}

//File :
func (my mpForm) File(name string, index int, limit ...int64) []byte {
	if index < 0 {
		dbg.Red("mpForm.File.index err :", index)
		return nil
	}
	mf := my.form

	limitSize := int64(0)
	if len(limit) > 0 && limit[0] > 0 {
		limitSize = limit[0]
	}

	var buf []byte
	if fhs, do := mf.File["name"]; do == false || index >= len(fhs) {
		dbg.Red("mpForm.File is not :", name)
		return nil
	} else {
		header := fhs[index]
		if limitSize != 0 {
			if header.Size > limitSize {
				dbg.Red("mpForm.File size over :", header.Size, "/", limitSize)
				return nil
			}
		}
		fp, err := header.Open()
		if err != nil {
			dbg.Red("mpForm.File open err", name, index)
			return nil
		}
		defer fp.Close()
		buf = make([]byte, header.Size)
		reader := bufio.NewReader(fp)
		if _, err := reader.Read(buf); err != nil {
			dbg.Red("mpForm.File read err", name, index)
			return nil
		}
	}

	return buf
}
