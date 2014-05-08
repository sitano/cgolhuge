package main

import (
	"container/list"
	"testing"
)

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
