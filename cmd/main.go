package main

import (
	"log"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
)

func main() {
	config, err := meshdoc.LoadFromFile("./examples/simple/config.toml")
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	meshdoc := meshdoc.NewMeshdoc(&meshdoc.MeshdocOptions{
		Config: config,
		Components: meshdoc.ComponentOptions{
			BookReader:    htmldoc.NewDefaultBookReader(),
			BookWriter:    htmldoc.NewDefaultBookWriter(),
			ParsedReader:  htmldoc.NewDefaultParsedReader,
			ParsedWriter:  htmldoc.NewDefaultParsedWriter(),
			Preprocessors: []meshdoc.Preprocessor{},
			Postprocessors: []meshdoc.Postprocessor{
				toc.NewTOCPostprocessor(),
			},
		},
	})

	err = meshdoc.Run()
	log.Printf("err = %+v\n", err)
}
