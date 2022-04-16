package main

import (
	"log"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
)

func main() {
	config, err := meshdoc.LoadFromFile("./examples/simple/config.toml")
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	counterHierarchy := counter.NewHierarchy()
	counterHierarchy.Add(counter.FileKey, counter.RootKey)
	counterHierarchy.Add("H1", counter.FileKey)
	counterHierarchy.Add("H2", "H1")
	counterOptions := &counter.Options{
		Hierarchy: counterHierarchy,
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
				counter.NewCounter(counterOptions),
			},
		},
	})

	err = meshdoc.Run()
	log.Printf("err = %+v\n", err)
}
