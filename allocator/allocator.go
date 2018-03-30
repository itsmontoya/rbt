package allocator

import (
	"unsafe"
)

const (
	// PairSize is the size (in bytes) of the pair struct
	PairSize = int64(unsafe.Sizeof(Pair{}))
)

// Allocator is a allocating interface
type Allocator interface {
	Get(offset, sz int64) []byte
	// Grow will grow the underlying bytes, but will not adjust the tail
	Grow(sz int64) (grew bool)
	// Ensures the tail is at the size or greater, will grow if necessary
	EnsureSize(sz int64) (grew bool)

	// Allocate will allocate a new section
	Allocate(sz int64) (s *Section, grew bool)
	// Release will release a section (and it's bytes)
	Release(*Section)

	// Function to be called on grow
	OnGrow(fn func())

	Close() (err error)
}

// Pair is a data section pair
type Pair struct {
	Offset int64
	Size   int64
}

// IsEmpty will return if a pair is empty
func (p *Pair) IsEmpty() bool {
	if p.Offset != 0 {
		return false
	}

	if p.Size != 0 {
		return false
	}

	return true
}
