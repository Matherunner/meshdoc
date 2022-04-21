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
