package model

import (
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo/tools/dbg"
	"txscheduler/brix/tools/database/mongo/tools/jmath"
)

func Erc20Balance(
	caller ebcm.CALLER,
	erc20 string,
	user string,
) (string, error) {
	balance := ZERO
	err := caller.Call(
		erc20,
		abi.Method{
			Name: "balanceOf",
			Params: abi.NewParams(
				abi.NewAddress(user),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		erc20,
		func(rs abi.RESULT) {
			balance = rs.Uint256(0)
		},
	)
	return balance, err
}

///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

type RpcTotalLocked struct {
	LockedAmount string `json:"locked_amount"`
	LockedCount  string `json:"locked_count"`
}

func (my RpcTotalLocked) Valid() bool    { return my.LockedAmount != "" }
func (my RpcTotalLocked) String() string { return dbg.ToJsonString(my) }

func (my RpcTotalLocked) GetLockedCount() int { return jmath.Int(my.LockedCount) }

///////////////////////////////////////////////////////////////////////

type RpcLockInfo struct {
	ReleaseTime unix.Time `json:"release_time"`
	BalanceLock string    `json:"balance_lock"`
}

func (my RpcLockInfo) Valid() bool    { return my.BalanceLock != "" }
func (my RpcLockInfo) String() string { return dbg.ToJsonString(my) }

func (my RpcLockInfo) CMP(release_time unix.Time, amount string) bool {
	return jmath.CMP(my.BalanceLock, amount) == 0 && my.ReleaseTime == release_time
}

func (RpcLockInfo) _type_list() abi.TypeList {
	return abi.TypeList{
		abi.Uint256,
		abi.Uint256,
	}
}
func (my RpcLockInfo) _abi_parse(rs abi.IRESULT) RpcLockInfo {
	my.ReleaseTime = unix.FromString(rs.Uint256(0))
	my.BalanceLock = rs.Uint256(1)
	return my
}

/*
RpcLockInfoList :

	type RpcLockInfo struct {
		ReleaseTime unix.Time `json:"release_time"`
		BalanceLock string    `json:"balance_lock"`
	}
*/
type RpcLockInfoList []RpcLockInfo

type RpcLockUserInfo struct {
	User             string          `json:"user"`
	Balance          string          `json:"balance"`
	BalanceLockTotal string          `json:"balance_lock_total"`
	LockInfo         RpcLockInfoList `json:"lock_info"`
}

func (my RpcLockUserInfo) String() string { return dbg.ToJSONString(my) }
func (RpcLockUserInfo) _type_lisT() abi.TypeList {
	return abi.TypeList{
		abi.Address,
		abi.Uint256,
		abi.Uint256,
		abi.TupleArray(RpcLockInfo{}._type_list()...),
	}
}
func (my RpcLockUserInfo) _abi_parse(rs abi.IRESULT) RpcLockUserInfo {
	my.User = rs.Address(0)
	my.Balance = rs.Uint256(1)
	my.BalanceLockTotal = rs.Uint256(2)
	for _, tuple := range rs.Tuple(3) {
		lock_info := RpcLockInfo{}._abi_parse(tuple)
		my.LockInfo = append(my.LockInfo, lock_info)
	} //for
	return my
}
func (my *RpcLockUserInfo) sort_release_time() {
	// sort.Slice(my.LockInfo, func(i, j int) bool {
	// 	return my.LockInfo[i].ReleaseTime < my.LockInfo[j].ReleaseTime
	// })
}

///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////

type LockTokenUtil struct{}

func (my LockTokenUtil) GetUserLockInfoAll(
	caller ebcm.CALLER,
	util_contract string,
	erc20 string,
	user string,
) (RpcLockUserInfo, error) {

	return my._raw_get_lock_user_info(
		caller,
		erc20,
		user,
	)

	// util_contract = strings.TrimSpace(util_contract)
	// if util_contract == "" || !ebcm.IsAddress(util_contract) {
	// 	return my._raw_get_lock_user_info(
	// 		caller,
	// 		erc20,
	// 		user,
	// 	)
	// }

	// //custom : util_contract
	// lock_user_info := RpcLockUserInfo{}
	// err := caller.Call(
	// 	util_contract,
	// 	abi.Method{
	// 		Name: "getLockInfo",
	// 		Params: abi.NewParams(
	// 			abi.NewAddress(erc20),
	// 			abi.NewAddress(user),
	// 		),
	// 		Returns: abi.NewReturns(
	// 			abi.Tuple(RpcLockUserInfo{}._type_lisT()...),
	// 		),
	// 	},
	// 	util_contract,
	// 	func(rs abi.RESULT) {
	// 		lock_user_info = RpcLockUserInfo{}._abi_parse(
	// 			rs.TupleOne(0),
	// 		)

	// 	},
	// 	false,
	// )

	// lock_user_info.sort_release_time()
	// return lock_user_info, err
}

///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////

func (LockTokenUtil) AvailableBalanceOf(
	caller ebcm.CALLER,
	erc20 string,
	holder string,
) (string, error) {
	balance := ebcm.ZERO

	err := caller.Call(
		erc20,
		abi.Method{
			Name: "availableBalanceOf",
			Params: abi.NewParams(
				abi.NewAddress(holder),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		erc20,
		func(rs abi.RESULT) {
			balance = rs.Uint256(0)
		},
	)

	return balance, err
}

func (my LockTokenUtil) AvailablePriceOf(
	caller ebcm.CALLER,
	erc20 string,
	holder string,
	decimals string,
) (string, error) {
	balance, err := my.AvailableBalanceOf(caller, erc20, holder)
	if err != nil {
		return "", err
	}
	if jmath.CMP(balance, 0) <= 0 {
		return balance, nil
	}
	return ebcm.WeiToToken(balance, decimals), nil
}

func (LockTokenUtil) TotalLocked(
	caller ebcm.CALLER,
	erc20 string,
	user string,
) (RpcTotalLocked, error) {
	total_locked := RpcTotalLocked{}

	err := caller.Call(
		erc20,
		abi.Method{
			Name: "totalLocked",
			Params: abi.NewParams(
				abi.NewAddress(user),
			),
			Returns: abi.NewReturns(
				abi.Uint256, //locked_amount
				abi.Uint256, //length
			),
		},
		erc20,
		func(rs abi.RESULT) {
			total_locked.LockedAmount = rs.Uint256(0)
			total_locked.LockedCount = rs.Uint256(1)
		},
	)

	return total_locked, err
}

func (LockTokenUtil) LockInfo(
	caller ebcm.CALLER,
	erc20 string,
	user string,
	position_index int,
) (RpcLockInfo, error) {
	lock_info := RpcLockInfo{}
	err := caller.Call(
		erc20,
		abi.Method{
			Name: "lockInfo",
			Params: abi.NewParams(
				abi.NewAddress(user),
				abi.NewUint256(position_index),
			),
			Returns: abi.NewReturns(
				abi.Uint256, //releaseTime
				abi.Uint256, //amount
			),
		},
		erc20,
		func(rs abi.RESULT) {
			lock_info = RpcLockInfo{
				ReleaseTime: unix.FromString(rs.Uint256(0)),
				BalanceLock: rs.Uint256(1),
			}
		},
	)

	return lock_info, err
}

func (my LockTokenUtil) _raw_get_lock_user_info(
	caller ebcm.CALLER,
	erc20 string,
	user string,
) (RpcLockUserInfo, error) {
	ebcm.IsAddressP(&user)

	lock_user_info := RpcLockUserInfo{
		User:             user,
		Balance:          ZERO,
		BalanceLockTotal: ZERO,
		LockInfo:         RpcLockInfoList{},
	}

	balance, err := Erc20Balance(
		caller,
		erc20,
		user,
	)
	if err != nil {
		return lock_user_info, err
	}
	lock_user_info.Balance = balance

	if jmath.CMP(lock_user_info.Balance, 0) <= 0 {
		return lock_user_info, nil
	}

	total_locked, err := my.TotalLocked(
		caller,
		erc20,
		user,
	)
	if err != nil {
		return lock_user_info, err
	}

	lock_user_info.BalanceLockTotal = total_locked.LockedAmount

	locked_length := total_locked.GetLockedCount()

	for i := 0; i < locked_length; i++ {
		lock_info, err := my.LockInfo(
			caller,
			erc20,
			user,
			i,
		)
		if err != nil {
			return lock_user_info, err
		}
		lock_user_info.LockInfo = append(lock_user_info.LockInfo, lock_info)

	} //for

	lock_user_info.sort_release_time()
	return lock_user_info, nil
}

///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////

type LockedStateInfo struct {
	PositionIndex int       `json:"position_index"`
	ReleaseTime   unix.Time `json:"release_time"`
	LockedPrice   string    `json:"locked_price"`
	IsTimeOver    bool      `json:"is_time_over"`
}

func (LockedStateInfo) TagString() []string {
	return []string{
		"position_index", "데이터 위치 정보( owner권한으로 unlock 트랜잭션 실행시 필요 )",
		"release_time", "언락시간 (10자리 UTC)",
		"lock_price", "락업된 수량 (10.1234)",
		"is_time_over", "락 해제 시간과 서버 응답시간을 비교해서 락업 기간이 지났을 경우 true (락업 해재 할 필요 없음.)",
	}
}

type LockAccountInfo struct {
	Account            string            `json:"account"`
	TotalPrice         string            `json:"total_price"`
	LockedTotalPrice   string            `json:"locked_total_price"`
	LockedCalcPrice    string            `json:"locked_calc_price"`
	LockedStateInfos   []LockedStateInfo `json:"locked_state_infos"`
	ResponseServerTime unix.Time         `json:"response_server_time"`
}

func (LockAccountInfo) TagString() []string {
	return []string{
		"account", "지갑 주소",
		"total_price", "총 보유 잔액",
		"locked_total_price", "락업된 총 잔액",
		"locked_calc_price", "응답시 서버 시간을 기준으로 계산된 락업된 총 잔액",
		"locked_state_infos", "락업 상태정보 리스트",
		"response_server_time", "응답시 서버 시간",
	}
}

func (my LockAccountInfo) String() string { return dbg.ToJSONString(my) }

func (my RpcLockUserInfo) LockUserInfo(decmals any) LockAccountInfo {
	get_price := func(wei string) string {
		return ebcm.WeiToToken(wei, decmals)
	}

	rs_time := unix.Now()
	info := LockAccountInfo{
		Account:            my.User,
		TotalPrice:         get_price(my.Balance),
		LockedTotalPrice:   get_price(my.BalanceLockTotal),
		LockedStateInfos:   make([]LockedStateInfo, len(my.LockInfo)),
		ResponseServerTime: rs_time,
	}

	locked_calc_price := ZERO
	for i, v := range my.LockInfo {
		info.LockedStateInfos[i].PositionIndex = i
		info.LockedStateInfos[i].ReleaseTime = v.ReleaseTime
		info.LockedStateInfos[i].LockedPrice = get_price(v.BalanceLock)
		info.LockedStateInfos[i].IsTimeOver = rs_time >= v.ReleaseTime
		if !info.LockedStateInfos[i].IsTimeOver {
			locked_calc_price = jmath.ADD(locked_calc_price, info.LockedStateInfos[i].LockedPrice)
		}
	} //for

	info.LockedCalcPrice = locked_calc_price
	return info
}
