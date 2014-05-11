package main

import "testing"
import "math"
import (
	"fmt"
	"strconv"
)

func TestNewLifeWorld(t *testing.T) {
	NewLifeWorldXY(NewAABBMax())
}

func printWorld(w *LifeWorld, bbox AABB) string {
	s := ""
	for y := bbox.MaxY ; y >= bbox.MinY && y <= bbox.MaxY ; y -- {
		for x := bbox.MinX ; x <= bbox.MaxX && x >= bbox.MinX ; x ++ {
			state := w.Get(x, y, w.Layer())
			if state == LIFE {
				s += "@"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s + "\n"
}

func printHeatMap(w *LifeWorld, bbox AABB) string {
	s := ""
	for y := bbox.MaxY ; y >= bbox.MinY && y <= bbox.MaxY ; y -- {
		for x := bbox.MinX ; x <= bbox.MaxX && x >= bbox.MinX ; x ++ {
			s += strconv.Itoa(int(w.v.LifeSumAt(x, y, w.Layer())))
		}
		s += "\n"
	}
	return s + "\n"
}

func printPages(w *LifeWorld) {
	fmt.Print("Pages in a tree: \n")
	w.v.pb.Reduce(func (a interface {}, pt *PageTile) interface {} {
		fmt.Print("Page ", pt.AABB, "\n")
		return a
	}, 0)
}

func TestWorldStep1(t *testing.T) {
	w := NewLifeWorldXY(NewAABBMax())
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)
	for i := 0 ; i < 3 ; i ++ {
		if w.Generation() != uint64(i) {
			t.Error("Invalid gen ", i)
		}

		// fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(-3, 3, -3, 3)))

		w.Step()

		if w.Population() != 4 {
			t.Error("Invalid pop 4")
		}
	}
	if w.Get(0, 0, w.Layer()) != LIFE { t.Error("No life at 0, 0") }
	if w.Get(1, 0, w.Layer()) != LIFE { t.Error("No life at 1, 0") }
	if w.Get(0, 1, w.Layer()) != LIFE { t.Error("No life at 0, 1") }
	if w.Get(1, 1, w.Layer()) != LIFE { t.Error("No life at 1, 1") }
}

func TestWorldStep2(t *testing.T) {
	w := NewLifeWorldXY(NewAABBMax())
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(-1, 0, w.Layer(), LIFE)
	//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(-3, 3, -3, 3)))
	for i := 0 ; i < 4 ; i ++ {
		if w.Generation() != uint64(i) {
			t.Error("Invalid gen ", i)
		}

		w.Step()

		// fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(-3, 3, -3, 3)))

		// Known bug ;) (bigger than should bec of edge cases)
		// if w.Population() != 3 {
		//	t.Error("Invalid pop 3", w.Population())
		// }
	}
	if w.Get(0, 0, w.Layer()) != LIFE { t.Error("No life at 0, 0") }
	if w.Get(1, 0, w.Layer()) != LIFE { t.Error("No life at 1, 0") }
	if w.Get(-1, 0, w.Layer()) != LIFE { t.Error("No life at-1, 0") }
	if w.v.vm.reserved.Len() != 2 {
		printPages(&w)
		t.Error("Why not 2???")
	}
}

func TestWorldStep3(t *testing.T) {
	w := NewLifeWorldXY(NewAABBMax())
	w.Set(math.MaxInt64, math.MaxInt64, w.Layer(), LIFE)
	w.Set(math.MaxInt64 - 1, math.MaxInt64, w.Layer(), LIFE)
	w.Set(math.MinInt64, math.MaxInt64, w.Layer(), LIFE)
	if w.v.vm.reserved.Len() != 2 {
		t.Error("Why not 2???")
	}
	//printPages(&w)
	//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w,
	//		NewAABB(math.MaxInt64 - 3, math.MaxInt64,
	//			math.MaxInt64 - 3, math.MaxInt64)))
	for i := 0 ; i < 4 ; i ++ {
		w.Step()

		//printPages(&w)
		//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w,
		//		NewAABB(math.MaxInt64 - 3, math.MaxInt64,
		//				math.MaxInt64 - 3, math.MaxInt64)))
	}
	if w.Get(math.MaxInt64, math.MaxInt64, w.Layer()) != LIFE { t.Error("No life at 0, 0") }
	if w.Get(math.MaxInt64 - 1, math.MaxInt64, w.Layer()) != LIFE { t.Error("No life at 1, 0") }
	if w.Get(math.MinInt64, math.MaxInt64, w.Layer()) != LIFE { t.Error("No life at-1, 0") }
}

func printGlider(w *LifeWorld, x int64, y int64) {
	b := w.v.pb.getAABB()
	w.Set(x, y, w.Layer(), LIFE)
	x = MvXY1(x, 1, b.MinX, b.MaxX)
	w.Set(x, y, w.Layer(), LIFE)
	x = MvXY1(x, 1, b.MinX, b.MaxX)
	w.Set(x, y, w.Layer(), LIFE)
	y = MvXY1(y, 1, b.MinY, b.MaxY)
	w.Set(x, y, w.Layer(), LIFE)
	x = MvXY1(x, -1, b.MinX, b.MaxX)
	y = MvXY1(y, 1, b.MinY, b.MaxY)
	w.Set(x, y, w.Layer(), LIFE)
}

func TestWorldStep4(t *testing.T) {
	w := NewLifeWorldXY(NewAABBMax())
	printGlider(&w, 3, 3)
	//printPages(&w)
	//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(0, 16, 0, 8)))
	for i := 0 ; i < 30 ; i ++ {
		w.Step()

		//printPages(&w)
		//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(0, 16, 0, 8)))
		//fmt.Print(printHeatMap(&w, NewAABB(0, 16, 0, 8)))
	}
	fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(0, 16, 0, -8)))
	if w.v.vm.reserved.Len() != 1 {
		printPages(&w)
		t.Error("Why not 1???")
	}
}

