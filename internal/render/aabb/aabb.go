package aabb

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type R struct {
	aabb  hyperrectangle.R
	color color.Color
}

func New(r hyperrectangle.R, color color.Color) *R {
	if k := r.Min().Dimension(); k != 2 {
		panic(fmt.Sprintf("invalid hyperrectangle dimension %v", k))
	}
	return &R{
		aabb:  r,
		color: color,
	}
}

func (r *R) Draw(img *image.Paletted) {
	xmin, xmax := int(math.Round(r.aabb.Min().X(vector.AXIS_X))), int(math.Round(r.aabb.Max().X(vector.AXIS_X)))
	ymin, ymax := int(math.Round(r.aabb.Min().X(vector.AXIS_Y))), int(math.Round(r.aabb.Max().X(vector.AXIS_Y)))

	for i := xmin; i <= xmax; i++ {
		img.Set(i, ymin, r.color)
		img.Set(i, ymax, r.color)
	}

	for i := ymin; i <= ymax; i++ {
		img.Set(xmin, i, r.color)
		img.Set(xmax, i, r.color)
	}
}
