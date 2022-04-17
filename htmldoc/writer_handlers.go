package htmldoc

import (
	"errors"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshdoc/htmldoc/processors/xref"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

var (
	ErrTargetNotFound = errors.New("target not found for xref")
)

type TitleHandler struct {
	Ctx *meshdoc.Context
}

func (h *TitleHandler) Name() string {
	return "TITLE"
}

func (t *TitleHandler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(t.Ctx)
	num := nums.ElementNumber(node)

	attrs := []html.Attr{{
		Name:  "class",
		Value: "page-header",
	}}

	if err = enc.Start("h1", attrs); err != nil {
		return
	}

	if err = enc.Start("span", nil); err != nil {
		return
	}

	if err = enc.Text(num); err != nil {
		return
	}

	if err = enc.End("span"); err != nil {
		return
	}

	if err = enc.DangerousText(`.&nbsp;`); err != nil {
		return
	}

	return
}

func (t *TitleHandler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("h1")
	return
}

type H1Handler struct {
	Ctx *meshdoc.Context
}

func (h *H1Handler) Name() string {
	return "H1"
}

func (h *H1Handler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(h.Ctx)
	num := nums.ElementNumber(node)

	if err = enc.Start("h1", nil); err != nil {
		return
	}

	if err = enc.Start("span", nil); err != nil {
		return
	}

	if err = enc.Text(num); err != nil {
		return
	}

	if err = enc.End("span"); err != nil {
		return
	}

	if err = enc.DangerousText(`.&nbsp;`); err != nil {
		return
	}

	return
}

func (h *H1Handler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("h1")
	return
}

type H2Handler struct {
	Ctx *meshdoc.Context
}

func (h *H2Handler) Name() string {
	return "H2"
}

func (h *H2Handler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(h.Ctx)
	num := nums.ElementNumber(node)

	var attrs []html.Attr

	refStore := xref.FromContext(h.Ctx)
	target, ok := refStore.IDByTargetNode(node)
	if ok {
		attrs = append(attrs, html.Attr{
			Name:  "id",
			Value: target.ID,
		})
	}

	if err = enc.Start("h2", attrs); err != nil {
		return
	}

	if err = enc.Start("span", nil); err != nil {
		return
	}

	if err = enc.Text(num); err != nil {
		return
	}

	if err = enc.End("span"); err != nil {
		return
	}

	if err = enc.DangerousText(`.&nbsp;`); err != nil {
		return
	}

	return
}

func (h *H2Handler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("h2")
	return
}

type TOCHandler struct {
}

func (h *TOCHandler) Name() string {
	return "TOC"
}

func (h *TOCHandler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	err = enc.Start("pre", nil)
	return
}

func (t *TOCHandler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("pre")
	return
}

type ParagraphHandler struct {
}

func (h *ParagraphHandler) Name() string {
	return "P"
}

func (t *ParagraphHandler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	err = enc.Start("p", nil)
	return
}

func (t *ParagraphHandler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("p")
	return
}

type EmphasisHandler struct {
}

func (h *EmphasisHandler) Name() string {
	return "E"
}

func (h *EmphasisHandler) Enter(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	err = enc.Start("em", nil)
	return
}

func (h *EmphasisHandler) Exit(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("em")
	return
}

type StrongHandler struct {
}

func (h *StrongHandler) Name() string {
	return "S"
}

func (h *StrongHandler) Enter(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	err = enc.Start("strong", nil)
	return
}

func (h *StrongHandler) Exit(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("strong")
	return
}

type CodeHandler struct {
}

func (h *CodeHandler) Name() string {
	return "C"
}

func (h *CodeHandler) Enter(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	err = enc.Start("code", nil)
	return
}

func (h *CodeHandler) Exit(enc *html.Encoder, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("code")
	return
}

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
