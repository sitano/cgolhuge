package main

//import "fmt"
import "testing"

func TestLoadLIF(t *testing.T) {
	vm := NewVM()
	p := vm.ReservePage()
	p.AABB = New00WH(PageSizeWidth, PageSizeHeight)

	LoadLIF(p, 1, 1, "pattern/glider.lif")
	//fmt.Print(p.Print(NewXYWH(0, 0, 5, 5)))

	if ! p.Match(NewXYWH(0, 0, 5, 5), GliderPattern5x5) {
		t.Errorf("Page(%v) invalid match for glider 5x5 ", p)
	}
}
