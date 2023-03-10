package agent

import (
	"fmt"
	"image"
	"image/color"

	"github.com/downflux/go-database/agent"
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
	ColorVelocity       = color.RGBA{0, 0, 255, 255}
	ColorHeading        = color.RGBA{255, 0, 0, 255}
	ColorAgent          = color.Black
	ColorTargetPosition = color.RGBA{0, 128, 0, 255}
	ColorTrail          = color.RGBA{192, 192, 192, 255}

	fontOffset = vector.V{-6, -8}
)

type A struct {
	agent agent.RO

	trail *trail.T

	target bool
	label  bool
}

func New(a agent.RO, label bool, target bool) *A {
	return &A{
		agent:  a,
		trail:  trail.New(ColorTrail),
		label:  label,
		target: target,
	}
}

func (r *A) Draw(img *image.Paletted) {
	r.trail.Push(r.agent.Position())
	r.trail.Draw(img)

	if r.target {
		circle.New(
			*hypersphere.New(
				r.agent.TargetPosition(),
				2,
			),
			ColorTargetPosition,
		).Draw(img)
	}

	circle.New(
		*hypersphere.New(
			r.agent.Position(),
			r.agent.Radius(),
		),
		ColorAgent,
	).Draw(img)

	segment.New(
		*s2d.New(*l2d.New(
			r.agent.Position(),
			polar.Cartesian(r.agent.Heading()),
		), 0, 2*r.agent.Radius()),
		ColorHeading,
	).Draw(img)

	if r.label {
		label.New(
			// Use hexadecimals here since the agents are small and we need
			// a clear designator.
			fmt.Sprintf("%02X", r.agent.ID()),
			vector.Add(fontOffset, r.agent.Position()),
		).Draw(img)
	}
}
