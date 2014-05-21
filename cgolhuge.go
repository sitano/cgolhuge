package main

import (
	"flag"
	"os"
	"log"
	"runtime/pprof"
	"runtime"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var loadfile = flag.String("load", "", "Load RLE/LIF file into world")
var loadx = flag.Uint64("x", uint64(0), "Start x position")
var loady = flag.Uint64("y", uint64(0), "Start y position")
var profiling  = false

func main() {
	flag.Parse()

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
	}

	if *loadfile != "" {
		// TODO: do something into x, y coord
	}

	screen.Reset()
	screen.PrintAt(5, 5, "Hi guys\n1\n2\n3\n4\n5")
	screen.Println()

/*	stepStart := int64(0)
	for i := 0 ; i < 1000 ; i ++ {
		if ! profiling {
			stepStart = time.Now().UnixNano()
		}
		w.Step()
		if ! profiling {
			stepEnd := time.Now().UnixNano()
			fmt.Print("\033[2J")
			PrintWorld(&w, viewNW, 2, 2)
			PrintWorld(&w, viewNE, 2, 4 + int(viewNW.SizeX()))
			PrintWorld(&w, viewSW, 4 + int(viewNW.SizeY()), 2)
			PrintWorld(&w, viewSE, 4 + int(viewNW.SizeY()), 4 + int(viewNW.SizeX()))
			fmt.Printf("\033[%d;%dH Gen: %d, Pop: %d, VMPages: %d, Elapsed: %.1fs, Avg/Step: %d ns", 0, 0,
				w.generation,
				w.population,
				w.v.vm.Reserved(),
				time.Now().Sub(start).Seconds(),
				stepEnd - stepStart)
		}
	}
    */

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
