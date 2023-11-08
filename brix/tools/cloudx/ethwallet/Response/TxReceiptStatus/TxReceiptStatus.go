package TxReceiptStatus

const (
	ResultOK      = "1"
	ResultPending = ""
	ResultFail    = "0"
)

//New :
func New() *StatusResult {
	ins := &StatusResult{
		Result: TxResult{},
	}
	return ins
}

//StatusResult :
type StatusResult struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  TxResult `json:"result"`
}

//IsStatusOk :
func (my StatusResult) IsStatusOk() bool {
	if my.Status == ResultOK {
		return true
	}
	return false
}

//ToString :
func (my StatusResult) ToString() string {
	return `{
	status : ` + my.Status + `
	message : ` + my.Message + `
	result : ` + my.Result.ToString() + `
}`
}

//TxResult :
type TxResult struct {
	Status string `json:"status"`
}

//IsConfirm : Result- "1"(success) , "0"(fail) , ""(pending)
func (my TxResult) IsConfirm() string {
	return my.Status
}

//ConfirmString :
func (my TxResult) ConfirmString() string {
	switch my.Status {
	case ResultOK:
		return "ResultOK"
	case ResultFail:
		return "ResultFail"
	}
	return "ResultPending"
}

//ToString :
func (my TxResult) ToString() string {
	return `{ status : ` + my.Status + `}`
}
