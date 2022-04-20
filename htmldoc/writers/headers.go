package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshforce/tree"
)

type genericHeaderHandler struct {
	blockName string
	tagName   string
}

func (h *genericHeaderHandler) Name() string {
	return h.blockName
}

func (h *genericHeaderHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	counters := counter.FromContext(ctx)
	num := counters.ElementNumber(node)

	items = append(items,
		NewHTMLItemTag(h.tagName, nil, StartTag),
		NewHTMLItemTag("span", nil, StartTag),
		NewHTMLItemText(num),
		NewHTMLItemTag("span", nil, EndTag),
		NewHTMLItemText(string(CharEmspace)),
	)

	return
}

func (h *genericHeaderHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag(h.tagName, nil, EndTag),
	)
	return
}

func NewH1Handler() HTMLBlockWriterHandler {
	return &genericHeaderHandler{blockName: "H1", tagName: "h2"}
}

func NewH2Handler() HTMLBlockWriterHandler {
	return &genericHeaderHandler{blockName: "H2", tagName: "h3"}
}

func NewH3Handler() HTMLBlockWriterHandler {
	return &genericHeaderHandler{blockName: "H3", tagName: "h4"}
}

func NewH4Handler() HTMLBlockWriterHandler {
	return &genericHeaderHandler{blockName: "H4", tagName: "h5"}
}

func NewH5Handler() HTMLBlockWriterHandler {
	return &genericHeaderHandler{blockName: "H5", tagName: "h6"}
}
