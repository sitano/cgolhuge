package main

type PageTree QuadTree

type PageTile struct {
	AABB

	// Page from vm
	p *Page

	// Page position in page-world coords
	px int64
	py int64

	// Game
	alive uint
}

// Make AABB implement the QuadElement interface
func (pt PageTile) getAABB() AABB {
	return pt.AABB
}

func NewPageTree(bbox AABB) PageTree {
	return PageTree(NewQuadTree(bbox))
}

func NewPageTile(p *Page, wsize uint, px int64, py int64) PageTile {
	return PageTile{
		AABB: NewAABBPXY(px, py, wsize),
		p: p,
		px: px,
		py: py,
		alive: 0,
	}
}

func NewAABBPXY(px int64, py int64, wsize uint) AABB {
	return NewAABB(
		PtoW(px, wsize), PtoW(px + 1, wsize),
		PtoW(py, wsize), PtoW(py + 1, wsize))
}

func PtoW(pc int64, wsize uint) int64 {
	return pc * int64(wsize)
}

func WtoP(wc int64, wsize uint) int64 {
	if (wc >= 0) {
		return wc / int64(wsize)
	}
	return wc / int64(wsize) - 1
}

func (pb *PageTree) Add(pt *PageTile) {
	(*QuadTree)(pb).Add(pt)
}

func (pb *PageTree) Remove(pt *PageTile) {
	panic("Not implemented")
}

func (pb *PageTree) Count() uint64 {
	return pb.count
}

func (pb *PageTree) getAABB() AABB {
	return pb.root.getAABB()
}

func (pb *PageTree) QueryPage(px int64, py int64, wsize uint) *PageTile {
	return pb.root.queryPage(px, py, NewAABBPXY(px, py, wsize))
}

func (tile *QuadTile) queryPage(px int64, py int64, qbox AABB) *PageTile {
	// end recursion if this tile does not intersect the query range
	if ! tile.Intersects(qbox) {
		return nil
	}

	// return candidates at this tile
	for _, v := range tile.contents {
		p := v.(*PageTile)
		if p.px == px && p.py == py {
			return p
		}
	}

	// recurse into childs (if any)
	if tile.childs[_TOPRIGHT] != nil {
		for _, child := range tile.childs {
			ret := child.queryPage(px, py, qbox)
			if ret != nil { return ret }
		}
	}

	return nil
}
