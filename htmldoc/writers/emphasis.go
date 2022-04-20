package writers

import (
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

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
