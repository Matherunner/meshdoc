package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type hyperlinkHandler struct {
}

func NewHyperlinkHandler() HTMLInlineWriterHandler {
	return &hyperlinkHandler{}
}

func (h *hyperlinkHandler) Name() string {
	return "A"
}

func (h *hyperlinkHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	href, ok := inline.Options().Get("href")
	if !ok {
		href = ""
	}

	items = append(items,
		NewHTMLItemTag("a", NewAttributes(html.Attr{
			Name:  "href",
			Value: href,
		}), StartTag),
	)
	return
}

func (h *hyperlinkHandler) Exit(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("a", nil, EndTag),
	)
	return
}
