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

	"github.com/downflux/go-boids/x/boids"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/database"
	"github.com/downflux/go-database/feature"
	"github.com/downflux/go-database/projectile"
	"github.com/downflux/go-geometry/2d/vector"

	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
	rfeature "github.com/downflux/pathfinding-demo/internal/render/feature"
	rlabel "github.com/downflux/pathfinding-demo/internal/render/label"
	rprojectile "github.com/downflux/pathfinding-demo/internal/render/projectile"
)

func sanitize(s string) string {
	return strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(s, "_"))
}

type O struct {
	Name        string
	Agents      []agent.O
	Projectiles []projectile.O
	Features    []feature.O
	Collider    collider.O
	Boids       boids.O

	EnableBoids bool

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
	projectileRenderers []*rprojectile.P
	featureRenderers    []*rfeature.F
	db                  *database.DB
	collider            *collider.C
	boids               *boids.B

	enableBoids bool

	minX         float64
	minY         float64
	maxX         float64
	maxY         float64
	tickDuration time.Duration

	tickTimer time.Duration
}

func New(o O) *S {
	db := database.New(database.DefaultO)
	var b *boids.B
	if o.EnableBoids {
		b = boids.New(db, o.Boids)
	}
	s := &S{
		nFrames: o.NFrames,

		db:           db,
		collider:     collider.New(db, o.Collider),
		boids:        b,
		enableBoids:  o.EnableBoids,
		tickDuration: o.TickDuration,
		minX:         o.MinX,
		minY:         o.MinY,
		maxX:         o.MaxX,
		maxY:         o.MaxY,
	}
	for _, opt := range o.Agents {
		s.agentRenderers = append(s.agentRenderers, ragent.New(db.InsertAgent(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Projectiles {
		s.projectileRenderers = append(s.projectileRenderers, rprojectile.New(db.InsertProjectile(opt), opt.Radius >= 10))
	}
	for _, opt := range o.Features {
		s.featureRenderers = append(s.featureRenderers, rfeature.New(db.InsertFeature(opt)))
	}
	return s
}

func (s *S) Tick(d time.Duration) {
	if s.enableBoids {
		s.boids.Tick(d)
	}
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

				rprojectile.ColorVelocity,
				rprojectile.ColorTrail,
				rprojectile.ColorProjectile,
				rprojectile.ColorHeading,

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
