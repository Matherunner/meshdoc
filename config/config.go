package config

import (
	"github.com/BurntSushi/toml"
	"github.com/Matherunner/meshdoc"
)

func LoadFromFile(path string) (config *meshdoc.MeshdocConfig, err error) {
	config = &meshdoc.MeshdocConfig{}
	_, err = toml.DecodeFile(path, config)
	return
}
