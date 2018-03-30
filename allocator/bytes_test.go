package allocator

import (
	"testing"
)

func TestBytes(t *testing.T) {
	bs := NewBytes(0)
	if bs.tail != 0 {
		t.Fatalf("invalid tail, expected %d and received %d", 0, bs.tail)
	}

	if bs.cap != 32 {
		t.Fatalf("invalid cap, expected %d and received %d", 32, bs.cap)
	}

	sec, grew := bs.Allocate(16)
	if grew {
		t.Fatal("unexpected grow")
	}

	if sec.Size != 16 {
		t.Fatalf("invalid size, expected %d and received %d", 16, sec.Size)
	}
}
