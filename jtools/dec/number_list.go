package dec

import (
	"jtools/dbg"
)

type NumberList []Number

// String : dbg.ToJsonString(my)
func (my NumberList) String() string { return dbg.ToJsonString(my) }

///////////////////////////////////////////

func (my NumberList) Int64() any {
	sl := make(Int64List, len(my))
	for i, v := range my {
		sl[i] = v.Int64()
	}
	return sl
}

func (my NumberList) Uint256() any {
	sl := make(Uint256List, len(my))
	for i, v := range my {
		sl[i] = v.Uint256()
	}
	return sl
}

func (my *NumberList) SetInt64() NumberList {
	for i := range *my {
		(*my)[i].SetInt64()
	}
	return *my
}

func (my *NumberList) SetUint256() NumberList {
	for i := range *my {
		(*my)[i].SetUint256()
	}
	return *my
}

func (my *NumberList) AppendInt64(v any) NumberList {
	(*my) = append((*my), INT64(v).Number())
	return *my
}

func (my *NumberList) AppendUint256(v any) NumberList {
	(*my) = append((*my), _UINT256(v).Number())
	return *my
}

func (my NumberList) Do(n any) bool {
	for _, v := range my {
		if v.Cmp(n) == 0 {
			return true
		}
	}
	return false
}

func (my *NumberList) Remove(n any) bool {
	for i, v := range *my {
		if v.Cmp(n) == 0 {
			*my = append((*my)[:i], (*my)[i+1:]...)
			return true
		}
	}
	return false
}
