package writers

import (
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

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
