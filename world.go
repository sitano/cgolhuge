package main

type World interface {
	Generation() uint64
	Population() uint64

	LifePoints(z byte) [][2]int64
	View() View
	ZLayer() byte

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

func (lw *LifeWorld) ZLayer() byte {
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

	cz ++
	panic("Not implemented")
	// TODO: traverse tree tiles -> recalc state per basis

	lw.z = nz
}
