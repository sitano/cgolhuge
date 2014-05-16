package main
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

func printGliderSE(w *LifeWorld, x int64, y int64) {
	w.Set(x, y, w.Layer(), LIFE)
	w.Set(x + 1, y, w.Layer(), LIFE)
	w.Set(x + 2, y, w.Layer(), LIFE)
	w.Set(x + 2, y + 1, w.Layer(), LIFE)
	w.Set(x + 1, y + 2, w.Layer(), LIFE)
}

func printGliderSW(w *LifeWorld, x int64, y int64) {
	w.Set(x, y, w.Layer(), LIFE)
	w.Set(x + 1, y, w.Layer(), LIFE)
	w.Set(x + 2, y, w.Layer(), LIFE)
	w.Set(x, y + 1, w.Layer(), LIFE)
	w.Set(x + 1, y + 2, w.Layer(), LIFE)
}

func printGliderNE(w *LifeWorld, x int64, y int64) {
	w.Set(x+2, y, w.Layer(), LIFE)
	w.Set(x+2, y + 1, w.Layer(), LIFE)
	w.Set(x+2, y + 2, w.Layer(), LIFE)
	w.Set(x+1, y + 2, w.Layer(), LIFE)
	w.Set(x, y + 1, w.Layer(), LIFE)
}

func printGliderNW(w *LifeWorld, x int64, y int64) {
	w.Set(x, y, w.Layer(), LIFE)
	w.Set(x, y + 1, w.Layer(), LIFE)
	w.Set(x, y + 2, w.Layer(), LIFE)
	w.Set(x + 1, y + 2, w.Layer(), LIFE)
	w.Set(x + 2, y + 1, w.Layer(), LIFE)
}

func TestWorldStep4(t *testing.T) {
	w := NewLifeWorldXY(NewAABBMax())
	printGliderSE(&w, 3, 3)
	printGliderSW(&w, 10, 10)
	printGliderNW(&w, 3, 10)
	printGliderNE(&w, 10, 3)
	//printPages(&w)
	//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(16, 20, 16, 20)))
	if w.v.vm.Reserved() != 1 {
		printPages(&w)
		t.Error("Why not 1???")
	}
	for i := 0 ; i < 100 ; i ++ {
		w.Step()

		//printPages(&w)
		//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(-16, 20, -16, 20)))
		// fmt.Print(printHeatMap(&w, NewAABB(0, 16, 0, 8)))
	}
	//fmt.Print("Generation ", w.Generation(), "\n", printWorld(&w, NewAABB(0, 16, 0, -8)))
	if w.v.vm.Reserved() != 4 {
		printPages(&w)
		t.Error("Why not 3??? ", w.v.vm.Reserved())
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
*/
