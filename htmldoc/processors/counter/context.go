package counter

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type contextKeyType struct{}

type ContextValue struct {
	fileByKey    map[meshdoc.GenericPath]string
	elementByKey map[*tree.Node]string
}

func (v *ContextValue) FileNumber(key meshdoc.GenericPath) string {
	if v == nil {
		return ""
	}
	return v.fileByKey[key]
}

func (v *ContextValue) ElementNumber(key *tree.Node) string {
	if v == nil {
		return ""
	}
	return v.elementByKey[key]
}

var ContextKey = contextKeyType{}

func setToContext(ctx *meshdoc.Context, value *ContextValue) {
	ctx.Set(ContextKey, value)
}

func FromContext(ctx *meshdoc.Context) *ContextValue {
	v, ok := ctx.Get(ContextKey)
	if !ok {
		return nil
	}
	return v.(*ContextValue)
}
