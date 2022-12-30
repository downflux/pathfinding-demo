package grid

import (
	"math"

	"image"
	"image/color"
)

var (
	ColorGrid = color.RGBA{200, 200, 200, 255}
)

type G struct {
	x     float64
	y     float64
	color color.Color
}

func New(x, y float64) *G {
	return &G{
		x:     x,
		y:     y,
		color: ColorGrid,
	}
}

func (g *G) Draw(img *image.Paletted) {
	for i := math.Ceil(float64(img.Bounds().Min.X) / g.x); i < math.Ceil(float64(img.Bounds().Max.X)/g.x); i++ {
		for j := img.Bounds().Min.Y; j <= img.Bounds().Max.Y; j++ {
			img.Set(int(i*g.x), j, g.color)
		}
	}

	for j := math.Ceil(float64(img.Bounds().Min.Y) / g.y); j < math.Ceil(float64(img.Bounds().Max.Y)/g.y); j++ {
		for i := img.Bounds().Min.X; i <= img.Bounds().Max.X; i++ {
			img.Set(i, int(j*g.y), g.color)
		}
	}
}
