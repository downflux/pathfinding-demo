package projectile

import (
	"fmt"
	"image"
	"image/color"

	"github.com/downflux/go-database/projectile"
	"github.com/downflux/go-geometry/2d/hypersphere"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/pathfinding-demo/internal/render/circle"
	"github.com/downflux/pathfinding-demo/internal/render/label"
	"github.com/downflux/pathfinding-demo/internal/render/segment"
	"github.com/downflux/pathfinding-demo/internal/render/trail"

	l2d "github.com/downflux/go-geometry/2d/line"
	s2d "github.com/downflux/go-geometry/2d/segment"
)

var (
	ColorVelocity   = color.RGBA{0, 0, 255, 255}
	ColorHeading    = color.RGBA{255, 0, 0, 255}
	ColorProjectile = color.Black
	ColorTrail      = color.RGBA{192, 192, 192, 255}

	fontOffset = vector.V{-6, -8}
)

type P struct {
	projectile projectile.RO

	trail *trail.T

	label bool
}

func New(a projectile.RO, label bool) *P {
	return &P{
		projectile: a,
		trail:      trail.New(ColorTrail),
		label:      label,
	}
}

func (r *P) Draw(img *image.Paletted) {
	r.trail.Push(r.projectile.Position())
	r.trail.Draw(img)

	circle.New(
		*hypersphere.New(
			r.projectile.Position(),
			r.projectile.Radius(),
		),
		ColorProjectile,
	).Draw(img)

	segment.New(
		*s2d.New(*l2d.New(
			r.projectile.Position(),
			polar.Cartesian(r.projectile.Heading()),
		), 0, 2*r.projectile.Radius()),
		ColorHeading,
	).Draw(img)

	if r.label {
		label.New(
			// Use hexadecimals here since the projectiles are small and we need
			// a clear designator.
			fmt.Sprintf("%02X", r.projectile.ID()),
			vector.Add(fontOffset, r.projectile.Position()),
		).Draw(img)
	}
}
