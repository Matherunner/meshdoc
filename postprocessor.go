package meshdoc

type Postprocessor interface {
	Process(ctx *Context, r ParsedReader) (ParsedReader, error)
}
