package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type tocHandler struct {
}

func NewTOCHandler() HTMLBlockWriterHandler {
	return &tocHandler{}
}

func (h *tocHandler) Name() string {
	return "TOC"
}

func (h *tocHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	items = append(items,
		NewHTMLItemTag("ul", nil, StartTag),
	)

	// TODO: put numberings!

	// TODO: actually, this whole thing has already been implemented in the writer!

	entries := toc.FromContext(ctx)
	for _, entry := range entries {
		items = append(items,
			NewHTMLItemTag("li", nil, StartTag),
			NewHTMLItemTag("a", NewAttributes(html.Attr{
				Name:  "href",
				Value: entry.Path.WebPath(), // FIXME: need the mapper?
			}), StartTag),
			NewHTMLItemText(entry.Path.Path()), // FIXME: put the actual title!
			NewHTMLItemTag("a", nil, EndTag),
			NewHTMLItemTag("li", nil, EndTag),
		)
	}

	instruction = tree.InstructionIgnoreChild

	return
}

func (h *tocHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("ul", nil, EndTag),
	)
	return
}
