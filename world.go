package main

import (
	"io"
	"math"
)

const (
	DEAD = byte(0)
	LIFE = byte(1)

	RULE_LIVE_MIN = 2
	RULE_LIVE_MAX = 3
	RULE_BORN = 3

	BITS10  = uint64(0x2) // 0b10

	BITS1   = uint64(0x1) //   0b1
	BITS11  = uint64(0x3) //  0b11
	BITS111 = uint64(0x7) // 0b111
	BITS101 = uint64(0x5) // 0b101

	BITS01  = BITS1

	PageStrideLast2 = PageStrideBits - 2
	PageStrideLast1 = PageStrideBits - 1
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

func NewLifeWorldWH(wbox AABB) *LifeWorld {
	return NewLifeWorld(NewWorldView(NewVM(), wbox))
}

func NewLifeWorldMax() *LifeWorld {
	return NewLifeWorldWH(New00WH(math.MaxUint64, math.MaxUint64))
}

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

	// TODO: reclaim pages
	// TODO: do not allocate next buf if there were no changes
	for pi := 0; pi < w.v.vm.Pages(); pi ++ {
		p := w.v.vm.reserved[pi]
		p.next = NewPageBuf()

		raw := p.raw
		next := p.next

		ap_nw := p.ap_nw
		ap_n  := p.ap_n
		ap_ne := p.ap_ne
		ap_w  := p.ap_w
		ap_e  := p.ap_e
		ap_sw := p.ap_sw
		ap_s  := p.ap_s
		ap_se := p.ap_se

		prev_w_line := uint64(0)
		prev_line := uint64(0)
		prev_e_line := uint64(0)

		if ap_nw != nil {
			prev_w_line = ap_nw.raw[PageStrides - 1]
		}
		if ap_n != nil {
			prev_line = ap_n.raw[PageStrides - 1]
		}
		if ap_ne != nil {
			prev_e_line = ap_ne.raw[PageStrides - 1]
		}

		curr_w_line := uint64(0)
		curr_line := p.raw[0]
		curr_e_line := uint64(0)

		if ap_w != nil {
			curr_w_line = ap_w.raw[0]
		}
		if ap_e != nil {
			curr_e_line = ap_e.raw[0]
		}

		next_w_line := uint64(0)
		next_line := uint64(0)
		next_e_line := uint64(0)

		ci := 0
		last_line := false
		for !last_line {
			new_line := uint64(0)

			if ci < PageStrides - 1 {
				next_line = raw[ci + 1]

				if ap_w != nil {
					next_w_line = ap_w.raw[ci + 1]
				}
				if ap_e != nil {
					next_e_line = ap_e.raw[ci + 1]
				}
			} else {
				// Last next line on the next page
				next_w_line = 0
				next_line = 0
				next_e_line = 0

				if ap_sw != nil {
					next_w_line = ap_sw.raw[0]
				}
				if ap_s != nil {
					next_line = ap_s.raw[0]
				}
				if ap_se != nil {
					next_e_line = ap_se.raw[0]
				}

				last_line = true
			}

			// Process 1 stride line if there are anything to process
			sum_west :=
				(prev_w_line >> PageStrideLast1) +
				(curr_w_line >> PageStrideLast1) +
				(next_w_line >> PageStrideLast1)
			sum_middle := prev_line | curr_line | next_line
			sum_east :=
				(prev_e_line & BITS1) +
				(curr_e_line & BITS1) +
				(next_e_line & BITS1)

			if sum_middle | sum_west != 0 {
				// First 2 bits with mask 0b011 (west edge)
				// Test: go build && ./cgolhuge -load pattern/glider_gun.lif -lx 17 -ly 5 -wait
				sum := sum_west + PopCnt(
						(prev_line & BITS11) |
						(curr_line & BITS10) << 4 |
						(next_line & BITS11) << 8)

				if sum >= RULE_LIVE_MIN {
					st := byte(curr_line & BITS1)

					if st == DEAD {
						if sum == RULE_BORN {
							new_line |= BITS1
							w.population ++
						}
					} else {
						if sum <= RULE_LIVE_MAX {
							new_line |= BITS1
							w.population ++
						}
					}

					// Outer west edge check
					if ci == 0 && ap_nw == nil && ap_n != nil && ap_w != nil &&
						p.px > w.v.pb.MinX && p.py > w.v.pb.MinY &&
						(prev_line & BITS1) +
						(curr_line & BITS1) +
						(curr_w_line >> PageStrideLast1) == RULE_BORN {
						w.v.Set(p.MinX - 1, p.MinY - 1, DEAD)
						ap_nw = p.ap_nw
					}
					if ap_w == nil &&
						p.px > w.v.pb.MinX &&
						(prev_line & BITS1) +
						(curr_line & BITS1) +
						(next_line & BITS1) == RULE_BORN {
						w.v.Set(p.MinX - 1, p.MinY, DEAD)
						ap_w = p.ap_w
					}
					if ci == PageStrides - 1 && ap_sw == nil && ap_s != nil && ap_w != nil &&
						p.px < w.v.pb.MaxX && p.py < w.v.pb.MaxY &&
						(curr_line & BITS1) +
						(next_line & BITS1) +
						(curr_w_line >> PageStrideLast1) == RULE_BORN {
						w.v.Set(p.MinX - 1, p.MaxY + 1, DEAD)
						ap_sw = p.ap_sw
					}
				}
			}

			if sum_middle != 0 {
				// Middle bits
				pl := prev_line
				cl := curr_line
				nl := next_line
				for bi := uint(1); bi < PageStrideBits - 1; bi ++ {
					sum := PopCnt((pl & BITS111) << 8 | (cl & BITS101) << 4 | (nl & BITS111))

					if sum >= RULE_LIVE_MIN {
						st := byte((cl >> 1) & BITS1)

						if st == DEAD {
							if sum == RULE_BORN {
								new_line |= BITS1 << bi
								w.population ++
							}
						} else {
							if sum <= RULE_LIVE_MAX {
								new_line |= BITS1 << bi
								w.population ++
							}
						}
					}

					// Outer north / south edge check
					// Can't inject this under common sum >= condition as middle bit on CL was skipped
					if ci == 0 && ap_n == nil && p.py > w.v.pb.MinY && PopCnt(cl & BITS111) == RULE_BORN {
						w.v.Set(p.MinX, p.MinY - 1, DEAD)
						ap_n = p.ap_n
					}
					if ci == PageStrides - 1 && ap_s == nil && p.py < w.v.pb.MaxY && PopCnt(cl & BITS111) == RULE_BORN {
						w.v.Set(p.MinX, p.MaxY + 1, DEAD)
						ap_s = p.ap_s
					}

					// Shift to next point (from west -> east)
					pl >>= 1
					cl >>= 1
					nl >>= 1
				}
			}

			if sum_middle | sum_east != 0 {
				// Last 2 bits with mask 0b011 (east edge)
				// Test: go build && ./cgolhuge -load pattern/glider_gun.lif -lx 45 -ly 5 -wait
				sum := sum_east + PopCnt(
							((prev_line >> PageStrideLast2) & BITS11) |
							((curr_line >> PageStrideLast2) & BITS01) << 4 |
							((next_line >> PageStrideLast2) & BITS11) << 8)

				if sum >= RULE_LIVE_MIN {
					st := byte(curr_line >> PageStrideLast1)

					if st == DEAD {
						if sum == RULE_BORN {
							new_line |= BITS1 << PageStrideLast1
							w.population ++
						}
					} else {
						if sum <= RULE_LIVE_MAX {
							new_line |= BITS1 << PageStrideLast1
							w.population ++
						}
					}

					// Outer east edge check
					if ci == 0 && ap_ne == nil && ap_n != nil && ap_e != nil &&
						p.px < w.v.pb.MaxX && p.py > w.v.pb.MinY &&
						(prev_line >> PageStrideLast1) +
						(curr_line >> PageStrideLast1) +
						(curr_e_line & BITS1) == RULE_BORN {
						w.v.Set(p.MaxX + 1, p.MinY - 1, DEAD)
						ap_ne = p.ap_ne
					}
					if ap_e == nil &&
						p.px < w.v.pb.MaxX &&
						(prev_line >> PageStrideLast1) +
						(curr_line >> PageStrideLast1) +
						(next_line >> PageStrideLast1) == RULE_BORN {
						w.v.Set(p.MaxX + 1, p.MinY, DEAD)
						ap_e = p.ap_e
					}
					if ci == PageStrides - 1 && ap_se == nil && ap_s != nil && ap_e != nil &&
						p.px < w.v.pb.MaxX && p.py < w.v.pb.MaxY &&
						(curr_line >> PageStrideLast1) +
						(next_line >> PageStrideLast1) +
						(curr_e_line & BITS1) == RULE_BORN {
						w.v.Set(p.MaxX + 1, p.MaxY + 1, DEAD)
						ap_se = p.ap_se
					}
				}
			}

			next[ci] = new_line

			prev_w_line = curr_w_line
			prev_line = curr_line
			prev_e_line = curr_e_line

			curr_w_line = next_w_line
			curr_line = next_line
			curr_e_line = next_e_line

			ci ++
		}
	}

	w.Swap()
}

func (w *LifeWorld) Swap() {
	for _, p := range w.v.vm.reserved {
		p.raw = p.next
		p.next = nil
	}
}

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
