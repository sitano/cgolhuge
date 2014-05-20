package main

import "testing"
import "io"
import "bytes"

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

func TestVM(t *testing.T) {
	vm := NewVM()

	if vm.Pages() != 0 { t.Error("No pages") }

	p1 := vm.ReservePage()

	if vm.Pages() != 1 { t.Error("Pages 1") }
	if p1 == nil { t.Error("Reserve failed") }
	if len(p1.raw) != PageStrideSize { t.Error("Invalid page size") }

	vm.ReclaimPage(p1)

	if vm.Pages() != 0 { t.Error("Pages 1") }

	p1 = vm.ReservePage()

	if vm.Pages() != 1 { t.Error("Pages 1") }

	p3 := vm.ReservePage()

	if vm.Pages() != 2 { t.Error("Pages 2") }
	if p3 == nil { t.Error("Reserve failed") }

	vm.ReclaimPage(p1)
	vm.ReclaimPage(p3)

	if vm.Pages() != 0 { t.Error("Pages 2") }
}

func TestPageMustBeClean(t *testing.T) {
	vm := NewVM()
	p := vm.ReservePage()
	p.raw[0] = 0xff
	vm.ReclaimPage(p)
	p2 := vm.ReservePage()
	if p2.raw[0] != 0 {
		t.Error("Reserved page after reclaimation must be clean")
	}
}

func TestPageView(t *testing.T) {
	if XtoPX(0) != 0 || XtoPX(1) != 0 || XtoPX(PageStrideByte) != 1 {
		t.Error("XtoPX error")
	}

	if YtoPY(PageStrideByte) != PageStrideByte {
		t.Error("YtoPY error")
	}

	if XYtoPI(0, 0) != 0 || XYtoPI(1, 0) != 0 || XYtoPI(PageStrideByte - 1, 0) != 0 ||
		XYtoPI(0, 1) != PageStrideWidth || XYtoPI(1, 1) != PageStrideWidth {
		t.Error("XYtoPI error")
	}

	if XtoSB(0) != 0 || XtoSB(5) != 5 || XtoSB(PageStrideByte) != 0 {
		t.Error("XtoSB error")
	}

	vm := NewVM()
	p := vm.ReservePage()
	p.raw[0] = 0x5 // 0b101
	if p.Get(0, 0) != 1 || p.Get(1, 0) != 0 || p.Get(2, 0) != 1 {
		t.Error("Get(0-2, 0) error")
	}

	p.raw[1] = 0x2 // 0b010
	if p.Get(0, 1) != 0 || p.Get(1, 1) != 1 || p.Get(2, 1) != 0 {
		t.Error("Get(0-2, 1) error")
	}

	p.Set(2, 1, 1)
	if p.Get(2, 1) != 1 || p.raw[1] != 0x6 {
		t.Error("Get(0-2, 1) error after set 1 to (2, 1)")
	}
}

var gliderPattern3x3 []byte = []byte{
	0, 1, 0,
	0, 0, 1,
	1, 1, 1,
}

var gliderPattern5x5 []byte = []byte{
	0, 0, 0, 0, 0,
	0, 0, 1, 0, 0,
	0, 0, 0, 1, 0,
	0, 1, 1, 1, 0,
	0, 0, 0, 0, 0,
}

func TestPageUtil(t *testing.T) {
	vm := NewVM()
	p := vm.ReservePage()
	p.AABB = New00WH(PageSizeWidth, PageSizeHeight)

	var v View = p
	var u ViewUtil = p

	if v == nil || u == nil {
		t.Error("Page(v, u) interface error")
	}

	n, err := u.Writer(NewXYWH(1, 1, 3, 3)).Write(gliderPattern3x3)
	if n != len(gliderPattern3x3) || err != io.EOF {
		t.Errorf("Page(%v, %v) invalid write op len/EOF (%v, %v)", v, p, n, err)
	}

	if ! u.Match(NewXYWH(0, 0, 5, 5), gliderPattern5x5) {
		t.Errorf("Page(%v, %v) invalid match for glider 5x5 ", v, p)
	}

	var m2 []byte = make([]byte, 3*3, 3*3)
	n, err = u.Reader(NewXYWH(1, 1, 3, 3)).Read(m2)
	if n != len(gliderPattern3x3) || err != io.EOF || ! bytes.Equal(m2, gliderPattern3x3) {
		t.Errorf("Page(%v, %v) invalid read op len/EOF (%v, %v)", v, p, n, err)
	}
}
