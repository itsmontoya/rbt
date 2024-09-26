package allocator

import (
	"unsafe"

	"github.com/itsmontoya/rbt/backend"
)

// NewSimple will return a new allocator
func NewSimple(b backend.Backend, sz int64) *Simple {
	var a Simple
	a.b = b
	a.grow(sz + headerSize)
	return &a
}

// Simple will allocate
type Simple struct {
	b  backend.Backend
	bs []byte
	*header

	preGrow  func()
	postGrow func()
}

func (s *Simple) setHeader() {
	s.header = (*header)(unsafe.Pointer(&s.bs[0]))
	if s.tail == 0 {
		s.tail = headerSize
	}

	s.hmm = 37
	s.cap = int64(len(s.bs))
}

func (s *Simple) allocate(sz int64) (offset int64, grew bool) {
	offset = s.tail
	if s.tail += sz; s.tail < s.cap {
		return
	}

	s.Grow(s.tail)
	grew = true
	return
}

// Grow will ensure the backing slice is at least a provided size
func (s *Simple) grow(sz int64) {
	s.bs = s.b.Grow(sz)
	s.setHeader()
}

// Allocate will allocate a set of bytes for a given size
func (s *Simple) Allocate(sz int64) (offset int64, grew bool) {
	offset, grew = s.allocate(sz)
	return
}

// Byte will return a pointer to the byte at a given offset
func (s *Simple) Byte(offset int64) *byte {
	return &s.bs[offset]
}

// First will return the first byte following the header
func (s *Simple) First() *byte {
	return &s.bs[headerSize]
}

// Grow will ensure the backing slice is at least a provided size
func (s *Simple) Grow(sz int64) (grew bool) {
	if s.cap > sz {
		return
	}

	if s.preGrow != nil {
		s.preGrow()
	}

	s.grow(sz)

	if s.postGrow != nil {
		s.postGrow()
	}

	return true
}

// Len will return the length
func (s *Simple) Len() int64 {
	return s.tail - headerSize
}

// Bytes will return a byteslice reference
func (s *Simple) Bytes(offset, sz int64) []byte {
	return s.bs[offset : offset+sz]
}

// Release will release a set of bytes at a given offset
func (s *Simple) Release(offset, sz int64) {
	return
}

// Reset will reset the allocator
func (s *Simple) Reset() {
	s.tail = headerSize
}

// OnPreGrow will set the pre-grow func
func (s *Simple) OnPreGrow(fn func()) {
	s.preGrow = fn
}

// OnPostGrow will set the post-grow func
func (s *Simple) OnPostGrow(fn func()) {
	s.postGrow = fn
}
