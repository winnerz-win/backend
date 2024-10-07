package itype

import "jtools/dbg"

type ResultError struct {
	Method  string `json:"method"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (my ResultError) String() string { return dbg.ToJsonString(my) }

func ParseResultError(err error) ResultError {
	r, _ := dbg.DecodeStruct[ResultError](err.Error())
	return r
}
