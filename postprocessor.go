package meshdoc

type Postprocessor interface {
	Process(ctx *Context, r ParsedReader) (ParsedReader, error)
	// TODO: add "Dependencies()", then use topological sorting to make sure postprocessors run in the correct order
}
