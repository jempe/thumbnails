package thumbnails

/*
#cgo LDFLAGS: -lwebp
#cgo darwin,arm64 CFLAGS: -I/opt/homebrew/include
#cgo darwin,arm64 LDFLAGS: -L/opt/homebrew/lib
#cgo darwin,amd64 CFLAGS: -I/usr/local/include
#cgo darwin,amd64 LDFLAGS: -L/usr/local/lib
#include <stdint.h>
#include <webp/encode.h>
*/
import "C"

import (
	"errors"
	"image"
	"image/draw"
	"unsafe"
)

// encodeWebp encodes img as a lossy webp image using libwebp, quality goes from 0 to 100
func encodeWebp(img image.Image, quality float32) ([]byte, error) {
	rgba, ok := img.(*image.NRGBA)
	if !ok || rgba.Rect.Min != (image.Point{}) {
		rgba = image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	}

	width := rgba.Rect.Dx()
	height := rgba.Rect.Dy()

	if width < 1 || height < 1 {
		return nil, errors.New("webp: empty image")
	}

	var output *C.uint8_t

	size := C.WebPEncodeRGBA((*C.uint8_t)(unsafe.Pointer(&rgba.Pix[0])), C.int(width), C.int(height), C.int(rgba.Stride), C.float(quality), &output)
	if size == 0 || output == nil {
		return nil, errors.New("webp: encoding failed")
	}

	defer C.WebPFree(unsafe.Pointer(output))

	return C.GoBytes(unsafe.Pointer(output), C.int(size)), nil
}
