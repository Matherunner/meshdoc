package math

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

// NewMathBatchRenderer returns a postprocessor that renders all the inline and block math
// elements in every document. Math rendering is technically the job of writers, but we want
// to batch render all elements using a single NodeJS execution.
func NewMathBatchRenderer() meshdoc.Postprocessor {
	return &mathBatchRenderer{}
}

type mathBatchRenderer struct {
}

func (p *mathBatchRenderer) Process(ctx *meshdoc.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	var katexInputs []katexInput

	for _, tr := range r.Files() {
		it := tree.NewIterator(tr.Root())
		for it.Next(tree.InstructionEnterChild) {
			node := it.Value()
			if it.Exit() {
				continue
			}

			switch val := node.Value.(type) {
			case *tree.InlineNode:
				if val.Name() == "MATH" {
					katexInputs = append(katexInputs, katexInput{
						Node:        node,
						DisplayMode: false,
					})
				}
			case *tree.BlockNode:
				if val.Name() == "MATH" {
					katexInputs = append(katexInputs, katexInput{
						Node:        node,
						DisplayMode: true,
					})
				}
			}
		}
	}

	nodePath, rendererPath, err := getRendererPaths(ctx)
	if err != nil {
		return nil, err
	}

	htmlByNode, err := renderByKatex(katexInputs, nodePath, rendererPath)
	if err != nil {
		return nil, err
	}

	setToContext(ctx, htmlByNode)

	return r, nil
}

type katexInput struct {
	Node        *tree.Node
	DisplayMode bool
}

func renderByKatex(nodes []katexInput, nodePath, rendererPath string) (htmlByNode map[*tree.Node]string, err error) {
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

		var inputs []mathInput
		for _, node := range nodes {
			content := tree.CoalesceStringContent(node.Node)
			inputs = append(inputs, mathInput{
				Input:   content,
				Display: node.DisplayMode,
				Macros:  map[string]string{},
			})
		}

		enc := json.NewEncoder(wc)
		enc.Encode(inputs)
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

	htmlByNode = map[*tree.Node]string{}
	for i, node := range nodes {
		htmlByNode[node.Node] = outHTMLs[i]
	}
	return
}

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
