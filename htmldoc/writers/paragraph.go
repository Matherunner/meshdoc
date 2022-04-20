package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type paragraphHandler struct {
}

func NewParagraphHandler() HTMLBlockWriterHandler {
	return &paragraphHandler{}
}

func (h *paragraphHandler) Name() string {
	return "P"
}

func (h *paragraphHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("p", nil, StartTag),
	)
	return
}

func (h *paragraphHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("p", nil, EndTag),
	)
	return
}
