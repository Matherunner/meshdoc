package meshdoc

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/Matherunner/meshforce"
)

type meshdocConfig struct {
	Path string
}

type MeshdocOptions struct {
	ConfigPath string
}

type Meshdoc struct {
	options *MeshdocOptions

	config meshdocConfig
}

func NewMeshdoc(options *MeshdocOptions) *Meshdoc {
	return &Meshdoc{
		options: options,
	}
}

func (m *Meshdoc) Run() (err error) {
	_, err = toml.DecodeFile(m.options.ConfigPath, &m.config)
	if err != nil {
		return
	}

	fileInfo, err := ioutil.ReadDir(m.config.Path)
	if err != nil {
		return
	}

	for _, file := range fileInfo {
		filePath := path.Join(m.config.Path, file.Name())
		log.Printf("path = %+v\n", filePath)

		parser := meshforce.NewParser()
		log.Printf("parser = %+v\n", parser)
	}

	return
}
