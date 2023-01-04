package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/feature"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/projectile"
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
			Min: vector.V{xmin - width, ymin - width},
			Max: vector.V{xmin, ymax + width},
		},
		{
			Min: vector.V{xmax, ymin - width},
			Max: vector.V{xmax + width, ymax + width},
		},
		{
			Min: vector.V{xmin, ymin - width},
			Max: vector.V{xmax, ymin},
		},
		{
			Min: vector.V{xmin, ymax},
			Max: vector.V{xmax, ymax + width},
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
					TargetVelocity: v,
					Velocity:       vector.V{0, 0},

					Radius:             r,
					MaxVelocity:        vector.Magnitude(v),
					MaxAngularVelocity: 2 * math.Pi,
					MaxAcceleration:    5,
					Flags:              flags.FSizeSmall,
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
				NFrames:      600,
			})
		}
	}

	opts = append(opts, simulation.O{
		Name: "Box_And_Ball",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    5,
				Flags:              flags.FSizeSmall,
			},
		},
		Features: []feature.O{
			{
				Min: vector.V{70, 20},
				Max: vector.V{90, 80},
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      600,
	}, simulation.O{
		// Make sure agents handle an inner corner by stopping fully.
		Name: "Box_And_Ball_Corner",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    5,
				Flags:              flags.FSizeSmall,
			},
		},
		Features: []feature.O{
			{
				Min: vector.V{70, 20},
				Max: vector.V{90, 80},
			},
			{
				Min: vector.V{50, 80},
				Max: vector.V{90, 100},
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      400,
	}, simulation.O{
		// Make sure that agents are rotation through the smallest angle
		// to their target heading.
		Name: "Rotation_Test",
		Agents: []agent.O{
			// (+X, +Y) to (+X, -Y)
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, math.Pi / 4},
				TargetVelocity:     vector.V{10, -10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Flags:              flags.FSizeSmall,
			},
			// (-X, -Y) to (-X, -Y)
			{
				Position:           vector.V{100, 50},
				Heading:            polar.V{1, 5 * math.Pi / 4},
				TargetVelocity:     vector.V{10, -10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Flags:              flags.FSizeSmall,
			},
			// (+X, -Y) to (+X, +Y)
			{
				Position:           vector.V{100, 100},
				Heading:            polar.V{1, 7 * math.Pi / 4},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Flags:              flags.FSizeSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      250,
	}, simulation.O{
		Name: "Acceleration_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{130, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{-30, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi,
				MaxAcceleration:    10,
				Flags:              flags.FSizeSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      250,
	}, simulation.O{
		Name: "Collision_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{100, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    10,
				Flags:              flags.FSizeSmall,
			},
			{
				Position:           vector.V{100, 45},
				Heading:            polar.V{1, math.Pi},
				TargetVelocity:     vector.V{-100, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    10,
				Flags:              flags.FSizeSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      400,
	}, simulation.O{
		// Agents should move smoothly alongside each other at shallow
		// angles of incidence.
		Name: "Collision_Slide_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{20, 2},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    10,
				Flags:              flags.FSizeSmall,
			},
			{
				Position:           vector.V{50, 80},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{20, -2},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    10,
				Flags:              flags.FSizeSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      250,
	}, simulation.O{
		Name: "Projectile_No_Collision",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{20, 0},
				Velocity:           vector.V{20, 0},
				Radius:             10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    10,
				Flags:              flags.FSizeSmall,
			},
		},
		Projectiles: []projectile.O{
			{
				Position:       vector.V{100, 50},
				Heading:        polar.V{1, math.Pi},
				TargetVelocity: vector.V{-30, 0},
				Velocity:       vector.V{-30, 0},
				Radius:         2,
				Flags:          flags.FSizeProjectile,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      250,
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
