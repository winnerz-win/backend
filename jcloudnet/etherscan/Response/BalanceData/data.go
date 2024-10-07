package BalanceData

import (
	"jtools/jmath"
)

// ValueResult :
type ValueResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// IsStatusOk :
func (my ValueResult) IsStatusOk() bool {
	if my.Status == "1" {
		return true
	}
	return false
}

// ResultString : balanceString
func (my ValueResult) ResultString() string {
	if my.IsStatusOk() {
		return my.Result
	}
	return "0"
}

// ResultValue : *jmath.BigDecimal
func (my ValueResult) ResultValue() string {
	return jmath.VALUE(my.ResultString())
}

// New :
func New() *ValueResult {
	return &ValueResult{
		Status: "0",
	}
}
