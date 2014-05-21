package main

import "io"
import "fmt"

type View interface {
	GetAABB() AABB
	Get(x uint64, y uint64) byte
	Set(x uint64, y uint64, v byte)
}

type ViewUtil interface {
	Print(b AABB) string
	Match(b AABB, matcher []byte) bool

	MirrorH(b AABB)
	MirrorV(b AABB)

	Writer(b AABB) io.Writer
	Reader(b AABB) io.Reader
}

type ViewIO struct {
	io.Writer
	io.Reader

	v View

	b AABB

	lx uint64
	ly uint64

	err error
}

type WorldView struct {
	// View
	// ViewIO

	vm *VM

	// pb *PageTree
}

func Print(v View, b AABB) string {
	r := ""

	if ! v.GetAABB().Intersects(b) {
		return r
	}

	bbview := v.GetAABB().Intersection(b)
	for iy := bbview.MinY; iy <= bbview.MaxY; iy ++ {
		for ix := bbview.MinX; ix <= bbview.MaxX; ix ++ {
			val := v.Get(ix, iy)
			if val == 0 {
				r += "."
			} else {
				r += "@"
			}
		}

		r += "\n"
	}

	return r
}

func Match(v View, b AABB, matcher []byte) bool {
	if ! v.GetAABB().Intersects(b) {
		return len(matcher) == 0
	}

	ii := 0
	bbview := v.GetAABB().Intersection(b)
	for iy := bbview.MinY; iy <= bbview.MaxY; iy ++ {
		for ix := bbview.MinX; ix <= bbview.MaxX; ix ++ {
			// fmt.Printf("x=%v, y=%v, v1=%v, v2=%v\n", ix, iy, v.Get(ix, iy), matcher[ii])
			if v.Get(ix, iy) != matcher[ii] {
				return false
			}
			ii ++
		}
	}

	return true
}

func MirrorH(v View, b AABB) {
	if ! v.GetAABB().Intersects(b) {
		return
	}

	bbview := v.GetAABB().Intersection(b)
	for iy := bbview.MinY; iy <= bbview.MaxY; iy ++ {
		for ix := bbview.MinX; ix <= bbview.MaxX / 2; ix ++ {
			ix2 := bbview.MaxX - (ix - bbview.MinX)
			val := v.Get(ix2, iy)
			v.Set(ix2, iy, v.Get(ix, iy))
			v.Set(ix2, iy, val)
		}
	}
}

func MirrorV(v View, b AABB) {
	if ! v.GetAABB().Intersects(b) {
		return
	}

	bbview := v.GetAABB().Intersection(b)
	for ix := bbview.MinX; ix <= bbview.MaxX; ix ++ {
		for iy := bbview.MinY; iy <= bbview.MaxY / 2; iy ++ {
			iy2 := bbview.MaxY - (iy - bbview.MinY)
			val := v.Get(ix, iy2)
			v.Set(ix, iy2, v.Get(ix, iy))
			v.Set(ix, iy2, val)
		}
	}
}

func Writer(v View, b AABB) io.Writer {
	if ! v.GetAABB().Intersects(b) {
		panic(fmt.Sprintf("Writer failed: %v does not intersect %v", v.GetAABB(), b))
	}
	bbox := v.GetAABB().Intersection(b)
	return &ViewIO{
		v: v,
		b: bbox,
		lx: bbox.MinX,
		ly: bbox.MinY,
		err: nil,
	}
}

func Reader(v View, b AABB) io.Reader {
	if ! v.GetAABB().Intersects(b) {
		panic(fmt.Sprintf("Reader failed: %v does not intersect %v", v.GetAABB(), b))
	}
	bbox := v.GetAABB().Intersection(b)
	return &ViewIO{
		v: v,
		b: bbox,
		lx: bbox.MinX,
		ly: bbox.MinY,
		err: nil,
	}
}

func (v *ViewIO) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, v.err
	}
	ii := 0
	for v.ly <= v.b.MaxY {
		for ; v.lx <= v.b.MaxX && ii < len(p); v.lx ++ {
			p[ii] = v.v.Get(v.lx, v.ly)
			ii ++
		}
		if v.lx > v.b.MaxX {
			v.lx = v.b.MinX
			v.ly ++
		}
		if ii == len(p) {
			break
		}
	}
	if v.err == nil && v.ly > v.b.MaxY {
		v.err = io.EOF
	}
	return ii, v.err
}

func (v *ViewIO) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, v.err
	}
	ii := 0
	for v.ly <= v.b.MaxY {
		for ; v.lx <= v.b.MaxX && ii < len(p); v.lx ++ {
			v.v.Set(v.lx, v.ly, p[ii])
			ii ++
		}
		if v.lx > v.b.MaxX {
			v.lx = v.b.MinX
			v.ly ++
		}
		if ii == len(p) {
			break
		}
	}
	if v.err == nil && v.ly > v.b.MaxY {
		v.err = io.EOF
	}
	return ii, v.err
}

