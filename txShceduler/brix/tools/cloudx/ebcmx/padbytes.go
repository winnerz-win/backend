package ebcmx

type Appender interface {
	/////////////////////////////////////////////////

	IndexSize() int
	SetParamHeader(count int) int
	//SetIndex(header int, index int)

	SetAddress(header int, address string) int
	SetAddressArray(header, i int, address ...string) int

	SetAmount(header int, value string) int
	SetAmountArray(header, i int, value ...string) int

	SetBool(header int, b bool) int
	SetBoolArray(header int, i int, bs ...bool) int

	SetText(header, i int, text string) int
	SetTextArray(header, i int, text ...string) int

	SetBytes(header, i int, Bytes []byte) int
	SetBytesHex(header, i int, hexBytes string) int

	SetBytesArray(header, i int, Bytes ...[]byte) int
	SetBytesHexArray(header, i int, hexBytes ...string) int

	SetBytes32(header int, Bytes []byte) int
	SetBytes32Hex(header int, hexBytes string) int

	SetBytes32Array(header, i int, Bytes ...[]byte) int
	SetBytes32HexArray(header, i int, hexBytes ...string) int

	HexDebug() string
}

type PADBYTES interface {
	Bytes() []byte
	Hex() string
	HexDebug() string
}

type XSenderData struct {
	MakePadBytes func(method string, callback func(Appender)) interface{}
}

type GasSpeed string

func (my GasSpeed) Value() GasSpeed { return my }
func (my GasSpeed) String() string  { return string(my) }

const (
	GasFastest = GasSpeed("fastest")
	GasFast    = GasSpeed("fast")
	GasAverage = GasSpeed("average")
	GasSafeLow = GasSpeed("safeLow")
	GasBegger  = GasSpeed("begger")
)
