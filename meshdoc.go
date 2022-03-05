package meshdoc

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/Matherunner/meshforce"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type lineScanner struct {
	scanner *bufio.Scanner
	i       int
	eof     bool
	peekBuf string
	peeked  bool
}

func newLineScanner(scanner *bufio.Scanner) *lineScanner {
	eof := !scanner.Scan()
	return &lineScanner{scanner: scanner, eof: eof}
}

func (s *lineScanner) Scan() bool {
	s.i++
	s.eof = !s.scanner.Scan()
	s.peeked = false
	return !s.eof
}

func (s *lineScanner) Line() string {
	if s.peeked {
		return s.peekBuf
	}
	return s.scanner.Text()
}

func (s *lineScanner) LineNumber() int {
	return s.i
}

func (s *lineScanner) EOF() bool {
	return s.eof
}

func (s *lineScanner) Peek() bool {
	if s.eof {
		return false
	}
	s.peekBuf = s.scanner.Text()
	s.peeked = true
	return true
}

type nodeDefinitions struct {
}

func newNodeDefinitions() *nodeDefinitions {
	return &nodeDefinitions{}
}

func (d *nodeDefinitions) Register(p *meshforce.Parser) {
	for _, def := range meshforce.BlockDefinitions {
		p.RegisterBlock(def)
	}
	for _, def := range meshforce.InlineDefinitions {
		p.RegisterInline(def)
	}
}

type treeVisitor struct {
}

func newTreeVisitor() *treeVisitor {
	return &treeVisitor{}
}

func (v *treeVisitor) Enter(cur *tree.Node, stack []*tree.Node) (instruction tree.VisitInstruction, err error) {
	return
}

func (v *treeVisitor) Exit(cur *tree.Node, stack []*tree.Node) (err error) {
	return
}

type meshdocConfig struct {
	Path string
}

type MeshdocOptions struct {
	ConfigPath string
}

type Meshdoc struct {
	options *MeshdocOptions

	config meshdocConfig
}

func NewMeshdoc(options *MeshdocOptions) *Meshdoc {
	return &Meshdoc{
		options: options,
	}
}

func (m *Meshdoc) Run() (err error) {
	_, err = toml.DecodeFile(m.options.ConfigPath, &m.config)
	if err != nil {
		return
	}

	fileInfo, err := ioutil.ReadDir(m.config.Path)
	if err != nil {
		return
	}

	treeByPath := map[string]*tree.Tree{}

	definitions := newNodeDefinitions()

	writer := html.NewWriter()
	writer.RegisterBlockHandler(&TitleHandler{})
	writer.RegisterBlockHandler(&H1Handler{})
	writer.RegisterBlockHandler(&ParagraphHandler{})
	writer.RegisterInlineHandler(&StrongHandler{})
	writer.RegisterInlineHandler(&EmphasisHandler{})

	for _, info := range fileInfo {
		filePath := path.Join(m.config.Path, info.Name())
		log.Printf("Parsing %s\n", filePath)

		var file *os.File
		file, err = os.Open(filePath)
		if err != nil {
			return
		}

		scanner := bufio.NewScanner(file)
		peeker := newLineScanner(scanner)
		parser := meshforce.NewParser()

		definitions.Register(parser)

		parser.Parse(peeker)

		msgs := parser.Messages().Messages()
		for _, msg := range msgs {
			if msg.Kind == meshforce.MessageKindError {
				return fmt.Errorf("parse error: %+v", msg)
			}
		}

		treeByPath[filePath] = parser.Tree()

		buf := bytes.Buffer{}
		err = writer.Write2(&buf, parser.Tree())
		if err != nil {
			return
		}

		content := buf.String()
		log.Printf("content = %+v\n", content)

		// visitor := newTreeVisitor()
		// parser.Tree().Visit(visitor)
	}

	return
}
