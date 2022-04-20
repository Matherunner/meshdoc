package meshdoc

import (
	"github.com/BurntSushi/toml"
)

type MeshdocConfig struct {
	SourcePath        string
	TemplatePath      string
	OutputPath        string
	KatexRendererPath string
	NodeExecPath      string
}

func LoadFromFile(path string) (config *MeshdocConfig, err error) {
	config = &MeshdocConfig{}
	_, err = toml.DecodeFile(path, config)
	return
}
