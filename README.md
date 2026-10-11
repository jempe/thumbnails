# Thumbnails

Thumbnails is a Go-based tool for generating thumbnails from images using a specified configuration file. This README provides instructions on how to install, configure, and use the tool.

## Installation

To install the Thumbnails tool, use the following command:

```sh
go get github.com/jempe/thumbnails
```

## Usage

The Thumbnails tool requires a configuration file and an image file to generate thumbnails. You can specify these files using command-line flags.

### Command-Line Flags

- `-config`:  Path to the configuration file (default: `thumbnails.json`).
- `-image`:  Path to the image file.

### Example

To generate thumbnails using a configuration file and an image file, run the following command:

```sh
go run cmd/gothumbnails/main.go -config=path/to/thumbnails.json -image=path/to/image.jpg
```

## Requirements

All thumbnails are WebP files, encoded with libwebp through cgo. You need a C compiler and the libwebp headers and library:

```sh
brew install webp              # macOS (Homebrew)
sudo port install webp         # macOS (MacPorts)
apt install libwebp-dev        # Debian/Ubuntu
```

The cgo flags in `webp.go` already look in the Homebrew locations (`/opt/homebrew`, `/usr/local`). MacPorts installs into `/opt/local`, so tell cgo where to find it when you build, test or run:

```sh
export CGO_CFLAGS="-I/opt/local/include"
export CGO_LDFLAGS="-L/opt/local/lib"
go build ./...
```

## Configuration

The configuration file is a JSON file:

```json
{
  "images_folder": "images",
  "thumbs_folder": "images/thumbs",
  "webp_quality": 80,
  "sizes": {
    "300w": "<300x300",
    "square": "=200x200"
  }
}
```

- `images_folder`: folder with the original JPEG, PNG, GIF and WebP images.
- `thumbs_folder`: folder where one subfolder per size is created.
- `webp_quality`: WebP quality from 0 to 100 (default 80).
- `sizes`: map of folder name to size. `<WxH` fits the image inside WxH without upscaling, `=WxH` crops the center of the image and scales it to exactly WxH.

## Output names

Every thumbnail is a `.webp` file. A suffix based on the type of the original image avoids collisions between files with the same name:

| Original | Thumbnail |
|---|---|
| `photo.gif` | `photo_g.webp` |
| `photo.jpg` | `photo_j.webp` |
| `photo.png` | `photo_p.webp` |
| `photo.webp` | `photo_w.webp` |

Animated WebP files can't be decoded and produce an error.

## Error Handling

If there are any errors during the configuration or thumbnail generation process, they will be printed to the console. Make sure to check the output for any error messages.

## Contributing

If you would like to contribute to the project, please fork the repository and submit a pull request. We appreciate your contributions!

## License

This project is licensed under the Apache 2.0 License. See the LICENSE file for details.
