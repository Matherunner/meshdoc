package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type strongHandler struct {
}

func NewStrongHandler() HTMLInlineWriterHandler {
	return &strongHandler{}
}

func (h *strongHandler) Name() string {
	return "S"
}

func (h *strongHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("strong", nil, StartTag),
	)
	return
}

func (h *strongHandler) Exit(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("strong", nil, EndTag),
	)
	return
}
