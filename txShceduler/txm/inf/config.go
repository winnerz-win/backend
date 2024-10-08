package inf

import (
	"jtools/cloud/ebcm"
	"strings"
	"sync"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jcfg"
	"txscheduler/brix/tools/jpath"
	"txscheduler/txm/pwd"
)

type IConfig struct {
	Mainnet    bool              `yaml:"-" json:"mainnet"`
	Version    string            `yaml:"version" json:"version"`
	Seed       string            `yaml:"seed" json:"seed"`
	DB         string            `yaml:"db" json:"db"`
	IPCheck    bool              `yaml:"ip_check" json:"ip_check"`
	ClientHost map[bool][]string `yaml:"client_host" json:"client_host"`
	AdminSalt  string            `yaml:"-" json:"-"`

	Confirms              int         `yaml:"confirms" json:"confirms"`
	IsLockTransferByOwner bool        `yaml:"is_lock_transfer_by_owner"  json:"is_lock_transfer_by_owner"`
	Owners                KeyPairList `yaml:"owners" json:"owners"`
	Masters               KeyPairList `yaml:"master" json:"master"`
	Chargers              KeyPairList `yaml:"charger" json:"charger"`

	Tokens TokenInfoList `yaml:"tokens" json:"tokens"`

	InfuraKeys []string `yaml:"infura" json:""`
	ESKeys     []string `yaml:"es" json:""`
}

func NewHost(ip, port string) []string { return []string{ip, port} }

func (my IConfig) String() string {
	viewJSON := struct {
		Mainnet    bool          `yaml:"-" json:"mainnet"`
		Version    string        `yaml:"version" json:"version"`
		Seed       string        `yaml:"seed" json:"seed"`
		DB         string        `yaml:"db" json:"db"`
		IPCheck    bool          `yaml:"ip_check" json:"ip_check"`
		ClientHost string        `yaml:"client_host" json:"client_host"`
		Confirms   int           `yaml:"confirms" json:"confirms"`
		Master     string        `json:"master_address"`
		Charger    string        `json:"charger_address"`
		Tokens     TokenInfoList `yaml:"tokens" json:"tokens"`
	}{
		Mainnet:    my.Mainnet,
		Version:    my.Version,
		Seed:       my.Seed,
		DB:         my.DB,
		IPCheck:    my.IPCheck,
		ClientHost: ClientAddress(),
		Confirms:   my.Confirms,
		Master:     Master().Address,
		Charger:    Charger().Address,
		Tokens:     my.Tokens,
	}
	return dbg.ToJSONString(viewJSON)
}

func (my IConfig) View() string {
	my.Masters[0].PrivateKey = "****"
	my.Chargers[0].PrivateKey = "****"
	my.InfuraKeys = []string{}
	my.ESKeys = []string{}
	return my.String()
}

type xConfig struct {
	*IConfig
	mu sync.RWMutex
}

var (
	config = &xConfig{}
	seed   = ""
)

func Config() IConfig {
	defer config.mu.RUnlock()
	config.mu.RLock()
	return *config.IConfig
}

func InitConfig(name string, mainnet bool) {
	item := &IConfig{}
	jcfg.LoadYAML(jpath.NowPath()+"/"+name, item)

	item.Mainnet = mainnet
	SetConfig(item)

}

func SetConfig(c *IConfig) {
	defer config.mu.Unlock()
	config.mu.Lock()

	config.IConfig = c

	if config.Mainnet {
		seed = "mainnet_" + config.DB + "_" + config.Seed
	} else {
		seed = "testnet_" + config.DB + "_" + config.Seed
	}
	config.Seed = seed

	owner := KeyPairList{}
	for i := 0; i < len(config.Owners); i++ {
		if config.Masters[i].Mainnet != config.Mainnet {

		} else {
			config.Owners[i].Refactory()
			owner = append(owner, config.Owners[i])
		}
	} //for
	config.Owners = owner

	master := KeyPairList{}
	for i := 0; i < len(config.Masters); i++ {
		if config.Masters[i].Mainnet != config.Mainnet {

		} else {
			config.Masters[i].Refactory()
			master = append(master, config.Masters[i])
		}
	} //for
	config.Masters = master

	charger := KeyPairList{}
	for i := 0; i < len(config.Chargers); i++ {
		if config.Chargers[i].Mainnet != config.Mainnet {

		} else {
			config.Chargers[i].Refactory()
			charger = append(charger, config.Chargers[i])
		}
	} //for
	config.Chargers = charger

	tokens := TokenInfoList{}
	for i := 0; i < len(config.Tokens); i++ {
		token := config.Tokens[i]
		if token.Mainnet != config.Mainnet {
		} else {

			var finder *ebcm.Sender = nil
			if strings.HasPrefix(config.Tokens[i].Contract, "0x") {
				finder = GetSender()
			}
			config.Tokens[i].Refactory(finder)
			tokens = append(tokens, config.Tokens[i])
		}
	} //for
	config.Tokens = tokens

	pwd.InitPWD(config.AdminSalt)
}
