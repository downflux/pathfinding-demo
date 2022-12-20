package circle

import (
	"image"
	"image/color"

	"github.com/downflux/go-geometry/2d/hypersphere"
)

type C struct {
	circle hypersphere.C
	color  color.Color
}

func New(c hypersphere.C, color color.Color) *C {
	return &C{
		circle: c,
		color:  color,
	}
}

func (c *C) Draw(img *image.Paletted) {
	r := int(c.circle.R())
	vx, vy := int(c.circle.P().X()), int(c.circle.P().Y())

	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(vx+x, vy+y, c.color)
		img.Set(vx+y, vy+x, c.color)
		img.Set(vx-y, vy+x, c.color)
		img.Set(vx-x, vy+y, c.color)
		img.Set(vx-x, vy-y, c.color)
		img.Set(vx-y, vy-x, c.color)
		img.Set(vx+y, vy-x, c.color)
		img.Set(vx+x, vy-y, c.color)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
