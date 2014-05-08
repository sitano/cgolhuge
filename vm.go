package main

import (
	"container/list"
)

const (
	// Page of size 2**7 * 2**7 = 128 * 128 = 16kb
	KSIZE_16K = 14
)

type Page []uint8

type VM struct {
	// Power of 2 to set page size
	ksize uint

	reserved *list.List
	reclaimed *list.List
}

func NewVM(ksize uint) VM {
	return VM{
		ksize: ksize,
		reserved: list.New(),
		reclaimed: list.New(),
	}
}

func (vm VM) NewPage() Page {
	size := pow2ui64(vm.ksize)
	return Page(make([]uint8, size, size))
}

func searchPage(l *list.List, p *Page) *list.Element {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*Page) == p {
			return e
		}
	}
	return nil
}

func pow2ui64(n uint) uint {
	if n > 0 {
		return 2 << (n - 1)
	}

	return 1
}

// Take any free page or create one and put it into reserved list
func (vm VM) ReservePage() *Page {
	var p *Page

	if vm.reclaimed.Len() > 0 {
		el := vm.reclaimed.Front()
		vm.reclaimed.Remove(el)
		p = el.Value.(*Page)
	} else {
		np := vm.NewPage()
		p = &np
	}

	vm.reserved.PushBack(p)

	return p
}

func (vm VM) ReclaimPage(p *Page) bool {
	el := searchPage(vm.reserved, p)
	if el == nil {
		return false
	}
	vm.reserved.Remove(el)
	vm.reclaimed.PushBack(p)
	return true
}

func (vm VM) Pages() int {
	return vm.reserved.Len() + vm.reclaimed.Len()
}

func (vm VM) Reserved() int {
	return vm.reserved.Len()
}

func (vm VM) Reclaimed() int {
	return vm.reclaimed.Len()
}

func (vm VM) PageSize() uint {
	return pow2ui64(vm.ksize)
}
