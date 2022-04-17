package xref

import (
	"errors"
	"fmt"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

var (
	ErrTargetNotExist = errors.New("xref ID does not point to a valid target")
)

type XRefTarget struct {
	ID   string
	Path meshdoc.GenericPath
	Node *tree.Node
}

type XRef struct {
}

func NewXRef() meshdoc.Postprocessor {
	return &XRef{}
}

func (x *XRef) Process(ctx *meshdoc.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	ctxValue := NewContextValue()

	// FIXME: the path written here needs to be converted to the output file path! Maybe we can write a postprocessor that injects the mapping into ctx?

	// Scan for all nodes with ID
	for p, t := range r.Files() {
		it := tree.NewIterator(t.Root())
		for it.Next(tree.InstructionEnterChild) {
			node := it.Value()
			if it.Exit() {
				continue
			}

			id := ""
			ok := false
			switch node := node.Value.(type) {
			case *tree.BlockNode:
				id, ok = node.Options().Get("ID")
			case *tree.InlineNode:
				if node.Name() != "XREF" {
					id, ok = node.Options().Get("ID")
				}
			}
			if ok {
				target := XRefTarget{
					ID:   id,
					Path: p,
					Node: node,
				}
				ctxValue.targetByID[id] = target
				ctxValue.idByTargetNode[node] = target
			}
		}
	}

	// Scan for XREF nodes
	for _, t := range r.Files() {
		it := tree.NewIterator(t.Root())
		for it.Next(tree.InstructionEnterChild) {
			node := it.Value()
			if it.Exit() {
				continue
			}

			if inline, ok := node.Value.(*tree.InlineNode); ok {
				if inline.Name() != "XREF" {
					continue
				}

				id, ok := inline.Options().Get("ID")
				if !ok {
					continue
				}

				target, existing := ctxValue.targetByID[id]
				if !existing {
					return nil, fmt.Errorf("%w: %s", ErrTargetNotExist, id)
				}

				ctxValue.targetByXRefNode[node] = target
			}
		}
	}

	setToContext(ctx, ctxValue)

	return r, nil
}
