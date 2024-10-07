package ebcm

import (
	"jtools/cc"
	"jtools/dbg"
	"jtools/jmath"
	"strings"
)

type EIP712DomainData struct {
	Name              string
	Version           string
	ChainId           int64
	VerifyingContract string
}

func EIP712Domain(name string, version string, chainId any, verifyingContract string) EIP712DomainData {
	return EIP712DomainData{
		Name:              name,
		Version:           version,
		ChainId:           jmath.Int64(chainId),
		VerifyingContract: verifyingContract,
	}
}

///////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////

type TYPE string

func (my TYPE) String() string { return string(my) }

const (
	ADDRESS = TYPE("address")
	UINT256 = TYPE("uint256")
	BOOL    = TYPE("bool")
	BYTES32 = TYPE("bytes32")
)

type IStructType interface {
	Name() string
	Type() string
}

type _structType struct {
	field_name string
	field_type TYPE
}

func (my _structType) Name() string { return my.field_name }
func (my _structType) Type() string { return my.field_type.String() }

type IStructDesign interface {
	Spec(field_name string, field_type TYPE, field_value ...any)
	Set(field_name string, field_value any)

	Types() []IStructType
	Values() map[string]any
}
type _structMap struct {
	type_list []_structType //
	_types    map[string]TYPE

	valueMap map[string]any // [field_name]value
}

func (my _structMap) debug_v(primaryType string) any {
	dt_name := "primaryType"
	if key := strings.ReplaceAll(primaryType, " ", ""); key != "" {
		dt_name = key
	}
	dt := []struct {
		Name string
		Type string
	}{}
	for _, v := range my.type_list {
		dt = append(dt,
			struct {
				Name string
				Type string
			}{
				Name: v.Name(),
				Type: v.Type(),
			},
		)
	}
	return map[string]any{
		dt_name:   dt,
		"message": my.valueMap,
	}
}
func (my _structMap) String() string {
	return dbg.ToJsonString(my.debug_v("primaryType"))
}

func (my _structMap) clone_type() _structMap {
	clone := _structMap{
		_types:   map[string]TYPE{},
		valueMap: map[string]any{},
	}
	clone.type_list = append(clone.type_list, my.type_list...)
	for k, v := range my._types {
		clone._types[k] = v
	} //for
	return clone
}

func (my _structMap) Types() []IStructType {
	sl := []IStructType{}
	for _, v := range my.type_list {
		sl = append(sl, v)
	}
	return sl
}
func (my _structMap) Values() map[string]any { return my.valueMap }

func (my _structMap) refresh_value(field_type TYPE, field_value any) any {
	switch field_type {
	case UINT256:
		val := jmath.HEX(field_value)
		if val == "0x" {
			val = "0x0"
		}
		//cc.Purple(val)
		field_value = val

	case ADDRESS:
		address := dbg.Void(field_value)
		IsAddressP(&address)
		field_value = address
	} //switch
	return field_value
}

func (my *_structMap) Spec(field_name string, field_type TYPE, field_value ...any) {
	if _, do := my._types[field_name]; do {
		cc.Red("ebcm.StructMap[", field_name, "] is already exist.")
		return
	}

	var val interface{} = nil
	if len(field_value) > 0 {
		if field_value[0] != nil {
			val = my.refresh_value(field_type, field_value[0])
		}
	}

	my._types[field_name] = field_type

	my.valueMap[field_name] = val
	my.type_list = append(my.type_list,
		_structType{
			field_name: field_name,
			field_type: field_type,
		},
	)
}
func (my *_structMap) Set(field_name string, field_value any) {
	field_type, do := my._types[field_name]
	if !do {
		cc.Red("ebcm.StructMap[", field_name, "] is not exist.")
		return
	}

	field_value = my.refresh_value(field_type, field_value)
	my.valueMap[field_name] = field_value
}

///////////////////////////////////////////////////////////////////////////////////

type IEIP712Struct interface {
	PrimaryType() string
	Design() IStructDesign
	Set(field_name string, field_value any)

	CloneType() IEIP712Struct
}
type eip712StructData struct {
	primaryType string
	design      *_structMap
}

func (my eip712StructData) String() string {
	v := map[string]any{
		"primaryType": my.primaryType,
		"design":      my.design.debug_v(my.primaryType),
	}
	return dbg.ToJsonString(v)
}

func (my eip712StructData) PrimaryType() string    { return my.primaryType }
func (my *eip712StructData) Design() IStructDesign { return my.design }

//	func (my *eip712StructData) Spec(field_name string, field_type TYPE, field_value ...any)  {
//		my.message.Spec(field_name, field_type, field_value)
//	}
func (my *eip712StructData) Set(field_name string, field_value any) {
	my.design.Set(field_name, field_value)
}

func (my eip712StructData) CloneType() IEIP712Struct {
	_design := my.design.clone_type()
	clone := &eip712StructData{
		primaryType: my.primaryType,
		design:      &_design,
	}
	return clone
}

func EIP712Struct(primaryType string, f func(sm IStructDesign)) IEIP712Struct {
	data := &eip712StructData{
		primaryType: primaryType,
		design: &_structMap{
			valueMap: map[string]any{},
			_types:   map[string]TYPE{},
		},
	}
	f(data.design)
	return data
}

type EIP712Validator func(
	domain EIP712DomainData,
	msg_data IEIP712Struct,
	signer_address string,
	signature string) (bool, error)

func MakeEIP712Validator(f EIP712Validator) EIP712Validator {
	return f
}
