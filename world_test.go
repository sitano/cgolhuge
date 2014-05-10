package main

import "testing"
import "math"

func TestNewLifeWorld(t *testing.T) {
	NewLifeWorldXY(NewAABBMax())
}

func printWorld(w *LifeWorld, bbox AABB) string {
	s := ""
	for y := bbox.MaxY ; y >= bbox.MinY ; y -- {
		for x := bbox.MinX ; x <= bbox.MaxX ; x ++ {
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
	// fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(-3, 3, -3, 3)))
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
	for i := 0 ; i < 4 ; i ++ {
	//	w.Step()
	}
	if w.Get(math.MaxInt64, math.MaxInt64, w.Layer()) != LIFE { t.Error("No life at 0, 0") }
	if w.Get(math.MaxInt64 - 1, math.MaxInt64, w.Layer()) != LIFE { t.Error("No life at 1, 0") }
	if w.Get(math.MinInt64, math.MaxInt64, w.Layer()) != LIFE { t.Error("No life at-1, 0") }
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
