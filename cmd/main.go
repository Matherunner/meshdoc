package main

import (
	"flag"
	"log"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/counter"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/xref"
)

var (
	configFlag = flag.String("config", "", "path to config file")
)

func main() {
	flag.Parse()

	if *configFlag == "" {
		log.Fatalf("config must be set")
	}

	config, err := meshdoc.LoadFromFile(*configFlag)
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	counterHierarchy := counter.NewHierarchy()
	counterHierarchy.Add(counter.FileKey, counter.RootKey)
	counterHierarchy.Add("H1", counter.FileKey)
	counterHierarchy.Add("H2", "H1")
	counterHierarchy.Add("H3", "H2")
	counterHierarchy.Add("H4", "H3")
	counterHierarchy.Add("H5", "H4")
	counterOptions := &counter.Options{
		Hierarchy: counterHierarchy,
	}

	meshdoc := meshdoc.NewMeshdoc(&meshdoc.MeshdocOptions{
		Config: config,
		Components: meshdoc.ComponentOptions{
			BookReader:   htmldoc.NewDefaultBookReader(),
			BookWriter:   htmldoc.NewDefaultBookWriter(),
			ParsedReader: htmldoc.NewDefaultParsedReader,
			// ParsedWriter:  htmldoc.NewDefaultParsedWriter, // TODO: not used?
			Preprocessors: []meshdoc.Preprocessor{},
			Postprocessors: []meshdoc.Postprocessor{
				toc.NewTOC(),
				counter.NewCounter(counterOptions),
				xref.NewXRef(),
			},
		},
	})

	err = meshdoc.Run()
	log.Printf("err = %+v\n", err)
}
