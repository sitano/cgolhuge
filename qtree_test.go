/*
Based on work of Volker Poplawski, 2013 (https://github.com/volkerp/goquadtree)
*/
package main

import "testing"

func TestAABB(t *testing.T) {
	a := NewAABB( 0, 10, 0, 10 )

	b := NewAABB( 4, 6, 4, 6 )    // b completely within a

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

	c := NewAABB( 10, 20, 0, 10 )

	if a.Intersects(c) {
		t.Errorf("%v does intersect %v", a, c)
	}
	if c.Intersects(a) {
		t.Errorf("%v does intersect %v", c, a)
	}

	if a.Contains(c) || c.Contains(a) {
		t.Errorf("%v contains %v (or vise versa)", a, c)
	}

	d := NewAABB( -10, 0, 0, 10 )

	if a.Intersects(d) {
		t.Errorf("%v does intersect %v", a, d)
	}
	if d.Intersects(a) {
		t.Errorf("%v does intersect %v", d, a)
	}

	e := NewAABB( 9, 15, 9, 15 )

	if ! a.Intersects(e) || ! e.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, e)
	}

	f := NewAABB( -10, 20, 4, 6 )

	if  ! a.Intersects(f) || ! f.Intersects(a) {
		t.Errorf("%v does not intersect %v", a, f)
	}
}
