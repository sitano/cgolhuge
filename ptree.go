package main

import "math"
import (
	"fmt"
)

type PageTree struct {
	QuadTree

	wsize uint
}

type PageTile struct {
	AABB

	// Page from vm
	p *Page

	// Page position in page-world coords
	px int64
	py int64

	// Game
	alive uint
}

// Make AABB implement the QuadElement interface
func (pt *PageTile) GetAABB() AABB {
	return pt.AABB
}

func (b PageTree) String() string {
	return fmt.Sprint(b.AABB)
}

func NewPageTree(bbox AABB, wsize uint) PageTree {
	szx := bbox.SizeX()
	szy := bbox.SizeY()
	if (szx != math.MaxUint64 && szx % uint64(wsize) != 0) ||
		(szy != math.MaxUint64 && szy % uint64(wsize) != 0) {
		panic("NewPageTree: bbox size does not fit page wsize")
	}

	return PageTree{NewQuadTree(bbox), wsize}
}

func NewPageTile(p *Page, wsize uint, px int64, py int64) PageTile {
	return PageTile{
		AABB: NewAABBPXY(px, py, wsize),
		p: p,
		px: px,
		py: py,
		alive: 0,
	}
}

func NewAABBPXY(px int64, py int64, wsize uint) AABB {
	bbox := NewAABB(
		PtoW(px, wsize), PtoW(px + 1, wsize),
		PtoW(py, wsize), PtoW(py + 1, wsize))
	// Check overflow
	if px > 0 && bbox.MinX < 0 && bbox.MaxX > 0 {
		maxX := int64(math.MaxInt64)
		bbox = NewAABB(bbox.MaxX, maxX, bbox.MinY, bbox.MaxY)
	}
	if py > 0 && bbox.MinY < 0 && bbox.MaxY > 0 {
		maxY := int64(math.MaxInt64)
		bbox = NewAABB(bbox.MinX, bbox.MaxX, bbox.MaxY, maxY)
	}
	return bbox
}

func PtoW(pc int64, wsize uint) int64 {
	return pc * int64(wsize)
}

func WtoP(wc int64, wsize uint) int64 {
	if (wc >= 0) {
		return wc / int64(wsize)
	}

	return (wc + 1) / int64(wsize) - 1
}

func (pb *PageTree) Add(pt *PageTile) {
	// TODO: test page wsize, aabb, etc
	pb.QuadTree.Add(pt)
}

func (pb *PageTree) Remove(pt *PageTile) bool {
	return pb.RemovePXY(pt.px, pt.py)
}

func (pb *PageTree) RemovePXY(px int64, py int64) bool {
	if pb.root.remove(px, py, PtoW(px, pb.wsize), PtoW(py, pb.wsize)) {
		pb.count --
		return true
	}
	return false
}

func (pb *PageTree) Count() uint64 {
	return pb.count
}

func (pb *PageTree) GetAABB() *AABB {
	return &pb.AABB
}

func (pb *PageTree) QueryPage(px int64, py int64) *PageTile {
	return pb.root.queryPage(px, py, PtoW(px, pb.wsize), PtoW(py, pb.wsize))
}

func (qb *PageTree) Reduce(f func(a interface{}, t *PageTile) interface{}, v interface{}) interface{} {
	return qb.root.reduce(func(a interface{}, t QuadElement) interface{} {
		return f(a, t.(*PageTile))
	}, v)
}

func (tile *QuadTile) queryPage(px int64, py int64, x int64, y int64) *PageTile {
	// end recursion if this tile does not intersect the query range
	if ! tile.ContainsPoint(x, y) {
		return nil
	}

	// return candidates at this tile
	for _, v := range tile.contents {
		p := v.(*PageTile)
		if p.px == px && p.py == py {
			return p
		}
	}

	// recurse into childs (if any)
	if tile.childs[_TOPRIGHT] != nil {
		for _, child := range tile.childs {
			ret := child.queryPage(px, py, x, y)
			if ret != nil { return ret }
		}
	}

	return nil
}

func (tile *QuadTile) remove(px int64, py int64, x int64, y int64) bool {
	// end recursion if this tile does not intersect the query range
	if ! tile.ContainsPoint(x, y) {
		return false
	}

	// return candidates at this tile
	for i, v := range tile.contents {
		p := v.(*PageTile)
		if p.px == px && p.py == py {
			tile.contents[i] = tile.contents[len(tile.contents)-1]
			tile.contents = tile.contents[0:len(tile.contents)-1]
			// TODO: merge parent tree node childs if can, but i dont need it RN
			return true
		}
	}

	// recurse into childs (if any)
	if tile.childs[_TOPRIGHT] != nil {
		for _, child := range tile.childs {
			ret := child.remove(px, py, x, y)
			if ret { return ret }
		}
	}

	return false
}

func (pb *PageTree) MaxPagesX() uint64 {
	bbox := pb.AABB
	ws := uint64(pb.wsize)
	return Abs(bbox.MinX) / ws + Abs(bbox.MaxX) / ws
}

func (pb *PageTree) MaxPagesY() uint64 {
	bbox := pb.AABB
	ws := uint64(pb.wsize)
	return Abs(bbox.MinY) / ws + Abs(bbox.MaxY) / ws
}

func (pt *PageTile) SetByte(i uint, v byte) {
	(*pt.p)[i] = v
}

func (pt *PageTile) GetByte(i uint) byte {
	return (*pt.p)[i]
}

func (pt *PageTile) GetAlive() uint {
	return pt.alive
}

func (pt *PageTile) SetAlive(alive uint) {
	pt.alive = alive
}
