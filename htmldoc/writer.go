package htmldoc

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"path"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

const (
	htmlExt = ".html"
)

type DefaultParsedWriter struct {
	writer *html.Writer
}

func NewDefaultParsedWriter() meshdoc.ParsedWriter {
	writer := html.NewWriter()

	writer.RegisterBlockHandler(&TitleHandler{})
	writer.RegisterBlockHandler(&H1Handler{})
	writer.RegisterBlockHandler(&H2Handler{})
	writer.RegisterBlockHandler(&TOCHandler{})
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

var (
	defaultPageTemplate *template.Template
)

type DefaultBookWriter struct {
}

func NewDefaultBookWriter() meshdoc.BookWriter {
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

func (w *DefaultBookWriter) tocToWebPath(entry meshdoc.GenericPath) string {
	return entry.SetExt(htmlExt).WebPath()
}

func (w *DefaultBookWriter) inputToOutputFileName(config *meshdoc.MeshdocConfig, input meshdoc.GenericPath) string {
	return path.Join(config.OutputPath, input.SetExt(htmlExt).Path())
}

func (w *DefaultBookWriter) writeFile(filePath string, tmpl *template.Template, parseTree *tree.Tree, navigations []navigation) (err error) {
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
		Navigations: navigations,
	})

	return
}

func (w *DefaultBookWriter) Write(ctx *meshdoc.Context, reader meshdoc.ParsedReader) (err error) {
	config := ctx.Config()

	counterValues := counter.FromContext(ctx)

	pageTmpl, err := w.parseTemplates(config.TemplatePath)
	if err != nil {
		return
	}

	err = os.MkdirAll(config.OutputPath, os.ModePerm)
	if err != nil {
		return
	}

	tableOfContents := toc.FromContext(ctx)
	navigations := make([]navigation, 0, len(tableOfContents))
	for _, entry := range tableOfContents {
		number := counterValues.FileNumber(entry.Path)
		path := w.tocToWebPath(entry.Path)
		navigations = append(navigations, navigation{
			Name: number + ". " + entry.Name,
			Path: path,
		})
	}

	for filePath, tree := range reader.Files() {
		filePath := w.inputToOutputFileName(config, filePath)

		log.Printf("writing to %s", filePath)

		err = w.writeFile(filePath, pageTmpl, tree, navigations)
		if err != nil {
			return
		}
	}

	return nil
}

type navigation struct {
	Name string
	Path string
}

type pageTemplateData struct {
	HTMLContent template.HTML
	Navigations []navigation
}
