package ecsx

import (
	"strings"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

// CustomInput :
type CustomInput struct {
	Input string `bson:"input" json:"input"`

	MethodID string   `bson:"method_id" json:"method_id"`
	Data     []string `bson:"data" json:"data"`
	Count    int      `bson:"count" json:"count"`
}

func (my CustomInput) String() string { return dbg.ToJSONString(my) }

func newCustomInput(input string) CustomInput {
	c := CustomInput{
		Input: input,
		Data:  []string{},
	}
	if len(input) < 8 {
		return c
	}

	c.MethodID = "0x" + input[:8]
	input = input[8:]
	for {
		if len(input) >= 64 {
			item := "0x" + input[:64]
			c.Data = append(c.Data, item)
			input = input[64:]
		} else {
			break
		}
	} //for
	c.Count = len(c.Data)
	return c
}

// IndexAddress :
func (my CustomInput) IndexAddress(index int) string {
	if index >= 0 && index < len(my.Data) {
		data := my.Data[index]
		if strings.HasPrefix(data, "0x000000000000000000000000") {
			return strings.ReplaceAll(data, "0x000000000000000000000000", "0x")
		}
	}
	return ""
}

// IndexValue :
func (my CustomInput) IndexValue(index int) string {
	if index >= 0 && index < len(my.Data) {
		data := my.Data[index]
		return jmath.VALUE(data)
	}
	return ""
}
