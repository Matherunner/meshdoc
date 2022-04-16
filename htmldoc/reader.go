package htmldoc

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

type DefaultParsedReader struct {
	m map[meshdoc.GenericPath]*tree.Tree
}

func NewDefaultParsedReader(m meshdoc.TreeByPath) meshdoc.ParsedReader {
	return &DefaultParsedReader{m: m}
}

func (r *DefaultParsedReader) Files() map[meshdoc.GenericPath]*tree.Tree {
	return r.m
}

type DefaultFileReader struct {
	path string
	*os.File
}

func NewDefaultFileReader(path string) (r meshdoc.FileReader, err error) {
	fr := &DefaultFileReader{
		path: path,
	}

	fr.File, err = os.Open(path)
	if err != nil {
		return
	}

	return fr, nil
}

type DefaultBookReader struct {
}

func NewDefaultBookReader() meshdoc.BookReader {
	return &DefaultBookReader{}
}

func (r *DefaultBookReader) Files(ctx *meshdoc.Context) (readers map[meshdoc.GenericPath]meshdoc.FileReader, err error) {
	config := ctx.Config()

	fileInfo, err := ioutil.ReadDir(config.SourcePath)
	if err != nil {
		return
	}

	readers = map[meshdoc.GenericPath]meshdoc.FileReader{}

	for _, fi := range fileInfo {
		if path.Ext(fi.Name()) != ".mf" {
			continue
		}

		// FIXME: not always a file name, can be in a subdirectory
		docPath := fi.Name()

		filePath := path.Join(config.SourcePath, docPath)
		readers[meshdoc.NewGenericPath(docPath)], err = NewDefaultFileReader(filePath)
		if err != nil {
			return
		}
	}
	return
}
