package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"

	"github.com/fogleman/gg"
)

func DownloadTemplate(url string) image.Image {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("template from %s failed because of %v", url, err)
	}
	defer res.Body.Close()
	image, _, err := image.Decode(res.Body)
	if err != nil {
		log.Fatalf("Could not decode %s because of %v", url, err)
	}
	return image
}

func addText(url string, text []string) string {
	length := len(text[1])
	if len(text[0]) > len(text[1]) {
		length = len(text[0])
	}
	path := "./meme.jpg"
	img := DownloadTemplate(url)
	r := img.Bounds()
	w := r.Dx()
	h := r.Dy()

	fontSize := w / length

	if fontSize > h/8 {
		fontSize = h / 8
	}

	fmt.Println(fontSize)

	m := gg.NewContext(w, h)
	m.DrawImage(img, 0, 0)
	m.LoadFontFace("impact.ttf", float64(fontSize))

	// Apply black stroke
	m.SetHexColor("#000")
	strokeSize := 6
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			// give it rounded corners
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := float64(w/2 + dx)
			y := float64(h - fontSize + dy)
			m.DrawStringAnchored(text[1], x, y, 0.5, 0.5)
		}
	}

	// Apply white fill
	m.SetHexColor("#FFF")
	m.DrawStringAnchored(text[1], float64(w)/2, float64(h-fontSize), 0.5, 0.5)

	// Apply black stroke
	m.SetHexColor("#000")
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			// give it rounded corners
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := float64(w/2 + dx)
			y := float64(fontSize - dy)
			m.DrawStringAnchored(text[0], x, y, 0.5, 0.5)
		}
	}

	// Apply white fill
	m.SetHexColor("#FFF")
	m.DrawStringAnchored(text[0], float64(w)/2, float64(fontSize), 0.5, 0.5)

	m.SavePNG(path)

	fmt.Printf("Saved to %s\n", path)
	return path
}
