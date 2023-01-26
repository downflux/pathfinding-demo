package aabb

import (
	"image"
	"image/color"
	"math"

	"github.com/downflux/go-geometry/2d/hyperrectangle"
)

type R struct {
	aabb  hyperrectangle.R
	color color.Color
}

func New(r hyperrectangle.R, color color.Color) *R {
	return &R{
		aabb:  r,
		color: color,
	}
}

func (r *R) Draw(img *image.Paletted) {
	xmin, xmax := int(math.Round(r.aabb.Min().X())), int(math.Round(r.aabb.Max().X()))
	ymin, ymax := int(math.Round(r.aabb.Min().Y())), int(math.Round(r.aabb.Max().Y()))

	for i := xmin; i <= xmax; i++ {
		img.Set(i, ymin, r.color)
		img.Set(i, ymax, r.color)
	}

	for i := ymin; i <= ymax; i++ {
		img.Set(xmin, i, r.color)
		img.Set(xmax, i, r.color)
	}
}
