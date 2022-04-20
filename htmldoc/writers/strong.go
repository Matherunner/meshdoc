package writers

import (
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

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
