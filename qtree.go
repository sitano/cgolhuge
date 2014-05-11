/*
Based on work of Volker Poplawski, 2013 (https://github.com/volkerp/goquadtree)
*/
package main

import "container/list"
import "math"
import "fmt"

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

// Use AABB() to construct a AABB object
// Min - inclusive, Max - exlusive
type AABB struct {
	MinX, MaxX, MinY, MaxY int64
}

func (b AABB) String() string {
	return fmt.Sprintf("[AABB x=%d, y=%d / %d, %d]", b.MinX, b.MinY, b.MaxX, b.MaxY)
}

func Min(x int64, y int64) int64 {
	if x <= y {
		return x
	}

	return y
}

func Max(x int64, y int64) int64 {
	if x >= y {
		return x
	}

	return y
}

func Abs(x int64) uint64 {
	if x >= 0 {
		return uint64(x)
	}
	return uint64(-1 * x)
}

func NewAABB(xa, xb, ya, yb int64) AABB {
	return AABB{ Min(xa, xb), Max(xa, xb), Min(ya, yb), Max(ya, yb) }
}

func NewAABBMax() AABB {
	return NewAABB(math.MinInt64, math.MaxInt64, math.MaxInt64, math.MinInt64)
}

// Make AABB implement the QuadElement interface
func (b AABB) GetAABB() AABB {
	return b
}

func (b *AABB) SizeX() uint64 {
	// Check overflow
	if  b.MinX == math.MinInt64 && b.MaxX == math.MaxInt64 {
		return math.MaxUint64
	}
	return Abs(b.MaxX - b.MinX)
}

func (b *AABB) SizeY() uint64 {
	// Check overflow
	if  b.MinY == math.MinInt64 && b.MaxY == math.MaxInt64 {
		return math.MaxUint64
	}
	return Abs(b.MaxY - b.MinY)
}

// Returns true if o intersects this
func (b *AABB) Intersects(o AABB) bool {
	return b.MinX < o.MaxX && b.MinY < o.MaxY &&
		b.MaxX > o.MinX && b.MaxY > o.MinY
}

// Returns true if o is within this
func (b *AABB) Contains(o AABB) bool {
	return b.MinX <= o.MinX && b.MinY <= o.MinY &&
		b.MaxX >= o.MaxX && b.MaxY >= o.MaxY
}

func (b *AABB) ContainsPoint(x int64, y int64) bool {
	return b.MinX <= x && b.MinY <= y &&
		b.MaxX >= x && b.MaxY >= y
}

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
