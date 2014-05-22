package main

import (
	"flag"
	"os"
	"log"
	"runtime/pprof"
	"runtime"
	"time"
	"fmt"
	"strings"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var loadfile = flag.String("load", "", "Load RLE/LIF file into world")
var loadx = flag.Uint64("lx", uint64(0), "Load into x position")
var loady = flag.Uint64("ly", uint64(0), "Load into y position")
var viewx = flag.Uint64("vx", uint64(0), "View port top-left x position")
var viewy = flag.Uint64("vy", uint64(0), "View port top-left y position")
var idle = flag.Int64("idle", int64(0), "Idle ms between steps")
var start time.Time
var profiling bool

func main() {
	flag.Parse()

	w := NewWorldView(NewVM())
	vx:= *viewx
	vy:= *viewy
	vp:= w.pb.GetAABB()

	runtime.GC()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		profiling = true
	}

	if *memprofile != "" {
		profiling = true
	}

	var screen *Screen
	if ! profiling {
		screen = NewScreen()
		vp = vp.Intersection(NewXYWH(vx, vy, uint64(screen.cols) - 2, uint64(screen.rows) - 2))
		start = time.Now()
	}

	if *loadfile != "" {
		if strings.HasSuffix(*loadfile, ".rle") {
			LoadRLE(w, *loadx, *loady, *loadfile)
		} else if strings.HasSuffix(*loadfile, ".lif") {
			LoadLIF(w, *loadx, *loady, *loadfile)
		} else {
			panic(fmt.Sprintf("Unknown file format to load ad (%d, %d): %s", vx, vy, *loadfile))
		}
	}

	stepStart := int64(0)
	for {
		if ! profiling {
			stepStart = time.Now().UnixNano()
		}
		// TODO: w.Step()
		if ! profiling {
			stepEnd := time.Now().UnixNano()
			screen.Reset()
			screen.PrintAt(1, 1, fmt.Sprintf("Gen: %d, Pop: %d, VMPages: %d, Elapsed: %.1fs, Avg/Step: %d ns, VP: %v",
				0, 0, w.vm.Pages(), time.Now().Sub(start).Seconds(), stepEnd - stepStart, vp))
			screen.PrintAt(2, 1, w.Print(vp))
			if *idle > 0 {
				time.Sleep(time.Millisecond * time.Duration(*idle))
			}
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
