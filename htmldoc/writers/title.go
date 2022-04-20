package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type TitleHandler struct {
}

func (h *TitleHandler) Name() string {
	return "TITLE"
}

// TODO: wait, then how do you xref to things that don't have numbers? how to display the xref?

func (h *TitleHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	counters := counter.FromContext(ctx)
	num := counters.ElementNumber(node)

	items = append(
		items,
		NewHTMLItemTag("h1", []html.Attr{{
			Name:  "class",
			Value: "page-header",
		}}, StartTag),
		NewHTMLItemTag("span", nil, StartTag),
		NewHTMLItemText(num),
		NewHTMLItemTag("span", nil, EndTag),
		NewHTMLItemText(string(CharEmspace)),
	)

	return
}

func (h *TitleHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(
		items,
		NewHTMLItemTag("h1", nil, EndTag),
	)
	return
}
