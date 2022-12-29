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
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-collider/feature"
	"github.com/downflux/go-geometry/2d/vector"

	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
	rfeature "github.com/downflux/pathfinding-demo/internal/render/feature"
	rlabel "github.com/downflux/pathfinding-demo/internal/render/label"
)

func sanitize(s string) string {
	return strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(s, "_"))
}

type O struct {
	Name        string
	Agents      []agent.O
	Projectiles []agent.O
	Features    []feature.O
	Collider    collider.O

	MinX         float64
	MinY         float64
	MaxX         float64
	MaxY         float64
	TickDuration time.Duration

	NFrames int
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
	nFrames int

	agentRenderers      []*ragent.A
	projectileRenderers []*ragent.A
	featureRenderers    []*rfeature.F
	collider            *collider.C

	minX         float64
	minY         float64
	maxX         float64
	maxY         float64
	tickDuration time.Duration

	tickTimer time.Duration

	tileSize float64
}

func New(o O) *S {
	s := &S{
		nFrames: o.NFrames,

		collider:     collider.New(o.Collider),
		tickDuration: o.TickDuration,
		minX:         o.MinX,
		minY:         o.MinY,
		maxX:         o.MaxX,
		maxY:         o.MaxY,

		tileSize: o.Collider.DigitizerTileSize,
	}
	for _, opt := range o.Agents {
		s.agentRenderers = append(s.agentRenderers, ragent.New(s.collider.Insert(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Projectiles {
		s.agentRenderers = append(s.agentRenderers, ragent.New(s.collider.Insert(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Features {
		s.featureRenderers = append(s.featureRenderers, rfeature.New(s.collider.InsertFeature(opt)))
	}
	return s
}

func (s *S) Tick(d time.Duration) {
	s.collider.Tick(d)
}

func (s *S) TickTimer() time.Duration { return s.tickTimer }

func (s *S) Execute() *gif.GIF {
	frames := make([]*image.Paletted, 0, s.nFrames)
	delays := make([]int, 0, s.nFrames)
	for f := 0; f < s.nFrames; f++ {
		img := image.NewPaletted(
			image.Rectangle{image.Point{int(s.minX), int(s.minY)}, image.Point{int(s.maxX), int(s.maxY)}},
			[]color.Color{
				color.White,
				rlabel.ColorText,
				ragent.ColorVelocity,
				ragent.ColorTrail,
				ragent.ColorAgent,
				ragent.ColorHeading,
				rfeature.ColorBox,
			},
		)

		rlabel.New(fmt.Sprintf("frame %v / %v", f, s.nFrames), vector.V{0, 0}).Draw(img)

		for _, a := range s.agentRenderers {
			a.Draw(img)
		}
		for _, p := range s.projectileRenderers {
			p.Draw(img)
		}
		for _, f := range s.featureRenderers {
			f.Draw(img)
		}
		frames = append(frames, img)
		delays = append(delays, 2) // 50fps

		start := time.Now()
		s.Tick(s.tickDuration)
		s.tickTimer = s.tickTimer + time.Now().Sub(start)
	}

	s.tickTimer = time.Duration(float64(s.tickTimer) / float64(s.nFrames))

	return &gif.GIF{
		Delay: delays,
		Image: frames,
	}
}
