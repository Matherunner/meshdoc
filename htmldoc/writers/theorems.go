package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type theoremHandler struct {
}

func NewTheoremHandler() HTMLBlockWriterHandler {
	return &theoremHandler{}
}

func (h *theoremHandler) Name() string {
	return "THEOREM"
}

func (h *theoremHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	title, ok := block.Options().Get("NAME")
	if !ok {
		title = "Theorem"
	}

	items = append(items,
		NewHTMLItemTag("div", NewAttributes(html.Attr{
			Name:  "class",
			Value: "theorem",
		}), StartTag),
		NewHTMLItemTag("span", NewAttributes(html.Attr{
			Name:  "class",
			Value: "theorem",
		}), StartTag),
		NewHTMLItemText(title+"."),
		NewHTMLItemTag("span", nil, EndTag),
	)

	return
}

func (h *theoremHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("div", nil, EndTag),
	)
	return
}
