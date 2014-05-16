package main
/*
import "testing"
import "math"
import "fmt"

func TestWtoP(t *testing.T) {
	ws := uint(128)
	cases := [][2]int64{
		[2]int64{0, 0},
		[2]int64{127, 0},
		[2]int64{128, 1},
		[2]int64{129, 1},
		[2]int64{-1, -1},
		[2]int64{-127, -1},
		[2]int64{-128, -1},
		[2]int64{-129, -2},
	}
	for _, v := range cases {
		pc :=  WtoP(v[0], ws)
		if pc != v[1] {
			t.Error("pc ", v, ", px = ", WtoP(v[0], ws), " but ", pc)
		}
	}
}

func TestPtoW(t *testing.T) {
	if PtoW(0, 128) != 0 { t.Error("Invalid PtoW coord at 0 / 128") }
	if PtoW(1, 128) != 128 { t.Error("Invalid PtoW coord at 1 / 128") }
	if PtoW(2, 128) != 256 { t.Error("Invalid PtoW coord at 2 /  128") }

	if PtoW(-1, 128) != -128 { t.Error("Invalid PtoW coord at -1 / 128: ", PtoW(-1, 128)) }
	if PtoW(-2, 128) != -256 { t.Error("Invalid PtoW coord at -2 / 128") }
}

func TestNewPageTree(t *testing.T) {
	NewPageTree(NewAABB(-128, 128, 128, -128), 128)
}

func TestOFInt64(t *testing.T) {
	fmt.Printf("MinInt64 = %d\n", math.MinInt64)
	// OF: fmt.Printf("-1 * MinInt64 = %d\n", math.MinInt64 * -1)
	fmt.Printf("MaxInt64 = %d\n", math.MaxInt64)
	// fmt.Printf("MaxUInt64 = %d\n", math.MaxUint64)
	fmt.Printf("+ = %d\n", uint64(-1 * math.MinInt64) + uint64(math.MaxInt64))
	// OF: fmt.Printf("1 + = %x\n", uint64(-1 * math.MinInt64) + uint64(math.MaxInt64) + uint64(1))
}

func TestLastPage(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	p := vm.ReservePage()

	px := int64((uint64(math.MaxInt64) + 1) / uint64(vm.wsize))
	py := px
	pt := NewPageTile(p, vm.PageWidth(), px, py)
	fmt.Print("pt next last TR: wsize = ", vm.wsize," px = ", px, " pt = ", pt,"\n")

	px = int64((uint64(math.MaxInt64) + 1) / uint64(vm.wsize)) - 1
	py = px
	pt = NewPageTile(p, vm.PageWidth(), px, py)
	fmt.Print("pt last +1, + 1: wsize = ", vm.wsize," px = ", px, " pt = ", pt,"\n")

	if pt.AABB.MinX < 0 || pt.AABB.MinY < 0{
		t.Error("Invalid AABB overflow")
	}

	px = WtoP(math.MinInt64, vm.wsize)
	py = WtoP(math.MaxInt64, vm.wsize)
	pt = NewPageTile(p, vm.PageWidth(), px, py)
	fmt.Print("pt last -1, +1: wsize = ", vm.wsize," px = ", px, " pt = ", pt,"\n")
}

func TestNewPageTreeMaxInt64(t *testing.T) {
	pt := NewPageTree(NewAABBMax(), 128)
	fmt.Print("pt max = ", pt.AABB, ", sizex = ", pt.AABB.SizeX(),"\n")
}

func TestNewPageTile(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	p := vm.ReservePage()

	pt := NewPageTile(p, vm.PageWidth(), 0, 0)
	if pt.px != 0 || pt.py != 0 { t.Error("New page tile have invalid pt coords") }
	if pt.MinX != 0 || pt.MinY != 0 || pt.MaxX != 128 || pt.MaxY != 128 {
		t.Error("New page have invalid rect")
	}
	if pt.AABB.SizeX() != 128 || pt.AABB.SizeY() != 128{
		t.Error("New page have invalid rect size")
	}

	pt = NewPageTile(p, vm.PageWidth(), 1, 1)
	if pt.px != 1 || pt.py != 1 { t.Error("New page tile have invalid pt coords") }
	if pt.MinX != 128 || pt.MinY != 128 || pt.MaxX != 128 + 128 || pt.MaxY != 128 + 128 {
		t.Error("New page have invalid rect")
	}
	if pt.AABB.SizeX() != 128 || pt.AABB.SizeY() != 128{
		t.Error("New page have invalid rect size")
	}

	pt = NewPageTile(p, vm.PageWidth(), -1, -1)
	if pt.px != -1 || pt.py != -1 { t.Error("New page tile have invalid pt coords") }
	if pt.MinX != -128 || pt.MinY != -128 || pt.MaxX != 0 || pt.MaxY != 0 {
		t.Error("New page have invalid rect", pt)
	}
	if pt.AABB.SizeX() != 128 || pt.AABB.SizeY() != 128{
		t.Error("New page have invalid rect size", pt.AABB)
	}

}

func TestAddPage(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	p := vm.ReservePage()
	pt := NewPageTile(p, vm.PageWidth(), 1, 1)
	pb := NewPageTree(NewAABB(-128, 128, 128, -128), 128)
	pb.Add(&pt)
	if pb.Count() != 1 { t.Error("Pages count in a tree must be 1") }
	// TODO: Must throw an error
	pb.Add(&pt)
	if pb.Count() != 2 { t.Error("Pages count in a tree must be 2") }
}

func fillPageTree(t *testing.T, vm *VM) PageTree {
	pb := NewPageTree(NewAABB(-128 * 21, 128 * 21, 128 * 21, -128 * 21), 128)
	for x := int64(-20) ; x <= 20; x ++ {
		for y := int64(-20) ; y <= 20 ; y++ {
			pt := NewPageTile(vm.ReservePage(), vm.PageWidth(), x, y)
			if ! pt.Intersects(pb.AABB) {
				t.Error("This tile ", pt, " does not intersect whole tree ", pb)
			}
			pb.Add(&pt)
		}
	}
	if pb.Count() != 41 * 41 { t.Error("Tree have invalid size after fill ", pb.Count()) }
	return pb
}

func TestPageTreeSearchPage(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	pb := fillPageTree(t, &vm)
	for x := int64(-20) ; x <= 20; x ++ {
		for y := int64(-20) ; y <= 20 ; y++ {
			q := pb.QueryPage(x, y)
			if q == nil {
				t.Error("Can't query page at pos ", x, "x", y)
			}
			if q.px != x || q.py != y {
				t.Error("Query ret another page ", q)
			}
		}
	}
}

func TestPageRemove(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	pb := fillPageTree(t, &vm)
	for x := int64(-20) ; x <= 20; x ++ {
		for y := int64(-20) ; y <= 20 ; y++ {
			if ! pb.RemovePXY(x, y) {
				t.Error("Can't remove page at pos ", x, "x", y)
			}
		}
	}
	if pb.Count() != 0 {
		t.Error("Tree must be empty rn")
	}
	for x := int64(-20) ; x <= 20; x ++ {
		for y := int64(-20) ; y <= 20 ; y++ {
			q := pb.QueryPage(x, y)
			if q != nil {
				t.Error("Query page after empty at pos ", x, "x", y)
			}
		}
	}
}
*/
