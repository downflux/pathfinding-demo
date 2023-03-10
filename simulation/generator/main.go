package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/downflux/go-boids/boids"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/feature"
	"github.com/downflux/go-database/flags/move"
	"github.com/downflux/go-database/flags/size"
	"github.com/downflux/go-database/projectile"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
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
			AABB: *hyperrectangle.New(vector.V{xmin - width, ymin - width}, vector.V{xmin, ymax + width}),
		},
		{
			AABB: *hyperrectangle.New(vector.V{xmax, ymin - width}, vector.V{xmax + width, ymax + width}),
		},
		{
			AABB: *hyperrectangle.New(vector.V{xmin, ymin - width}, vector.V{xmax, ymin}),
		},
		{
			AABB: *hyperrectangle.New(vector.V{xmin, ymax}, vector.V{xmax, ymax + width}),
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
				p := vector.V{
					float64(i%int(cols)) + 0.5,
					math.Floor(float64(i)/cols) + 0.5,
				}
				v := vector.Scale(5*r, vector.V{
					rn(-1, 1),
					rn(-1, 1),
				})
				agents = append(agents, agent.O{
					Position:       vector.Scale(grid, p),
					TargetPosition: vector.V{0, 0},
					Heading: polar.Normalize(
						polar.V{1, rn(0, 2*math.Pi)},
					),
					TargetVelocity:     v,
					Velocity:           vector.V{0, 0},
					Radius:             r,
					Mass:               10,
					MaxVelocity:        vector.Magnitude(v),
					MaxAngularVelocity: 2 * math.Pi,
					MaxAcceleration:    5,
					Size:               size.FSmall,
				})
			}

			opts = append(opts, simulation.O{
				Name:         fmt.Sprintf("Random/N=%v/??=%v", n, density),
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
	for _, n := range []int{100, 200, 1000} {
		for _, density := range []float64{0.01, 0.1} {
			agents := make([]agent.O, 0, n)

			min, max := 0.0, math.Sqrt(float64(n)*math.Pi*r*r/density)
			cols := math.Floor(math.Sqrt(float64(n)))
			grid := (max - min) / cols

			targets := make([]vector.V, 0, n)
			for i := 0; i < n; i++ {
				targets = append(targets, vector.Scale(grid, vector.V{
					float64(i%int(cols)) + 0.5,
					math.Floor(float64(i)/cols) + 0.5,
				}))
			}
			rand.Shuffle(len(targets), func(i, j int) { targets[i], targets[j] = targets[j], targets[i] })
			// Allow for some duplicate goals to see artificial
			// flocking better.
			targets = targets[:n/10]

			for i := 0; i < n; i++ {
				p := vector.V{
					float64(i%int(cols)) + 0.5,
					math.Floor(float64(i)/cols) + 0.5,
				}
				v := vector.Scale(5*r, vector.V{
					rn(-1, 1),
					rn(-1, 1),
				})
				agents = append(agents, agent.O{
					Position:       vector.Scale(grid, p),
					TargetPosition: targets[i%len(targets)],
					Heading: polar.Normalize(
						polar.V{1, rn(0, 2*math.Pi)},
					),
					TargetVelocity:     v,
					Velocity:           vector.V{0, 0},
					Radius:             r,
					Mass:               5,
					MaxVelocity:        50,
					MaxAngularVelocity: math.Pi / 4,
					MaxAcceleration:    40,
					Size:               size.FSmall,
					Move:               move.FAvoidance | move.FArrival | move.FFlocking,
				})
			}

			opts = append(opts, simulation.O{
				Name:         fmt.Sprintf("Random_Boids/N=%v/??=%v", n, density),
				Agents:       agents,
				Features:     borders(min, max, min, max),
				Collider:     collider.DefaultO,
				Boids:        boids.DefaultO,
				EnableBoids:  true,
				MinX:         min,
				MinY:         min,
				MaxX:         max,
				MaxY:         max,
				TickDuration: 20 * time.Millisecond,
				NFrames:      1200,
			})
		}
	}

	opts = append(opts, simulation.O{
		Name: "Box_And_Ball",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    5,
				Size:               size.FSmall,
			},
		},
		Features: []feature.O{
			{
				AABB: *hyperrectangle.New(vector.V{70, 20}, vector.V{90, 80}),
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
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    5,
				Size:               size.FSmall,
			},
		},
		Features: []feature.O{
			{
				AABB: *hyperrectangle.New(vector.V{70, 20}, vector.V{90, 80}),
			},
			{
				AABB: *hyperrectangle.New(vector.V{50, 80}, vector.V{90, 100}),
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
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, math.Pi / 4},
				TargetVelocity:     vector.V{10, -10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Size:               size.FSmall,
			},
			// (-X, -Y) to (-X, -Y)
			{
				Position:           vector.V{100, 50},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 5 * math.Pi / 4},
				TargetVelocity:     vector.V{10, -10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Size:               size.FSmall,
			},
			// (+X, -Y) to (+X, +Y)
			{
				Position:           vector.V{100, 100},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 7 * math.Pi / 4},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    5,
				Size:               size.FSmall,
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
				Position:           vector.V{50, 50},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{30, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi,
				MaxAcceleration:    10,
				Size:               size.FSmall,
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
		Name: "Brake_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{0, 0},
				Velocity:           vector.V{30, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi,
				MaxAcceleration:    10,
				Size:               size.FSmall,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      200,
	}, simulation.O{
		Name: "Boids_Box_And_Ball",
		Agents: []agent.O{
			{
				Position:           vector.V{25, 50},
				TargetPosition:     vector.V{100, 100},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    50,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FArrival | move.FFlocking,
			},
		},
		Features: []feature.O{
			{
				AABB: *hyperrectangle.New(vector.V{70, 20}, vector.V{90, 80}),
			},
		},
		Collider:     collider.DefaultO,
		Boids:        boids.DefaultO,
		EnableBoids:  true,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      600,
	}, simulation.O{
		Name: "Boids_Box_And_Ball_Convex_Corner",
		Agents: []agent.O{
			{
				Position:           vector.V{25, 25},
				TargetPosition:     vector.V{100, 100},
				Heading:            polar.V{1, math.Pi / 4},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi,
				MaxAcceleration:    50,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FArrival | move.FFlocking,
			},
		},
		Features: []feature.O{
			{
				AABB: *hyperrectangle.New(vector.V{50, 50}, vector.V{80, 80}),
			},
		},
		Collider:     collider.DefaultO,
		Boids:        boids.DefaultO,
		EnableBoids:  true,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      600,
	}, simulation.O{
		Name: "Boids_Box_And_Ball_Corner",
		Agents: []agent.O{
			{
				Position:           vector.V{25, 50},
				TargetPosition:     vector.V{100, 100},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{10, 10},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi,
				MaxAcceleration:    50,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FArrival | move.FFlocking,
			},
		},
		Features: []feature.O{
			{
				AABB: *hyperrectangle.New(vector.V{70, 20}, vector.V{90, 80}),
			},
			{
				AABB: *hyperrectangle.New(vector.V{50, 80}, vector.V{90, 100}),
			},
		},
		Collider:     collider.DefaultO,
		Boids:        boids.DefaultO,
		EnableBoids:  true,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      1200,
	}, simulation.O{
		Name: "Arrival_Boids_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{150, 50},
				TargetPosition:     vector.V{250, 50},
				Heading:            polar.V{1, math.Pi / 2},
				TargetVelocity:     vector.V{0, 0},
				Velocity:           vector.V{0, 0},
				Radius:             5,
				Mass:               1,
				MaxVelocity:        50,
				MaxAngularVelocity: math.Pi / 2,
				MaxAcceleration:    50,
				Size:               size.FSmall,
				Move:               move.FArrival,
			},
		},
		Collider:     collider.DefaultO,
		Boids:        boids.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         300,
		MaxY:         100,
		TickDuration: 20 * time.Millisecond,
		NFrames:      600,
		EnableBoids:  true,
	}, simulation.O{
		Name: "Avoidance_Boids_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{150, 50},
				TargetPosition:     vector.V{250, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{0, 0},
				Velocity:           vector.V{0, 0},
				Radius:             5,
				Mass:               1,
				MaxVelocity:        50,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    50,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FSeek,
			},
			{
				Position:           vector.V{200, 45},
				TargetPosition:     vector.V{100, 45},
				Heading:            polar.V{1, math.Pi},
				TargetVelocity:     vector.V{0, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        30,
				MaxAngularVelocity: math.Pi / 8,
				MaxAcceleration:    10,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FSeek,
			},
		},
		Collider:     collider.DefaultO,
		Boids:        boids.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         300,
		MaxY:         100,
		TickDuration: 20 * time.Millisecond,
		NFrames:      600,
		EnableBoids:  true,
	}, simulation.O{
		Name: "Avoidance_Boids_Direct_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{150, 50},
				TargetPosition:     vector.V{250, 50},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{0, 0},
				Velocity:           vector.V{0, 0},
				Radius:             5,
				Mass:               1,
				MaxVelocity:        50,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    50,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FSeek,
			},
			{
				Position:           vector.V{200, 50},
				TargetPosition:     vector.V{100, 50},
				Heading:            polar.V{1, math.Pi},
				TargetVelocity:     vector.V{0, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        30,
				MaxAngularVelocity: math.Pi / 8,
				MaxAcceleration:    10,
				Size:               size.FSmall,
				Move:               move.FAvoidance | move.FSeek,
			},
		},
		Collider:     collider.DefaultO,
		Boids:        boids.DefaultO,
		EnableBoids:  true,
		MinX:         0,
		MinY:         0,
		MaxX:         300,
		MaxY:         100,
		TickDuration: 20 * time.Millisecond,
		NFrames:      600,
	}, simulation.O{
		Name: "Collision_Test",
		Agents: []agent.O{
			{
				Position:           vector.V{50, 50},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{100, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    10,
				Size:               size.FSmall,
			},
			{
				Position:           vector.V{100, 45},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, math.Pi},
				TargetVelocity:     vector.V{-100, 0},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: math.Pi / 4,
				MaxAcceleration:    10,
				Size:               size.FSmall,
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
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{20, 2},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    10,
				Size:               size.FSmall,
			},
			{
				Position:           vector.V{50, 80},
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{20, -2},
				Velocity:           vector.V{0, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    10,
				Size:               size.FSmall,
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
				TargetPosition:     vector.V{0, 0},
				Heading:            polar.V{1, 0},
				TargetVelocity:     vector.V{20, 0},
				Velocity:           vector.V{20, 0},
				Radius:             10,
				Mass:               10,
				MaxVelocity:        100,
				MaxAngularVelocity: 2 * math.Pi,
				MaxAcceleration:    10,
				Size:               size.FSmall,
			},
		},
		Projectiles: []projectile.O{
			{
				Position:       vector.V{100, 50},
				TargetPosition: vector.V{0, 0},
				Heading:        polar.V{1, math.Pi},
				TargetVelocity: vector.V{-30, 0},
				Velocity:       vector.V{-30, 0},
				Radius:         2,
			},
		},
		Collider:     collider.DefaultO,
		MinX:         0,
		MinY:         0,
		MaxX:         150,
		MaxY:         150,
		TickDuration: 20 * time.Millisecond,
		NFrames:      250,
	}, func() simulation.O {
		agents := flock(vector.V{50, 150}, vector.V{450, 350}, 10, 7)
		agents = append(agents, flock(vector.V{50, 450}, vector.V{300, 150}, 10, 7)...)
		agents = append(agents, flock(vector.V{150, 50}, vector.V{50, 450}, 10, 7)...)
		return simulation.O{
			Name:         "Flocking_Small",
			Agents:       agents,
			Collider:     collider.DefaultO,
			Boids:        boids.DefaultO,
			MinX:         0,
			MinY:         0,
			MaxX:         500,
			MaxY:         500,
			TickDuration: 20 * time.Millisecond,
			EnableBoids:  true,
			NFrames:      1200,
		}

	}())

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

func flock(p vector.V, g vector.V, m float64, r float64) []agent.O {
	agents := []agent.O{
		{
			Position:           p,
			TargetPosition:     g,
			Heading:            polar.V{1, 0},
			TargetVelocity:     vector.V{0, 0},
			Velocity:           vector.V{20, 0},
			MaxVelocity:        50,
			MaxAngularVelocity: math.Pi,
			MaxAcceleration:    50,
			Size:               size.FSmall,
			Move:               move.FAvoidance | move.FArrival | move.FFlocking,
		},
		{
			Position:           vector.Add(p, vector.V{-15, -15}),
			TargetPosition:     vector.Add(g, vector.V{-15, -15}),
			Heading:            polar.V{1, 0},
			TargetVelocity:     vector.V{0, 0},
			Velocity:           vector.V{20, 0},
			MaxVelocity:        50,
			MaxAngularVelocity: math.Pi,
			MaxAcceleration:    50,
			Size:               size.FSmall,
			Move:               move.FAvoidance | move.FArrival | move.FFlocking,
		},
		{
			Position:           vector.Add(p, vector.V{-30, 0}),
			TargetPosition:     vector.Add(g, vector.V{-30, 0}),
			Heading:            polar.V{1, 0},
			TargetVelocity:     vector.V{0, 0},
			Velocity:           vector.V{20, 0},
			MaxVelocity:        50,
			MaxAngularVelocity: math.Pi,
			MaxAcceleration:    50,
			Size:               size.FSmall,
			Move:               move.FAvoidance | move.FArrival | move.FFlocking,
		},
		{
			Position:           vector.Add(p, vector.V{-15, 15}),
			TargetPosition:     vector.Add(g, vector.V{-15, 15}),
			Heading:            polar.V{1, 0},
			TargetVelocity:     vector.V{0, 0},
			Velocity:           vector.V{20, 0},
			MaxVelocity:        50,
			MaxAngularVelocity: math.Pi,
			MaxAcceleration:    50,
			Size:               size.FSmall,
			Move:               move.FAvoidance | move.FArrival | move.FFlocking,
		},
		{
			Position:           vector.Add(p, vector.V{30, 10}),
			TargetPosition:     vector.Add(g, vector.V{30, 10}),
			Heading:            polar.V{1, 0},
			TargetVelocity:     vector.V{0, 0},
			Velocity:           vector.V{20, 0},
			MaxVelocity:        50,
			MaxAngularVelocity: math.Pi,
			MaxAcceleration:    50,
			Size:               size.FSmall,
			Move:               move.FAvoidance | move.FArrival | move.FFlocking,
		},
	}
	for i := 0; i < len(agents); i++ {
		f := rn(0.25, 2)
		mass := m * f
		radius := r * math.Sqrt(f)
		a := &agents[i]
		a.Radius = radius
		a.Mass = mass

	}
	return agents
}
