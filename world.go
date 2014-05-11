package main

import "math"
import "container/list"

const (
	RULE_LIVE_MIN = 2
	RULE_LIVE_MAX = 3
	RULE_BORN = 3
)

type World interface {
	View

	Generation() uint64
	Population() uint64

	LifePoints(z byte) [][2]int64
	View() View
	Layer() byte

	Step()
}

type LifeWorld struct {
	// World state
	v *WorldView

	// Current life layer
	z byte

	// Life stats
	generation uint64
	population uint64
}

func NewLifeWorld(v *WorldView) LifeWorld {
	return LifeWorld{v, 0, 0, 0}
}

func NewLifeWorldXY(bbox AABB) LifeWorld {
	vm := NewVM(KSIZE_16K)
	pb := NewPageTree(bbox, vm.wsize)
	wv := NewWorldView(&vm, &pb)
	return NewLifeWorld(&wv)
}

func (lw *LifeWorld) Generation() uint64 {
	return lw.generation
}

func (lw *LifeWorld) Population() uint64 {
	return lw.population
}

func (lw *LifeWorld) LifePoints() [][2]int64 {
	panic("Not implemented")
}

func (lw *LifeWorld) View() View {
	return lw.v
}

func (lw *LifeWorld) Layer() byte {
	return lw.z
}

func (lw *LifeWorld) NextZLayer() byte {
	return (lw.z + 1) % 2
}

func (lw *LifeWorld) Step() {
	lw.generation ++
	lw.population = 0
	cz := lw.z
	nz := lw.NextZLayer()
	ws := lw.v.pb.wsize
	ks := lw.v.vm.ksize

	// Process life
	new := list.New()
	lw.population = lw.v.pb.Reduce(func(a interface{}, pt *PageTile) interface{} {
		population := a.(uint64)

		// Prevent reclamation during processing
		pt.alive = 1

		for i := uint(0) ; i < ks; i ++ {
			x := POtoWX(i, pt.px, ws)
			y := POtoWY(i, pt.py, ws)

			sum := lw.LifeSumAt(x, y, cz)
			st := lw.Get(x, y, cz)
			nst := DEAD

			if st == DEAD {
				if sum == RULE_BORN {
					nst = LIFE
					pt.alive ++
					population ++
				}
			}

			if st == LIFE {
				if sum >= RULE_LIVE_MIN && sum <= RULE_LIVE_MAX {
					nst = LIFE
					pt.alive ++
					population ++
				}
			}

			lw.Set(x, y, nz, nst)
		}

		pt.alive --

		// Check page edges (special case when there is no page)
		lw.TryEdgeLines(new, pt.getAABB())

		return population
	}, lw.population).(uint64)

	// Fill in additional life (egde case)
	lw.population += uint64(new.Len())
	lw.PurgePoints(new, nz)

	// Reclaim memory
	lw.TryReclaim()

	lw.z = nz
}

func (lw *LifeWorld) PurgePoints(ll *list.List, z byte) {
	for ll.Len() > 0 {
		np := ll.Remove(ll.Front()).([2]int64)
		lw.Set(np[0], np[1], z, LIFE)
	}
}

func (lw *LifeWorld) TryEdgeLines(ll *list.List, pbb AABB) {
	maxX := pbb.MaxX - 1
	maxY := pbb.MaxY - 1
	if pbb.MaxX == math.MaxInt64 { maxX = math.MaxInt64 }
	if pbb.MaxY == math.MaxInt64 { maxY = math.MaxInt64 }
	lw.TryEdgePoint(ll, pbb.MinX, pbb.MinY, -1, -1)
	lw.TryEdgePoint(ll, pbb.MinX, maxY,     -1, +1)
	lw.TryEdgePoint(ll, maxX, pbb.MinY,     +1, -1)
	lw.TryEdgePoint(ll, maxX, maxY,         +1, +1)
	for x := pbb.MinX ; x <= maxX && x >= pbb.MinX; x++ {
		lw.TryEdgePoint(ll, x, pbb.MinY, 0, -1)
		lw.TryEdgePoint(ll, x, maxY, 0, +1)
	}
	for y := pbb.MinY ; y <= maxY && y >= pbb.MinY; y++ {
		lw.TryEdgePoint(ll, pbb.MinX, y, -1, 0)
		lw.TryEdgePoint(ll, maxX,     y, +1, 0)
	}
}

func (lw *LifeWorld) TryEdgePoint(ll *list.List, x int64, y int64, dx int64, dy int64) {
	gbb := lw.v.pb.getAABB()
	ws := lw.v.vm.wsize
	tx := MvXY1(x, dx, gbb.MinX, gbb.MaxX)
	ty := MvXY1(y, dy, gbb.MinY, gbb.MaxY)
	if lw.v.pb.QueryPage(WtoP(tx, ws), WtoP(ty, ws)) == nil {
		ts := lw.v.LifeSumAt(tx, ty, lw.Layer())
		if ts == RULE_BORN {
			ll.PushBack([2]int64{tx, ty})
		}
	}
}

func (lw *LifeWorld) TryReclaim() {
	lw.v.pb.Reduce(func(a interface{}, pt *PageTile) interface{} {
		if pt.alive == 0 {
			lw.v.TryReclaim(pt)
		}
		return nil
	}, nil)
}

func (lw *LifeWorld) Set(x int64, y int64, z byte, t byte) {
	lw.v.Set(x, y, z, t)
}

func (lw *LifeWorld) Get(x int64, y int64, z byte) byte {
	return lw.v.Get(x, y, z)
}

func (lw *LifeWorld) NextTo(x int64, y int64, z byte, dx int64, dy int64) byte {
	return lw.v.NextTo(x, y, z, dx, dy)
}

func (lw *LifeWorld) LifeSumAt(x int64, y int64, z byte) byte {
	return lw.v.LifeSumAt(x, y, z)
}
