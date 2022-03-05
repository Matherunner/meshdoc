package meshdoc

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

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

type pageTemplateData struct {
	HTMLContent template.HTML
}

type parsedDocument struct {
	tree       *tree.Tree
	content    string
	outputPath string
}

type meshdocConfig struct {
	SourcePath   string
	TemplatePath string
	OutputPath   string
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

func (m *Meshdoc) parseTemplates(dir string) (pageTmpl *template.Template, err error) {
	pagePath := path.Join(dir, "page.tmpl")
	contentPath := path.Join(dir, "content.tmpl")

	pageTmpl, err = template.ParseFiles(pagePath, contentPath)
	if err != nil {
		return
	}

	return
}

func (m *Meshdoc) Run() (err error) {
	_, err = toml.DecodeFile(m.options.ConfigPath, &m.config)
	if err != nil {
		return
	}

	log.Printf("config = %+v\n", m.config)

	docByPath := map[string]parsedDocument{}

	definitions := newNodeDefinitions()

	writer := html.NewWriter()
	writer.RegisterBlockHandler(&TitleHandler{})
	writer.RegisterBlockHandler(&H1Handler{})
	writer.RegisterBlockHandler(&ParagraphHandler{})
	writer.RegisterInlineHandler(&StrongHandler{})
	writer.RegisterInlineHandler(&EmphasisHandler{})

	pageTmpl, err := m.parseTemplates(m.config.TemplatePath)
	if err != nil {
		return
	}

	fileInfo, err := ioutil.ReadDir(m.config.SourcePath)
	if err != nil {
		return
	}

	for _, info := range fileInfo {
		if path.Ext(info.Name()) != ".mf" {
			continue
		}

		filePath := path.Join(m.config.SourcePath, info.Name())
		log.Printf("Parsing %s\n", filePath)

		outputPath := path.Join(m.config.OutputPath, info.Name())
		outputPath = strings.TrimSuffix(outputPath, path.Ext(outputPath))
		outputPath += ".html"

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

		doc := parsedDocument{
			tree:       parser.Tree(),
			outputPath: outputPath,
		}

		buf := bytes.Buffer{}
		err = writer.Write2(&buf, parser.Tree())
		if err != nil {
			return
		}

		doc.content = buf.String()
		docByPath[filePath] = doc
	}

	err = os.MkdirAll(m.config.OutputPath, os.ModePerm)
	if err != nil {
		return
	}

	for _, doc := range docByPath {
		var outFile *os.File
		outFile, err = os.Create(doc.outputPath)
		if err != nil {
			return
		}
		defer outFile.Close()
		err = pageTmpl.Execute(outFile, &pageTemplateData{
			HTMLContent: template.HTML(doc.content),
		})
		if err != nil {
			return
		}
	}

	return
}
