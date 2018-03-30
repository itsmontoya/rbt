package backend

import (
	"github.com/itsmontoya/rbt/allocator"
	"github.com/missionMeteora/journaler"
	"github.com/missionMeteora/toolkit/errors"
)

// New will return a new backend
func New(m *Multi) *Backend {
	var b Backend
	b.m = m
	b.s = &allocator.Section{}
	m.a.OnGrow(b.SetBytes)
	return &b
}

// Backend represents a data Backend
type Backend struct {
	m *Multi
	s *allocator.Section
}

// SetBytes will refresh the bytes reference
func (b *Backend) SetBytes() {
	offset := b.s.Offset
	journaler.Debug("Got offset!")
	size := b.s.Size
	journaler.Debug("Got size!")
	b.s.Bytes = b.m.a.Get(offset, size)

	journaler.Debug("Uhh ya.")
}

// Bytes are the current bytes
func (b *Backend) Bytes() []byte {
	return b.s.Bytes
}

// Section will return a section
func (b *Backend) Section() allocator.Section {
	return *b.s
}

func (b *Backend) allocateSection(sz int64) (ns *allocator.Section, grew bool) {
	ns, grew = b.m.a.Allocate(sz)
	if b.s.Size == 0 {
		return
	}

	if grew {
		b.SetBytes()
	}

	// Copy old bytes to new byteslice
	copy(ns.Bytes, b.s.Bytes)
	// Release old bytes to allocator
	b.m.a.Release(b.s)
	return
}

// Grow will grow the backend
func (b *Backend) Grow(sz int64) (bs []byte) {
	if !b.s.IsEmpty() {
		bs = b.s.Bytes
	}

	if cap := nextCap(b.s.Size, sz); cap > -1 {
		b.s, _ = b.allocateSection(cap)
	}

	bs = b.s.Bytes
	return
}

// Notify will notify the parent
func (b *Backend) Notify() {
	b.m.Set(b)
}

// Dup will duplicate a backend
func (b *Backend) Dup() (out *Backend) {
	out = New(b.m)
	out.Grow(b.s.Size)
	b.SetBytes()
	copy(out.s.Bytes, b.s.Bytes)
	return
}

// Destroy will destroy a backend and it's contents
func (b *Backend) Destroy() (err error) {
	if b.m == nil {
		return errors.ErrIsClosed
	}

	b.m.a.Release(b.s)
	b.m = nil
	b.s = nil
	return
}

// Close will close an Backend
func (b *Backend) Close() (err error) {
	if b.m.a == nil {
		return errors.ErrIsClosed
	}

	b.m = nil
	b.s = nil
	return
}
