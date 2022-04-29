package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type codeHandler struct {
}

func NewCodeHandler() HTMLInlineWriterHandler {
	return &codeHandler{}
}

func (h *codeHandler) Name() string {
	return "C"
}

func (h *codeHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("code", nil, StartTag),
	)
	return
}

func (h *codeHandler) Exit(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("code", nil, EndTag),
	)
	return
}

type codeBlockHandler struct {
}

func NewCodeBlockHandler() HTMLBlockWriterHandler {
	return &codeBlockHandler{}
}

func (h *codeBlockHandler) Name() string {
	return "C"
}

func (h *codeBlockHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("pre", nil, StartTag),
		NewHTMLItemTag("code", nil, StartTag),
	)
	return
}

func (h *codeBlockHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("code", nil, EndTag),
		NewHTMLItemTag("pre", nil, EndTag),
	)
	return
}
