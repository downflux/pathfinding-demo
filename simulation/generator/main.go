package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/downflux/go-collider/agent"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-collider/feature"
	"github.com/downflux/go-collider/mask"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/pathfinding-demo/simulation"
)

const (
	r = 5.0
)

var (
	output = flag.String("output", "/dev/null", "output config directory")
)

func rn(min, max float64) float64 { return min + rand.Float64()*(max-min) }

func borders(xmin, xmax, ymin, ymax float64) []feature.O {
	width := 2 * r
	return []feature.O{
		{
			Min:  vector.V{xmin - width, ymin - width},
			Max:  vector.V{xmin, ymax + width},
			Mask: mask.MSizeSmall,
		},
		{
			Min:  vector.V{xmax, ymin - width},
			Max:  vector.V{xmax + width, ymax + width},
			Mask: mask.MSizeSmall,
		},
		{
			Min:  vector.V{xmin, ymin - width},
			Max:  vector.V{xmax, ymin},
			Mask: mask.MSizeSmall,
		},
		{
			Min:  vector.V{xmin, ymax},
			Max:  vector.V{xmax, ymax + width},
			Mask: mask.MSizeSmall,
		},
	}
}

func main() {
	flag.Parse()

	var opts []simulation.O
	for _, n := range []int{100, 200, 1000} {
		for _, density := range []float64{0.01, 0.1} {
			agents := make([]agent.O, 0, n)

			min, max := 0.0, math.Sqrt(float64(n)*math.Pi*r*r/density)
			cols := math.Floor(math.Sqrt(float64(n)))
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

			opts = append(opts, simulation.O{
				Name:         fmt.Sprintf("Random/N=%v/ρ=%v", n, density),
				Agents:       agents,
				Features:     borders(min, max, min, max),
				Collider:     collider.DefaultO,
				MinX:         min,
				MinY:         min,
				MaxX:         max,
				MaxY:         max,
				TickDuration: 20 * time.Millisecond,
			})
		}
	}

	opts = append(opts, simulation.O{
		Name: "Box_And_Ball",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				Velocity:           vector.V{10, 10},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 2,
				MaxAcceleration:    5,
				Mask:               mask.MSizeSmall,
			},
		},
		Features: []feature.O{
			{
				Min:  vector.V{70, 20},
				Max:  vector.V{90, 80},
				Mask: mask.MSizeSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
	},
	simulation.O{
		Name: "Box_And_Ball_Corner",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				Velocity:           vector.V{10, 10},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 2,
				MaxAcceleration:    5,
				Mask:               mask.MSizeSmall,
			},
		},
		Features: []feature.O{
			{
				Min:  vector.V{70, 20},
				Max:  vector.V{90, 80},
				Mask: mask.MSizeSmall,
			},
			{
				Min:  vector.V{50, 80},
				Max:  vector.V{90, 100},
				Mask: mask.MSizeSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
	})

	for _, o := range opts {
		fn := path.Join(*output, fmt.Sprintf("%v.json", o.Filename()))

		func() {
			if *output == "/dev/null" {
				return
			}

			w, err := os.Create(fn)
			if err != nil {
				panic(fmt.Sprintf("cannot write to file %v: %v", fn, err))
			}
			defer w.Close()

			w.Write(o.Marshal())
		}()
	}
}
