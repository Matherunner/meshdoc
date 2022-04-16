package meshdoc

import (
	"io"

	"github.com/Matherunner/meshforce/tree"
)

type FileReader interface {
	io.ReadCloser
	// Context() Context
}

type BookReader interface {
	Files(ctx *Context) (readers map[GenericPath]FileReader, err error)
}

type ParsedReader interface {
	// TODO: also messages
	Files() map[GenericPath]*tree.Tree
	// Context() Context
}
