package main

import (
	"flag"
	"os"
	"os/signal"
	"log"
	"runtime/pprof"
	"runtime"
	"time"
	"fmt"
	"strings"
	"syscall"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var loadfile = flag.String("load", "", "Load RLE/LIF file into world")
var loadx = flag.Uint64("lx", uint64(0), "Load into x position")
var loady = flag.Uint64("ly", uint64(0), "Load into y position")
var viewx = flag.Uint64("vx", uint64(0), "View port top-left x position")
var viewy = flag.Uint64("vy", uint64(0), "View port top-left y position")
var vieww = flag.Uint64("vw", uint64(70), "View port top-left width position")
var viewh = flag.Uint64("vh", uint64(70), "View port top-left height position")
var idle = flag.Int64("idle", int64(0), "Idle ms between steps")
var wait = flag.Bool("wait", false, "Wait <ENTER> for every step")
var start time.Time
var profiling bool

func main() {
	flag.Parse()

	w := NewLifeWorldMax()
	vx:= *viewx
	vy:= *viewy
	vp:= w.v.AABB

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
		screen.DisableInputBuffering()
		screen.HideInputChars()
		// Set view port
		vp = vp.Intersection(NewXYWH(vx, vy,
			Max(1, Min(uint64(screen.cols) - 2, *vieww)),
			Max(1, Min(uint64(screen.rows) - 2, *viewh))))
		// Program start
		start = time.Now()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<- c
		// sig is a ^C, handle it
		if (profiling) {
			pprof.StopCPUProfile()
		} else {
			screen.ShowInputChars()
		}
		os.Exit(1)
	}()

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
	stepEnd := int64(0)
	for {
		if ! profiling {
			screen.Reset()
			screen.PrintAt(1, 1, fmt.Sprintf("Gen: %d, Pop: %d, VMPages: %d, Elapsed: %.1fs, Avg/Step: %d ns, VP: %v",
					w.Generation(), w.Population(), w.v.vm.Pages(), time.Now().Sub(start).Seconds(), stepEnd - stepStart, vp))
			screen.PrintAt(2, 1, w.Print(vp))
			if *idle > 0 {
				time.Sleep(time.Millisecond * time.Duration(*idle))
			}
			if *wait {
				var b []byte = make([]byte, 3)
				for {
					os.Stdin.Read(b)
					if b[0] == 10 {
						break
					}
					if b[0] == 'q' {
						c <- syscall.SIGINT
					}
					if b[0] == 27 {
						if b[1] == 91 {
							switch b[2] {
							case 65: /* up */
								if vy > 0 {
									vy --
								}
							case 66: /* down */
								vy ++
							case 67: /* right */
								vx ++
							case 68: /* left */
								if vx > 0 {
									vx --
								}
							}
							// New view port
							vp = w.v.AABB.Intersection(NewXYWH(vx, vy,
								Max(1, Min(uint64(screen.cols) - 2, *vieww)),
								Max(1, Min(uint64(screen.rows) - 2, *viewh))))
							// Redraw
							screen.Reset()
							screen.PrintAt(1, 1, fmt.Sprintf("Gen: %d, Pop: %d, VMPages: %d, Elapsed: %.1fs, Avg/Step: %d ns, VP: %v",
									w.Generation(), w.Population(), w.v.vm.Pages(), time.Now().Sub(start).Seconds(), stepEnd - stepStart, vp))
							screen.PrintAt(2, 1, w.Print(vp))
						}
					}
				}
			}
			stepStart = time.Now().UnixNano()
		}
		w.Step()
		if ! profiling {
			stepEnd = time.Now().UnixNano()
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
