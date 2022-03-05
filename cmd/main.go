package main

import (
	"log"

	"github.com/Matherunner/meshdoc"
)

func main() {
	meshdoc := meshdoc.NewMeshdoc(&meshdoc.MeshdocOptions{
		ConfigPath: "./examples/simple/config.toml",
	})
	err := meshdoc.Run()
	log.Printf("err = %+v\n", err)
}
