package writers

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshdoc/htmldoc/processors/xref"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type H1Handler struct {
	Ctx *meshdoc.Context
}

func (h *H1Handler) Name() string {
	return "H1"
}

func (h *H1Handler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(h.Ctx)
	num := nums.ElementNumber(node)

	if err = enc.Start("h1", nil); err != nil {
		return
	}

	if err = enc.Start("span", nil); err != nil {
		return
	}

	if err = enc.Text(num); err != nil {
		return
	}

	if err = enc.End("span"); err != nil {
		return
	}

	if err = enc.DangerousText(`.&nbsp;`); err != nil {
		return
	}

	return
}

func (h *H1Handler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("h1")
	return
}

type H2Handler struct {
	Ctx *meshdoc.Context
}

func (h *H2Handler) Name() string {
	return "H2"
}

// FIXME: the handling of "ID" should be common to all, so maybe have another layer to wrap this Enter,
//  instead of writing to enc directly, return a HTML node "model", maybe a list of []HTMLItem, like NewHTMLItem(name, attrs),
//  then the bottom layer could add common attributes like ID or other things

func (h *H2Handler) Enter(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	nums := counter.FromContext(h.Ctx)
	num := nums.ElementNumber(node)

	var attrs []html.Attr

	refStore := xref.FromContext(h.Ctx)
	target, ok := refStore.IDByTargetNode(node)
	if ok {
		attrs = append(attrs, html.Attr{
			Name:  "id",
			Value: target.ID,
		})
	}

	if err = enc.Start("h2", attrs); err != nil {
		return
	}

	if err = enc.Start("span", nil); err != nil {
		return
	}

	if err = enc.Text(num); err != nil {
		return
	}

	if err = enc.End("span"); err != nil {
		return
	}

	if err = enc.DangerousText(`.&nbsp;`); err != nil {
		return
	}

	return
}

func (h *H2Handler) Exit(enc *html.Encoder, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (err error) {
	err = enc.End("h2")
	return
}