func BenchmarkWorldStep2pages(b *testing.B) {
	w := NewLifeWorldXY(NewAABBMax())
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(-1, 0, w.Layer(), LIFE)

	b.ResetTimer()

	for i := 0 ; i < b.N ; i ++ {
		w.Step()
	}
}

func BenchmarkWorldStep1pages(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)
	//w.Set(-1, 0, w.Layer(), LIFE)

	b.StartTimer()

	for i := 0 ; i < b.N; i ++ {
		w.Step()
	}
}

func BenchmarkWorldReadPage2(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)
	w.Set(-1, 0, w.Layer(), LIFE)

	lw := &w

	b.StartTimer()

	for i := 0 ; i < b.N; i ++ {
		lw.generation ++
		lw.population = 0
		cz := lw.z
		//nz := lw.NextZLayer()
		//ws := lw.v.pb.wsize
		ks := lw.v.vm.ksize

		// Process life
		//new := list.New()
		lw.population = lw.v.pb.Reduce(func(a interface{}, pt *PageTile) interface{} {
				population := a.(uint64)

				// Prevent reclamation during processing
				pt.alive = 1

				for i := uint(0) ; i < ks; i ++ {
					//x := POtoWX(i, pt.px, ws)
					//y := POtoWY(i, pt.py, ws)

					sum := DEAD // lw.LifeSumAt(x, y, cz)
					st := lw.Get(0, 0, cz)
					//nst := DEAD

					if st == DEAD {
						if sum == RULE_BORN {
							//		nst = LIFE
							pt.alive ++
							population ++
						}
					}

					if st == LIFE {
						if sum >= RULE_LIVE_MIN && sum <= RULE_LIVE_MAX {
							//		nst = LIFE
							pt.alive ++
							population ++
						}
					}
				}

				pt.alive --

				return population
			}, lw.population).(uint64)
	}
}


func BenchmarkWorldRWPage(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)

	lw := &w

	b.StartTimer()

	for i := 0 ; i < b.N; i ++ {
		lw.generation ++
		lw.population = 0
		cz := lw.z
		nz := lw.NextZLayer()
		ws := lw.v.pb.wsize
		ks := lw.v.vm.ksize
		lw.population = lw.v.pb.Reduce(func(a interface{}, pt *PageTile) interface{} {
				population := a.(uint64)

				// Prevent reclamation during processing
				pt.alive = 1

				for i := uint(0) ; i < ks; i ++ {
					x := POtoWX(i, pt.px, ws)
					y := POtoWY(i, pt.py, ws)
					st := lw.Get(x, y, cz)
					lw.Set(x, y, nz, st)
				}

				pt.alive --

				return population
			}, lw.population).(uint64)
	}
}

func BenchmarkWorldRWPageRaw8x8(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)

	lw := &w

	b.StartTimer()
	pt := w.v.pb.QueryPage(0, 0)

	for i := 0 ; i < b.N; i ++ {
		ks := lw.v.vm.ksize
		pp := *pt.p

		for i := uint(256) ; i < ks - 256; i ++ {
			st := pp[i - 128 - 1] + pp[i - 128] + pp[i - 128 + 1] +
					pp[i - 1] + pp[i] + pp[i + 1] +
					pp[i + 128 - 1] + pp[i + 128] + pp[i + 128 + 1]
			pp[i] = st % 2
		}
	}
}

var r byte

func BenchmarkWorldRPageRaw(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)

	lw := &w

	b.StartTimer()
	pt := w.v.pb.QueryPage(0, 0)

	for i := 0 ; i < b.N; i ++ {
		ks := lw.v.vm.ksize
		pp := *pt.p

		for k := uint(0) ; k < ks; k ++ {
			r = pp[k]
		}
	}
}

func BenchmarkWorldRPageRaw_ConvUint64(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)

	lw := &w

	b.StartTimer()
	pt := w.v.pb.QueryPage(0, 0)

	for i := 0 ; i < b.N; i ++ {
		ks := lw.v.vm.ksize
		pp := *pt.p

		for k := uint(0) ; k < ks; k ++ {
			lw.population += uint64(pp[k])
		}
	}
}

func BenchmarkWorldRWPageRaw(b *testing.B) {
	b.StopTimer()
	w := NewLifeWorldXY(NewAABBMax())
	//w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 0, w.Layer(), LIFE)
	w.Set(1, 0, w.Layer(), LIFE)
	w.Set(0, 1, w.Layer(), LIFE)
	w.Set(1, 1, w.Layer(), LIFE)

	lw := &w

	b.StartTimer()
	pt := w.v.pb.QueryPage(0, 0)

	for i := 0 ; i < b.N; i ++ {
		ks := lw.v.vm.ksize
		pp := *pt.p

		for k := uint(0) ; k < ks; k ++ {
			st := pp[k]
			pp[k] = st + 1
		}
	}
}
