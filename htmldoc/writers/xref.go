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

type XRefHandler struct {
	Ctx           *meshdoc.Context
	MapOutputFile OutputFileMapper
}

func (h *XRefHandler) Name() string {
	return "XREF"
}

func (h *XRefHandler) Enter(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	refStore := xref.FromContext(h.Ctx)
	target, ok := refStore.TargetByXRefNode(node)
	if !ok {
		return 0, ErrTargetNotFound
	}

	outputFile := h.MapOutputFile(target.Path)

	attrs := []html.Attr{{
		Name:  "href",
		Value: outputFile + "#" + target.ID,
	}}

	err = enc.Start("a", attrs)
	return
}

func (h *XRefHandler) Exit(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("a")
	return
}
