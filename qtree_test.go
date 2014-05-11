/*
Based on work of Volker Poplawski, 2013 (https://github.com/volkerp/goquadtree)
*/
package main

import "math/rand"
import "testing"
import _ "fmt"

func TestAABB(t *testing.T) {
	a := NewAABB( 0, 10, 0, 10 )

	b := NewAABB( 4, 6, 4, 6 )    // b completely within a

	if ! a.Intersects(b) || ! b.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, b)
	}

	if ! a.Intersects(a) {
		t.Errorf("%v does not intersect itself", a)
	}

	if ! a.Contains(b) {
		t.Errorf("%v does not contain %v", a, b)
	}

	if ! a.Contains(a) {
		t.Errorf("%v does not contain itself", a)
	}

	if b.Contains(a) {
		t.Errorf("%v contains %v", b, a)
	}

	c := NewAABB( 10, 20, 0, 10 )

	if a.Intersects(c) {
		t.Errorf("%v does intersect %v", a, c)
	}
	if c.Intersects(a) {
		t.Errorf("%v does intersect %v", c, a)
	}

	if a.Contains(c) || c.Contains(a) {
		t.Errorf("%v contains %v (or vise versa)", a, c)
	}

	d := NewAABB( -10, 0, 0, 10 )

	if a.Intersects(d) {
		t.Errorf("%v does intersect %v", a, d)
	}
	if d.Intersects(a) {
		t.Errorf("%v does intersect %v", d, a)
	}

	e := NewAABB( 9, 15, 9, 15 )

	if ! a.Intersects(e) || ! e.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, e)
	}

	f := NewAABB( -10, 20, 4, 6 )

	if  ! a.Intersects(f) || ! f.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, f)
	}
}

// Generates n AABBes in the range of frame with average width and height avgSize
func randomAABBes(n int, frame AABB, avgSize int64) []AABB {
	ret := make([]AABB, n)

	for i := 0; i < len(ret); i++ {
		w := int64(rand.NormFloat64() * float64(avgSize))
		h := int64(rand.NormFloat64() * float64(avgSize))
		x := int64(rand.Float64() * float64(frame.SizeX()) + float64(frame.MinX))
		y := int64(rand.Float64() * float64(frame.SizeY()) + float64(frame.MinY))
		ret[i] = NewAABB(x, Min(frame.MaxX, x+w), y, Min(frame.MaxY, y+h))
	}

	return ret
}


// Returns all elements of data which intersect query
func queryLinear(data []AABB, query AABB) (ret []QuadElement) {
	for _, v := range data {
		if query.Intersects(v.GetAABB()) {
			ret = append(ret, v)
		}
	}

	return ret
}


func compareQuadElement(v1, v2 QuadElement) bool {
	b1 := v1.GetAABB()
	b2 := v2.GetAABB()

	return b1.MinX == b2.MinX && b1.MaxX == b2.MaxX &&
		b1.MinY == b2.MinY && b2.MaxY == b2.MaxY
}


func lookupResults(r1, r2 []QuadElement) int {
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

// World-space extends from -1000..1000 in X and Y direction
var world AABB = NewAABB(-1000, 1000, -1000, 1000)


// Compary correctness of quad-tree results vs simple look-up on set of random rectangles
func TestQuadTreeRects(t *testing.T) {
	var rects []AABB = randomAABBes(100, world, 5)
	qt := NewQuadTree(world)

	for _, v := range rects {
		qt.Add(v)
	}

	queries := randomAABBes(100, world, 10)

	for _, q := range queries {
		r1 := queryLinear(rects, q)
		r2 := qt.Query(q)

		if len(r1) != len(r2) {
			t.Errorf("r1 and r2 differ: %v   %v\n", r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("%v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("%v was not returned by brute-force\n", r2[i])
		}

	}
}


// Compary correctness of quad-tree results vs simple look-up on set of random points
func TestQuadTreePoints(t *testing.T) {
	var points []AABB = randomAABBes(100, world, 0)
	qt := NewQuadTree(world)

	for _, v := range points {
		qt.Add(v)
	}

	queries := randomAABBes(100, world, 10)

	for _, q := range queries {
		r1 := queryLinear(points, q)
		r2 := qt.Query(q)

		if len(r1) != len(r2) {
			t.Errorf("r1 and r2 differ: %v   %v\n", r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("%v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("%v was not returned by brute-force\n", r2[i])
		}

	}
}

func TestQuadTreeReduce(t *testing.T) {
	var rects []AABB = randomAABBes(100, world, 5)
	qt := NewQuadTree(world)

	for _, v := range rects {
		qt.Add(v)
	}

	count := qt.Reduce(func(c interface{}, e QuadElement) interface{} {
		return c.(int) + 1
	}, 0).(int)

	if count != 100 {
		t.Error("Reduce wrong abount 100 inner rects != ", count)
	}
}


// A set of 10 million randomly distributed rectangles of avg size 5
var boxes10M []AABB

func BenchmarkInitBoxes(b *testing.B) {
	boxes10M = randomAABBes(10*1000*1000, world, 5)
}

// Benchmark insertion into quad-tree
func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	var values []AABB = randomAABBes(b.N, world, 5)
	qt := NewQuadTree(world)

	b.StartTimer()

	for _, v := range values {
		qt.Add(v)
	}
}


// Benchmark quad-tree on set of rectangles
func BenchmarkRectsQuadtree(b *testing.B) {
	b.StopTimer()
	rand.Seed(1)
	qt := NewQuadTree(world)

	for _, v := range boxes10M {
		qt.Add(v)
	}

	queries := randomAABBes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		qt.Query(q)
	}
}


// Benchmark simple look up on set of rectangles
func BenchmarkRectsLinear(b *testing.B) {
	b.StopTimer()
	rand.Seed(1)
	queries := randomAABBes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		queryLinear(boxes10M, q)
	}
}

// A set of 10 million randomly distributed points
var points10M []AABB

func BenchmarkInitPoints(b *testing.B) {
	points10M = randomAABBes(10*1000*1000, world, 0)
}

// Benchmark quad-tree on set of points
func BenchmarkPointsQuadtree(b *testing.B) {
	b.StopTimer()
	rand.Seed(1)
	qt := NewQuadTree(world)

	for _, v := range points10M {
		qt.Add(v)
	}

	queries := randomAABBes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		qt.Query(q)
	}
}


// Benchmark simple look-up on set of points
func BenchmarkPointsLinear(b *testing.B) {
	b.StopTimer()
	rand.Seed(1)
	queries := randomAABBes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		queryLinear(points10M, q)
	}
}
