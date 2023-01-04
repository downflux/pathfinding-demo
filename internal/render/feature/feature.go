package feature

import (
	"image"
	"image/color"

	"github.com/downflux/go-database/feature"
	"github.com/downflux/pathfinding-demo/internal/render/aabb"
)

var (
	ColorBox = color.Black
)

type F struct {
	feature feature.RO
}

func New(f feature.RO) *F {
	return &F{
		feature: f,
	}
}

func (f *F) Draw(img *image.Paletted) {
	aabb.New(f.feature.AABB(), ColorBox).Draw(img)
}
