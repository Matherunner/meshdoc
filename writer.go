package meshdoc

import (
	"io"

	"github.com/Matherunner/meshforce/tree"
)

type FileWriter interface {
	io.Writer
}

type BookWriter interface {
	Write(ctx *Context, reader ParsedReader) error
}

type ParsedWriter interface {
	// TODO: context!
	Write(w io.Writer, tree *tree.Tree) (err error)
}
