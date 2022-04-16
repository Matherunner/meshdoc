package meshdoc

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/Matherunner/meshdoc/context"
)

type configContextKeyType struct{}

var configContextKey = configContextKeyType{}

var ErrConfigNotSet = errors.New("config not set in context")

func LoadFromFile(path string) (config *MeshdocConfig, err error) {
	config = &MeshdocConfig{}
	_, err = toml.DecodeFile(path, config)
	return
}

func ConfigToContext(ctx context.Context, config *MeshdocConfig) {
	ctx.Set(configContextKey, config)
}

func ConfigFromContext(ctx context.Context) (config *MeshdocConfig, err error) {
	v, ok := ctx.Get(configContextKey)
	if !ok {
		return nil, ErrConfigNotSet
	}
	return v.(*MeshdocConfig), nil
}
