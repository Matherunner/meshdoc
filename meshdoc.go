package meshdoc

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Matherunner/meshforce"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type GenericPath string

type Context interface {
	Set(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

type DefaultContext struct {
	kv map[interface{}]interface{}
}

func NewDefaultContext() Context {
	return &DefaultContext{
		kv: map[interface{}]interface{}{},
	}
}

func (c *DefaultContext) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = c.kv[key]
	return
}

func (c *DefaultContext) Set(key, value interface{}) {
	c.kv[key] = value
}

type FileReader interface {
	io.ReadCloser
	// Context() Context
}

type DefaultFileReader struct {
	path string
	*os.File
}

func NewDefaultFileReader(path string) (r FileReader, err error) {
	fr := &DefaultFileReader{
		path: path,
	}

	fr.File, err = os.Open(path)
	if err != nil {
		return
	}

	return fr, nil
}

type BookReader interface {
	Files(config *MeshdocConfig) (readers map[GenericPath]FileReader, err error)
}

type FileWriter interface {
	io.Writer
}

type BookWriter interface {
	Write(ctx Context, config *MeshdocConfig, reader ParsedReader) error
}

var (
	defaultPageTemplate *template.Template
)

type DefaultBookWriter struct {
}

func NewDefaultBookWriter() BookWriter {
	return &DefaultBookWriter{}
}

func (w *DefaultBookWriter) parseTemplates(dir string) (tmpl *template.Template, err error) {
	if defaultPageTemplate != nil {
		return defaultPageTemplate, nil
	}

	pagePath := path.Join(dir, "page.tmpl")
	contentPath := path.Join(dir, "content.tmpl")

	defaultPageTemplate, err = template.ParseFiles(pagePath, contentPath)
	if err != nil {
		return
	}

	return defaultPageTemplate, nil
}

func (w *DefaultBookWriter) writeFile(ctx Context, filePath string, parseTree *tree.Tree, tmpl *template.Template) (err error) {
	var buf bytes.Buffer
	renderer := NewDefaultParsedWriter()
	err = renderer.Write(&buf, parseTree)
	if err != nil {
		return err
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	htmlContent := template.HTML(buf.String())
	err = tmpl.Execute(outFile, &pageTemplateData{
		HTMLContent: htmlContent,
		// TODO: fill in navigations from ctx
		Navigations: []navigation{},
	})

	return
}

func (w *DefaultBookWriter) Write(ctx Context, config *MeshdocConfig, reader ParsedReader) (err error) {
	pageTmpl, err := w.parseTemplates(config.TemplatePath)
	if err != nil {
		return
	}

	err = os.MkdirAll(config.OutputPath, os.ModePerm)
	if err != nil {
		return
	}

	for filePath, tree := range reader.Files() {
		filePath := string(filePath)
		filePath = strings.TrimSuffix(filePath, path.Ext(filePath))
		filePath += ".html"

		err = w.writeFile(ctx, filePath, tree, pageTmpl)
		if err != nil {
			return
		}
	}

	return nil
}

type Preprocessor interface {
	Process(r FileReader) FileReader
}

type ParsedReader interface {
	// TODO: also messages
	Files() map[GenericPath]*tree.Tree
	// Context() Context
}

type DefaultParsedReader struct {
	m map[GenericPath]*tree.Tree
}

func NewDefaultParsedReader(m map[GenericPath]*tree.Tree) ParsedReader {
	return &DefaultParsedReader{m: m}
}

func (r *DefaultParsedReader) Files() map[GenericPath]*tree.Tree {
	return r.m
}

type Postprocessor interface {
	Process(ctx Context, r ParsedReader) ParsedReader
}

type ParsedWriter interface {
	// TODO: context!
	Write(w io.Writer, tree *tree.Tree) (err error)
}

type DefaultParsedWriter struct {
	writer *html.Writer
}

func NewDefaultParsedWriter() ParsedWriter {
	writer := html.NewWriter()
	writer.RegisterBlockHandler(&TitleHandler{})
	writer.RegisterBlockHandler(&H1Handler{})
	writer.RegisterBlockHandler(&ParagraphHandler{})
	writer.RegisterInlineHandler(&StrongHandler{})
	writer.RegisterInlineHandler(&EmphasisHandler{})

	return &DefaultParsedWriter{
		writer: writer,
	}
}

func (p *DefaultParsedWriter) Write(w io.Writer, tree *tree.Tree) (err error) {
	err = p.writer.Write2(w, tree)
	return
}

type ConfigProvider interface {
	Config(options *MeshdocOptions) (config *MeshdocConfig, err error)
}

type DefaultConfigProvider struct {
}

func NewDefaultConfigProvider() ConfigProvider {
	return &DefaultConfigProvider{}
}

func (p *DefaultConfigProvider) Config(options *MeshdocOptions) (config *MeshdocConfig, err error) {
	config = &MeshdocConfig{}
	_, err = toml.DecodeFile(options.ConfigPath, config)
	return
}

type DefaultBookReader struct {
}

func NewDefaultBookReader() BookReader {
	return &DefaultBookReader{}
}

func (r *DefaultBookReader) Files(config *MeshdocConfig) (readers map[GenericPath]FileReader, err error) {
	fileInfo, err := ioutil.ReadDir(config.SourcePath)
	if err != nil {
		return
	}

	readers = map[GenericPath]FileReader{}

	for _, fi := range fileInfo {
		if path.Ext(fi.Name()) != ".mf" {
			continue
		}

		// FIXME: not always a file name, can be in a subdirectory
		docPath := fi.Name()

		filePath := path.Join(config.SourcePath, docPath)
		readers[GenericPath(filePath)], err = NewDefaultFileReader(filePath)
		if err != nil {
			return
		}
	}
	return
}

type Meshdoc2 struct {
	configProvider ConfigProvider

	bookReader     BookReader
	preprocessors  []Preprocessor
	parser         *meshforce.Parser
	postprocessors []Postprocessor
	writer         ParsedWriter
	bookWriter     BookWriter

	options *MeshdocOptions
}

func NewMeshdoc2(options *MeshdocOptions) *Meshdoc2 {
	return &Meshdoc2{
		options:        options,
		configProvider: NewDefaultConfigProvider(),
		bookReader:     NewDefaultBookReader(),
		writer:         NewDefaultParsedWriter(),
		bookWriter:     NewDefaultBookWriter(),
	}
}

func (t *Meshdoc2) AddPreprocessor(p Preprocessor) {
	t.preprocessors = append(t.preprocessors, p)
}

func (t *Meshdoc2) AddPostprocessor(p Postprocessor) {
	t.postprocessors = append(t.postprocessors, p)
}

func (t *Meshdoc2) resetParser() {
	definitions := newNodeDefinitions()
	t.parser = meshforce.NewParser()
	definitions.Register(t.parser)
}

func (t *Meshdoc2) Run() (err error) {
	ctx := NewDefaultContext()

	config, err := t.configProvider.Config(t.options)
	if err != nil {
		return
	}

	files, err := t.bookReader.Files(config)
	if err != nil {
		return
	}

	for path, r := range files {
		for _, p := range t.preprocessors {
			files[path] = p.Process(r)
		}
	}

	treeByPath := map[GenericPath]*tree.Tree{}

	for path, r := range files {
		scanner := bufio.NewScanner(r)
		lineScanner := NewLineScanner(scanner)

		t.resetParser()
		t.parser.Parse(lineScanner)

		for _, msg := range t.parser.Messages().Messages() {
			if msg.Kind == meshforce.MessageKindError {
				// TODO: print all errors instead of just the first one?
				return fmt.Errorf("error when parsing %s: %s", path, msg.Name)
			}
		}

		treeByPath[path] = t.parser.Tree()
	}

	fmt.Printf("tree by path = %+v\n", treeByPath)

	parsedReader := NewDefaultParsedReader(treeByPath)
	for _, p := range t.postprocessors {
		parsedReader = p.Process(ctx, parsedReader)
	}

	err = t.bookWriter.Write(ctx, config, parsedReader)
	if err != nil {
		return
	}

	return nil
}

type LineScanner struct {
	scanner *bufio.Scanner
	i       int
	eof     bool
	peekBuf string
	peeked  bool
}

func NewLineScanner(scanner *bufio.Scanner) *LineScanner {
	eof := !scanner.Scan()
	return &LineScanner{scanner: scanner, eof: eof}
}

func (s *LineScanner) Scan() bool {
	s.i++
	s.eof = !s.scanner.Scan()
	s.peeked = false
	return !s.eof
}

func (s *LineScanner) Line() string {
	if s.peeked {
		return s.peekBuf
	}
	return s.scanner.Text()
}

func (s *LineScanner) LineNumber() int {
	return s.i
}

func (s *LineScanner) EOF() bool {
	return s.eof
}

func (s *LineScanner) Peek() bool {
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

type navigation struct {
	Name string
	Path string
}

type pageTemplateData struct {
	HTMLContent template.HTML
	Navigations []navigation
}

type MeshdocConfig struct {
	SourcePath   string
	TemplatePath string
	OutputPath   string
}

type MeshdocOptions struct {
	ConfigPath string
}
