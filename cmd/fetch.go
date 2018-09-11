package main

import (
	"fmt"
  "log"
  "./internal/app" as "internal"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "LyricFetcher"
	app.Usage = "Add lyrics to your music library files"
	app.Action = func(c *cli.Context) error {
		internal.HandleFiles(os.Args[1:])
		fmt.Println("boom! I say!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
