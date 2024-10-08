package model

import (
	"jtools/jmath"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

type AA_GAS_PRICE struct {
	Key         string `bson:"key" json:"key"`
	GasMulValue string `bson:"gas_mul_value" json:"gas_mul_value"`

	LT_TransferPendingMin string `bson:"lt_transfer_pending_min" json:"lt_transfer_pending_min"`
	LT_LockPendingMin     string `bson:"lt_lock_pending_min" json:"lt_lock_pending_min"`

	NN_LockUnlockPendingMin string `bson:"nn_lock_unlock_pending_min" json:"nn_lock_unlock_pending_min"`
}

func (my *AA_GAS_PRICE) _selector() mongo.Bson {
	my.Key = inf.AA_GAS_PRICE
	return mongo.Bson{"key": my.Key}
}
func (my AA_GAS_PRICE) _GetDB(db mongo.DATABASE) AA_GAS_PRICE {
	db.C(inf.AA_GAS_PRICE).Find(my._selector()).One(&my)
	return my
}

// GetGasMulValue : gas_mul_value < 1 ? 1 : gas_mul_value
func (my AA_GAS_PRICE) GetGasMulValue() string {
	if jmath.CMP(my.GasMulValue, 1) < 0 {
		return "1"
	}
	return my.GasMulValue
}
func (my AA_GAS_PRICE) GetLtTransferPendingMin() time.Duration {
	return time.Duration(jmath.Int(my.LT_TransferPendingMin)) * time.Minute
}
func (my AA_GAS_PRICE) GetLtLockPendingMin() time.Duration {
	return time.Duration(jmath.Int(my.LT_LockPendingMin)) * time.Minute
}
func (my AA_GAS_PRICE) GetLockUnLockPendingMin() time.Duration {
	return time.Duration(jmath.Int(my.NN_LockUnlockPendingMin)) * time.Minute
}

func GetGAS(db mongo.DATABASE) AA_GAS_PRICE {
	return AA_GAS_PRICE{}._GetDB(db)
}

func (my AA_GAS_PRICE) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.AA_GAS_PRICE)
		c.EnsureIndex(mongo.SingleIndex("key", 1, true))

		if cnt, _ := c.Find(nil).Count(); cnt == 0 {
			my._selector()
			my.GasMulValue = "1.5"
			my.LT_TransferPendingMin = "4"
			my.LT_LockPendingMin = "4"
			my.NN_LockUnlockPendingMin = "4"
			c.Insert(my)
		}
	})
}
