package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type commentHandler struct {
}

func NewCommentHandler() HTMLBlockWriterHandler {
	return &commentHandler{}
}

func (h *commentHandler) Name() string {
	return "COMMENT"
}

func (h *commentHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	instruction = tree.InstructionIgnoreChild
	return
}

func (h *commentHandler) Exit(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, err error) {
	return
}
