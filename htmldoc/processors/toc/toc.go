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

type Entry struct {
	Title *tree.Node
	Path  meshdoc.GenericPath
}

type TOC struct {
}

func NewTOC() meshdoc.Postprocessor {
	return &TOC{}
}

func (t *TOC) extractTitle(path meshdoc.GenericPath, r meshdoc.ParsedReader) *tree.Node {
	parseTree, ok := r.Files()[path]
	if !ok {
		return nil
	}

	it := tree.NewIterator(parseTree.Root())
	for it.Next(tree.InstructionEnterChild) {
		node := it.Value()
		if !it.Exit() {
			if block, ok := node.Value.(*tree.BlockNode); ok {
				if block.Name() == "TITLE" {
					return node.Child
				}
			}
		}
	}

	return nil
}

func (t *TOC) parseTOCList(ctx *meshdoc.Context, content string, r meshdoc.ParsedReader) (toc []Entry, err error) {
	inputFiles := make([]string, 0, len(ctx.InputFiles()))
	for _, file := range ctx.InputFiles() {
		inputFiles = append(inputFiles, file.WithoutExt())
	}

	toc = make([]Entry, 0, 8)
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

		toc = append(toc, Entry{
			Title: t.extractTitle(p, r),
			Path:  p,
		})
	}
	return
}

func (t *TOC) Process(ctx *meshdoc.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	// TODO: extracting text seems like a very common operation, maybe can add a method to Node?

	builder := strings.Builder{}
	foundTOC := false
	for _, t := range r.Files() {
		it := tree.NewIterator(t.Root())
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
		toc, err := t.parseTOCList(ctx, builder.String(), r)
		if err != nil {
			return nil, err
		}

		setToContext(ctx, &contextValue{
			toc: toc,
		})
	}

	return r, nil
}
