package dec

import "jtools/jmath"

////////////////////////////////////////////////////////

type Uint256 Number

func (my Uint256) _string() string { return string(my) }
func (my Uint256) _v() string      { return jmath.VALUE(my) }

func (my Uint256) Type() Type { return my.Number().Type() }

func (my Uint256) Number() Number { return Number(my) }

func (my Uint256) Value() string {
	return my.Number().Uint256Value()
}
func (my Uint256) Int64() Int64 {
	return my.Number().Int64()
}
func (my Uint256) Uint256() Uint256 { return my }

func (my Uint256) Cmp(v any) int {
	return my.Number().Cmp(v)
}

////////////////////////////////////////////////////////

type Uint256List []Uint256

func (my Uint256List) Int64List() Int64List {
	list := make(Int64List, len(my))
	for i, v := range my {
		list[i] = v.Int64()
	}
	return list
}

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////

type Int64 Number

func (my Int64) _string() string { return string(my) }
func (my Int64) _v() string      { return jmath.VALUE(my) }

func (my Int64) Type() Type { return my.Number().Type() }

func (my Int64) Number() Number { return Number(my) }

func (my Int64) Value() string {
	return my.Number().Int64Value()
}

func (my Int64) Uint256() Uint256 {
	return my.Number().Uint256()
}
func (my Int64) Int64() Int64 { return my }

func (my Int64) Cmp(v any) int {
	return my.Number().Cmp(v)
}

// //////////////////////////////////////////////////////

type Int64List []Int64

func (my Int64List) Uint256List() Uint256List {
	list := make(Uint256List, len(my))
	for i, v := range my {
		list[i] = v.Uint256()
	}
	return list
}
