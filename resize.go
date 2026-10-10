package thumbnails

import (
	"image"
	"math"

	xdraw "golang.org/x/image/draw"
)

var lanczos3 = &xdraw.Kernel{
	Support: 3,
	At: func(t float64) float64 {
		if t == 0 {
			return 1
		}
		if t <= -3 || t >= 3 {
			return 0
		}
		pt := math.Pi * t
		return 3 * math.Sin(pt) * math.Sin(pt/3) / (pt * pt)
	},
}

// scale resamples the src rectangle of img into a new width x height image
func scale(img image.Image, src image.Rectangle, width int, height int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	lanczos3.Scale(dst, dst.Bounds(), img, src, xdraw.Src, nil)

	return dst
}

// fitResize downscales img to fit inside max_width x max_height keeping the aspect ratio.
// The original image is returned if it already fits.
func fitResize(img image.Image, max_width int, max_height int) image.Image {
	bounds := img.Bounds()
	img_width := bounds.Dx()
	img_height := bounds.Dy()

	if img_width <= max_width && img_height <= max_height {
		return img
	}

	ratio := math.Min(float64(max_width)/float64(img_width), float64(max_height)/float64(img_height))

	width := max(1, int(math.Round(float64(img_width)*ratio)))
	height := max(1, int(math.Round(float64(img_height)*ratio)))

	return scale(img, bounds, width, height)
}

// coverResize returns an image of exactly width x height pixels.
// The centered area of img with the thumbnail aspect ratio is cropped and scaled to that size.
func coverResize(img image.Image, width int, height int) image.Image {
	bounds := img.Bounds()
	img_width := bounds.Dx()
	img_height := bounds.Dy()

	crop_width := img_width
	crop_height := img_height

	if img_width*height > img_height*width {
		crop_width = max(1, int(math.Round(float64(img_height)*float64(width)/float64(height))))
	} else {
		crop_height = max(1, int(math.Round(float64(img_width)*float64(height)/float64(width))))
	}

	x := bounds.Min.X + (img_width-crop_width)/2
	y := bounds.Min.Y + (img_height-crop_height)/2

	return scale(img, image.Rect(x, y, x+crop_width, y+crop_height), width, height)
}
