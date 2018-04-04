package rbt

import "testing"

func TestCounter(t *testing.T) {
	var arr [8]byte
	c := newCounter(arr[:])
	c.Set(0)
	if n := c.Increment(); n != 1 {
		t.Fatalf("invalid value, expected %d and received %d", 1, n)
	}

	if n := c.Increment(); n != 2 {
		t.Fatalf("invalid value, expected %d and received %d", 2, n)
	}

	if n := c.Increment(); n != 3 {
		t.Fatalf("invalid value, expected %d and received %d", 3, n)
	}

	if n := c.Decrement(); n != 2 {
		t.Fatalf("invalid value, expected %d and received %d", 2, n)
	}

	if n := c.Decrement(); n != 1 {
		t.Fatalf("invalid value, expected %d and received %d", 1, n)
	}

	if n := c.Decrement(); n != 0 {
		t.Fatalf("invalid value, expected %d and received %d", 0, n)
	}

	c.Set(3)

	c.Close()
	c = newCounter(arr[:])

	if n := c.Get(); n != 3 {
		t.Fatalf("invalid value, expected %d and received %d", 3, n)
	}

	if n := c.Decrement(); n != 2 {
		t.Fatalf("invalid value, expected %d and received %d", 2, n)
	}

	if n := c.Decrement(); n != 1 {
		t.Fatalf("invalid value, expected %d and received %d", 1, n)
	}

	if n := c.Decrement(); n != 0 {
		t.Fatalf("invalid value, expected %d and received %d", 0, n)
	}

}
