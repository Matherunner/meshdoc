package meshdoc

type Preprocessor interface {
	Process(r FileReader) FileReader
}
