package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type orderedListHandler struct {
}

func NewOrderedListHandler() HTMLBlockWriterHandler {
	return &orderedListHandler{}
}

func (h *orderedListHandler) Name() string {
	return "OL"
}

func (h *orderedListHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("ol", nil, StartTag),
	)
	return
}

func (h *orderedListHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("ol", nil, EndTag),
	)
	return
}
