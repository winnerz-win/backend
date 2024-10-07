package etherscan

//CallTimeLimitError :
type CallTimeLimitError struct {
	OrginalType string `json:"-"`
	Status      string `json:"status"` //  1, 0
	Message     string `json:"message"`
	Result      string `json:"result"`
}

//NewCallTimeLimitError :
func NewCallTimeLimitError(orgType string) *CallTimeLimitError {
	return &CallTimeLimitError{
		OrginalType: orgType,
	}
}

//GetOrignalClass : cFailCallbackData
func (my CallTimeLimitError) GetOrignalClass() interface{} {
	switch my.OrginalType {
	case "ContractTransactional":
		return my.GetContractTransactional()
	} //switch
	return nil
}

//GetContractTransactional :
func (my CallTimeLimitError) GetContractTransactional() *ContractTransactional {
	return &ContractTransactional{
		Status:            my.Status,
		Message:           my.Message,
		_isLimitTimeError: true,
		_limitTImeMessage: my.Result,
		_isRefactData:     true,
		_lastBlockNumber:  0,
	}
}
