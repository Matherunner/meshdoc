package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type emphasisHandler struct {
}

func NewEmphasisHandler() HTMLInlineWriterHandler {
	return &emphasisHandler{}
}

func (h *emphasisHandler) Name() string {
	return "E"
}

func (h *emphasisHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("em", nil, StartTag),
	)
	return
}

func (h *emphasisHandler) Exit(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("em", nil, EndTag),
	)
	return
}
