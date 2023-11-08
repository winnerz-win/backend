package TxExecuteStatus

const (
	ResultOK   = "0"
	ResultFail = "1"
)

//New :
func New() *StatusResult {
	ins := &StatusResult{
		Result: TxStatus{
			IsError: "1",
		},
	}
	return ins
}

//StatusResult :
type StatusResult struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  TxStatus `json:"result"`
}

//IsStatusOk : true is "1" //(StatusResult).Result.IsConfirm()
func (my StatusResult) IsStatusOk() bool {
	if my.Status == "1" {
		return true
	}
	return false
}

//ToString :
func (my StatusResult) ToString() string {
	return `{
	status  : ` + my.Status + `
	message : ` + my.Message + `
	result  : ` + my.Result.ToString() + `
}`
}

//TxStatus :
type TxStatus struct {
	IsError        string `json:"isError"`
	ErrDescription string `json:"errDescription"`
}

//IsConfirm : ResultOK = "0", ResultFail = "1"
func (my TxStatus) IsConfirm() string {
	return my.IsError
}

//ConfirmString :ResultOK , ResultFail
func (my TxStatus) ConfirmString() string {
	switch my.IsError {
	case ResultOK:
		return "ResultOK"
	case ResultFail:
		return "ResultFail"
	default:
		return "??"
	}
}

//ToString :
func (my TxStatus) ToString() string {
	return `{
		isError : ` + my.IsError + `
		errDescription : ` + my.ErrDescription + `
	}`
}
