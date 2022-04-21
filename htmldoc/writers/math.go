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

type mathBlockHandler struct {
}

func NewMathBlockHandler() HTMLBlockWriterHandler {
	return &mathBlockHandler{}
}

func (h *mathBlockHandler) Name() string {
	return "MATH"
}

func (h *mathBlockHandler) Enter(ctx *meshdoc.Context, block *tree.BlockNode, node *tree.Node, stack []*tree.Node) (items []HTMLItem, instruction tree.VisitInstruction, err error) {
	type mathInput struct {
		Input   string            `json:"input"`
		Display bool              `json:"display"`
		Macros  map[string]string `json:"macros"`
	}

	config := ctx.Config()
	if config.KatexRendererPath == "" {
		err = ErrKatexRendererPathNotSet
		return
	}
	if config.NodeExecPath == "" {
		err = ErrNodeExecNotSet
		return
	}

	items = append(items,
		NewHTMLItemTag("div", nil, StartTag),
	)

	entryJS := path.Join(config.KatexRendererPath, "index.js")
	cmd := exec.Command(config.NodeExecPath, entryJS)
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
			Display: true,
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

	items = append(items,
		NewHTMLItemDangerousText(outHTMLs[0]),
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
