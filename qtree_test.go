package main

import "math"
import "math/rand"
import "testing"
import _ "fmt"

// Generates n AABBes in the range of frame with average width and height avgSize
func randomPages(n int, frame AABB) []*Page {
	ret := make([]*Page, n)

	for i := 0; i < len(ret); i++ {
		ret[i] = NewPage()
		ret[i].px = uint64(rand.Float64() * float64(frame.SizeX()) + float64(frame.MinX))
		ret[i].py = uint64(rand.Float64() * float64(frame.SizeY()) + float64(frame.MinY))
	}

	return ret
}

// Returns all elements of data which intersect query
func queryLinear(data []*Page, query *Page) (ret []*Page) {
	for _, v := range data {
		if compareQuadElement(v, query) {
			ret = append(ret, v)
		}
	}

	return ret
}

func compareQuadElement(v1, v2 *Page) bool {
	return v1.px == v2.px && v1.py == v2.py
}

func lookupResults(r1, r2 []*Page) int {
	for i, v1 := range r1 {
		found := false

		for _, v2 := range r2 {
			if compareQuadElement(v1, v2) {
				found = true
				break
			}
		}

		if ! found {
			return i
		}
	}

	return -1
}

// World-space extends from 0..1000 in X and Y direction
var world AABB = NewAABB(0, 0, 1000, 1000)
var worldMax AABB = NewAABB(0, 0, math.MaxUint64, math.MaxUint64)

// Compary correctness of quad-tree results vs simple look-up on set of random rectangles
/*func TestQuadTreeRects(t *testing.T) {
	var rects []*Page = randomPages(100, world)
	qt := NewQuadTree(world)

	for _, v := range rects {
		qt.Add(v)
	}

	queries := randomPages(100, world)

	for _, q := range queries {
		r1 := queryLinear(rects, q)
		r2 := qt.QueryBox(q.AABB)
		r3 := qt.QueryPoint(q.px, q.py)

		if len(r1) != len(r2) {
			t.Errorf("r1 and r2 differ: %v   %v\n", r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("%v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("%v was not returned by brute-force\n", r2[i])
		}

		if r3 == nil {
			t.Errorf("%v was not returned by QT\n", r3)
		}
	}
}
  */

