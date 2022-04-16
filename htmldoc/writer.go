package htmldoc

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/context"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

type DefaultParsedWriter struct {
	writer *html.Writer
}

func NewDefaultParsedWriter() meshdoc.ParsedWriter {
	writer := html.NewWriter()

	writer.RegisterBlockHandler(&TitleHandler{})
	writer.RegisterBlockHandler(&H1Handler{})
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

func (w *DefaultBookWriter) writeFile(ctx context.Context, filePath string, parseTree *tree.Tree, tmpl *template.Template) (err error) {
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

	tableOfContents := toc.FromContext(ctx)
	fmt.Printf("TOC = %+v\n", tableOfContents)

	htmlContent := template.HTML(buf.String())
	err = tmpl.Execute(outFile, &pageTemplateData{
		HTMLContent: htmlContent,
		// TODO: fill in navigations from ctx
		Navigations: []navigation{},
	})

	return
}

func (w *DefaultBookWriter) Write(ctx context.Context, config *meshdoc.MeshdocConfig, reader meshdoc.ParsedReader) (err error) {
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
		filePath = path.Join(config.OutputPath, filePath)

		log.Printf("writing to %s", filePath)

		err = w.writeFile(ctx, filePath, tree, pageTmpl)
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
