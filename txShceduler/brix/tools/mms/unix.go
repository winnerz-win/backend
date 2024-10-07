package mms

import "txscheduler/brix/tools/unix"

func FromUNIX(v unix.Time) MMS {
	return MMS(v.Int64() * 1000)
}

func (my MMS) ToUNIX() unix.Time {
	return unix.Time(my.Int64() / 1000)
}
