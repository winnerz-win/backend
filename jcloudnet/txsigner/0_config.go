package txsigner

import (
	"jtools/dbg"
	"jtools/jlog"
	"jtools/jyml"
)

type Config struct {
	TITLE  string             `yaml:"title"`
	Logger jlog.ConfigLogYAML `yaml:"logger"`
}

func (my Config) String() string { return dbg.ToYamlString(my) }

func LoadConfig(filename string) *Config {
	config := &Config{}

	if err := jyml.LoadFile(filename, config); err != nil {
		_Exit(err)
	}

	jlog.Init(config.Logger)
	jlog.Info("LoadConfig(", filename, ")")

	return config
}

func _Exit(err ...interface{}) {
	if len(err) > 0 {
		jlog.Panic(err...)
	}

	dbg.Exit()
}
