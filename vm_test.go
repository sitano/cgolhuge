package main

import (
	"container/list"
	"testing"
)

func TestBitsShift(t *testing.T) {
	if 2 >> 1 != 1 { t.Error("2 >> 1 != 1") }
	if 2 >> 0 != 2 { t.Error("2 >> 0 != 2") }
	if 2 >> 2 != 0 { t.Error("2 >> 2 != 0") }
	if 2 >> 3 != 0 { t.Error("2 >> 3 != 0") }
	if -8 << 1 != -16 { t.Error("-8 << 1 != -16: ", -8 << 1) }
	if -8 >> 1 != -4 { t.Error("-8 >> 1 != -4: ", -8 >> 1) }
	if -8 >> 2 != -2 { t.Error("-8 >> 2 != -2: ", -8 >> 2) }
	if -8 >> 3 != -1 { t.Error("-8 >> 3 != -1: ", -8 >> 1) }
	if -8 >> 4 != -1 { t.Error("-8 >> 4 != 1: ", -8 >> 4) }
	if -8 >> 5 != -1 { t.Error("-8 >> 5 != 1: ", -8 >> 5) }
}

func TestBits(t *testing.T) {
	if bits(0) != 0 { t.Error("bits(0) != 0, but", bits(0)) }
	if bits(2) != 1 { t.Error("bits(2) != 2, but", bits(1)) }
	if bits(3) != 2 { t.Error("bits(3) != 3, but", bits(2)) }
	if bits(4) != 1 { t.Error("bits(4) != 4, but", bits(1)) }
}

func TestLg2(t *testing.T) {
	if lg2(1) != 0 { t.Error("lg2(1) != 0, but", lg2(1)) }
	if lg2(2) != 1 { t.Error("lg2(2) != 1, but", lg2(2)) }
	if lg2(4) != 2 { t.Error("lg2(4) != 2, but", lg2(4)) }
	if lg2(8) != 3 { t.Error("lg2(8) != 3, but", lg2(8)) }
}

func TestIsPowerOf2(t *testing.T) {
	if isPowerOf2(0) { t.Error("isPowerOf2(0) != true, but", isPowerOf2(0)) }
	if !isPowerOf2(1) { t.Error("isPowerOf2(1) != false, but", isPowerOf2(1)) }
	if !isPowerOf2(2) { t.Error("isPowerOf2(2) != false, but", isPowerOf2(2)) }
	if isPowerOf2(3) { t.Error("isPowerOf2(3) != true, but", isPowerOf2(3)) }
	if !isPowerOf2(4) { t.Error("isPowerOf2(4) != false, but", isPowerOf2(4)) }
}

func TestPow2UInt64(t *testing.T) {
	if pow2ui64(0) != 1 { t.Error("pow2ui64(0) != 1, but", pow2ui64(0)) }
	if pow2ui64(1) != 2 { t.Error("pow2ui64(1) != 2, but", pow2ui64(1)) }
	if pow2ui64(2) != 4 { t.Error("pow2ui64(2) != 4, but", pow2ui64(2)) }
	if pow2ui64(3) != 8 { t.Error("pow2ui64(3) != 8, but", pow2ui64(3)) }
}

func TestSearchPage(t *testing.T) {
	l := list.New()
	p1 := Page(make([]byte, 1, 1))
	p2 := Page(make([]byte, 1, 1))
	p3 := Page(make([]byte, 1, 1))
	l.PushBack(&p1)
	l.PushBack(&p2)
	l.PushBack(&p3)
	e1 := searchPage(l, &p1)
	e2 := searchPage(l, &p2)
	e3 := searchPage(l, &p3)
	if (e1 == nil) {
		t.Error("Can't find page 1")
	}
	if (e2 == nil) {
		t.Error("Can't find page 2")
	}
	if (e3 == nil) {
		t.Error("Can't find page 3")
	}
}

func TestVM(t *testing.T) {
	vm := NewVM(KSIZE_16K)

	if vm.Pages() != 0 { t.Error("No pages") }
	if vm.Reclaimed() != 0 { t.Error("No reclaimed") }
	if vm.Reserved() != 0 { t.Error("No reserved") }

	p1 := vm.ReservePage()

	if vm.Pages() != 1 { t.Error("Pages 1") }
	if vm.Reclaimed() != 0 { t.Error("No reclaimed") }
	if vm.Reserved() != 1 { t.Error("Reserved 1") }
	if p1 == nil { t.Error("Reserve failed") }
	if len(*p1) != KSIZE_16K { t.Error("Invalid page size") }

	vm.ReclaimPage(p1)

	if vm.Pages() != 1 { t.Error("Pages 1") }
	if vm.Reclaimed() != 1 { t.Error("Reclaimed 1") }
	if vm.Reserved() != 0 { t.Error("Reserved 0") }

	p2 := vm.ReservePage()

	if vm.Pages() != 1 { t.Error("Pages 1") }
	if vm.Reclaimed() != 0 { t.Error("No reclaimed") }
	if vm.Reserved() != 1 { t.Error("Reserved 1") }
	if p1 != p2 { t.Error("Reserve reclaimed failed") }

	p3 := vm.ReservePage()

	if vm.Pages() != 2 { t.Error("Pages 2") }
	if vm.Reclaimed() != 0 { t.Error("No reclaimed") }
	if vm.Reserved() != 2 { t.Error("Reserved 2") }
	if p3 == nil { t.Error("Reserve failed") }

	vm.ReclaimPage(p1)
	vm.ReclaimPage(p3)

	if vm.Pages() != 2 { t.Error("Pages 2") }
	if vm.Reclaimed() != 2 { t.Error("Reclaimed 2") }
	if vm.Reserved() != 0 { t.Error("Reserved 0") }
}

func TestPageWidth(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	if (vm.PageWidth() != 128) {
		t.Error("Invalid page side calc for ps = 16k")
	}
}

func TestPageLen(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	p := vm.ReservePage()
	r := ([]byte)(*p)
	if len(r) != KSIZE_16K {
		t.Error("New page have invalid size")
	}
}

func TestPageMustBeClean(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	p := vm.ReservePage()
	r := ([]byte)(*p)
	r[0] = 0xff
	vm.ReclaimPage(p)
	p2 := vm.ReservePage()
	r2 := ([]byte)(*p2)
	if r2[0] != 0 {
		t.Error("Reserved page after reclaimation must be clean")
	}
}
