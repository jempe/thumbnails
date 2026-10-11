package thumbnails

import (
	"encoding/json"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateWebp(t *testing.T) {
	dir := t.TempDir()
	images := filepath.Join(dir, "images")
	thumbs := filepath.Join(dir, "thumbs")
	os.Mkdir(images, 0755)
	os.Mkdir(thumbs, 0755)

	src := image.NewRGBA(image.Rect(0, 0, 400, 200))
	for x := 0; x < 400; x++ {
		for y := 0; y < 200; y++ {
			src.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	save := func(name string, encode func(f *os.File) error) {
		f, err := os.Create(filepath.Join(images, name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := encode(f); err != nil {
			t.Fatal(err)
		}
	}

	save("a.jpg", func(f *os.File) error { return jpeg.Encode(f, src, nil) })
	save("a.png", func(f *os.File) error { return png.Encode(f, src) })
	save("a.gif", func(f *os.File) error { return gif.Encode(f, src, nil) })
	save("a.webp", func(f *os.File) error {
		data, err := encodeWebp(src, 90)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		return err
	})

	conf, _ := json.Marshal(map[string]any{
		"images_folder": images,
		"thumbs_folder": thumbs,
		"webp_quality":  80,
		"sizes":         map[string]string{"fit": "<100x100", "exact": "=60x60", "big": "<1000x1000"},
	})
	conf_file := filepath.Join(dir, "thumbnails.json")
	os.WriteFile(conf_file, conf, 0644)

	if err := Config(conf_file); err != nil {
		t.Fatal(err)
	}

	if err := Generate("", true); err != nil {
		t.Fatal(err)
	}

	expected := map[string][2]int{"fit": {100, 50}, "exact": {60, 60}, "big": {400, 200}}

	for folder, size := range expected {
		for _, name := range []string{"a_j.webp", "a_p.webp", "a_g.webp", "a_w.webp"} {
			path := filepath.Join(thumbs, folder, name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("missing thumbnail %s: %s", path, err)
			}

			if mime := http.DetectContentType(data); mime != "image/webp" {
				t.Errorf("%s has type %s", path, mime)
			}

			// VP8 lossy images store the dimensions in the frame header
			width := int(data[26]) | int(data[27]&0x3f)<<8
			height := int(data[28]) | int(data[29]&0x3f)<<8
			if width != size[0] || height != size[1] {
				t.Errorf("%s is %dx%d, expected %dx%d", path, width, height, size[0], size[1])
			}
		}
	}
}
