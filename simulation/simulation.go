package simulation

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"time"

	"github.com/downflux/go-collider/agent"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	vnd "github.com/downflux/go-geometry/nd/vector"
	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
	rlabel "github.com/downflux/pathfinding-demo/internal/render/label"
)

type O struct {
	Agents      []agent.O
	Projectiles []agent.O
	Collider    collider.O

	Dimensions   hyperrectangle.R
	TickDuration time.Duration
}

type S struct {
	agentRenderers      []*ragent.A
	projectileRenderers []*ragent.A
	collider            *collider.C

	dimensions   hyperrectangle.R
	tickDuration time.Duration
}

func New(o O) *S {
	s := &S{
		collider:     collider.New(o.Collider),
		tickDuration: o.TickDuration,
		dimensions:   o.Dimensions,
	}
	for _, opt := range o.Agents {
		s.agentRenderers = append(s.agentRenderers, ragent.New(s.collider.Insert(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Projectiles {
		s.agentRenderers = append(s.agentRenderers, ragent.New(s.collider.Insert(opt), opt.Radius >= 10))
	}
	return s
}

func (s *S) Tick(d time.Duration) {
	s.collider.Tick(d)
}

func (s *S) Execute(nFrames int) *gif.GIF {
	minX, minY := s.dimensions.Min().X(vnd.AXIS_X), s.dimensions.Min().X(vnd.AXIS_Y)
	maxX, maxY := s.dimensions.Max().X(vnd.AXIS_X), s.dimensions.Max().X(vnd.AXIS_Y)
	frames := make([]*image.Paletted, 0, nFrames)
	delays := make([]int, 0, nFrames)
	for f := 0; f < nFrames; f++ {
		img := image.NewPaletted(
			image.Rectangle{image.Point{int(minX), int(minY)}, image.Point{int(maxX), int(maxY)}},
			[]color.Color{
				color.White,
				rlabel.ColorText,
				ragent.ColorVelocity,
				ragent.ColorTrail,
				ragent.ColorAgent,
				ragent.ColorHeading,
			},
		)
		rlabel.New(fmt.Sprintf("frame %v / %v", f, nFrames), vector.V{0, 0}).Draw(img)

		for _, a := range s.agentRenderers {
			a.Draw(img)
		}
		frames = append(frames, img)
		delays = append(delays, 2) // 50fps

		s.Tick(s.tickDuration)
	}

	return &gif.GIF{
		Delay: delays,
		Image: frames,
	}
}
