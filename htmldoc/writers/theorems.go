package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
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
	counters := counter.FromContext(ctx)
	num := counters.ElementNumber(node)

	title, ok := block.Options().Get("NAME")
	if !ok {
		title = "Theorem"
	}

	if num != "" {
		title += " " + num
	}

	title += "."

	items = append(items,
		NewHTMLItemTag("div", NewAttributes(html.Attr{
			Name:  "class",
			Value: "theorem",
		}), StartTag),
		NewHTMLItemTag("span", NewAttributes(html.Attr{
			Name:  "class",
			Value: "theorem",
		}), StartTag),
		NewHTMLItemText(title),
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

type proofHandler struct {
}

func NewProofHandler() HTMLBlockWriterHandler {
	return &proofHandler{}
}

func (h *proofHandler) Name() string {
	return "PROOF"
}

func (h *proofHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("div", NewAttributes(html.Attr{
			Name:  "class",
			Value: "proof",
		}), StartTag),
		NewHTMLItemTag("span", NewAttributes(html.Attr{
			Name:  "class",
			Value: "proof",
		}), StartTag),
		NewHTMLItemText("Proof."),
		NewHTMLItemTag("span", nil, EndTag),
	)
	return
}

func (h *proofHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("div", nil, EndTag),
	)
	return
}