// Compary correctness of quad-tree results vs simple look-up on set of random points
func TestQuadTreePointsSmall(t *testing.T) {
	var points []*Page = randomPages(100, world)
	qt := NewQuadTree(world)

	for _, v := range points {
		qt.Add(v)
	}

	for _, q := range points {
		r1 := queryLinear(points, q)
		r2 := qt.QueryBox(NewAABBW2P(q.AABB))
		r3 := qt.QueryPoint(q.px, q.py)

		if len(r1) != len(r2) {
			t.Errorf("q for r1 and r2 differ: %v    %v   %v\n", q, r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("r1[i] = %v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("r2[i] = %v was not returned by brute-force\n", r2[i])
		}

		if r3 == nil {
			t.Errorf("r3 = %v was not returned by QT\n", r3)
		}
	}

	for _, q := range points {
		if qt.RemoveAt(q.px, q.py) == nil {
			t.Errorf("Failed to remove %v\n", q)
		}
		r3 := qt.QueryPoint(q.px, q.py)
		if r3 != nil {
			t.Errorf("r3 = %v was not returned by QT\n", r3)
		}
	}

	for _, v := range points {
		qt.Add(v)
	}

	for _, q := range points {
		r1 := queryLinear(points, q)
		r2 := qt.QueryBox(NewAABBW2P(q.AABB))
		r3 := qt.QueryPoint(q.px, q.py)

		if len(r1) != len(r2) {
			t.Errorf("q for r1 and r2 differ: %v    %v   %v\n", q, r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("r1[i] = %v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("r2[i] = %v was not returned by brute-force\n", r2[i])
		}

		if r3 == nil {
			t.Errorf("r3 = %v was not returned by QT\n", r3)
		}
	}
}

func TestQuadTreePointsMax(t *testing.T) {
	var points []*Page = randomPages(1000, worldMax)
	qt := NewQuadTree(worldMax)

	for _, v := range points {
		qt.Add(v)
	}

	for _, q := range points {
		r1 := queryLinear(points, q)
		r2 := qt.QueryBox(NewAABBW2P(q.AABB))
		r3 := qt.QueryPoint(q.px, q.py)

		if len(r1) != len(r2) {
			t.Errorf("q for r1 and r2 differ: %v    %v   %v\n", q, r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("r1[i] = %v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("r2[i] = %v was not returned by brute-force\n", r2[i])
		}

		if r3 == nil {
			t.Errorf("r3 = %v was not returned by QT\n", r3)
		}
	}
}

func TestQuadTreeReduce(t *testing.T) {
	var rects []*Page = randomPages(100, world)
	qt := NewQuadTree(world)

	for _, v := range rects {
		qt.Add(v)
	}

	count := 0
	qt.Reduce(func(p *Page) {
		count += 1
	})

	if count != 100 {
		t.Error("Reduce wrong abount 100 inner rects != ", count)
	}
}

func fillQuadTreeWithPageMatrix(qt *QuadTree, ps [][]*Page, x uint64, y uint64) {
	for r, ps1 := range ps {
		for c, ps2 := range ps1 {
			if ps2 != nil {
				qt.AddTo(ps2, x + uint64(c - 1), y + uint64(r - 1))
			}
		}
	}
}

func validateQuadTreePageRelations(t *testing.T, ps [][]*Page) {
	for r, ps1 := range ps {
		for c, ps2 := range ps1 {
			if ps2 != nil {
				if c > 0 {
					if ps2.ap_w != ps1[c - 1] {
						t.Errorf("Error for (%d, %d) of %v/w, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_w, ps1[c - 1])
					}
					if r > 0 {
						if ps2.ap_nw != ps[r - 1][c - 1] {
							t.Errorf("Error for (%d, %d) of %v/nw, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_nw, ps[r - 1][c - 1])
						}
					}
					if r < len(ps) - 1 {
						if ps2.ap_sw != ps[r + 1][c - 1] {
							t.Errorf("Error for (%d, %d) of %v/w, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_sw, ps[r + 1][c - 1])
						}
					}
				}
				if c < len(ps1) - 1 {
					if ps2.ap_e != ps1[c + 1] {
						t.Errorf("Error for (%d, %d) of %v/w, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_e, ps1[c + 1])
					}
					if r > 0 {
						if ps2.ap_ne != ps[r - 1][c + 1] {
							t.Errorf("Error for (%d, %d) of %v/nw, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_ne, ps[r - 1][c + 1])
						}
					}
					if r < len(ps) - 1 {
						if ps2.ap_se != ps[r + 1][c + 1] {
							t.Errorf("Error for (%d, %d) of %v/w, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_se, ps[r + 1][c + 1])
						}
					}
				}
				if r > 0 {
					if ps2.ap_n != ps[r - 1][c] {
						t.Errorf("Error for (%d, %d) of %v/nw, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_n, ps[r - 1][c])
					}
				}
				if r < len(ps) - 1 {
					if ps2.ap_s != ps[r + 1][c] {
						t.Errorf("Error for (%d, %d) of %v/w, have %v, wait %v\n", c - 1, r - 1, ps2, ps2.ap_s, ps[r + 1][c])
					}
				}
			}
		}
	}
}

func TestQuadTreeAdjacent(t *testing.T) {
	qt := NewQuadTree(world)
	vm := NewVM()
	// Init matrix for 0, 0
	ps := [][]*Page{
		{ nil, nil, nil },
		{ nil, vm.ReservePage(), vm.ReservePage() },
		{ nil, vm.ReservePage(), vm.ReservePage() },
	}
	// Fill 0, 0
	fillQuadTreeWithPageMatrix(qt, ps, 0, 0)
	// Check 0, 0
	validateQuadTreePageRelations(t, ps)

	// Init matrix for MaxX, MaxY
	ps = [][]*Page{
		{ vm.ReservePage(), vm.ReservePage(), nil },
		{ vm.ReservePage(), vm.ReservePage(), nil },
		{ nil, nil, nil },
	}
	// Fill
	fillQuadTreeWithPageMatrix(qt, ps, world.MaxX, world.MaxY)
	// Check MaxX, MaxY
	validateQuadTreePageRelations(t, ps)

	// Init matrix for center
	ps = [][]*Page{
		{ vm.ReservePage(), vm.ReservePage(), vm.ReservePage() },
		{ vm.ReservePage(), vm.ReservePage(), vm.ReservePage() },
		{ vm.ReservePage(), vm.ReservePage(), vm.ReservePage() },
	}
	// Fill
	fillQuadTreeWithPageMatrix(qt, ps, world.MinX + world.SizeX() / 2, world.MinY + world.SizeY() / 2)
	// Check MaxX, MaxY
	validateQuadTreePageRelations(t, ps)
}

// A set of 10 million randomly distributed rectangles of avg size 5
var boxes10M []*Page

func BenchmarkTreeInitBoxes(b *testing.B) {
	boxes10M = randomPages(10*1000*1000, world)
}

// Benchmark insertion into quad-tree
func BenchmarkTreeInsert(b *testing.B) {
	var values []*Page = randomPages(b.N, world)
	qt := NewQuadTree(world)

	b.ResetTimer()

	for _, v := range values {
		qt.Add(v)
	}
}


// Benchmark quad-tree on set of rectangles
func BenchmarkTreeRectsQuadtree(b *testing.B) {
	rand.Seed(1)
	qt := NewQuadTree(world)

	for _, v := range boxes10M {
		qt.Add(v)
	}

	queries := randomPages(b.N, world)

	b.ResetTimer()
	for _, q := range queries {
		qt.QueryPoint(q.px, q.py)
	}
}


// Benchmark simple look up on set of rectangles
func BenchmarkTreeRectsLinear(b *testing.B) {
	rand.Seed(1)
	queries := randomPages(b.N, world)

	b.ResetTimer()
	for _, q := range queries {
		queryLinear(boxes10M, q)
	}
}

// A set of 10 million randomly distributed points
var points10M []*Page

func BenchmarkTreePointsInit(b *testing.B) {
	points10M = randomPages(1000*1000, world)
}

// Benchmark quad-tree on set of points
func BenchmarkTreePointsQuadtree(b *testing.B) {
	rand.Seed(1)
	qt := NewQuadTree(world)

	for _, v := range points10M {
		qt.Add(v)
	}

	queries := randomPages(b.N, world)

	b.ResetTimer()
	for _, q := range queries {
		qt.QueryPoint(q.px, q.py)
	}
}

// Benchmark simple look-up on set of points
func BenchmarkTreePointsLinear(b *testing.B) {
	rand.Seed(1)
	queries := randomPages(b.N, world)

	b.ResetTimer()
	for _, q := range queries {
		queryLinear(points10M, q)
	}
}
