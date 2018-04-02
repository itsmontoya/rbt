package backend

import (
	"unsafe"

	"github.com/Path94/atoms"
	"github.com/itsmontoya/rbt/allocator"
	"github.com/missionMeteora/toolkit/errors"
)

// NewMulti will return a new Multi
func NewMulti(a allocator.Allocator) *Multi {
	var m Multi
	m.a = a
	m.a.EnsureSize(allocator.SectionSize)
	m.a.OnGrow(m.onGrow)
	m.onGrow()
	return &m
}

// Multi is a multiple backend manager
type Multi struct {
	a allocator.Allocator
	s *allocator.Section

	closed atoms.Bool
}

func (m *Multi) onGrow() (end bool) {
	if m.closed.Get() {
		return true
	}

	m.s = (*allocator.Section)(unsafe.Pointer(&m.a.Get(0, 1)[0]))
	return
}

// New will return a new backend
func (m *Multi) New() (b *Backend) {
	return New(m)
}

// Get will get the current backend
func (m *Multi) Get() (b *Backend) {
	if m.s == nil {
		return m.New()
	}

	b = New(m)
	b.s = *m.s
	b.SetBytes()
	return
}

// Set will set the primary backend
func (m *Multi) Set(b *Backend) {
	// This is fairly cheap, and it's safe to ensure the bytes are correct
	// Segfaults are bad, and I should feel bad. *sad panda*
	m.onGrow()
	// Set pair values
	m.s.Offset = b.s.Offset
	m.s.Size = b.s.Size
}

// Close will close multi
func (m *Multi) Close() (err error) {
	if !m.closed.Set(true) {
		return errors.ErrIsClosed
	}

	return
}
