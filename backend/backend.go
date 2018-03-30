package backend

import (
	"github.com/itsmontoya/rbt/allocator"
	"github.com/missionMeteora/toolkit/errors"
)

// New will return a new backend
func New(m *Multi) *Backend {
	var b Backend
	b.m = m
	m.a.OnGrow(b.SetBytes)
	return &b
}

// Backend represents a data Backend
type Backend struct {
	m  *Multi
	s  allocator.Section
	bs []byte
}

// SetBytes will refresh the bytes reference
func (b *Backend) SetBytes() {
	b.bs = b.m.a.Get(b.s.Offset, b.s.Size)
}

// Bytes are the current bytes
func (b *Backend) Bytes() []byte {
	return b.bs
}

// Section will return a section
func (b *Backend) Section() allocator.Section {
	return b.s
}

func (b *Backend) allocate(sz int64) (bs []byte) {
	var (
		ns   allocator.Section
		grew bool
	)

	if ns, grew = b.m.a.Allocate(sz); grew {
		b.SetBytes()
	}

	bs = b.m.a.Get(ns.Offset, ns.Size)

	if b.s.Size > 0 {
		// Copy old bytes to new byteslice
		copy(bs, b.bs)
		// Release old bytes to allocator
		b.m.a.Release(b.s)
	}

	b.s = ns
	b.bs = bs
	return
}

// Grow will grow the backend
func (b *Backend) Grow(sz int64) (bs []byte) {
	if !b.s.IsEmpty() {
		bs = b.m.a.Get(b.s.Offset, b.s.Size)
	}

	var cap int64
	if cap = nextCap(b.s.Size, sz); cap == -1 {
		return
	}

	bs = b.allocate(cap)
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
	copy(out.bs, b.bs)
	return
}

// Destroy will destroy a backend and it's contents
func (b *Backend) Destroy() (err error) {
	if b.m == nil {
		return errors.ErrIsClosed
	}

	b.m.a.Release(b.s)
	b.m = nil
	return
}

// Close will close an Backend
func (b *Backend) Close() (err error) {
	if b.m.a == nil {
		return errors.ErrIsClosed
	}

	b.m = nil
	return
}
