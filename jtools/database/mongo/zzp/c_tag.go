package zzp

/*
	() brackets / parentheses
	{} braces
	[] square brackets

	<> angle brackets
	<<>> double angle brackets
*/

type Tag int

const (
	//TagSmall : ()
	TagSmall = Tag(1)
	//TagMiddle: {}
	TagMiddle = Tag(2)
	//TagBig : []
	TagBig = Tag(3)
	//TagAngle : <>
	TagAngle = Tag(4)
)

var (
	tagMap = map[Tag][2]string{
		TagSmall:  {"(", ")"},
		TagMiddle: {"{", "}"},
		TagBig:    {"[", "]"},
		TagAngle:  {"<", ">"},
	}
)

type ArgType string

const (
	argEmpty     = "empty"
	ArgJson      = "json"
	ArgJsonArray = "jsonArray"
	ArgText      = "text"
	ArgTextComma = "textComma" // [,,,,]
	ArgTextDot   = "textDot"   // [....]
)

type TagPair [2]string

func (my TagPair) String() string {
	return my.Front() + "" + my.Back()
}
func (my TagPair) Front() string { return my[0] }
func (my TagPair) Back() string  { return my[1] }
func (my TagPair) IsSame(v string) bool {
	return my.Front() == v || my.Back() == v
}

func tagPair(elem string) string {
	switch elem {
	case ")":
		return "("
	case "(":
		return ")"

	case "}":
		return "{"
	case "{":
		return "}"

	case "]":
		return "["
	case "[":
		return "]"

	case ">":
		return "<"
	case "<":
		return ">"
	}
	return ""
}
