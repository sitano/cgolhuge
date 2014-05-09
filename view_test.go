package main

import "math"
import "testing"

func TestNewWorldView(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	pb := NewPageTree(NewAABB(-128, 128, 128, -128), 128)
	NewWorldView(&vm, &pb)
}

func TestReadStateZ(t *testing.T) {
	b := byte(135) // 10 00 01 11
	if ReadStateZ(b, 0) != 3 { t.Error("read z0") }
	if ReadStateZ(b, 1) != 1 { t.Error("read z1") }
	if ReadStateZ(b, 2) != 0 { t.Error("read z2") }
	if ReadStateZ(b, 3) != 2 { t.Error("read z3", ReadStateZ(b, 3)) }
}

func TestWriteStateZ(t *testing.T) {
	b := byte(135) // 10000111
	if WriteStateZ(b, 0, 1) & 0x3 != 1 { t.Error("write z0") }
	if (WriteStateZ(b, 1, 2) >> 2) & 0x3 != 2 { t.Error("write z2", ReadStateZ(WriteStateZ(b, 1, 2), 1)) }
	if (WriteStateZ(b, 2, 3) >> 4) & 0x3 != 3 { t.Error("write z3") }
	if (WriteStateZ(b, 3, 2) >> 6) & 0x3 != 2 { t.Error("write z4", ReadStateZ(WriteStateZ(b, 3, 4), 3)) }
}

func TestClearStateZ(t *testing.T) {
	b := byte(135) // 10000111
	if ReadStateZ(b, 3) != 2 { t.Error("clear z4") }
	b = ClearStateZ(b, 3)
	if ReadStateZ(b, 3) != 0 { t.Error("clear z3") }
	if b != byte(135-128) { t.Error("clear z3/2") }
}

func TestGetPPageOfftest(t *testing.T) {
	ws := uint(128)
	cases := [][5]int64{
		[5]int64{0,0,0,0, 128 * 128 - 128},
		[5]int64{127,0,0,0, KSIZE_16K - 1},
		[5]int64{0,127,0,0, 0},
		[5]int64{127,127,0,0, 127},
	}
	for _, v := range cases {
		offset :=  GetPPageOffset(v[0], v[1], v[2], v[3], ws)
		if offset != uint(v[4]) {
			t.Error("offset ", v, " but ", offset)
		}
		if offset >= KSIZE_16K {
			t.Error("offset exceeded max ksize at ", v)
		}
	}
}

func fillPageDead(pb *PageTree, p *PageTile, z byte) {
	for i := uint(0) ; i < pb.wsize * pb.wsize ; i ++ {
		p.SetByte(i, WriteStateZ(p.GetByte(i), z, DEAD))
	}
}

func fillPageLife(pb *PageTree, p *PageTile, z byte) {
	for i := uint(0) ; i < pb.wsize * pb.wsize ; i ++ {
		p.SetByte(i, WriteStateZ(p.GetByte(i), z, LIFE))
	}
}

func TestWorldViewSetGet(t *testing.T) {
	vm := NewVM(KSIZE_16K)
	pb := NewPageTree(NewAABBMax(), vm.wsize)
	wv := NewWorldView(&vm, &pb)
	ws := vm.wsize
	z := byte(0)
	//fmt.Print("px = ", WtoP(0, wv.pb.wsize), "\n")
	if wv.Get(0, 0, z) != DEAD { t.Error("0, 0 d") }
	wv.Set(0, 0, z, DEAD)
	if wv.Get(0, 0, z) != DEAD { t.Error("0, 0 d") }
	wv.Set(0, 0, z, LIFE)
	if wv.Get(0, 0, z) != LIFE { t.Error("0, 0 l") }
	wv.Set(128, 0, 0, LIFE)
	q1280 := pb.QueryPage(1, 0)
	if q1280.alive != 1 { t.Error("128, 0 alive 1") }
	wv.Set(128, 0, 1, DEAD)
	if q1280.alive != 1 { t.Error("128, 0 alive 1") }
	wv.Set(128, 0, 2, LIFE)
	if q1280.alive != 2 { t.Error("128, 0 alive 2") }
	if wv.Get(128, 0, 0) != LIFE { t.Error("128, 0, 0 l") }
	if wv.Get(128, 0, 1) != DEAD { t.Error("128, 0, 1 d") }
	if wv.Get(128, 0, 2) != LIFE { t.Error("128, 0, 2 l") }
	wv.Set(128, 0, 0, DEAD)
	if q1280.alive != 1 { t.Error("128, 0 alive 1") }
	wv.Set(128, 0, 2, DEAD)
	if q1280.alive != 0 { t.Error("128, 0 alive 0") }
	if wv.Get(128, 0, 0) != DEAD { t.Error("128, 0, 0 d") }
	if wv.Get(128, 0, 1) != DEAD { t.Error("128, 0, 1 d") }
	if wv.Get(128, 0, 2) != DEAD { t.Error("128, 0, 2 d") }

	wv.Set(math.MaxInt64, math.MaxInt64, 0, LIFE)
	if wv.Get(math.MaxInt64, math.MaxInt64, 0) != LIFE {
		t.Error("math.MaxInt64, math.MaxInt64, 0 l")
	}
	wv.Set(math.MaxInt64, math.MaxInt64, 0, DEAD)
	wv.Set(math.MaxInt64, math.MaxInt64, 0, DEAD)
	qMax := pb.QueryPage(WtoP(math.MaxInt64, ws), WtoP(math.MaxInt64, ws))
	if qMax != nil {
		if qMax.alive != 0 {
			t.Error("Empty page must have 0 alive")
		}
		t.Error("Empty page must be reclaimed out of tree")
	}
}
