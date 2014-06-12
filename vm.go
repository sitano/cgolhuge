package main

import (
	"io"
	"fmt"
)

const (
	PageSizeWidth = 64
	PageSizeHeight = 64
	PageStrideBits = 64
	PageSizeByte = (8 * PageSizeWidth / PageStrideBits) * PageSizeHeight // 512
	PageStridePO2 = 6
	PageStrideMod = 63 // 0b111111
	PageStrideWidth = PageSizeWidth / PageStrideBits // 1
	PageStrideWidthPO2 = 0
	PageStrides = PageSizeByte / (PageStrideBits / 8)
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
	next []uint64

	px, py uint64

	// Adjacent Pages
	ap_nw, ap_n, ap_ne *Page
	ap_w,        ap_e  *Page
	ap_sw, ap_s, ap_se *Page

	// TODO: rect of life presence
	// TODO: changes since last step
	// TODO: life total count in page
	// TODO: count life on edges sep
}

type VM struct {
	reserved []*Page
}


func NewVM() *VM {
	return &VM{
		reserved: make([]*Page, 0, 16),
	}
}

func NewPageBuf() []uint64 {
   return make([]uint64, PageStrides, PageStrides)
}

func NewPage() *Page {
	return &Page{
		AABB: New00WH(PageSizeWidth, PageSizeHeight),
		raw: NewPageBuf(),
		px: 0,
		py: 0,
	}
}

func (p *Page) String() string {
	return fmt.Sprintf("(Page %d, %d / %v)", p.px, p.py, p.AABB)
}

func (vm *VM) Search(p *Page) int {
	for i, pi := range vm.reserved  {
		if pi == p {
			return i
		}
	}
	return -1
}

func (vm *VM) SearchPXY(px uint64, py uint64) int {
	for i, pi := range vm.reserved  {
		if pi.px == px && pi.py == py {
			return i
		}
	}
	return -1
}

// Take any free page or create one and put it into reserved list
func (vm *VM) ReservePage() *Page {
	p := NewPage()
	vm.reserved = append(vm.reserved, p)
	return p
}

func (vm *VM) ReclaimPage(p *Page) bool {
	i := vm.Search(p)
	if i < 0 {
		return false
	}
	vm.reserved[i], vm.reserved = vm.reserved[len(vm.reserved)-1], vm.reserved[:len(vm.reserved)-1]
	return true
}

func (vm *VM) Pages() int {
	return len(vm.reserved)
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
	sb := XtoSB(x)
	pi := XYtoPI(x, y)
	// Set
	p.raw[pi] = (p.raw[pi] & (^(uint64(1) << sb))) | (uint64(v) << sb)
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

