package meshdoc

import (
	"io"

	"github.com/Matherunner/meshforce/tree"
)

type TitleHandler struct {
}

func (h *TitleHandler) Name() string {
	return "TITLE"
}

func (t *TitleHandler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	_, err = io.WriteString(w, `<h1 class="page-header">`)
	return
}

func (t *TitleHandler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</h1>`)
	return
}

type H1Handler struct {
}

func (h *H1Handler) Name() string {
	return "H1"
}

func (t *H1Handler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	_, err = io.WriteString(w, `<h1>`)
	return
}

func (t *H1Handler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</h1>`)
	return
}

type ParagraphHandler struct {
}

func (h *ParagraphHandler) Name() string {
	return "P"
}

func (t *ParagraphHandler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	_, err = io.WriteString(w, `<p>`)
	return
}

func (t *ParagraphHandler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</p>`)
	return
}

type EmphasisHandler struct {
}

func (h *EmphasisHandler) Name() string {
	return "E"
}

func (h *EmphasisHandler) Enter(w io.Writer, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	_, err = io.WriteString(w, "<em>")
	return
}

func (h *EmphasisHandler) Exit(w io.Writer, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, "</em>")
	return
}

type StrongHandler struct {
}

func (h *StrongHandler) Name() string {
	return "S"
}

func (h *StrongHandler) Enter(w io.Writer, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	_, err = io.WriteString(w, "<strong>")
	return
}

func (h *StrongHandler) Exit(w io.Writer, block *tree.InlineNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, "</strong>")
	return
}
