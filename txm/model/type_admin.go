package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
	"txscheduler/txm/pwd"
)

const (
	V2 = "/v2"

	Platform         = "TXM"
	HeaderAdminToken = "Txm-Adn-Tk"

	AdminDefaultName = "admin"
	AdminDefaultPWD  = "admin1234"
)

//Admin :
type Admin struct {
	Name         string  `bson:"name" json:"name"`
	Pwd          string  `bson:"pwd" json:"pwd"`
	IsRoot       bool    `bson:"is_root" json:"is_root"`
	Salt         string  `bson:"salt" json:"salt"`                   //salt
	RefreshToken string  `bson:"refresh_token" json:"refresh_token"` //갱신 토큰
	CreateAt     mms.MMS `bson:"create_at" json:"create_at"`
	CreateYMD    int     `bson:"create_ymd" json:"create_ymd"`
	Timestamp    mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD          int     `bson:"ymd" json:"ymd"`
}

func (my Admin) String() string { return dbg.ToJSONString(my) }

func (my Admin) Selector() mongo.Bson { return mongo.Bson{"name": my.Name} }

func (my Admin) UpdateDB(db mongo.DATABASE) {
	db.C(inf.COLAdmin).Update(my.Selector(), my)
}

//IndexingDB :
func (my Admin) IndexingDB() {
	inf.DB().Run(inf.DBName, inf.COLAdmin, func(c mongo.Collection) {
		c.EnsureIndex(mongo.SingleIndex("name", "1", true))
		c.EnsureIndex(mongo.SingleIndex("is_root", "1", false))

		if cnt, _ := c.Find(nil).Count(); cnt == 0 {
			nowAt := mms.Now()
			clientPWD := pwd.Hex(AdminDefaultPWD)
			dbpassword, salt := pwd.MakeDBPWD(clientPWD)
			adminIsabel := Admin{
				Name:      AdminDefaultName,
				Pwd:       dbpassword,
				Salt:      salt,
				CreateAt:  nowAt,
				CreateYMD: nowAt.YMD(),
				Timestamp: nowAt,
				YMD:       nowAt.YMD(),
				IsRoot:    true,
			}
			c.Insert(adminIsabel)
		} //

	})
}
