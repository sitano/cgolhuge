package main

const (
	DEAD = byte(0)
	LIFE = byte(1)
	ZSTEP = 2
	ZMASK = 0x3
	ZMAX  = 3
)

type View interface {
	// io.Writer
	// io.Reader

	Set(x int64, y int64, z byte, t byte)
	Get(x int64, y int64, z byte) byte

	NextTo(x int64, y int64, z byte, dx int64, dy int64) byte
	LifeSumAt(x int64, y int64, z byte) byte
}

type WorldView struct {
	vm *VM
	pb *PageTree
	autoReclaim bool
}

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
