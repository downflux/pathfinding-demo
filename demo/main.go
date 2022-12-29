package main

import (
	"flag"
	"fmt"
	"image/gif"
	"os"
	"path/filepath"

	"github.com/downflux/pathfinding-demo/simulation"
)

var (
	configs = flag.String("configs", "", "config directory")
	output  = flag.String("output", "/dev/null", "output GIF directory")
)

func main() {
	flag.Parse()

	matches, err := filepath.Glob(filepath.Join(*configs, "*.json"))
	if err != nil {
		panic(fmt.Sprintf("cannot match input config directory: %v", err))
	}

	var opts []simulation.O
	for _, fn := range matches {
		func() {
			data, err := os.ReadFile(fn)
			if err != nil {
				panic(fmt.Sprintf("cannot read file: %v", err))
			}
			opts = append(opts, simulation.Unmarshal(data))
		}()
	}

	for _, o := range opts {
		fn := filepath.Join(*output, fmt.Sprintf("%v.gif", o.Filename()))
		fmt.Printf("running %v (%v)\n", o.Name, fn)
		s := simulation.New(o)
		anim := s.Execute()

		fmt.Printf("  average tick time = %v\n", s.TickTimer())

		func() {
			if *output == "/dev/null" {
				return
			}

			w, err := os.Create(fn)
			if err != nil {
				panic(fmt.Sprintf("cannot write to file %v: %v", fn, err))
			}
			defer w.Close()

			if err := gif.EncodeAll(w, anim); err != nil {
				panic(fmt.Sprintf("cannot write to file %v: %v", fn, err))
			}
		}()
	}
}
