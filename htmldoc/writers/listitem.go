package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type listItemHandler struct {
}

func NewListItemHandler() HTMLBlockWriterHandler {
	return &listItemHandler{}
}

func (h *listItemHandler) Name() string {
	return "LI"
}

func (h *listItemHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("li", nil, StartTag),
	)
	return
}

func (h *listItemHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("li", nil, EndTag),
	)
	return
}
