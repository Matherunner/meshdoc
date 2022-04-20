package writers

import (
	"errors"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshdoc/htmldoc/processors/xref"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type OutputFileMapper func(input meshdoc.GenericPath) string

var (
	ErrTargetNotFound = errors.New("target not found for xref")
)

const (
	CharEnspace = '\u2002'
	CharEmspace = '\u2003'
)

type TagType int

const (
	StartTag TagType = iota + 1
	EndTag
)

type HTMLItemType int

const (
	HTMLItemTypeTag HTMLItemType = iota + 1
	HTMLItemTypeText
)

type HTMLItem interface {
	Type() HTMLItemType
}

type HTMLItemText struct {
	content string
}

func NewHTMLItemText(content string) *HTMLItemText {
	return &HTMLItemText{content: content}
}

func (t *HTMLItemText) Type() HTMLItemType {
	return HTMLItemTypeText
}

func (t *HTMLItemText) Content() string {
	return t.content
}

type HTMLItemTag struct {
	tag     string
	attrs   []html.Attr
	tagType TagType
}

func NewHTMLItemTag(tag string, attrs []html.Attr, tagType TagType) *HTMLItemTag {
	return &HTMLItemTag{
		tag:     tag,
		attrs:   attrs,
		tagType: tagType,
	}
}

func (t *HTMLItemTag) Type() HTMLItemType {
	return HTMLItemTypeTag
}

func (t *HTMLItemTag) Tag() string {
	return t.tag
}

func (t *HTMLItemTag) Attrs() []html.Attr {
	return t.attrs
}

func (t *HTMLItemTag) TagType() TagType {
	return t.tagType
}

type HTMLBlockWriterHandler interface {
	Name() string
	Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error)
	Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error)
}

type BlockWriterHandler struct {
	ctx     *meshdoc.Context
	wrapped HTMLBlockWriterHandler
}

func WithBlockWriterHandler(ctx *meshdoc.Context, wrapped HTMLBlockWriterHandler) *BlockWriterHandler {
	return &BlockWriterHandler{
		ctx:     ctx,
		wrapped: wrapped,
	}
}

func (h *BlockWriterHandler) Name() string {
	return h.wrapped.Name()
}

func (h *BlockWriterHandler) writeItems(enc *html.Encoder, items []HTMLItem) {
	for _, item := range items {
		switch item := item.(type) {
		case *HTMLItemTag:
			if item.TagType() == StartTag {
				enc.Start(item.Tag(), item.Attrs())
			} else {
				enc.End(item.Tag())
			}
		case *HTMLItemText:
			enc.Text(item.Content())
		}
	}
}

func (h *BlockWriterHandler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	items, instruction, err := h.wrapped.Enter(h.ctx, block, node, stack)
	if err != nil {
		return
	}
	h.writeItems(enc, items)
	return
}

func (h *BlockWriterHandler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	items, err := h.wrapped.Exit(h.ctx, block, node, stack)
	if err != nil {
		return
	}
	h.writeItems(enc, items)
	return
}

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

// FIXME: the handling of "ID" should be common to all, so maybe have another layer to wrap this Enter,
//  instead of writing to enc directly, return a HTML node "model", maybe a list of []HTMLItem, like NewHTMLItem(name, attrs),
//  then the bottom layer could add common attributes like ID or other things

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
