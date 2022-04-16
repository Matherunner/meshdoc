package toc

import (
	"errors"
	"sort"
	"strings"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

var (
	ErrUnknownFile = errors.New("unknown file defined in toc")
)

type contextKeyType struct{}

var (
	contextKey = contextKeyType{}
)

type tocContextValue struct {
	toc []meshdoc.GenericPath
}

func FromContext(ctx *meshdoc.Context) (toc []meshdoc.GenericPath) {
	value, ok := ctx.Get(contextKey)
	if !ok {
		return nil
	}
	return value.(*tocContextValue).toc
}

type TOCPostprocessor struct {
}

func NewTOCPostprocessor() meshdoc.Postprocessor {
	return &TOCPostprocessor{}
}

func (p *TOCPostprocessor) parseTOCList(ctx *meshdoc.Context, content string) (toc []meshdoc.GenericPath, err error) {
	inputFiles := make([]string, 0, len(ctx.InputFiles()))
	for _, file := range ctx.InputFiles() {
		inputFiles = append(inputFiles, file.WithoutExt())
	}

	toc = make([]meshdoc.GenericPath, 0, 8)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fileName := strings.TrimSpace(line)

		n := sort.Search(len(inputFiles), func(i int) bool {
			return inputFiles[i] >= fileName
		})
		if n == len(inputFiles) {
			return nil, ErrUnknownFile
		}

		// Assume the entries all point to .mf files
		p := meshdoc.NewGenericPath(fileName).SetExt(".mf")
		toc = append(toc, p)
	}
	return
}

func (p *TOCPostprocessor) Process(ctx *meshdoc.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	builder := strings.Builder{}
	foundTOC := false
	for _, t := range r.Files() {
		it := tree.NewIterator(t)
		for it.Next(tree.InstructionEnterChild) {
			node := it.Value()

			if foundTOC && !it.Exit() {
				if node, ok := node.Value.(*tree.TextNode); ok {
					builder.WriteString(node.Content())
				}
			}

			if node, ok := node.Value.(*tree.BlockNode); ok {
				if node.Name() == "TOC" {
					if it.Exit() {
						break
					}
					foundTOC = true
				}
			}
		}

		if foundTOC {
			break
		}
	}

	if foundTOC {
		toc, err := p.parseTOCList(ctx, builder.String())
		if err != nil {
			return nil, err
		}

		ctx.Set(contextKey, &tocContextValue{
			toc: toc,
		})
	}

	return r, nil
}
