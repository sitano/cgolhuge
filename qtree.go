package main

import "fmt"

const MAX_ENTRIES_PER_TILE = 16
const MAX_LEVELS = 64

type QuadTile struct {
	AABB

	level    int
	contents []*Page

	nw, ne, sw, se *QuadTile
}

type QuadTree struct {
	AABB

	root *QuadTile
	count uint64
}

func NewQuadTree(bbox AABB) QuadTree {
	qt := QuadTree{ bbox, &QuadTile{AABB:bbox}, 0 }
	return qt
}

func (qb QuadTree) String() string {
	return fmt.Sprintf("(Tree %v)", qb.root)
}

func (qt *QuadTile) String() string {
	return fmt.Sprintf("(Tile/%d/%v c=%v, nw=%v, ne=%v, sw=%v, se=%v)", qt.level, qt.AABB, qt.contents, qt.nw, qt.ne, qt.sw, qt.se)
}

func (qb *QuadTree) Add(v *Page) {
	qb.AddTo(v, v.px, v.py)
}

func (qb *QuadTree) AddTo(v *Page, px uint64, py uint64) {
	v.px = px
	v.py = py
	v.AABB = NewWAABB4PXY(px, py)
	qb.root.add(v)
	qb.count ++
}

func NewWAABB4PXY(px uint64, py uint64) AABB {
	return NewXYWH(px * PageSizeWidth, py * PageSizeHeight, PageSizeWidth, PageSizeHeight)
}

func NewWAABB4PAABB(pbb AABB) AABB {
	return NewXYWH(pbb.MinX >> PageStridePO2, pbb.MinY >> PageStridePO2, 1, 1)
}

func (qb *QuadTree) RemoveAt(px uint64, py uint64) bool {
	if qb.root.remove(px, py) {
		qb.count--
		return true
	}
	return false
}

func (qb *QuadTree) QueryBox(bbox AABB) (values []*Page) {
	return qb.root.queryBox(bbox, values)
}

func (qb *QuadTree) QueryPoint(px uint64, py uint64) *Page {
	return qb.root.queryPoint(px, py)
}

func (qb *QuadTree) Reduce(f func(p *Page)) {
	qb.root.reduce(f)
}

func (qt *QuadTile) add(v *Page) {
	// look for sub-tile directly below this tile to accomodate value.
	if c := qt.getChild(v.px, v.py); c == nil {
		// no suitable sub-tile for value found.
		// either this tile has no childs or
		// value does not fit in any subtile.
		// store value at this level.
		qt.contents = append(qt.contents, v)

		// tile is split if exceeds it max number of entries and
		// has not childs already and max tree depth for this sub-tree not reached.
		if len(qt.contents) > MAX_ENTRIES_PER_TILE && qt.nw == nil && qt.level < MAX_LEVELS {
			qt.split()
		}
	} else {
		// suitable sub-tile for value found at index i.
		// recursivly add value.
		c.add(v)
	}
}

func (qt *QuadTile) getChild(px uint64, py uint64) *QuadTile {
	if qt.nw != nil {
		if qt.nw.ContainsPoint(px, py) {
			return qt.nw
		}

		if qt.ne.ContainsPoint(px, py) {
			return qt.ne
		}

		if qt.sw.ContainsPoint(px, py) {
			return qt.sw
		}

		if qt.se.ContainsPoint(px, py) {
			return qt.se
		}
	}

	return nil
}

// create four child quads.
// distribute contents of this tiles on newly created childs.
func (qt *QuadTile) split() {
	w2 := qt.SizeX() / 2
	h2 := qt.SizeY() / 2

	qt.nw = &QuadTile{ AABB:NewAABB(qt.MinX, qt.MinY, qt.MinX + w2 - 1, qt.MinY + h2 - 1), level:qt.level+1 }
	qt.ne = &QuadTile{ AABB:NewAABB(qt.MinX + w2, qt.MinY, qt.MaxX, qt.MinY + h2 - 1), level:qt.level+1 }
	qt.sw = &QuadTile{ AABB:NewAABB(qt.MinX, qt.MinY + h2, qt.MinX + w2 - 1, qt.MaxY), level:qt.level+1 }
	qt.se = &QuadTile{ AABB:NewAABB(qt.MinX + w2, qt.MinY + h2, qt.MaxX, qt.MaxY), level:qt.level+1 }

	// copy values to temporary slice
	contentsBak := append([]*Page{}, qt.contents...)

	// clear values on this tile
	qt.contents = []*Page{}

	// reinsert from temporary slice
	for _,v := range contentsBak {
		qt.add(v)
	}
}

func (qt *QuadTile) queryBox(qbox AABB, ret []*Page) []*Page {
	if ! qt.Intersects(qbox) {
		return ret
	}

	for _, v := range qt.contents {
		if qbox.ContainsPoint(v.px, v.py) {
			ret = append(ret, v)
		}
	}

	if qt.nw != nil {
		ret = qt.nw.queryBox(qbox, ret)
		ret = qt.ne.queryBox(qbox, ret)
		ret = qt.sw.queryBox(qbox, ret)
		ret = qt.se.queryBox(qbox, ret)
	}

	return ret
}

func (qt *QuadTile) queryPoint(px uint64, py uint64) *Page {
	if ! qt.ContainsPoint(px, py) {
		return nil
	}

	for _, v := range qt.contents {
		if v.px == px && v.py == py {
			return v
		}
	}

	if qt.nw != nil {
		if v := qt.nw.queryPoint(px, py) ; v != nil {
			return v
		}
		if v := qt.ne.queryPoint(px, py) ; v != nil {
			return v
		}
		if v := qt.sw.queryPoint(px, py) ; v != nil {
			return v
		}
		if v := qt.se.queryPoint(px, py) ; v != nil {
			return v
		}
	}

	return nil
}

func (qt *QuadTile) remove(px uint64, py uint64) bool {
	if ! qt.ContainsPoint(px, py) {
		return false
	}

	for i, v := range qt.contents {
		if v.px == px && v.py == py {
			qt.contents[i], qt.contents = qt.contents[len(qt.contents)-1], qt.contents[:len(qt.contents)-1]
			return true
		}
	}

	if qt.nw != nil {
		if qt.nw.remove(px, py) {
			return true
		}
		if qt.ne.remove(px, py) {
			return true
		}
		if qt.sw.remove(px, py) {
			return true
		}
		if qt.se.remove(px, py) {
			return true
		}
	}

	return false
}

func (qt *QuadTile) reduce(f func(p *Page)) {
	for _, t := range qt.contents {
		f(t)
	}

	if qt.nw != nil {
		qt.nw.reduce(f)
		qt.ne.reduce(f)
		qt.sw.reduce(f)
		qt.se.reduce(f)
	}
}
