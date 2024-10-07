package mongo

import (
	"txscheduler/brix/tools/database/mongo/tools/dbg"
)

//Config :
type Config struct {
	Address    string   `yaml:"address" json:"address"`
	IsAuth     bool     `yaml:"auth" json:"auth"`
	ID         string   `yaml:"id" json:"id"`
	PW         string   `yaml:"pw" json:"pw"`
	AuthSource string   `yaml:"auth_source" json:"auth_source"`
	List       []string `yaml:"list" json:"list"`
}

//ToString :
func (my Config) ToString() string {
	return dbg.ToJSONString(my)
}

//NewConfig :
func NewConfig(c *Config) *CDB {
	return New(c.Address, c.IsAuth, c.AuthSource, c.ID, c.PW)
}

type IConfig interface {
	GetList() []string
	GetIsAuth() bool
	GetAuthSource() string
	GetID() string
	GetPWD() string
}

func NewIConfig(c IConfig) *CDB {
	return NewList(
		c.GetList(),
		c.GetIsAuth(),
		c.GetAuthSource(),
		c.GetID(),
		c.GetPWD(),
	)
}
