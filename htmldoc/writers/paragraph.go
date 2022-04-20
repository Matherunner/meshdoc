package writers

import (
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

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
