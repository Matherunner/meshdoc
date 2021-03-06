package meshdoc

import (
	"bufio"
	"fmt"

	"github.com/Matherunner/meshforce"
	"github.com/Matherunner/meshforce/tree"
)

type TreeByPath map[GenericPath]*tree.Tree

type ComponentOptions struct {
	BookReader     BookReader
	BookWriter     BookWriter
	ParsedReader   func(treeByPath TreeByPath) ParsedReader
	ParsedWriter   func(ctx *Context) ParsedWriter
	Preprocessors  []Preprocessor
	Postprocessors []Postprocessor
}

type MeshdocOptions struct {
	Config     *MeshdocConfig
	Components ComponentOptions
}

type Meshdoc struct {
	options *MeshdocOptions
}

func NewMeshdoc(options *MeshdocOptions) *Meshdoc {
	return &Meshdoc{
		options: options,
	}
}

func (t *Meshdoc) newParser() (parser *meshforce.Parser) {
	definitions := newNodeDefinitions()
	parser = meshforce.NewParser()
	definitions.Register(parser)
	return parser
}

func (t *Meshdoc) setInputFilesToContext(ctx *Context, files map[GenericPath]FileReader) {
	inputFiles := make([]GenericPath, 0, len(files))
	for k := range files {
		inputFiles = append(inputFiles, k)
	}
	ctx.SetInputFiles(inputFiles)
}

func (t *Meshdoc) parseFiles(files map[GenericPath]FileReader) (treeByPath map[GenericPath]*tree.Tree, err error) {
	treeByPath = map[GenericPath]*tree.Tree{}

	for path, r := range files {
		scanner := bufio.NewScanner(r)
		lineScanner := NewLineScanner(scanner)

		parser := t.newParser()
		parser.Parse(lineScanner)

		for _, msg := range parser.Messages().Messages() {
			if msg.Kind == meshforce.MessageKindError {
				// TODO: print all errors instead of just the first one?
				err = fmt.Errorf("error when parsing %s: %s", path, msg.Name)
				return
			}
		}

		treeByPath[path] = parser.Tree()
	}

	return
}

func (t *Meshdoc) Run() (err error) {
	ctx := NewContext()

	ctx.SetConfig(t.options.Config)

	files, err := t.options.Components.BookReader.Files(ctx)
	if err != nil {
		return
	}

	t.setInputFilesToContext(ctx, files)

	for path, r := range files {
		for _, p := range t.options.Components.Preprocessors {
			files[path] = p.Process(r)
		}
	}

	treeByPath, err := t.parseFiles(files)
	if err != nil {
		return
	}

	parsedReader := t.options.Components.ParsedReader(treeByPath)
	for _, p := range t.options.Components.Postprocessors {
		parsedReader, err = p.Process(ctx, parsedReader)
		if err != nil {
			return
		}
	}

	err = t.options.Components.BookWriter.Write(ctx, parsedReader)
	if err != nil {
		return
	}

	return nil
}

// TODO: should make this configurable?
type nodeDefinitions struct {
	blockDefs  []meshforce.BlockDefinition
	inlineDefs []meshforce.InlineDefinition
}

func newNodeDefinitions() *nodeDefinitions {
	return &nodeDefinitions{
		blockDefs: []meshforce.BlockDefinition{{
			Name:   "TOC",
			Struct: meshforce.BlockStructLiteral,
			Policy: meshforce.BlockPolicyReadUntilClose,
		}, {
			Name:   "C",
			Struct: meshforce.BlockStructLiteral,
			Policy: meshforce.BlockPolicyReadUntilClose,
		}, {
			Name:   "DL",
			Struct: meshforce.BlockStructCompound,
			Policy: meshforce.BlockPolicyReadUntilClose,
		}, {
			Name:   "DT",
			Struct: meshforce.BlockStructSimple,
			Policy: meshforce.BlockPolicyReadOne,
		}, {
			Name:   "DD",
			Struct: meshforce.BlockStructSimple,
			Policy: meshforce.BlockPolicyReadGreedy,
		}, {
			Name:   "MATH",
			Struct: meshforce.BlockStructLiteral,
			Policy: meshforce.BlockPolicyReadUntilClose,
		}, {
			Name:   "THEOREM",
			Struct: meshforce.BlockStructSimple,
			Policy: meshforce.BlockPolicyReadUntilClose,
		}, {
			Name:   "PROOF",
			Struct: meshforce.BlockStructSimple,
			Policy: meshforce.BlockPolicyReadUntilClose,
		}},
		inlineDefs: []meshforce.InlineDefinition{{
			Name:   "XREF",
			Struct: meshforce.InlineStructSimple,
		}, {
			Name:   "MATH",
			Struct: meshforce.InlineStructLiteral,
		}, {
			Name:   "A",
			Struct: meshforce.InlineStructSimple,
		}},
	}
}

func (d *nodeDefinitions) Register(p *meshforce.Parser) {
	for _, def := range meshforce.BlockDefinitions {
		p.RegisterBlock(def)
	}
	for _, def := range d.blockDefs {
		p.RegisterBlock(def)
	}
	for _, def := range meshforce.InlineDefinitions {
		p.RegisterInline(def)
	}
	for _, def := range d.inlineDefs {
		p.RegisterInline(def)
	}
}
