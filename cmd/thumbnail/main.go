package main

import (
	"flag"
	"fmt"
	"github.com/jempe/thumbnails"
)

var config_file = flag.String("config", "/home/kastro/Documents/thumbnails/thumbnails.json", "config file location")

func main() {
	flag.Parse()
	config := *config_file

	err := thumbnails.Config(config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = thumbnails.Generate("", true)

	if err != nil {
		panic(fmt.Errorf("Error generating thumbnails: %s \n", err))
	}
}
