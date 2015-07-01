package main

import (
	"flag"
	"fmt"
	"github.com/jempe/thumbnails"
)

var config_file = flag.String("config", "thumbnails.json", "config file location")
var image_file = flag.String("image", "", "image file name")

func main() {
	flag.Parse()
	config := *config_file
	image := *image_file

	err := thumbnails.Config(config)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		err = thumbnails.Generate(image, true)

		if err != nil {
			fmt.Printf("Generate thumbnails error: %s\n", err)
		}
	}
}
