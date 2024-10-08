package rpc

import (
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/unix"
)

type BLOCK_HISTORY_FUNC func(caller *ebcm.Sender, reader IReader, idx string, f func(no string)) error

type dBlockRecord struct{}

func (dBlockRecord) BlockHistory(
	caller *ebcm.Sender,
	reader IReader,
	blockIdx string,
	f func(_history_number string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "blockHistory",
			Params: abi.NewParams(
				abi.NewUint256(blockIdx),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(
				rs.Uint256(0),
			)
		},
		_is_debug_call,
	)
}

func (dBlockRecord) GetLastBlockHistory(
	caller *ebcm.Sender,
	reader IReader,
	f func(_history_number string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name:   "getLastBlockHistory",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		_is_debug_call,
	)
}

func (dBlockRecord) GetLastTimeHistory(
	caller *ebcm.Sender,
	reader IReader,
	f func(_time_number string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name:   "getLastTimeHistory",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		_is_debug_call,
	)
}

func (dBlockRecord) GetLastBlockIdx(
	caller *ebcm.Sender,
	reader IReader,
	f func(_blockIdx string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name:   "getLastBlockIdx",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		_is_debug_call,
	)
}

func (dBlockRecord) GetLastBlockRecord(
	caller *ebcm.Sender,
	reader IReader,
	f func(_history_number, _time_number, _blockIdx string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name:   "getLastBlockRecord",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256, //blockHistory_
				abi.Uint256, //timeHistory
				abi.Uint256, //blockIdx
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(
				rs.Uint256(0),
				rs.Uint256(1),
				rs.Uint256(2),
			)
		},
		_is_debug_call,
	)
}

const ERROR_GBH01 = "GBH01" //BlockRecord.getBlockHistory -> require(idx + 1 <= _blockIdx, "GBH01");

type BLOCK_HISTORY_FUNC_3 func(
	caller *ebcm.Sender,
	reader IReader,
	idx string,
	f func(
		historyNumber string, //v
		currentNumber string,
		currentTimestamp unix.Time, //v
	),
) error

func (dBlockRecord) GetBlockHistory(
	caller *ebcm.Sender,
	reader IReader,
	blockIdx string,
	f func(
		historyNumber string, //v
		currentNumber string,
		currentTimestamp unix.Time, //v
	),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "getBlockHistory",
			Params: abi.NewParams(
				abi.NewUint256(blockIdx),
			),
			Returns: abi.NewReturns(
				abi.Uint256, //historyNumber 	[0]-> _blockHistory[idx]
				abi.Uint256, //currentNumber	[1]-> _timeHistory[idx]
				abi.Uint256, //currentTimestamp	[2]-> block.number
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(
				rs.Uint256(0),                  //historyNumber
				rs.Uint256(1),                  //currentNumber
				unix.FromString(rs.Uint256(2)), //currentTimestamp
			)
		},
		//_is_debug_call,
		false,
	)
}
