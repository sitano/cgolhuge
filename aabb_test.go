package main

import "testing"

func TestAABB(t *testing.T) {
	a := NewAABB( 0, 0, 10, 10 )

	if ! a.ContainsPoint(0, 0) || ! a.ContainsPoint(10, 0) || ! a.ContainsPoint(0, 10) || ! a.ContainsPoint(10, 10) {
		t.Errorf("%v does not contain edge points")
	}

	b := NewAABB( 4, 4, 6, 6 )    // b completely within a

	if a.Intersection(b) != b {
		t.Errorf("%v does not intersect %v", a, b)
	}

	if ! a.Intersects(b) || ! b.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, b)
	}

	if ! a.Intersects(a) {
		t.Errorf("%v does not intersect itself", a)
	}

	if ! a.Contains(b) {
		t.Errorf("%v does not contain %v", a, b)
	}

	if ! a.Contains(a) {
		t.Errorf("%v does not contain itself", a)
	}

	if b.Contains(a) {
		t.Errorf("%v contains %v", b, a)
	}

	c := NewAABB( 10, 0, 20, 10 )

	if ! a.Intersects(c) {
		t.Errorf("%v does intersect %v", a, c)
	}
	if ! c.Intersects(a) {
		t.Errorf("%v does intersect %v", c, a)
	}

	if a.Contains(c) || c.Contains(a) {
		t.Errorf("%v contains %v (or vise versa)", a, c)
	}

	d := NewAABB( 0, 0, 0, 10 )

	if ! a.Intersects(d) {
		t.Errorf("%v does intersect %v", a, d)
	}
	if ! d.Intersects(a) {
		t.Errorf("%v does intersect %v", d, a)
	}

	e := NewAABB( 10, 10, 15, 15 )

	if ! a.Intersects(e) || ! e.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, e)
	}

	f := NewAABB( 0, 4, 20, 6 )

	if  ! a.Intersects(f) || ! f.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, f)
	}
}
