package model

import (
	"time"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmap"
)

var (
	acccount_lock_map = jmap.New[int64, any]()
)

func lock_member_uid(uid int64, f func()) {
	defer acccount_lock_map.Delete(uid)
	for {
		if _, loaded := acccount_lock_map.LoadOrStore(uid, struct{}{}); !loaded {
			f()
			break
		}
		dbg.Gray("LockMemberUID[", uid, "] wait...")
		time.Sleep(time.Millisecond)
	} //for
}

func LockMemberUID(db mongo.DATABASE, uid int64, f func(member Member)) bool {
	ok := false
	lock_member_uid(uid, func() {
		member := LoadMember(db, uid)
		if member.Valid() {
			f(member)
			ok = true
		}
	})
	return ok
}

func LockMember(db mongo.DATABASE, address string, f func(member Member)) bool {
	ok := false
	if v := LoadMemberAddress(db, address); v.Valid() {
		lock_member_uid(v.UID, func() {
			member := LoadMember(db, v.UID)
			f(member)
			ok = true
		})
	}
	return ok
}

func SyncMemberCoinUID(db mongo.DATABASE, uid int64, iSender any) {
	if iSender == nil {
		return
	}
	lock_member_uid(uid, func() {
		member := LoadMember(db, uid)
		if member.Valid() {
			switch sender := iSender.(type) {
			case *ecsx.Sender:
				member.UpdateCoinDB_Legacy(db, sender)

			case *ebcm.Sender:
				member.UpdateCoinDB(db, sender)
			} //switch
		}
	})
}

func SyncMemberCoin(db mongo.DATABASE, address string, iSender any) {
	if iSender == nil {
		return
	}

	if v := LoadMemberAddress(db, address); v.Valid() {
		lock_member_uid(v.UID, func() {
			member := LoadMember(db, v.UID)
			switch sender := iSender.(type) {
			case *ecsx.Sender:
				member.UpdateCoinDB_Legacy(db, sender)

			case *ebcm.Sender:
				member.UpdateCoinDB(db, sender)
			} //switch
		})
	}
}

// var (
// 	muCloud sync.Mutex
// )

// func CloudLock(f func()) {
// 	defer muCloud.Unlock()
// 	muCloud.Lock()
// 	f()
// }
