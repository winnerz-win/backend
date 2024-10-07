package doc

const (
	NSize = "18"
)

//Weight :
type Weight string

func (my Weight) String() string {
	return string(my)
}

//Color :
type Color string

func (my Color) String() string {
	return string(my)
}

//Normal :
const Normal Weight = "normal"

//Bold :
const Bold Weight = "bold"

//Black :
const Black Color = "#000000"

//Gray :
const Gray Color = "#7F7F7F"

//Red :
const Red Color = "#FF0000"

//Blue :
const Blue Color = "#0000FF"

//Green :
const Green Color = "#00FF00"

//Pink :
const Pink Color = "#FF00DD"

//Purple :
const Purple Color = "#5F00FF"
