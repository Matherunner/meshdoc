package writers

import (
	"errors"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/xref"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type OutputFileMapper func(input meshdoc.GenericPath) string

var (
	ErrTargetNotFound = errors.New("target not found for xref")
)

type xrefHandler struct {
	mapOutputFile OutputFileMapper
}

func NewXRefHandler(mapOutputFile OutputFileMapper) HTMLInlineWriterHandler {
	return &xrefHandler{mapOutputFile: mapOutputFile}
}

func (h *xrefHandler) Name() string {
	return "XREF"
}

func (h *xrefHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	refStore := xref.FromContext(ctx)
	target, ok := refStore.TargetByXRefNode(node)
	if !ok {
		err = ErrTargetNotFound
		return
	}

	outputFile := h.mapOutputFile(target.Path)
	items = append(items,
		NewHTMLItemTag("a", NewAttributes(html.Attr{
			Name:  "href",
			Value: outputFile + "#" + target.ID,
		}), StartTag),
	)
	return
}

func (h *xrefHandler) Exit(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("a", nil, EndTag),
	)
	return
}
