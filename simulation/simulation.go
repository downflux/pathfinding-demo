package simulation

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"regexp"
	"strings"
	"time"

	"github.com/downflux/go-collider/agent"
	"github.com/downflux/go-collider/feature"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-geometry/2d/vector"

	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
	rlabel "github.com/downflux/pathfinding-demo/internal/render/label"
)

func sanitize(s string) string {
	return strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(s, "_"))
}

type O struct {
	Name        string
	Agents      []agent.O
	Projectiles []agent.O
	Features []feature.O
	Collider    collider.O

	MinX         float64
	MinY         float64
	MaxX         float64
	MaxY         float64
	TickDuration time.Duration
}

func (o *O) Filename() string { return sanitize(o.Name) }

func (o O) Marshal() []byte {
	data, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("cannot marshal JSON: %v", err))
	}
	return data
}

func Unmarshal(data []byte) O {
	var o O
	if err := json.Unmarshal(data, &o); err != nil {
		panic(fmt.Sprintf("cannot unmarshal JSON: %v", err))
	}
	return o
}

type S struct {
	agentRenderers      []*ragent.A
	projectileRenderers []*ragent.A
	collider            *collider.C

	minX         float64
	minY         float64
	maxX         float64
	maxY         float64
	tickDuration time.Duration
}

func New(o O) *S {
	s := &S{
		collider:     collider.New(o.Collider),
		tickDuration: o.TickDuration,
		minX:         o.MinX,
		minY:         o.MinY,
		maxX:         o.MaxX,
		maxY:         o.MaxY,
	}
	for _, opt := range o.Agents {
		s.agentRenderers = append(s.agentRenderers, ragent.New(s.collider.Insert(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Projectiles {
		s.agentRenderers = append(s.agentRenderers, ragent.New(s.collider.Insert(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Features {
		s.collider.InsertFeature(opt)
	}
	return s
}

func (s *S) Tick(d time.Duration) {
	s.collider.Tick(d)
}

func (s *S) Execute(nFrames int) *gif.GIF {
	frames := make([]*image.Paletted, 0, nFrames)
	delays := make([]int, 0, nFrames)
	for f := 0; f < nFrames; f++ {
		img := image.NewPaletted(
			image.Rectangle{image.Point{int(s.minX), int(s.minY)}, image.Point{int(s.maxX), int(s.maxY)}},
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
