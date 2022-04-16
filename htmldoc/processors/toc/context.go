package toc

import "github.com/Matherunner/meshdoc"

type contextKeyType struct{}

var (
	contextKey = contextKeyType{}
)

type contextValue struct {
	toc []Entry
}

func FromContext(ctx *meshdoc.Context) (toc []Entry) {
	value, ok := ctx.Get(contextKey)
	if !ok {
		return nil
	}
	return value.(*contextValue).toc
}

func setToContext(ctx *meshdoc.Context, value *contextValue) {
	ctx.Set(contextKey, value)
}
