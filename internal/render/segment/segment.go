package segment

import (
	"image"
	"image/color"

	"github.com/downflux/go-geometry/2d/segment"
)

const (
	resolution = 500
)

type S struct {
	segment segment.S
	color   color.Color
}

func New(s segment.S, c color.Color) *S {
	return &S{
		segment: s,
		color:   c,
	}
}

func (s *S) Draw(img *image.Paletted) {
	grain := (s.segment.TMax() - s.segment.TMin()) / resolution
	for t := s.segment.TMin(); t <= s.segment.TMax(); t += grain {
		img.Set(int(s.segment.L().L(t).X()), int(s.segment.L().L(t).Y()), s.color)
	}
}
