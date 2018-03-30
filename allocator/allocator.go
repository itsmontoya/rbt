package allocator

import (
	"unsafe"
)

const (
	// SectionSize is the size (in bytes) of the pair struct
	SectionSize = int64(unsafe.Sizeof(Section{}))
)

// Allocator is a allocating interface
type Allocator interface {
	Get(offset, sz int64) []byte
	// Grow will grow the underlying bytes, but will not adjust the tail
	Grow(sz int64) (grew bool)
	// Ensures the tail is at the size or greater, will grow if necessary
	EnsureSize(sz int64) (grew bool)

	// Allocate will allocate a new section
	Allocate(sz int64) (s Section, grew bool)
	// Release will release a section (and it's bytes)
	Release(Section)

	// Function to be called on grow
	OnGrow(fn func())

	Close() (err error)
}
