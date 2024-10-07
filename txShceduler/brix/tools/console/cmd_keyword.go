package console

const (
	KEYWORD_HELP_LO  = "help"
	KEYWORD_HELP_UP  = "HELP"
	KEYWORD_CLS      = "cls"
	KEYWORD_NANOTIME = "nano.time"
	KEYWORD_MMSTIME  = "mms.time"

	KEYWORD_REMOTE            = "remote"
	KEYWORD_REMOTE_SWITCH_CMD = "remote."

	KEYWORD_CONSOLE_SET   = "console.set"
	KEYWORD_CONSOLE_CLEAR = "console.clear"

	KEYWORD_FILE_START = "console.file.set"
	KEYWORD_FILE_CLEAR = "console.file.clear"
)

var prefixKeywordList = []string{
	KEYWORD_CLS,
	KEYWORD_CONSOLE_SET,
	KEYWORD_CONSOLE_CLEAR,
	KEYWORD_NANOTIME,
	KEYWORD_MMSTIME,
	KEYWORD_REMOTE,
	KEYWORD_REMOTE_SWITCH_CMD,
	KEYWORD_FILE_START,
	KEYWORD_FILE_CLEAR,
}

var removePrefixKeywords = []string{
	KEYWORD_FILE_START,
	KEYWORD_FILE_CLEAR,
}
