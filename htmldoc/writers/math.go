package writers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

var (
	ErrKatexRendererPathNotSet = errors.New("katex renderer path not set in config")
	ErrNodeExecNotSet          = errors.New("node executable path not set in config")
)

func getRendererPaths(ctx *meshdoc.Context) (nodePath, rendererPath string, err error) {
	config := ctx.Config()
	nodePath = config.NodeExecPath
	rendererPath = config.KatexRendererPath
	if rendererPath == "" {
		err = ErrKatexRendererPathNotSet
		return
	}
	if nodePath == "" {
		err = ErrNodeExecNotSet
		return
	}
	return
}

func renderByKatex(node *tree.Node, displayMode bool, nodePath, rendererPath string) (html string, err error) {
	type mathInput struct {
		Input   string            `json:"input"`
		Display bool              `json:"display"`
		Macros  map[string]string `json:"macros"`
	}

	entryJS := path.Join(rendererPath, "index.js")
	cmd := exec.Command(nodePath, entryJS)
	wc, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	go func() {
		defer wc.Close()

		content := tree.CoalesceStringContent(node)
		enc := json.NewEncoder(wc)
		enc.Encode([]mathInput{{
			Input:   content,
			Display: displayMode,
			Macros:  map[string]string{},
		}})
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("failed to render math: %w\noutput = %s", err, out)
		return
	}

	var outHTMLs []string
	err = json.Unmarshal(out, &outHTMLs)
	if err != nil {
		return
	}

	html = outHTMLs[0]
	return
}

type mathInlineHandler struct {
}

func NewMathInlineHandler() HTMLInlineWriterHandler {
	return &mathInlineHandler{}
}

func (h *mathInlineHandler) Name() string {
	return "MATH"
}

func (h *mathInlineHandler) Enter(ctx *meshdoc.Context, inline *tree.InlineNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	nodePath, rendererPath, err := getRendererPaths(ctx)
	if err != nil {
		return
	}

	html, err := renderByKatex(node, false, nodePath, rendererPath)
	if err != nil {
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
	nodePath, rendererPath, err := getRendererPaths(ctx)
	if err != nil {
		return
	}

	html, err := renderByKatex(node, true, nodePath, rendererPath)
	if err != nil {
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
