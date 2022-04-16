package meshdoc

import (
	"bufio"
	"fmt"

	"github.com/Matherunner/meshdoc/context"
	"github.com/Matherunner/meshforce"
	"github.com/Matherunner/meshforce/tree"
)

type GenericPath string

type TreeByPath map[GenericPath]*tree.Tree

type ComponentOptions struct {
	BookReader     BookReader
	BookWriter     BookWriter
	ParsedReader   func(treeByPath TreeByPath) ParsedReader
	ParsedWriter   ParsedWriter
	Preprocessors  []Preprocessor
	Postprocessors []Postprocessor
}

type MeshdocOptions struct {
	Config     *MeshdocConfig
	Components ComponentOptions
}

type MeshdocConfig struct {
	SourcePath   string
	TemplatePath string
	OutputPath   string
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

func (t *Meshdoc) Run() (err error) {
	ctx := context.NewDefaultContext()

	ConfigToContext(ctx, t.options.Config)

	files, err := t.options.Components.BookReader.Files(ctx)
	if err != nil {
		return
	}

	for path, r := range files {
		for _, p := range t.options.Components.Preprocessors {
			files[path] = p.Process(r)
		}
	}

	treeByPath := map[GenericPath]*tree.Tree{}

	for path, r := range files {
		scanner := bufio.NewScanner(r)
		lineScanner := NewLineScanner(scanner)

		parser := t.newParser()
		parser.Parse(lineScanner)

		for _, msg := range parser.Messages().Messages() {
			if msg.Kind == meshforce.MessageKindError {
				// TODO: print all errors instead of just the first one?
				return fmt.Errorf("error when parsing %s: %s", path, msg.Name)
			}
		}

		treeByPath[path] = parser.Tree()
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

type nodeDefinitions struct {
	blockDefs []meshforce.BlockDefinition
}

func newNodeDefinitions() *nodeDefinitions {
	return &nodeDefinitions{
		blockDefs: []meshforce.BlockDefinition{{
			Name:   "TOC",
			Struct: meshforce.BlockStructLiteral,
			Policy: meshforce.BlockPolicyReadUntilClose,
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
}
