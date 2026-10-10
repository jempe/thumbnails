package thumbnails

import (
	"encoding/json"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	ImagesFolder string            `json:"images_folder"`
	ThumbsFolder string            `json:"thumbs_folder"`
	WebpQuality  float32           `json:"webp_quality"`
	Sizes        map[string]string `json:"sizes"`
}

var configFile string
var folder_path string
var thumbs_path string
var thumb_sizes map[string]string
var webp_quality float32

// suffix added to the thumbnail name to avoid collisions between files with the same name and different type
var mime_suffixes = map[string]string{
	"image/gif":  "_g",
	"image/jpeg": "_j",
	"image/png":  "_p",
}

func Config(file string) (err error) {
	if Exists(file) {
		configFile = file
	} else {
		err = errors.New("Config file doesn't exists")
	}

	return err
}

// Generate thumbnails
// if image is empty generate all thumbnails
func Generate(image string, overwrite bool) (err error) {
	if configFile == "" {
		return errors.New("Config file not set")
	}

	config_content, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	conf := config{WebpQuality: 80}
	err = json.Unmarshal(config_content, &conf)
	if err != nil {
		return err
	}

	folder_path = strings.TrimRight(conf.ImagesFolder, "/")
	thumbs_path = strings.TrimRight(conf.ThumbsFolder, "/")
	thumb_sizes = conf.Sizes
	webp_quality = conf.WebpQuality

	if !Exists(folder_path) {
		err = errors.New("Images folder doesn't exists")
		return err
	} else if !IsDirectory(folder_path) {
		err = errors.New("Images folder is not a directory")
		return err
	}

	if image == "" {
		images, _ := ioutil.ReadDir(folder_path)
		for _, f := range images {
			if IsImage(folder_path + "/" + f.Name()) {
				err = generateThumbnail(f.Name(), overwrite)
				if err != nil {
					return err
				}
			}
		}
	} else {
		image_file_path := folder_path + "/" + image
		if Exists(image_file_path) {
			if IsImage(image_file_path) {
				err = generateThumbnail(image, overwrite)
				if err != nil {
					return err
				}
			} else {
				return errors.New("That file is not an image")
			}
		} else {
			return errors.New("That file doesn't exist")
		}
	}

	return err
}

// generate one thumbnail for each configured size, all of them in webp format
func generateThumbnail(image_file string, overwrite bool) error {
	if !IsDirectory(thumbs_path) {
		return errors.New("Thumbs folder doesn't exist or is not a Folder")
	}

	source_image := folder_path + "/" + image_file

	mime, err := getContentType(source_image)
	if err != nil {
		return err
	}

	thumb_name := strings.TrimSuffix(image_file, filepath.Ext(image_file)) + mime_suffixes[mime] + ".webp"

	var img image.Image

	for thumb_folder, thumb_size := range thumb_sizes {
		thumb_folder_path := thumbs_path + "/" + thumb_folder
		if !Exists(thumb_folder_path) {
			log.Println("Creating folder" + thumb_folder_path)
			err := os.Mkdir(thumb_folder_path, 0755)
			if err != nil {
				return err
			}
		}

		if !IsDirectory(thumb_folder_path) {
			return errors.New("Can't create thumbnails. " + thumb_folder_path + " must be a directory")
		}

		thumb_file_path := thumb_folder_path + "/" + thumb_name

		width, height, exact_size, err := parseSize(thumb_size)
		if err != nil {
			return err
		}

		if Exists(thumb_file_path) && !overwrite {
			log.Printf("Nothing to do, thumb %s already exists\n", thumb_file_path)
			continue
		}

		if img == nil {
			img, err = decodeImage(source_image, mime)
			if err != nil {
				return err
			}
		}

		var resized_image image.Image

		if exact_size {
			resized_image = coverResize(img, width, height)
		} else {
			resized_image = fitResize(img, width, height)
		}

		err = writeWebp(thumb_file_path, resized_image)
		if err != nil {
			return err
		}
	}

	return nil
}

// decode an image file of one of the supported types
func decodeImage(path string, mime string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	switch mime {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/gif":
		img, err = gif.Decode(file)
	case "image/png":
		img, err = png.Decode(file)
	default:
		err = errors.New("Unsupported image type: " + mime)
	}

	return img, err
}

// encode an image as webp and save it in path
func writeWebp(path string, img image.Image) error {
	data, err := encodeWebp(img, webp_quality)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func CopyImage(image string, thumb_file string) error {
	log.Println("copy " + image + " to " + thumb_file)
	reader, err := os.Open(image)
	if err != nil {
		return err
	}

	defer reader.Close()

	writer, err := os.Create(thumb_file)
	if err != nil {
		return err
	}

	defer writer.Close()

	_, err = io.Copy(writer, reader)

	return err
}

func parseSize(size_string string) (crop_width int, crop_height int, exact bool, err error) {
	if !(strings.HasPrefix(size_string, "=") || strings.HasPrefix(size_string, "<")) {
		err = errors.New("Invalid thumbnail size It should start with a < or = :" + size_string)
		return
	} else if !strings.Contains(size_string, "x") {
		err = errors.New("Invalid thumbnail size: " + size_string)
		return
	} else {
		width_height := strings.Split(strings.TrimLeft(size_string, "<="), "x")

		crop_width, err = strconv.Atoi(width_height[0])
		if err != nil || crop_width < 1 {
			err = errors.New("width must be a positive integer")
			return
		}

		crop_height, err = strconv.Atoi(width_height[1])
		if err != nil || crop_height < 1 {
			err = errors.New("height must be a positive integer")
			return
		}

		if strings.HasPrefix(size_string, "=") {
			exact = true
		} else {
			exact = false
		}
		return
	}
}

// check if file is an image
func IsImage(image_path string) bool {
	mime, err := getContentType(image_path)

	if err != nil {
		return false
	} else {
		if mime == "image/gif" || mime == "image/jpeg" || mime == "image/png" {
			return true
		} else {
			return false
		}
	}
}

// get content type of file
func getContentType(path string) (mime string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return mime, err
	}

	defer file.Close()

	file_buffer := make([]byte, 512)
	n, err := file.Read(file_buffer)
	if err != nil {
		return mime, err
	}

	mime = http.DetectContentType(file_buffer[:n])

	return mime, err
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err == nil && fileInfo.IsDir() {
		return true
	} else {
		return false
	}
}

// check if file exists
func Exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}
