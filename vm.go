package main

import (
	"container/list"
)

const (
	PageSizeWidth = 64
	PageSizeHeight = 64
	PageSizeByte = PageSizeWidth * PageSizeHeight // 4096
	PageStrideByte = 8
	PageStrideWidth = PageSizeWidth / PageStrideByte // 8
	PageStrideHeight = PageSizeHeight / PageStrideByte // 8
	PageStrideSize = PageSizeWidth * PageStrideHeight // 64
)

// Page coordinates
// 0 1 2 3 4 | y=0
// 1 b b b b | y=1
// 2 b b b b | y=2
// 3 b b b b | y=2
// 4 b b b b | y=3

type Page struct {
	View

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
