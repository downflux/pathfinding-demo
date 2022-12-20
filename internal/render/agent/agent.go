package agent

import (
	"image"
	"image/color"

	"github.com/downflux/game-db/agent"
	"github.com/downflux/go-geometry/2d/hypersphere"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/pathfinding-demo/internal/render/circle"
	"github.com/downflux/pathfinding-demo/internal/render/segment"
	"github.com/downflux/pathfinding-demo/internal/render/trail"

	l2d "github.com/downflux/go-geometry/2d/line"
	s2d "github.com/downflux/go-geometry/2d/segment"
)

var (
	ColorVelocity = color.RGBA{0, 0, 255, 255}
	ColorHeading  = color.RGBA{255, 0, 0, 255}
	ColorAgent    = color.Black
	ColorTrail    = color.RGBA{192, 192, 192, 255}
)

type R struct {
	agent *agent.A

	trail *trail.T
}

func New(a *agent.A) *R {
	return &R{
		agent: a,
		trail: trail.New(ColorTrail),
	}
}

func (r *R) Draw(img *image.Paletted) {
	r.trail.Push(r.agent.Position())
	r.trail.Draw(img)

	circle.New(
		*hypersphere.New(
			r.agent.Position(),
			r.agent.Radius(),
		),
		ColorAgent,
	).Draw(img)

	segment.New(
		*s2d.New(*l2d.New(r.agent.Position(), polar.Cartesian(r.agent.Heading())), 0, 2*r.agent.Radius()),
		ColorHeading,
	).Draw(img)
}
