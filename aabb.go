package main

import (
	"fmt"
	"math"
)

// Min - inclusive, Max - inclusive
type AABB struct {
	MinX, MinY uint64
	MaxX, MaxY uint64
}

func (b AABB) String() string {
	return fmt.Sprintf("(AABB %d, %d -> %d, %d)", b.MinX, b.MinY, b.MaxX, b.MaxY)
}

func Min(x uint64, y uint64) uint64 {
	if x <= y { return x }
	return y
}

func Max(x uint64, y uint64) uint64 {
	if x >= y { return x }
	return y
}

func NewAABB(xa, ya, xb, yb uint64) AABB {
	return AABB{ Min(xa, xb), Min(ya, yb), Max(xa, xb), Max(ya, yb) }
}

func NewXYWH(x, y, w, h uint64) AABB {
	return AABB{ x, y, x + w - 1, y + h - 1 }
}

func New00WH(w, h uint64) AABB {
	return AABB{ 0, 0, w - 1, h - 1 }
}

func NewAABBMax() AABB {
	return NewAABB(0, 0, math.MaxUint64, math.MaxUint64)
}

func (b AABB) SizeX() uint64 {
	return b.MaxX - b.MinX + 1
}

func (b AABB) SizeY() uint64 {
	return b.MaxY - b.MinY + 1
}

func (b AABB) Intersects(o AABB) bool {
	return b.MinX <= o.MaxX && b.MinY <= o.MaxY && b.MaxX >= o.MinX && b.MaxY >= o.MinY
}

func (b AABB) Intersection(o AABB) AABB {
	return NewAABB(Max(b.MinX, o.MinX), Max(b.MinY, o.MinY),
		Min(b.MaxX, o.MaxX), Min(b.MaxY, o.MaxY))
}

func (b AABB) Contains(o AABB) bool {
	return b.MinX <= o.MinX && b.MinY <= o.MinY && b.MaxX >= o.MaxX && b.MaxY >= o.MaxY
}

func (b AABB) ContainsPoint(x uint64, y uint64) bool {
	return b.MinX <= x && b.MinY <= y && b.MaxX >= x && b.MaxY >= y
}
