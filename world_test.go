package main

import "testing"

func BenchmarkWorldGliderGun(t *testing.B) {
	w := NewLifeWorld(NewWorldView(NewVM()))
	LoadLIF(w, 20, 10, "pattern/glider_gun.lif")
	t.ResetTimer()
	for i := 0; i < t.N; i ++ {
		w.Step()
	}
}

/*
import "testing"
import "math"
import "fmt"
import "strconv"

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

*/
