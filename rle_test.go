package main

import "testing"

func TestLoadRLE(t *testing.T) {
	vm := NewVM()
	p := vm.ReservePage()
	p.AABB = New00WH(PageSizeWidth, PageSizeHeight)

	LoadRLE(p, 1, 1, "pattern/glider.rle")

	/*
	fmt.Print(p.Print(NewXYWH(0, 0, 5, 5)))
	fmt.Printf("%5b\n", p.raw[0])
	fmt.Printf("%5b\n", p.raw[1])
	fmt.Printf("%5b\n", p.raw[2])
	fmt.Printf("%5b\n", p.raw[3])
	fmt.Printf("%5b\n", p.raw[4])
	*/

	if ! p.Match(NewXYWH(0, 0, 5, 5), GliderPattern5x5) {
		t.Errorf("Page(%v) invalid match for glider 5x5 ", p)
	}
}
