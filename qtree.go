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

func NewQuadTree(bbox AABB) *QuadTree {
	return &QuadTree{ bbox, &QuadTile{AABB:bbox}, 0 }
}

func (qb *QuadTree) String() string {
	return fmt.Sprintf("(Tree %v)", qb.root)
}

func (qt *QuadTile) String() string {
	return fmt.Sprintf("(Tile/%d/%v c=%v, nw=%v, ne=%v, sw=%v, se=%v)", qt.level, qt.AABB, qt.contents, qt.nw, qt.ne, qt.sw, qt.se)
}

func (qb *QuadTree) Add(v *Page) {
	qb.AddTo(v, v.px, v.py)
}

func (qb *QuadTree) AddTo(p *Page, px uint64, py uint64) {
	// Update coords
	p.px = px
	p.py = py
	p.AABB = NewAABBPXY2W(px, py)
	// Add
	qb.root.add(p)
	// Update stats
	qb.count ++
	// Update my adjacent pages
	if px > qb.MinX {
		if py > qb.MinY {
			p.ap_nw = qb.QueryPoint(px - 1, py - 1)
		}
		p.ap_w = qb.QueryPoint(px - 1, py)
		if py < qb.MaxY {
			p.ap_sw = qb.QueryPoint(px - 1, py + 1)
		}
	}
	if py > qb.MinY {
		p.ap_n = qb.QueryPoint(px, py - 1)
	}
	if py < qb.MaxY {
		p.ap_s = qb.QueryPoint(px, py + 1)
	}
	if px < qb.MaxX {
		if py > qb.MinY {
			p.ap_ne = qb.QueryPoint(px + 1, py - 1)
		}
		p.ap_e = qb.QueryPoint(px + 1, py)
		if py < qb.MaxY {
			p.ap_se = qb.QueryPoint(px + 1, py + 1)
		}
	}
	// Update adjacent adjacent pages
	if p.ap_nw != nil { p.ap_nw.ap_se = p }
	if p.ap_w != nil { p.ap_w.ap_e = p }
	if p.ap_sw != nil { p.ap_sw.ap_ne = p }

	if p.ap_n != nil { p.ap_n.ap_s = p }
	if p.ap_s != nil { p.ap_s.ap_n = p }

	if p.ap_ne != nil { p.ap_ne.ap_sw = p }
	if p.ap_e != nil { p.ap_e.ap_w = p }
	if p.ap_se != nil { p.ap_se.ap_nw = p }
}

func NewAABBPXY2W(px uint64, py uint64) AABB {
	return NewXYWH(px * PageSizeWidth, py * PageSizeHeight, PageSizeWidth, PageSizeHeight)
}

func NewAABBW2P(wbox AABB) AABB {
	return NewXYWH(wbox.MinX >> PageStridePO2, wbox.MinY >> PageStridePO2, wbox.MaxX >> PageStridePO2, wbox.MaxY >> PageStridePO2)
}

func WXY2PXY(x uint64, y uint64) (px uint64, py uint64) {
	px = x >> PageStridePO2
	py = y >> PageStridePO2
	return
}

func (qb *QuadTree) RemoveAt(px uint64, py uint64) *Page {
	// Remove
	v := qb.root.remove(px, py)
	if v == nil {
		return nil
	}
	// Update stats
	qb.count--
	// Update adjacent adjacent pages
	if v.ap_nw != nil { v.ap_nw.ap_se = nil }
	if v.ap_w != nil { v.ap_w.ap_e = nil }
	if v.ap_sw != nil { v.ap_sw.ap_ne = nil }

	if v.ap_n != nil { v.ap_n.ap_s = nil }
	if v.ap_s != nil { v.ap_s.ap_n = nil }

	if v.ap_ne != nil { v.ap_ne.ap_sw = nil }
	if v.ap_e != nil { v.ap_e.ap_w = nil }
	if v.ap_se != nil { v.ap_se.ap_nw = nil }
	// Update my adjacent pages
	v.ap_nw = nil
	v.ap_w = nil
	v.ap_sw = nil

	v.ap_n = nil
	v.ap_s = nil

	v.ap_ne = nil
	v.ap_e = nil
	v.ap_se = nil
	return v
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

func (qt *QuadTile) add(p *Page) {
	// look for sub-tile directly below this tile to accomodate value.
	if c := qt.getChild(p.px, p.py); c == nil {
		// no suitable sub-tile for value found.
		// either this tile has no childs or
		// value does not fit in any subtile.
		// store value at this level.
		qt.contents = append(qt.contents, p)

		// tile is split if exceeds it max number of entries and
		// has not childs already and max tree depth for this sub-tree not reached.
		if len(qt.contents) > MAX_ENTRIES_PER_TILE && qt.nw == nil && qt.level < MAX_LEVELS {
			qt.split()
		}
	} else {
		// suitable sub-tile for value found at index i.
		// recursivly add value.
		c.add(p)
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

func (qt *QuadTile) remove(px uint64, py uint64) *Page {
	if ! qt.ContainsPoint(px, py) {
		return nil
	}

	for i, p := range qt.contents {
		if p.px == px && p.py == py {
			qt.contents[i], qt.contents = qt.contents[len(qt.contents)-1], qt.contents[:len(qt.contents)-1]
			return p
		}
	}

	if qt.nw != nil {
		if p := qt.nw.remove(px, py) ; p != nil {
			return p
		}
		if p := qt.ne.remove(px, py) ; p != nil {
			return p
		}
		if p := qt.sw.remove(px, py) ; p != nil {
			return p
		}
		if p := qt.se.remove(px, py) ; p != nil {
			return p
		}
	}

	return nil
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
