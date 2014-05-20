package main

import (
	"container/list"
	"io"
)

const (
	PageSizeWidth = 64
	PageSizeHeight = 64
	PageSizeByte = PageSizeWidth * PageSizeHeight // 4096
	PageStrideByte = 64
	PageStridePO2 = 6
	PageStrideMod = 63 // 0b111111
	PageStrideWidth = PageSizeWidth / PageStrideByte // 1
	PageStrideWidthPO2 = 0
	PageStrideHeight = PageSizeHeight // 64
	PageStrideSize = PageStrideWidth * PageStrideHeight // 64
)

// Page coordinates
//   y
// x 0 1 2 3 4
//   1 b b b b
//   2 b b b b
//   3 b b b b
//   4 b b b b

// Life packing 1 bit per life
// uint64 = 8 byte
// So, x = 0, byte & 1
//     x = 1, byte >> 1 ) & 1
// x = ...6543210
// v = ...0000000

type Page struct {
	AABB
	View
	ViewUtil

	raw []uint64

	px, py uint64

	// TODO: adjacent pages
	// TODO: rect of life presence
	// TODO: changes since last step
	// TODO: life total count in page
	// TODO: count life on edges sep
}

type VM struct {
	reserved *list.List
}


func NewVM() *VM {
	return &VM{
		reserved: list.New(),
	}
}

func NewPage() *Page {
	return &Page{
		raw: make([]uint64, PageStrideSize, PageStrideSize),
		px: 0,
		py: 0,
	}
}

func (vm *VM) Search(p *Page) *list.Element {
	for e := vm.reserved.Front(); e != nil; e = e.Next() {
		if e.Value.(*Page) == p {
			return e
		}
	}
	return nil
}

func (vm *VM) SearchPXY(px uint64, py uint64) *list.Element {
	for e := vm.reserved.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Page)
		if p.px == px && p.py == py {
			return e
		}
	}
	return nil
}

// Take any free page or create one and put it into reserved list
func (vm *VM) ReservePage() *Page {
	p := NewPage()
	vm.reserved.PushBack(p)
	return p
}

func (vm *VM) ReclaimPage(p *Page) bool {
	el := vm.Search(p)
	if el == nil {
		return false
	}
	vm.reserved.Remove(el)
	return true
}

func (vm *VM) Pages() int {
	return vm.reserved.Len()
}

func XtoPX(x uint64) uint64 {
	return x >> PageStridePO2
}

func YtoPY(y uint64) uint64 {
	return y
}

func XYtoPI(x uint64, y uint64) uint64 {
	return uint64((y << PageStrideWidthPO2) + (x >> PageStridePO2))
}

func XtoSB(x uint64) uint64 {
	return x & PageStrideMod
}

// View implementation

func (p *Page) GetAABB() AABB {
	return p.AABB
}

func (p *Page) Get(x uint64, y uint64) byte {
	return byte((p.raw[XYtoPI(x, y)] >> XtoSB(x)) & 0x1)
}

func (p *Page) Set(x uint64, y uint64, v byte) {
	// Mask
	sb := XtoSB(x)
	mask := ^(uint64(1) << sb)
	// Unset
	pi := XYtoPI(x, y)
	a := p.raw[pi] & mask
	// Set
	mask = uint64(v) << sb
	p.raw[pi] = a | mask
}

// ViewUtil implementation

func (p *Page) Print(b AABB) string {
	return Print(p, b)
}

func (p *Page) Match(b AABB, matcher []byte) bool {
	return Match(p, b, matcher)
}

func (p *Page) MirrorH(b AABB) {
	MirrorH(p, b)
}

func (p *Page) MirrorV(b AABB) {
	MirrorV(p, b)
}

func (p *Page) Writer(b AABB) io.Writer {
	return Writer(p, b)
}

func (p *Page) Reader(b AABB) io.Reader {
	return Reader(p, b)
}

