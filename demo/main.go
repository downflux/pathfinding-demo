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
	"github.com/downflux/game-db/agent/mask"
	"github.com/downflux/game-db/db"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"

	ragent "github.com/downflux/pathfinding-demo/internal/render/agent"
	rlabel "github.com/downflux/pathfinding-demo/internal/render/label"
)

const (
	n = 100
	r = 10

	density = 0.1
	nFrames = 600

	fnOut = "/home/mzhang/downflux/pathfinding-demo/demo/out2.gif"
)

func rn(min, max float64) float64 { return min + rand.Float64()*(max-min) }

func main() {
	world := db.New(db.O{
		LeafSize:  db.DefaultO.LeafSize,
		PoolSize:  1,
		Tolerance: db.DefaultO.Tolerance,
	})
	agents := make([]*ragent.A, 0, n)

	min, max := 0.0, math.Sqrt(n*math.Pi*r*r/density)
	agents = append(agents, ragent.New(
		world.Insert(agent.O{
			Position:           vector.V{50, 50},
			Heading:            polar.V{1, 0},
			Velocity:           vector.V{20, 0},
			Radius:             10,
			Mass:               rn(0, 100),
			MaxVelocity:        10,
			MaxAngularVelocity: math.Pi / 4,
			MaxAcceleration:    5,
			Mask:               mask.MSizeSmall,
		}),
	), ragent.New(
		world.Insert(agent.O{
			Position:           vector.V{200, 50},
			Heading:            polar.V{1, math.Pi},
			Velocity:           vector.V{-20, 0},
			Radius:             10,
			Mass:               rn(0, 100),
			MaxVelocity:        10,
			MaxAngularVelocity: math.Pi / 4,
			MaxAcceleration:    5,
			Mask:               mask.MSizeSmall,
		}),
	), ragent.New(
		world.Insert(agent.O{
			Position: vector.V{300, 50},
			Heading:  polar.V{1, math.Pi},
			Velocity: vector.V{-20, 0},
			Radius:   10,
			Mass:     rn(0, 100),
			// TODO(minkezhang: Check why this isn't being
			// respected by non-projectiles.
			MaxVelocity:        10,
			MaxAngularVelocity: math.Pi / 4,
			MaxAcceleration:    5,
			Mask:               mask.MSizeProjectile,
		}),
	))
	cols := math.Floor(math.Sqrt(n))
	grid := (max - min) / cols
	for i := 0; i < n; i++ {
		v := vector.Scale(5*r, vector.V{
			rn(-1, 1),
			rn(-1, 1),
		})
		agents = append(agents, ragent.New(
			world.Insert(agent.O{
				Position: vector.Scale(grid, vector.V{float64(i % int(cols)), float64(i) / cols}),
				Heading: polar.Normalize(
					polar.Polar(vector.Unit(v)),
				),
				Velocity: v,

				Radius:             r,
				Mass:               rn(0, 100),
				MaxVelocity:        vector.Magnitude(v),
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Mask:               mask.MSizeSmall,
			})),
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

	w, err := os.Create(fnOut)
	if err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", fnOut, err))
	}
	defer w.Close()

	if err := gif.EncodeAll(w, anim); err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", fnOut, err))
	}
}
