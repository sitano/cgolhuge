/*
Based on work of Volker Poplawski, 2013 (https://github.com/volkerp/goquadtree)
*/
package main
/*
import "container/list"

// Number of entries until a quad is split
const MAX_ENTRIES_PER_TILE = 16

// Maximum depth of quad-tree (not counting the root node)
const MAX_LEVELS = 64

// AABB orientation inside tree
const (
	_TOPRIGHT    = 0
	_TOPLEFT     = 1
	_BOTTOMLEFT  = 2
	_BOTTOMRIGHT = 3
)

// QuadTree expects its values to implement the QuadElement interface.
type QuadElement interface {
	GetAABB() AABB
}

// quad-tile / node of the quad-tree
type QuadTile struct {
	AABB

	level    int           // level this tile is at. root is level 0
	contents []QuadElement // values stored in this tile
	childs   [4]*QuadTile  // sub-tiles. none or four.
}

type QuadTree struct {
	AABB

	root QuadTile
	count uint64
}

// Constructs an empty quad-tree
// bbox specifies the extends of the coordinate system.
func NewQuadTree(bbox AABB) QuadTree {
	qt := QuadTree{ bbox, QuadTile{AABB:bbox}, 0 }
	return qt
}

// Adds a value to the quad-tree by trickle down from the root node.
func (qb *QuadTree) Add(v QuadElement) {
	qb.root.add(v)
	qb.count ++
}

func (qb *QuadTree) Count() uint64 {
	return qb.count
}

func (tile *QuadTile) Contents() []QuadElement {
	return tile.contents
}

func (tile *QuadTile) Childs() [4]*QuadTile {
	return tile.childs
}

// Returns all values which intersect the query box
func (qb *QuadTree) Query(bbox AABB) (values []QuadElement) {
	return qb.root.query(bbox, values)
}

func (qb *QuadTree) Reduce(f func(a interface{}, t QuadElement) interface{}, v interface{}) interface{} {
	return qb.root.reduce(f, v)
}

func (tile *QuadTile) add(v QuadElement) {
	// look for sub-tile directly below this tile to accomodate value.
	if i := tile.findChildIndex(v.GetAABB()); i < 0 {
		// no suitable sub-tile for value found.
		// either this tile has no childs or
		// value does not fit in any subtile.
		// store value at this level.
		tile.contents = append(tile.contents, v)

		// tile is split if exceeds it max number of entries and
		// has not childs already and max tree depth for this sub-tree not reached.
		if len(tile.contents) > MAX_ENTRIES_PER_TILE && tile.childs[_TOPRIGHT] == nil && tile.level < MAX_LEVELS {
			tile.split()
		}
	} else {
		// suitable sub-tile for value found at index i.
		// recursivly add value.
		tile.childs[i].add(v)
	}
}


// return child index for AABB
// returns -1 if quad has no children or AABB does not fit into any child
func (tile *QuadTile) findChildIndex(bbox AABB) int {
	if tile.childs[_TOPRIGHT] == nil {
		return -1
	}

	for i, child := range tile.childs {
		if child.Contains(bbox) {
			return i
		}
	}

	return -1
}


// create four child quads.
// distribute contents of this tiles on newly created childs.
func (tile *QuadTile) split() {
	mx := tile.MaxX/2.0 + tile.MinX/2.0
	my := tile.MaxY/2.0 + tile.MinY/2.0

	tile.childs[_TOPRIGHT]    = &QuadTile{ AABB:NewAABB(mx, tile.MaxX, my, tile.MaxY), level:tile.level+1 }
	tile.childs[_TOPLEFT]     = &QuadTile{ AABB:NewAABB(tile.MinX, mx, my, tile.MaxY), level:tile.level+1 }
	tile.childs[_BOTTOMLEFT]  = &QuadTile{ AABB:NewAABB(tile.MinX, mx, tile.MinY, my), level:tile.level+1 }
	tile.childs[_BOTTOMRIGHT] = &QuadTile{ AABB:NewAABB(mx, tile.MaxX, tile.MinY, my), level:tile.level+1 }

	// copy values to temporary slice
	var contentsBak []QuadElement
	contentsBak = append(contentsBak, tile.contents...)

	// clear values on this tile
	tile.contents = []QuadElement{}

	// reinsert from temporary slice
	for _,v := range contentsBak {
		tile.add(v)
	}
}


func (tile *QuadTile) query(qbox AABB, ret []QuadElement) []QuadElement {
	// end recursion if this tile does not intersect the query range
	if ! tile.Intersects(qbox) {
		return ret
	}

	// return candidates at this tile
	for _, v := range tile.contents {
		if qbox.Intersects(v.GetAABB()) {
			ret = append(ret, v)
		}
	}

	// recurse into childs (if any)
	if tile.childs[_TOPRIGHT] != nil {
		for _, child := range tile.childs {
			ret = child.query(qbox, ret)
		}
	}

	return ret
}

func (qt *QuadTile) reduce(f func(a interface{}, t QuadElement) interface{}, v interface{}) interface{} {
	stack := list.New()
	stack.PushFront(qt)
	for stack.Len() > 0 {
		tile := stack.Front().Value.(*QuadTile)
		stack.Remove(stack.Front())

		for _, t := range tile.contents {
			v = f(v, t)
		}

		if tile.childs[_TOPRIGHT] != nil {
			for _, child := range tile.childs {
				stack.PushFront(child)
			}
		}
	}
	return v
}
*/
