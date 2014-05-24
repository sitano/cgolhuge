package main

import (
	"io"
)

const (
	DEAD = byte(0)
	LIFE = byte(1)

	RULE_LIVE_MIN = 2
	RULE_LIVE_MAX = 3
	RULE_BORN = 3

	BITS1   = uint64(0x1) //   0b1
	BITS11  = uint64(0x3) //  0b11
	BITS111 = uint64(0x7) // 0b111
	BITS101 = uint64(0x5) // 0b101
)

type World interface {
	View
	ViewUtil

	Generation() uint64
	Population() uint64

	Step()
	View() View
}

type LifeWorld struct {
	World

	// World state
	v *WorldView

	// Life stats
	generation uint64
	population uint64
}

func NewLifeWorld(v *WorldView) *LifeWorld {
	return &LifeWorld{
		v: v,
		generation: 0,
		population: 0,
	}
}

/*
func NewLifeWorldXY(bbox AABB) LifeWorld {
	vm := NewVM(KSIZE_16K)
	pb := NewPageTree(bbox, vm.wsize)
	wv := NewWorldView(&vm, &pb)
	return NewLifeWorld(&wv)
}
*/

func (w *LifeWorld) Generation() uint64 {
	return w.generation
}

func (w *LifeWorld) Population() uint64 {
	return w.population
}

func (w *LifeWorld) View() View {
	return w.v
}

func (w *LifeWorld) Step() {
	w.generation ++
	w.population = 0

	if w.v.vm.Pages() == 0 {
		return
	}

	p := w.v.vm.reserved[0]
	p_len := len(p.raw)
	prev_line := uint64(0)
	curr_line := p.raw[0]
	next_line := uint64(0)
	for _, p = range w.v.vm.reserved {
		p.next = NewPageBuf()
		raw := p.raw
		next := p.next

		for ni := 1; ni < p_len; ni ++ {
			ci := ni - 1
			next_line = raw[ni]
			new_line := uint64(0)

			if prev_line != 0 || curr_line != 0 || next_line != 0 {
				// TODO: first bit 0x2

				for bi := uint(1); bi < PageStrideBits - 2; bi ++ {
					bi_1 := bi - 1
					sum := PopCnt((prev_line & (BITS111 << bi_1)) >> bi_1 +
							((curr_line & (BITS101 << bi_1)) >> bi_1) << 3 +
							((next_line & (BITS111 << bi_1)) >> bi_1) << 6)

					if sum < RULE_LIVE_MIN {
						continue
					}

					st := byte(PopCnt(curr_line & (BITS1 << bi)))
					// fmt.Printf("                                                                 line=%d, pl=%b, cl=%b, nl=%b\n", ci, prev_line, curr_line, next_line)
					// fmt.Printf("                                                                 line=%d, st=%d, sum=%d, pl=%b, cl=%b, nl=%b\n", ci, st, sum, prev_line & (BITS111 << (bi - 1)), curr_line & (BITS101 << (bi - 1)), next_line & (BITS111 << (bi - 1)))

					if st == DEAD {
						if sum == RULE_BORN {
							new_line = new_line | (BITS1 << bi)
							w.population ++
						}
					} else {
						if sum >= RULE_LIVE_MIN && sum <= RULE_LIVE_MAX {
							new_line = new_line | (BITS1 << bi)
							w.population ++
						}
					}
				}

				// TODO: last bit 0x2

				next[ci] = new_line
			}

			prev_line = curr_line
			curr_line = next_line
		}

		next_line = 0
		// TODO: scan last line
	}

	w.Swap()
}

func (w *LifeWorld) Swap() {
	for _, p := range w.v.vm.reserved {
		p.raw = p.next
		p.next = nil
	}
}

