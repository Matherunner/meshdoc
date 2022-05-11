package htmldoc

import (
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

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

	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewCodeBlockHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewCommentHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewDefinitionBodyHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewDefinitionListHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewDefinitionTermHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH1Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH2Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH3Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH4Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewH5Handler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewListItemHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewMathBlockHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewParagraphHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewProofHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewTheoremHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewTitleHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewTOCHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewUnorderedListHandler()))
	writer.RegisterBlockHandler(writers.WithBlockWriterHandler(ctx, writers.NewOrderedListHandler()))

	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewCodeHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewEmphasisHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewHyperlinkHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewMathInlineHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewStrongHandler()))
	writer.RegisterInlineHandler(writers.WithInlineWriterHandler(ctx, writers.NewXRefHandler(mapOutputFile)))

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
	var buf strings.Builder
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
	var buf strings.Builder
	renderer := NewDefaultParsedWriter(ctx, w.getOutputFilePath)
	err = renderer.Write(&buf, root)
	html = buf.String()
	return
}

func (w *DefaultBookWriter) Write(ctx *meshdoc.Context, reader meshdoc.ParsedReader) (err error) {
	config := ctx.Config()

	err = os.MkdirAll(config.OutputPath, os.ModePerm)
	if err != nil {
		return
	}

	navigations, err := w.generateNavigtations(ctx)
	if err != nil {
		return
	}

	err = w.parallelRenderFiles(ctx, reader, navigations)
	if err != nil {
		return
	}

	return nil
}

func (w *DefaultBookWriter) generateNavigtations(ctx *meshdoc.Context) (navigations []navigation, err error) {
	counterValues := counter.FromContext(ctx)
	tableOfContents := toc.FromContext(ctx)
	navigations = make([]navigation, 0, len(tableOfContents))
	for _, entry := range tableOfContents {
		number := counterValues.FileNumber(entry.Path)
		path := w.tocToWebPath(entry.Path)
		title, err := w.renderSnippet(ctx, entry.Title)
		if err != nil {
			return nil, err
		}
		navigations = append(navigations, navigation{
			Number: number,
			Title:  template.HTML(title),
			Path:   path,
		})
	}
	return
}

func (w *DefaultBookWriter) parallelRenderFiles(ctx *meshdoc.Context, reader meshdoc.ParsedReader, navigations []navigation) (err error) {
	config := ctx.Config()

	pageTmpl, err := w.parseTemplates(config.TemplatePath)
	if err != nil {
		return
	}

	var writeErr error
	wg := sync.WaitGroup{}

	for filePath, tree := range reader.Files() {
		filePath := filePath
		tree := tree

		wg.Add(1)

		go func() {
			defer wg.Done()

			outputFile := w.getOutputFilePath(filePath)
			physicalPath := path.Join(config.OutputPath, outputFile)

			log.Printf("writing to %s", filePath)

			err = w.writeFile(ctx, physicalPath, pageTmpl, tree, navigations)
			if err != nil {
				log.Printf("writing %s failed: %+v", filePath, err)
				writeErr = err
				return
			}

			log.Printf("writing %s done", filePath)
		}()
	}

	wg.Wait()

	err = writeErr

	return
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
