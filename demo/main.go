package main

import (
	"flag"
	"fmt"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/downflux/go-collider/agent"
	"github.com/downflux/go-collider/agent/mask"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/pathfinding-demo/simulation"

	vnd "github.com/downflux/go-geometry/nd/vector"
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

	agents := make([]agent.O, 0, n)
	min, max := 0.0, math.Sqrt(n*math.Pi*r*r/density)
	cols := math.Floor(math.Sqrt(n))
	grid := (max - min) / cols
	for i := 0; i < n; i++ {
		v := vector.Scale(5*r, vector.V{
			rn(-1, 1),
			rn(-1, 1),
		})
		agents = append(agents, agent.O{
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
		})
	}

	s := simulation.New(simulation.O{
		Agents:   agents,
		Collider: collider.DefaultO,
		Dimensions: *hyperrectangle.New(
			vnd.V{min, min},
			vnd.V{max, max},
		),
		TickDuration: 20 * time.Millisecond,
	})
	anim := s.Execute(nFrames)

	w, err := os.Create(*fnOut)
	if err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", *fnOut, err))
	}
	defer w.Close()

	if err := gif.EncodeAll(w, anim); err != nil {
		panic(fmt.Sprintf("cannot write to file %v: %v", *fnOut, err))
	}
}
