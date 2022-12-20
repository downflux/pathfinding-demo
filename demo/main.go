package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/downflux/game-db/agent"
	"github.com/downflux/game-db/db"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"

	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
)

const (
	n = 1000
	r = 10

	density = 0.15
	nFrames = 300

	fnOut = "/dev/stdout"
)

func rn(min, max float64) float64 { return min + rand.Float64()*(max-min) }

func main() {
	world := db.New(db.DefaultO)
	agents := make([]*ragent.A, 0, n)

	min, max := 0.0, math.Sqrt(n*math.Pi*r*r/density)
	for i := 0; i < n; i++ {
		v := vector.Scale(3*r, vector.V{
			rn(-1, 1),
			rn(-1, 1),
		})
		agents = append(agents, ragent.New(
			world.Insert(agent.O{
				Position: vector.V{rn(min, max), rn(min, max)},
				Heading: polar.Normalize(
					polar.Polar(vector.Unit(v)),
				),
				Velocity: v,

				Radius:             r,
				Mass:               rn(0, 100),
				MaxVelocity:        10,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    1,
			})),
		)
	}

	frames := make([]*image.Paletted, 0, nFrames)
	for i := 0; i < nFrames; i++ {
		img := image.NewPaletted(
			image.Rectangle{image.Point{int(min), int(min)}, image.Point{int(max), int(max)}},
			[]color.Color{
				color.White,
				ragent.ColorVelocity,
				ragent.ColorTrail,
				ragent.ColorAgent,
				ragent.ColorHeading,
			},
		)

		for _, a := range agents {
			a.Draw(img)
		}
		frames = append(frames, img)

		world.Tick(20 * time.Millisecond)
	}

	delays := make([]int, nFrames)
	for i := 0; i < nFrames; i++ {
		delays[i] = 2
	}
	anim := &gif.GIF{
		Delay: delays,
		Image: frames,
	}

	w, err := os.Create(fnOut)
	if err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", fnOut, err))
	}
	defer w.Close()

	if err := gif.EncodeAll(w, anim); err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", fnOut, err))
	}
}
