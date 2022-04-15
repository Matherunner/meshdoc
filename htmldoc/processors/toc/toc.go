package toc

import (
	"strings"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/context"
	"github.com/Matherunner/meshforce/tree"
)

var tocContextKey = &struct{}{}

type tocContextValue struct {
	toc []meshdoc.GenericPath
}

func FromContext(ctx context.Context) (toc []meshdoc.GenericPath) {
	value, ok := ctx.Get(tocContextKey)
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

func (p *TOCPostprocessor) Process(ctx context.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	toc := make([]meshdoc.GenericPath, 0, 8)

	for _, t := range r.Files() {
		builder := strings.Builder{}
		foundTOC := false
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
			content := builder.String()
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				toc = append(toc, meshdoc.GenericPath(strings.TrimSpace(line)))
			}
			break
		}
	}

	ctx.Set(tocContextKey, &tocContextValue{
		toc: toc,
	})

	return r, nil
}
