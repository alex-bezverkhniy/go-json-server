package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type Config struct {
	FileName string
	Port     string
}

const DefalutFileName = "db.json"

func main() {
	app := &cli.App{
		Name:  "go-json-server",
		Usage: "simple full REST API server with zero coding",
		Action: func(c *cli.Context) error {
			if c.NArg() > 0 {
				return cli.ShowAppHelp(c)
			}

			fileName := "db.json"
			if c.NArg() > 0 {
				fileName = c.Args().Get(0)
			}

			var config = Config{
				FileName: fileName,
			}

			log.Println(config)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// // Get path to .json file
	// argsWithoutProg := os.Args[1:]
	// if len(argsWithoutProg) == 0 {
	// 	fmt.Printf("Please provide path to .json file as a first argument\n\n")
	// 	fmt.Printf("e.g: go-json-server mydb.json\n\n")
	// 	os.Exit(-1)
	// }

	// config.FileName = argsWithoutProg[0]

	// jsonServer, err := NewJsonServer(config)
	// if err != nil {
	// 	os.Exit(-1)
	// }

	// jsonServer.Start()
}
