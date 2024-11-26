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

## Configuration

The configuration file should be a JSON file that specifies the settings for generating thumbnails. Here is an example of what the configuration file might look like:

```json
{
    "thumbnail_sizes": [
        {"width": 100, "height": 100},
        {"width": 200, "height": 200}
    ],
    "output_directory": "thumbnails"
}
```

## Error Handling

If there are any errors during the configuration or thumbnail generation process, they will be printed to the console. Make sure to check the output for any error messages.

## Contributing

If you would like to contribute to the project, please fork the repository and submit a pull request. We appreciate your contributions!

## License

This project is licensed under the Apache 2.0 License. See the LICENSE file for details.
