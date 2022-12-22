package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/downflux/game-db/agent"
	"github.com/downflux/game-db/agent/mask"
	"github.com/downflux/game-db/db"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"

	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
	rlabel "github.com/downflux/pathfinding-demo/internal/render/label"
)

const (
	n = 200
	r = 5

	density = 0.1
	nFrames = 600
)

var (
	fnOut = flag.String("out", "/dev/null", "GIF output path")

	// label indicates if agent IDs are labeled in the final output.
	label = r >= 10.0
)

func rn(min, max float64) float64 { return min + rand.Float64()*(max-min) }

func main() {
	flag.Parse()

	world := db.New(db.DefaultO)
	agents := make([]*ragent.A, 0, n)

	min, max := 0.0, math.Sqrt(n*math.Pi*r*r/density)

	cols := math.Floor(math.Sqrt(n))
	grid := (max - min) / cols
	for i := 0; i < n; i++ {
		v := vector.Scale(5*r, vector.V{
			rn(-1, 1),
			rn(-1, 1),
		})
		agents = append(agents, ragent.New(
			world.Insert(agent.O{
				Position: vector.Scale(grid, vector.V{
					float64(i%int(cols)) + 0.5,
					math.Floor(float64(i)/cols) + 0.5,
				}),
				Heading: polar.Normalize(
					polar.V{1, rn(0, 2*math.Pi)},
				),
				Velocity: v,

				Radius:             r,
				MaxVelocity:        vector.Magnitude(v),
				MaxAngularVelocity: math.Pi / 2,
				MaxAcceleration:    5,
				Mask:               mask.MSizeSmall,
			}), label),
		)
	}
	frames := make([]*image.Paletted, 0, nFrames)
	for i := 0; i < nFrames; i++ {
		img := image.NewPaletted(
			image.Rectangle{image.Point{int(min), int(min)}, image.Point{int(max), int(max)}},
			[]color.Color{
				color.White,
				rlabel.ColorText,
				ragent.ColorVelocity,
				ragent.ColorTrail,
				ragent.ColorAgent,
				ragent.ColorHeading,
			},
		)
		rlabel.New(fmt.Sprintf("frame %v / %v", i, nFrames), vector.V{0, 0}).Draw(img)

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

	w, err := os.Create(*fnOut)
	if err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", *fnOut, err))
	}
	defer w.Close()

	if err := gif.EncodeAll(w, anim); err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", *fnOut, err))
	}
}
