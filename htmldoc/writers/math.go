package writers

import (
	"errors"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/math"
	"github.com/Matherunner/meshforce/tree"
)

var (
	ErrContextValueNotFound = errors.New("rendered math not found in context")
	ErrMathNotPrerendered   = errors.New("math node not pre-rendered")
)

type mathInlineHandler struct {
}

func NewMathInlineHandler() HTMLInlineWriterHandler {
	return &mathInlineHandler{}
}

func (h *mathInlineHandler) Name() string {
	return "MATH"
}

func (h *mathInlineHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	htmlByNode := math.FromContext(ctx)
	if htmlByNode == nil {
		err = ErrContextValueNotFound
		return
	}

	html, ok := htmlByNode[node]
	if !ok {
		err = ErrMathNotPrerendered
		return
	}

	items = append(items,
		NewHTMLItemTag("span", nil, StartTag),
		NewHTMLItemDangerousText(html),
	)

	instruction = tree.InstructionIgnoreChild

	return
}

func (h *mathInlineHandler) Exit(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("span", nil, EndTag),
	)
	return
}

type mathBlockHandler struct {
}

func NewMathBlockHandler() HTMLBlockWriterHandler {
	return &mathBlockHandler{}
}

func (h *mathBlockHandler) Name() string {
	return "MATH"
}

func (h *mathBlockHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	htmlByNode := math.FromContext(ctx)
	if htmlByNode == nil {
		err = ErrContextValueNotFound
		return
	}

	html, ok := htmlByNode[node]
	if !ok {
		err = ErrMathNotPrerendered
		return
	}

	items = append(items,
		NewHTMLItemTag("div", nil, StartTag),
		NewHTMLItemDangerousText(html),
	)

	instruction = tree.InstructionIgnoreChild

	return
}

func (h *mathBlockHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	items = append(items,
		NewHTMLItemTag("div", nil, EndTag),
	)
	return
}
