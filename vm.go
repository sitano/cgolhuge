package main

import (
	"container/list"
	"fmt"
)

const (
	// Page of size 2**7 * 2**7 = 128 * 128 = 16kb
	KSIZE_16K = 128 * 128
)

type Page []byte

type VM struct {
	// Page size in bytes
	ksize uint
	// Width == Height of the page
	wsize uint
	wsbits uint

	reserved *list.List
}

// lg(ksize) must be div by 2
func NewVM(ksize uint) VM {
	if (!isPowerOf2(ksize) || lg2(ksize) & 1 != 0) {
		panic(fmt.Sprintf("Page size must fit equal w/h sizes, sz=%d, lg2=%d", ksize, lg2(ksize)))
	}
	ws := pow2ui64(lg2(ksize) >> 1)
	return VM{
		ksize: ksize,
		wsize: ws,
		wsbits: bits(ws),
		reserved: list.New(),
	}
}

func (vm VM) NewPage() Page {
	return make([]byte, vm.ksize, vm.ksize)
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
	p := vm.NewPage()
	pp:= &p
	vm.reserved.PushBack(pp)
	return pp
}

func (vm VM) ReclaimPage(p *Page) bool {
	el := searchPage(vm.reserved, p)
	if el == nil {
		return false
	}
	vm.reserved.Remove(el)
	return true
}

func (vm VM) Pages() int {
	return vm.reserved.Len()
}

func (vm VM) Reserved() int {
	return vm.reserved.Len()
}

func (vm VM) PageSize() uint {
	return pow2ui64(vm.ksize)
}

func bits(v uint) uint {
	k := uint(0)
	t := int64(v)
	for t != 0 {
		t &= t - 1
		k ++
	}
	return k
}

func isPowerOf2(v uint) bool {
	return v != 0 && v & (v - 1) == 0
}

func lg2(v uint) uint {
	k := uint(0)
	for ; v > 0 ; v >>= 1 {
		k++
	}
	return k - 1
}

func (vm VM) PageWidth() uint {
	return vm.wsize
}
