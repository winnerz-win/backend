package mongo

type AdminCommandResult struct {
	data MAP
}

func (my AdminCommandResult) Data() MAP {
	return my.data
}

func (my AdminCommandResult) Valid() bool {
	return my.data != nil
}

func (my AdminCommandResult) String() string {
	if !my.Valid() {
		return "AdminCommandResult is null"
	}
	return toString(my.data)
}

func (my *CDB) AdminCommand(cmd interface{}) AdminCommandResult {

	sr := my.RunAdmin(cmd)
	if sr == nil {
		return AdminCommandResult{}
	}

	ar := AdminCommandResult{
		data: MAP{},
	}
	sr.Decode(&ar.data)
	return ar
}
