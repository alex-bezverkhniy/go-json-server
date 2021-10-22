package main

import (
	"fmt"
	"os"
)

type Config struct {
	FileName string
	Port     string
}

var config = Config{}

func main() {
	// Get path to .json file
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Printf("Please provide path to .json file as a first argument\n\n")
		fmt.Printf("e.g: go-json-server mydb.json\n\n")
		os.Exit(-1)
	}

	config.FileName = argsWithoutProg[0]

	jsonServer, err := NewJsonServer(config)
	if err != nil {
		os.Exit(-1)
	}

	jsonServer.Start()
}
