package label

import (
	"image"
	"image/color"

	"github.com/downflux/go-geometry/2d/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	ColorText = color.Black

	// Face7x13 is 13 pixels high, and the font drawer uses the input Y
	// coordinate as the bottom most pixel.
	offset = vector.V{0, 13}
)

type L struct {
	label    string
	position vector.V
}

func New(label string, p vector.V) *L {
	return &L{
		label:    label,
		position: p,
	}
}

func (l *L) Draw(img *image.Paletted) {
	p := vector.Add(offset, l.position)
	x, y := int(p.X()), int(p.Y())

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(ColorText),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{fixed.I(x), fixed.I(y)},
	}
	d.DrawString(l.label)
}