/*
func NewWorldView(vm *VM, pb *PageTree) WorldView {
	if vm.wsize != pb.wsize {
		panic("VM wsize must match PageTree wsize")
	}
	return WorldView{vm, pb, true}
}

func (wv WorldView) Set(x int64, y int64, z byte, t byte) {
	wsize := wv.pb.wsize

	// Page coord
	px := WtoP(x, wsize)
	py := WtoP(y, wsize)

	// Page Tile
	pt := wv.pb.QueryPage(px, py)

	// Skip casual case
	if pt == nil && t == DEAD {
		return
	}

	// Reserve page if life needed
	if pt == nil && t == LIFE {
		np := NewPageTile(wv.vm.ReservePage(), wsize, px, py)
		pt = &np
		wv.pb.Add(pt)
	}

	// Coord inside of page
	pboffset := WPtoPO(x, y, px, py, wsize)

	// Get data
	data := pt.GetByte(pboffset)
	state := ReadStateZ(data, z)

	if state != t {
		data = WriteStateZ(data, z, t)
		pt.SetByte(pboffset, data)

		if wv.autoReclaim {
			switch t {
			case DEAD:
				pt.SetAlive(pt.GetAlive() - 1)
			case LIFE:
				pt.SetAlive(pt.GetAlive() + 1)
			}

			wv.TryReclaim(pt)
		}
	}
}

func (wv *WorldView) TryReclaim(pt *PageTile) {
	if pt != nil && pt.GetAlive() == 0 {
		wv.pb.Remove(pt)
		wv.vm.ReclaimPage(pt.p)
		pt.p = nil
		pt = nil
	}
}

func (wv *WorldView) Get(x int64, y int64, z byte) byte {
	wsize := wv.pb.wsize

	// Page coord
	px := WtoP(x, wsize)
	py := WtoP(y, wsize)

	// Page Tile
	pt := wv.pb.QueryPage(px, py)

	// Skip casual case
	if pt == nil {
		return DEAD
	}

	// Coord inside of page
	pboffset := WPtoPO(x, y, px, py, wsize)

	// Get data
	return ReadStateZ(pt.GetByte(pboffset), z)
}

func WPtoPO(x int64, y int64, px int64, py int64, wsize uint) uint {
	pxoffset := Abs(x - px * int64(wsize))
	pyoffset := Abs(y - py * int64(wsize))
	pydiff   := uint64(wsize) - pyoffset - 1
	return uint(pydiff * uint64(wsize) + pxoffset)
}

func POtoWX(offset uint, px int64, wsize uint) int64 {
	return px * int64(wsize) + int64(offset % wsize)
}

func POtoWY(offset uint, py int64, wsize uint) int64 {
	return py * int64(wsize) + int64(wsize) - int64(offset / wsize) - 1
}

func ClearStateZ(b byte, z byte) byte {
	return b & ^(ZMASK << (z * ZSTEP))
}

func WriteStateZ(b byte, z byte, v byte) byte {
	b = ClearStateZ(b, z)
	return b | (v << (z * ZSTEP))
}

func ReadStateZ(b byte, z byte) byte {
	return (b >> (z * ZSTEP)) & ZMASK
}

// dx, dy = + / - 1
func (vw *WorldView) NextTo(x int64, y int64, z byte, dx int64, dy int64) byte {
	bb := vw.pb.GetAABB()
	return vw.Get(MvXY1(x, dx, bb.MinX, bb.MaxX), MvXY1(y, dy, bb.MinY, bb.MaxY), z)
}

// dx, dy = + / - 1
func MvXY1(x int64, dx int64, min int64, max int64) int64 {
	nx := x + dx

	if max > 0 && max + 1 < 0 {
		if x >= max && dx > 0 {
			nx = min
		}
	} else {
		if x >= max - 1 && dx > 0 {
			nx = min
		}
	}

	if x <= min && dx < 0 {
		if max > 0 && max + 1 < 0 {
			nx = max
		} else {
			nx = max - 1
		}
	}

	return nx
}

func (vw *WorldView) LifeSumAt(x int64, y int64, z byte) byte {
	// Just sum them up as DEAD = 0, LIFE = 1
	return vw.NextTo(x, y, z, -1, +1) + vw.NextTo(x, y, z, 0, +1) + vw.NextTo(x, y, z, +1, +1) +
		vw.NextTo(x, y, z, -1, 0) + 0 +                           vw.NextTo(x, y, z, +1, 0) +
		vw.NextTo(x, y, z, -1, -1) + vw.NextTo(x, y, z, 0, -1) + vw.NextTo(x, y, z, +1, -1)
}
*/
