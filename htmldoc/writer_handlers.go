package htmldoc

import (
	"fmt"
	"io"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshforce/tree"
)

type TitleHandler struct {
	Ctx *meshdoc.Context
}

func (h *TitleHandler) Name() string {
	return "TITLE"
}

func (t *TitleHandler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(t.Ctx)
	num := nums.ElementNumber(node)
	_, err = fmt.Fprintf(w, `<h1 class="page-header"><span>%s</span>.&nbsp;`, num)
	return
}

func (t *TitleHandler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</h1>`)
	return
}

type H1Handler struct {
	Ctx *meshdoc.Context
}

func (h *H1Handler) Name() string {
	return "H1"
}

func (h *H1Handler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(h.Ctx)
	num := nums.ElementNumber(node)
	_, err = fmt.Fprintf(w, `<h1><span>%s</span>.&nbsp;`, num)
	return
}

func (h *H1Handler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</h1>`)
	return
}

type H2Handler struct {
	Ctx *meshdoc.Context
}

func (h *H2Handler) Name() string {
	return "H2"
}

func (h *H2Handler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(h.Ctx)
	num := nums.ElementNumber(node)
	_, err = fmt.Fprintf(w, `<h2><span>%s</span>.&nbsp;`, num)
	return
}

func (h *H2Handler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</h2>`)
	return
}

type TOCHandler struct {
}

func (h *TOCHandler) Name() string {
	return "TOC"
}

func (h *TOCHandler) Enter(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	_, err = io.WriteString(w, `<pre>TODO!`)
	return
}

func (t *TOCHandler) Exit(w io.Writer, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	_, err = io.WriteString(w, `</pre>`)
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