/*
func (w *LifeWorld) Step() {
	w.generation ++
	w.population = 0
	cz := w.Layer()
	nz := w.NextZLayer()
	ws := w.v.pb.wsize
	ks := w.v.vm.ksize

	// Process life
	w.v.autoReclaim = false
	new := list.New()
	w.population = w.v.pb.Reduce(func(a interface{}, pt *PageTile) interface{} {
		pt.alive = 0

		for i := uint(0) ; i < ks; i ++ {
			x := POtoWX(i, pt.px, ws)
			y := POtoWY(i, pt.py, ws)

			sum := w.LifeSumAt(x, y, cz)
			st := w.Get(x, y, cz)
			nst := DEAD

			if st == DEAD {
				if sum == RULE_BORN {
					nst = LIFE
					pt.alive ++
				}
			}

			if st == LIFE {
				if sum >= RULE_LIVE_MIN && sum <= RULE_LIVE_MAX {
					nst = LIFE
					pt.alive ++
				}
			}

			w.Set(x, y, nz, nst)
		}

		// Check page edges (special case when there is no page)
		w.TryEdgeLines(new, pt.GetAABB())

		return a.(uint64) + uint64(pt.alive)
	}, w.population).(uint64)

	// Restore before purge
	w.v.autoReclaim = true

	// Fill in additional life (egde case)
	w.population += uint64(new.Len())
	w.PurgePoints(new, nz)

	// Reclaim memory
	w.TryReclaim()

	w.z = nz
}

func (w *LifeWorld) PurgePoints(ll *list.List, z byte) {
	for ll.Len() > 0 {
		np := ll.Remove(ll.Front()).([2]int64)
		w.Set(np[0], np[1], z, LIFE)
	}
}

func (w *LifeWorld) TryEdgeLines(ll *list.List, pbb AABB) {
	maxX := pbb.MaxX - 1
	maxY := pbb.MaxY - 1
	if pbb.MaxX == math.MaxInt64 { maxX = math.MaxInt64 }
	if pbb.MaxY == math.MaxInt64 { maxY = math.MaxInt64 }
	w.TryEdgePoint(ll, pbb.MinX, pbb.MinY, -1, -1)
	w.TryEdgePoint(ll, pbb.MinX, maxY,     -1, +1)
	w.TryEdgePoint(ll, maxX, pbb.MinY,     +1, -1)
	w.TryEdgePoint(ll, maxX, maxY,         +1, +1)
	for x := pbb.MinX ; x <= maxX && x >= pbb.MinX; x++ {
		w.TryEdgePoint(ll, x, pbb.MinY, 0, -1)
		w.TryEdgePoint(ll, x, maxY, 0, +1)
	}
	for y := pbb.MinY ; y <= maxY && y >= pbb.MinY; y++ {
		w.TryEdgePoint(ll, pbb.MinX, y, -1, 0)
		w.TryEdgePoint(ll, maxX,     y, +1, 0)
	}
}

func (w *LifeWorld) TryEdgePoint(ll *list.List, x int64, y int64, dx int64, dy int64) {
	gbb := w.v.pb.GetAABB()
	ws := w.v.vm.wsize
	tx := MvXY1(x, dx, gbb.MinX, gbb.MaxX)
	ty := MvXY1(y, dy, gbb.MinY, gbb.MaxY)
	if w.v.pb.QueryPage(WtoP(tx, ws), WtoP(ty, ws)) == nil {
		ts := w.v.LifeSumAt(tx, ty, w.Layer())
		if ts == RULE_BORN {
			ll.PushBack([2]int64{tx, ty})
		}
	}
}

func (w *LifeWorld) TryReclaim() {
	w.v.pb.Reduce(func(a interface{}, pt *PageTile) interface{} {
		if pt.alive == 0 {
			w.v.TryReclaim(pt)
		}
		return nil
	}, nil)
}

func (w *LifeWorld) Set(x int64, y int64, z byte, t byte) {
	w.v.Set(x, y, z, t)
}

func (w *LifeWorld) Get(x int64, y int64, z byte) byte {
	return w.v.Get(x, y, z)
}

func (w *LifeWorld) NextTo(x int64, y int64, z byte, dx int64, dy int64) byte {
	return w.v.NextTo(x, y, z, dx, dy)
}

func (w *LifeWorld) LifeSumAt(x int64, y int64, z byte) byte {
	return w.v.LifeSumAt(x, y, z)
}
*/

// View implementation

func (w *LifeWorld) GetAABB() AABB {
	return w.v.AABB
}

func (w *LifeWorld) Get(x uint64, y uint64) byte {
	return w.v.Get(x, y)
}

func (w *LifeWorld) Set(x uint64, y uint64, v byte) {
	w.v.Set(x, y, v)
}

// ViewUtil implementation

func (w *LifeWorld) Print(b AABB) string {
	return Print(w, b)
}

func (w *LifeWorld) Match(b AABB, matcher []byte) bool {
	return Match(w, b, matcher)
}

func (w *LifeWorld) MirrorH(b AABB) {
	MirrorH(w, b)
}

func (w *LifeWorld) MirrorV(b AABB) {
	MirrorV(w, b)
}

func (w *LifeWorld) Writer(b AABB) io.Writer {
	return Writer(w, b)
}

func (w *LifeWorld) Reader(b AABB) io.Reader {
	return Reader(w, b)
}
