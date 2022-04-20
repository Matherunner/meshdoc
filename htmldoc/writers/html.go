package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
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

type Attributes struct {
	attrs []html.Attr
}

func NewAttributes(attrs ...html.Attr) *Attributes {
	return &Attributes{attrs: attrs}
}

func (a *Attributes) Add(name, value string) {
	a.attrs = append(a.attrs, html.Attr{
		Name:  name,
		Value: value,
	})
}

func (a *Attributes) Has(name string) bool {
	for _, attr := range a.attrs {
		if attr.Name == name {
			return true
		}
	}
	return false
}

func (a *Attributes) AddIfNotExists(name, value string) {
	if !a.Has(name) {
		a.Add(name, value)
	}
}

func (a *Attributes) Slice() []html.Attr {
	return a.attrs
}

type HTMLItem interface {
	Type() HTMLItemType
	Attrs() *Attributes
}

type HTMLItemText struct {
	content string
	attrs   *Attributes
}

func NewHTMLItemText(content string) *HTMLItemText {
	return &HTMLItemText{content: content, attrs: NewAttributes()}
}

func (t *HTMLItemText) Type() HTMLItemType {
	return HTMLItemTypeText
}

func (t *HTMLItemText) Attrs() *Attributes {
	return t.attrs
}

func (t *HTMLItemText) Content() string {
	return t.content
}

type HTMLItemTag struct {
	tag     string
	attrs   *Attributes
	tagType TagType
}

func NewHTMLItemTag(tag string, attrs *Attributes, tagType TagType) *HTMLItemTag {
	if attrs == nil {
		attrs = NewAttributes()
	}
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

func (t *HTMLItemTag) Attrs() *Attributes {
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
				enc.Start(item.Tag(), item.Attrs().Slice())
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

	// Add ID to the first child by default
	if len(items) != 0 {
		first := items[0]
		if id, ok := block.Options().Get("ID"); ok {
			first.Attrs().AddIfNotExists("id", id)
		}
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
