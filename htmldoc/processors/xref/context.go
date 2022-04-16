package xref

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

type ContextValue struct {
	targetByID       map[string]XRefTarget
	targetByXRefNode map[*tree.Node]XRefTarget
	idByTargetNode   map[*tree.Node]XRefTarget
}

func NewContextValue() *ContextValue {
	return &ContextValue{
		targetByID:       map[string]XRefTarget{},
		targetByXRefNode: map[*tree.Node]XRefTarget{},
		idByTargetNode:   map[*tree.Node]XRefTarget{},
	}
}

func (v *ContextValue) TargetByID(id string) (target XRefTarget, ok bool) {
	target, ok = v.targetByID[id]
	return
}

func (v *ContextValue) TargetByXRefNode(node *tree.Node) (target XRefTarget, ok bool) {
	target, ok = v.targetByXRefNode[node]
	return
}

func (v *ContextValue) IDByTargetNode(node *tree.Node) (target XRefTarget, ok bool) {
	target, ok = v.idByTargetNode[node]
	return
}

func setToContext(ctx *meshdoc.Context, value *ContextValue) {
	ctx.Set(contextKey, value)
}

func FromContext(ctx *meshdoc.Context) *ContextValue {
	v, ok := ctx.Get(contextKey)
	if !ok {
		return nil
	}
	return v.(*ContextValue)
}
