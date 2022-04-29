package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type definitionListHandler struct {
}

func NewDefinitionListHandler() HTMLBlockWriterHandler {
	return &definitionListHandler{}
}

func (h *definitionListHandler) Name() string {
	return "DL"
}

func (h *definitionListHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("dl", nil, StartTag),
	)
	return
}

func (h *definitionListHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("dl", nil, EndTag),
	)
	return
}

type definitionTermHandler struct {
}

func NewDefinitionTermHandler() HTMLBlockWriterHandler {
	return &definitionTermHandler{}
}

func (h *definitionTermHandler) Name() string {
	return "DT"
}

func (h *definitionTermHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("dt", nil, StartTag),
	)
	return
}

func (h *definitionTermHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("dt", nil, EndTag),
	)
	return
}

type definitionBodyHandler struct {
}

func NewDefinitionBodyHandler() HTMLBlockWriterHandler {
	return &definitionBodyHandler{}
}

func (h *definitionBodyHandler) Name() string {
	return "DD"
}

func (h *definitionBodyHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("dd", nil, StartTag),
	)
	return
}

func (h *definitionBodyHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("dd", nil, EndTag),
	)
	return
}
