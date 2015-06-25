package thumbnails

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
)

var configFile string

func Config(file string) (err error) {
	if exists(file) {
		viper.SetConfigFile(file)
	} else {
		err = errors.New("Config file doesn't exists")
	}

	return err
}

// Generate thumbnails
// if image is empty generate all thumbnails
func Generate(image string, overwrite bool) (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	folder_path := viper.GetString("images_folder")

	if !exists(folder_path) {
		err = errors.New("Images folder doesn't exists")
		return err
	} else if !isDirectory(folder_path) {
		err = errors.New("Images folder is not a directory")
		return err
	}

	if image == "" {
		images, _ := ioutil.ReadDir(folder_path)
		for _, f := range images {
			mime, err := getContentType(folder_path + "/" + f.Name())

			if err != nil {
				return err
			}

			fmt.Println(f.Name() + ":" + mime)
		}
	}

	fmt.Println(viper.GetString("thumbs_folder"))

	sizes := viper.GetStringMapString("sizes")

	for name, size := range sizes {
		fmt.Println(name + ":" + size)
	}

	return err
}

func getContentType(path string) (mime string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return mime, err
	}

	file_buffer := make([]byte, 512)
	_, err = file.Read(file_buffer)
	if err != nil {
		return mime, err
	}

	mime = http.DetectContentType(file_buffer)

	return mime, err
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err == nil && fileInfo.IsDir() {
		return true
	} else {
		return false
	}
}

func exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}
