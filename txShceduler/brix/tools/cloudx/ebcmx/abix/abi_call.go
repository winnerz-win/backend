package abix

type ISender interface {
	CallContract(from, to string, data []byte) ([]byte, error)
}

type Method struct {
	Name string

	//Params interface{}
	//Returns interface{}
	Params  TypeList
	Returns TypeList
}

// DelegateCall : func(finder ISender, contract string, method Method, caller string, f func(rs RESULT), isLogs ...bool) error
type DelegateCall func(finder ISender, contract string, method Method, caller string, f func(rs RESULT), isLogs ...bool) error

type Caller struct {
	// NewAddress func(data string) interface{}
	// NewBytes32 func(data string) interface{}
	// NewUint    func(data interface{}) interface{}
	// NewUint256 func(data interface{}) interface{}
	// NewUint128 func(data interface{}) interface{}
	// NewUint112 func(data interface{}) interface{}
	// NewUint64  func(data uint64) interface{}
	// NewUint32  func(data uint32) interface{}
	// NewUint16  func(data uint16) interface{}
	// NewUint8   func(data uint8) interface{}
	// NewBool    func(data bool) interface{}
	// NewString  func(data string) interface{}

	// NewAddressArray func(data ...string) interface{}
	// NewBytes32Array func(data ...string) interface{}
	// NewUint256Array func(data ...interface{}) interface{}
	// NewUint128Array func(data ...interface{}) interface{}

	// NewUint112Array func(data ...interface{}) interface{}
	// NewUint64Array  func(data ...uint64) interface{}
	// NewUint32Array  func(data ...uint32) interface{}
	// NewUint16Array  func(data ...uint16) interface{}
	// NewUint8Array   func(data ...uint8) interface{}
	// NewBoolArray    func(data ...bool) interface{}
	// NewStringArray  func(data ...string) interface{}

	// Tuple      func(params ...interface{}) interface{}
	// TupleArray func(params ...interface{}) interface{}

	// NewParams  func(ps ...interface{}) interface{}
	// NewReturns func(ps ...interface{}) interface{}

	// None             interface{}
	// Address          interface{}
	// Uint256          interface{}
	// Uint             interface{}
	// Uint128          interface{}
	// Uint112          interface{}
	// Uint64           interface{}
	// Uint32           interface{}
	// Uint16           interface{}
	// Uint8            interface{}
	// Bool             interface{}
	// String           interface{}
	// Bytes            interface{}
	// Bytes32          interface{}
	// AddressArray     interface{}
	// Uint256Array     interface{}
	// UintArray        interface{}
	// Uint128Array     interface{}
	// Uint112Array     interface{}
	// Uint64Array      interface{}
	// Uint32Array      interface{}
	// Uint16Array      interface{}
	// Uint8Array       interface{}
	// BoolArray        interface{}
	// StringArray      interface{}
	// BytesArray       interface{}
	// Bytes32Array     interface{}
	// ITuple           interface{}
	// ITupleFixedArray interface{}
	// ITupleFlexArray  interface{}

	/*
		Call : func(finder ISender, contract string, method Method, caller string, f func(rs RESULT), isLogs ...bool) error
		Call DelegateCall
	*/
	Call func(finder ISender, contract string, method Method, caller string, f func(rs RESULT), isLogs ...bool) error

	InputDataPure DelegateInputDataPure
}

type DelegateInputDataPure func(data string, prarmResturns TypeList) RESULT
