package inf

import "txscheduler/brix/tools/jargs"

var (
	args jargs.ArgAction
)

func SetArgs(a jargs.ArgAction) { args = a }

func Args() jargs.ArgAction { return args }
