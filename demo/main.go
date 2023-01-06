package main

import (
	"flag"
	"fmt"
	"image/gif"
	"log"
	"os"
	"path/filepath"

	"github.com/downflux/pathfinding-demo/simulation"
)

var (
	configs = flag.String("configs", "", "config directory")
	output  = flag.String("output", "/dev/null", "output GIF directory")
	logger  = flag.String("log", "/dev/null", "output log file")
)

func main() {
	flag.Parse()

	matches, err := filepath.Glob(*configs)
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
			o := simulation.Unmarshal(data)
			if *logger != "/dev/null" {
				(&o).Collider.Debug = true
			}
			opts = append(opts, o)
		}()
	}

	for _, o := range opts {
		if *logger != "/dev/null" {
			lfn, err := os.OpenFile(filepath.Join(*logger, fmt.Sprintf("%v.log", o.Filename())), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				panic(fmt.Sprintf("cannot open log file: %v", err))
			}
			log.SetOutput(lfn)
			log.SetFlags(log.Flags() | log.Lshortfile)
		}

		fn := filepath.Join(*output, fmt.Sprintf("%v.gif", o.Filename()))
		fmt.Printf("running %v (%v)\n", o.Name, fn)
		s := simulation.New(o)
		anim := s.Execute()

		fmt.Printf("  average tick time = %v / frame\n", s.TickTimer())

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
