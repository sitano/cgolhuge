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
var profiling  = false

func main() {
	flag.Parse()

	/*
	w := NewLifeWorldXY(NewAABBMax())
	PrintGliderSE(&w, 3, 3)
	PrintGliderSW(&w, 10, 10)
	PrintGliderNW(&w, 3, 10)
	PrintGliderNE(&w, 10, 3)

	viewNW := NewAABB(0, -20, 0, 20)
	viewNE := NewAABB(0,  20, 0, 20)
	viewSW := NewAABB(0, -20,-20, 0)
	viewSE := NewAABB(0,  20,-20, 0)
	*/

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

	/*
	if ! profiling {
		fmt.Print("\033[2J")
		PrintWorld(&w, viewNW, 2, 2)
		PrintWorld(&w, viewNE, 2, 4 + int(viewNW.SizeX()))
		PrintWorld(&w, viewSW, 4 + int(viewNW.SizeY()), 2)
		PrintWorld(&w, viewSE, 4 + int(viewNW.SizeY()), 4 + int(viewNW.SizeX()))
		start = time.Now()
	}
	stepStart := int64(0)
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
/*
func PrintGliderSE(w *LifeWorld, x int64, y int64) {
	w.Set(x, y, w.Layer(), LIFE)
	w.Set(x + 1, y, w.Layer(), LIFE)
	w.Set(x + 2, y, w.Layer(), LIFE)
	w.Set(x + 2, y + 1, w.Layer(), LIFE)
	w.Set(x + 1, y + 2, w.Layer(), LIFE)
}

func PrintGliderSW(w *LifeWorld, x int64, y int64) {
	w.Set(x, y, w.Layer(), LIFE)
	w.Set(x + 1, y, w.Layer(), LIFE)
	w.Set(x + 2, y, w.Layer(), LIFE)
	w.Set(x, y + 1, w.Layer(), LIFE)
	w.Set(x + 1, y + 2, w.Layer(), LIFE)
}

func PrintGliderNE(w *LifeWorld, x int64, y int64) {
	w.Set(x+2, y, w.Layer(), LIFE)
	w.Set(x+2, y + 1, w.Layer(), LIFE)
	w.Set(x+2, y + 2, w.Layer(), LIFE)
	w.Set(x+1, y + 2, w.Layer(), LIFE)
	w.Set(x, y + 1, w.Layer(), LIFE)
}

func PrintGliderNW(w *LifeWorld, x int64, y int64) {
	w.Set(x, y, w.Layer(), LIFE)
	w.Set(x, y + 1, w.Layer(), LIFE)
	w.Set(x, y + 2, w.Layer(), LIFE)
	w.Set(x + 1, y + 2, w.Layer(), LIFE)
	w.Set(x + 2, y + 1, w.Layer(), LIFE)
}

func PrintWorld(w *LifeWorld, bbox AABB, row int, col int) {
	for y := bbox.MaxY ; y >= bbox.MinY && y <= bbox.MaxY ; y -- {
		fmt.Printf("\033[%d;%dH", row, col)
		for x := bbox.MinX ; x <= bbox.MaxX && x >= bbox.MinX ; x ++ {
			state := w.Get(x, y, w.Layer())
			if state == LIFE {
				fmt.Print("@")
			} else {
				fmt.Print(".")
			}
		}
		row ++
	}
}                 */
