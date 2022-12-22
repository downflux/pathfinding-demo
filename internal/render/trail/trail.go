package trail

import (
	"image"
	"image/color"

	"github.com/downflux/go-geometry/2d/vector"
)

const (
	trailbufLen = 50
)

type T struct {
	frame    int
	trailbuf []vector.V

	color color.Color
}

func New(c color.Color) *T {
	return &T{
		trailbuf: make([]vector.V, 0, trailbufLen),
		color:    c,
	}
}

func (t *T) Push(v vector.V) {
	if t.frame < trailbufLen {
		u := vector.M{0, 0}
		u.Copy(v)
		t.trailbuf = append(t.trailbuf, u.V())
	} else {
		t.trailbuf[t.frame%trailbufLen].M().Copy(v)
	}
	t.frame += 1
}

func (t *T) Draw(img *image.Paletted) {
	for _, p := range t.trailbuf {
		img.Set(int(p.X()), int(p.Y()), t.color)
	}
}
