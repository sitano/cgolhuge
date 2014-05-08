/*
Based on work of Volker Poplawski, 2013 (https://github.com/volkerp/goquadtree)
*/
package main

// Use AABB() to construct a AABB object
type AABB struct {
	MinX, MaxX, MinY, MaxY int64
}

func Min(x int64, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

func Max(x int64, y int64) int64 {
	if x > y {
		return x
	}

	return y
}

func NewAABB(xa, xb, ya, yb int64) AABB {
	return AABB{ Min(xa, xb), Max(xa, xb), Min(ya, yb), Max(ya, yb) }
}

// Make AABB implement the AABBer interface
func (b AABB) AABB() AABB {
	return b
}

func (b AABB) SizeX() int64 {
	return b.MaxX - b.MinX
}

func (b AABB) SizeY() int64 {
	return b.MaxY - b.MinY
}

// Returns true if o intersects this
func (b AABB) Intersects(o AABB) bool {
	return b.MinX < o.MaxX && b.MinY < o.MaxY &&
		b.MaxX > o.MinX && b.MaxY > o.MinY
}

// Returns true if o is within this
func (b AABB) Contains(o AABB) bool {
	return b.MinX <= o.MinX && b.MinY <= o.MinY &&
		b.MaxX >= o.MaxX && b.MaxY >= o.MaxY
}
