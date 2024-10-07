package mongo

import (
	"encoding/json"
	"fmt"
	"jtools/cc"
	"jtools/dbg"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VOID :
type VOID map[string]interface{}

// VOIDS :
type VOIDS []VOID

func (my VOID) Valid() bool { return len(my) != 0 }
func (my VOID) Count() int  { return len(my) }

func (my VOID) IDString() string {
	id := ""
	switch v := my["_id"].(type) {
	case primitive.ObjectID:
		id = v.Hex()
	case string:
		id = v
	}
	return id
}

// JID : IDString()
func (my VOID) JID() string { return my.IDString() }

func (my VOID) ID() interface{} {
	var id interface{}
	switch v := my["_id"].(type) {
	case primitive.ObjectID:
		id = v
	case string:
		id = ObjectIDFromHex(v)
	}
	return id
}

func (my VOID) IDSelector() Bson {
	return Bson{"_id": my.ID()}
}

func StringToVOID(str string) VOID {
	void := VOID{}
	if err := json.Unmarshal([]byte(str), &void); err != nil {
		fmt.Println("err :", err)
	}
	return void
}
func (my VOID) StringWithoutID() string {
	data := my.Data()
	b, _ := json.MarshalIndent(data, "", "    ")
	return string(b)
}
func (my VOID) String() string {
	b, err := json.MarshalIndent(my, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
func (my VOID) StringTite() string {
	b, err := json.Marshal(my)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (my VOID) Bytes(pretty ...bool) []byte {
	clone := map[string]interface{}{}
	for k, v := range my {
		clone[k] = v
	}
	var buf []byte
	if len(pretty) > 0 && pretty[0] {
		buf, _ = json.MarshalIndent(clone, "", "    ")
	} else {
		buf, _ = json.Marshal(clone)
	}

	return buf
}

func (my VOIDS) String() string { return toString(my) }
func (my *VOIDS) Append(void ...VOID) {
	*my = append(*my, void...)
}
func (my VOIDS) Bytes(pretty ...bool) []byte {
	list := []interface{}{}
	for _, a := range my {
		clone := map[string]interface{}{}
		for k, v := range a {
			clone[k] = v
		}
		list = append(list, clone)
	}
	var buf []byte
	if len(pretty) > 0 && pretty[0] {
		buf, _ = json.MarshalIndent(list, "", "    ")
	} else {
		buf, _ = json.Marshal(list)
	}

	return buf
}

// Injection :
func (my VOID) Injection(p interface{}) {
	var buf []byte
	if v, do := p.([]byte); do {
		buf = v
	} else {
		b, err := json.Marshal(p)
		if err != nil {
			cc.Red("VOID.Injection.Marshal :", err)
			return
		}
		buf = b
	}
	my._injection(buf)
}
func (my VOID) _injection(buf []byte) {
	if err := json.Unmarshal(buf, &my); err != nil {
		cc.Red("VOID._injection.Unmarshal :", err)
	}
}

func (my *VOIDS) Injection(pp interface{}) {
	if b, do := pp.([]byte); do {
		if err := json.Unmarshal(b, &my); err != nil {
			cc.Red("VOIDS.Injection.Unmarshal :", err)
		}
		return
	}

	if b, err := json.Marshal(pp); err == nil {
		if err := json.Unmarshal(b, &my); err != nil {
			cc.Red("VOIDS.Injection.Unmarshal :", err)
		}
	} else {
		cc.Red("VOIDS.Injection.Marshal :", err)
	}
}

func (my VOIDS) ChangeStruct(p interface{}) error {
	return dbg.ParseStruct(my, p)
}

// Parse:
func (my VOID) Parse(p interface{}, isID ...bool) {
	void := my.Data(isID...)
	targetPointer := p
	if void == nil || targetPointer == nil {
		fmt.Println("Parse.param.nil")
	}
	b, err := json.Marshal(void)
	if err != nil {
		fmt.Println("void.Parse : ", err)
	}
	if err := json.Unmarshal(b, targetPointer); err != nil {
		fmt.Println(err)
	}
}

// Data :
func (my VOID) Data(isID ...bool) map[string]interface{} {
	isid := false
	if len(isID) > 0 && isID[0] {
		isid = true
	}

	data := map[string]interface{}{}
	for key, val := range my {
		if key == "_id" {
			if isid {
				switch v := val.(type) {
				case primitive.ObjectID:
					data["_id"] = v.Hex()
				default:
					data["_id"] = v
				}
			}
			continue
		}
		data[key] = val
	}
	return data
}

// TextJSON :
func (my VOID) TextJSON(isID ...bool) string {
	data := my.Data(isID...)
	b, _ := json.MarshalIndent(data, "", "    ")
	return string(b)
}

func ViewType(p interface{}) {
	switch p.(type) {
	case primitive.D:
		cc.Gray("D")
	case primitive.A:
		cc.Gray("A")
	case primitive.E:
		cc.Gray("E")
	case primitive.M:
		cc.Gray("M")
	default:
		cc.Gray("?", reflect.TypeOf(p).Name())
	}
}
