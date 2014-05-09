package main

type PageTree struct {
	QuadTree

	wsize uint
}

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

func NewPageTree(bbox AABB, wsize uint) PageTree {
	if bbox.SizeX() % uint64(wsize) != 0 || bbox.SizeY() % uint64(wsize) != 0 {
		panic("NewPageTree: bbox size does not fit page wsize")
	}
	return PageTree{NewQuadTree(bbox), wsize}
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
	// TODO: test page wsize, aabb, etc
	pb.QuadTree.Add(pt)
}

func (pb *PageTree) Remove(pt *PageTile) bool {
	if pb.root.remove(pt.px, pt.py, pt.AABB) {
		pb.count --
		return true
	}
	return false
}

func (pb *PageTree) RemovePXY(px int64, py int64) bool {
	if pb.root.remove(px, py, NewAABBPXY(px, py, pb.wsize)) {
		pb.count --
		return true
	}
	return false
}

func (pb *PageTree) Count() uint64 {
	return pb.count
}

func (pb *PageTree) getAABB() AABB {
	return pb.root.getAABB()
}

func (pb *PageTree) QueryPage(px int64, py int64) *PageTile {
	return pb.root.queryPage(px, py, NewAABBPXY(px, py, pb.wsize))
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

func (tile *QuadTile) remove(px int64, py int64, qbox AABB) bool {
	// end recursion if this tile does not intersect the query range
	if ! tile.Intersects(qbox) {
		return false
	}

	// return candidates at this tile
	for i, v := range tile.contents {
		p := v.(*PageTile)
		if p.px == px && p.py == py {
			tile.contents[i] = tile.contents[len(tile.contents)-1]
			tile.contents = tile.contents[0:len(tile.contents)-1]
			// TODO: merge parent tree node childs if can, but i dont need it RN
			return true
		}
	}

	// recurse into childs (if any)
	if tile.childs[_TOPRIGHT] != nil {
		for _, child := range tile.childs {
			ret := child.remove(px, py, qbox)
			if ret { return ret }
		}
	}

	return false
}
