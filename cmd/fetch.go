package main

import (
	"log"
	"os"

	"github.com/BernhardWebstudio/LyricFetcher/pkg/handler"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "LyricFetcher"
	app.Usage = "Add lyrics to your music library files"
	app.Action = func(c *cli.Context) error {
		handler.HandleFiles(os.Args[1:])
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
