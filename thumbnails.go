package thumbnails

import (
	"errors"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var configFile string
var folder_path string
var thumbs_path string
var thumb_sizes map[string]string

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

	folder_path = strings.TrimRight(viper.GetString("images_folder"), "/")
	thumbs_path = strings.TrimRight(viper.GetString("thumbs_folder"), "/")
	thumb_sizes = viper.GetStringMapString("sizes")

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
			if isImage(folder_path + "/" + f.Name()) {
				err = generateThumbnail(f.Name(), overwrite)
				if err != nil {
					return err
				}
			}
		}
	} else {
		image_file_path := folder_path + "/" + image
		if exists(image_file_path) {
			if isImage(image_file_path) {
				err = generateThumbnail(image_file_path, overwrite)
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

// generate one thumbnail
func generateThumbnail(image_file string, overwrite bool) error {
	var err error
	if isDirectory(thumbs_path) {
		source_image := folder_path + "/" + image_file

		mime, err := getContentType(source_image)
		if err != nil {
			return err
		}

		for thumb_folder, thumb_size := range thumb_sizes {
			thumb_folder_path := thumbs_path + "/" + thumb_folder
			if !exists(thumb_folder_path) {
				log.Println("Creating folder" + thumb_folder_path)
				err := os.Mkdir(thumb_folder_path, 0755)
				if err != nil {
					return err
				}
			}

			if isDirectory(thumb_folder_path) {
				thumb_file_path := thumb_folder_path + "/" + image_file

				width, height, exact_size, err := parseSize(thumb_size)
				if err != nil {
					return err
				}

				if exists(thumb_file_path) && overwrite == false {
					log.Printf("Nothing to do, thumb %s already exists\n", thumb_file_path)
				} else {
					var img image.Image

					file, err := os.Open(source_image)
					if err != nil {
						return err
					}

					if mime == "image/jpeg" {
						img, err = jpeg.Decode(file)
						if err != nil {
							return err
						}
					} else if mime == "image/gif" {
						img, err = gif.Decode(file)
						if err != nil {
							return err
						}
					} else if mime == "image/png" {
						img, err = png.Decode(file)
						if err != nil {
							return err
						}
					}

					file.Close()

					var resized_image image.Image

					if exact_size {
						// img_width, img_height, err := getImageDimension(source_image)
						// if err != nil {
						// 	return err
						// }
						resized_image = resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)
					} else {
						resized_image = resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)
					}

					out, err := os.Create(thumb_file_path)
					if err != nil {
						return err
					}

					defer out.Close()

					if mime == "image/jpeg" {
						jpeg.Encode(out, resized_image, nil)
					} else if mime == "image/gif" {
						var opt gif.Options
						opt.NumColors = 256

						gif.Encode(out, resized_image, &opt)
					} else if mime == "image/png" {
						png.Encode(out, resized_image)
					}
				}
			} else {
				return errors.New("Can't create thumbnails. " + thumb_folder_path + " must be a directory")
			}
		}
	} else {
		return errors.New("Thumbs folder doesn't exist or is not a Folder")
	}

	return err
}

func getImageDimensions(imagePath string) (width int, height int, err error) {
	file, file_err := os.Open(imagePath)
	if file_err != nil {
		err = file_err
		return
	}

	defer file.Close()

	image, _, image_err := image.DecodeConfig(file)
	if image_err != nil {
		err = image_err
		return
	}

	width = image.Width
	height = image.Height

	return
}

func copyImage(image string, thumb_file string) error {
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
		if err != nil {
			err = errors.New("width must be an integer")
			return
		}

		crop_height, err = strconv.Atoi(width_height[1])
		if err != nil {
			err = errors.New("height must be an integer")
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
func isImage(image_path string) bool {
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

// check if file exists
func exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}
