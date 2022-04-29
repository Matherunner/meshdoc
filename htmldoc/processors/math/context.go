package math

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type contextKeyType struct{}

var (
	contextKey = contextKeyType{}
)

func FromContext(ctx *meshdoc.Context) (htmlByNode map[*tree.Node]string) {
	value, ok := ctx.Get(contextKey)
	if !ok {
		return nil
	}
	return value.(map[*tree.Node]string)
}

func setToContext(ctx *meshdoc.Context, htmlByNode map[*tree.Node]string) {
	ctx.Set(contextKey, htmlByNode)
}
