package label

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	ColorText = color.Black
)

type L struct {
	label string
}

func New(label string) *L {
	return &L{
		label: label,
	}
}

func (l *L) Draw(img *image.Paletted) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(ColorText),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{fixed.I(0), fixed.I(13)},
	}
	d.DrawString(l.label)
}
