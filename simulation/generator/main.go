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
	"github.com/downflux/go-collider/agent/mask"
	"github.com/downflux/go-collider/collider"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/pathfinding-demo/simulation"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

const (
	r = 1.0
)

var (
	dir = flag.String("directory", "/dev/null", "")
)

func rn(min, max float64) float64 { return min + rand.Float64()*(max-min) }

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
				Name:     fmt.Sprintf("Random/N=%v/ρ=%v", n, density),
				Agents:   agents,
				Collider: collider.DefaultO,
				Dimensions: *hyperrectangle.New(
					vnd.V{min, min},
					vnd.V{max, max},
				),
				TickDuration: 20 * time.Millisecond,
			})
		}
	}

	for _, o := range opts {
		fn := path.Join(*dir, fmt.Sprintf("%v.json", o.Filename()))

		func() {
			if *dir == "/dev/null" {
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
