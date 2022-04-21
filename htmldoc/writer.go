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
	"github.com/Matherunner/meshdoc/htmldoc/writers"
	"github.com/Matherunner/meshforce/tree"
	"github.com/Matherunner/meshforce/writer/html"
)

const (
	htmlExt = ".html"
)

type DefaultParsedWriter struct {
	writer *html.Writer
}

func NewDefaultParsedWriter(ctx *meshdoc.Context, mapOutputFile writers.OutputFileMapper) meshdoc.ParsedWriter {
	writer := html.NewWriter()

	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewTitleHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH1Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH2Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH3Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH4Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH5Handler()))

	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewTOCHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewParagraphHandler()))

	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewMathBlockHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewTheoremHandler()))

	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewStrongHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewEmphasisHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewCodeHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewXRefHandler(mapOutputFile)))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewMathInlineHandler()))

	return &DefaultParsedWriter{
		writer: writer,
	}
}

func (p *DefaultParsedWriter) Write(w io.Writer, root *tree.Node) (err error) {
	err = p.writer.Write2(w, root)
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

func (w *DefaultBookWriter) getOutputFilePath(input meshdoc.GenericPath) string {
	return input.SetExt(htmlExt).Path()
}

func (w *DefaultBookWriter) writeFile(ctx *meshdoc.Context, filePath string, tmpl *template.Template, parseTree *tree.Tree, navigations []navigation) (err error) {
	var buf bytes.Buffer
	renderer := NewDefaultParsedWriter(ctx, w.getOutputFilePath)
	err = renderer.Write(&buf, parseTree.Root())
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

func (w *DefaultBookWriter) renderSnippet(ctx *meshdoc.Context, root *tree.Node) (html string, err error) {
	var buf bytes.Buffer
	renderer := NewDefaultParsedWriter(ctx, w.getOutputFilePath)
	err = renderer.Write(&buf, root)
	html = buf.String()
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
		title, err := w.renderSnippet(ctx, entry.Title)
		if err != nil {
			return err
		}
		navigations = append(navigations, navigation{
			Number: number,
			Title:  template.HTML(title),
			Path:   path,
		})
	}

	for filePath, tree := range reader.Files() {
		outputFile := w.getOutputFilePath(filePath)
		physicalPath := path.Join(config.OutputPath, outputFile)

		log.Printf("writing to %s", filePath)

		err = w.writeFile(ctx, physicalPath, pageTmpl, tree, navigations)
		if err != nil {
			return
		}
	}

	return nil
}

type navigation struct {
	Number string
	Title  template.HTML
	Path   string
}

type pageTemplateData struct {
	HTMLContent template.HTML
	Navigations []navigation
}
