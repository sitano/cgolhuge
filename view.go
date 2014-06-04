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
	AABB
	View
	ViewUtil

	vm *VM
	pb *QuadTree
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

func NewWorldView(vm *VM, bbox AABB) *WorldView {
	return &WorldView{
		AABB: bbox,
		vm: vm,
		pb: NewQuadTree(bbox),
	}
}

// View implementation

func (wv *WorldView) GetAABB() AABB {
	return wv.AABB
}

func (wv *WorldView) Get(x uint64, y uint64) byte {
	px, py := WXY2PXY(x, y)
	p := wv.pb.QueryPoint(px, py)
	if p == nil { return DEAD }
	return p.Get(x - p.MinX, y - p.MinY)
}

func (wv *WorldView) Set(x uint64, y uint64, v byte) {
	px, py := WXY2PXY(x, y)
	p := wv.pb.QueryPoint(px, py)
	if p == nil {
	    p = wv.vm.ReservePage()
		wv.pb.AddTo(p, px, py)
	}
	p.Set(x - p.MinX, y - p.MinY, v)
}

// ViewUtil implementation

func (wv *WorldView) Print(b AABB) string {
	return Print(wv, b)
}

func (wv *WorldView) Match(b AABB, matcher []byte) bool {
	return Match(wv, b, matcher)
}

func (wv *WorldView) MirrorH(b AABB) {
	MirrorH(wv, b)
}

func (wv *WorldView) MirrorV(b AABB) {
	MirrorV(wv, b)
}

func (wv *WorldView) Writer(b AABB) io.Writer {
	return Writer(wv, b)
}

func (wv *WorldView) Reader(b AABB) io.Reader {
	return Reader(wv, b)
}
