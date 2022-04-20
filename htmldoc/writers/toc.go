package writers

import (
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

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
