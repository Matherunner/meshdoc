package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type unorderedListHandler struct {
}

func NewUnorderedListHandler() HTMLBlockWriterHandler {
	return &unorderedListHandler{}
}

func (h *unorderedListHandler) Name() string {
	return "UL"
}

func (h *unorderedListHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("ul", nil, StartTag),
	)
	return
}

func (h *unorderedListHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("ul", nil, EndTag),
	)
	return
}
