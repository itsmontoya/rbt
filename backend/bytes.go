package backend

// NewBytes will return a new byteslice backend
func NewBytes(sz int64) *Bytes {
	var b Bytes
	if sz == 0 {
		sz = 32
	}

	b.bs = make([]byte, sz)
	b.cap = sz
	return &b
}

// Bytes manages a byteslice backend
type Bytes struct {
	bs   []byte
	tail int64
	cap  int64

	listeners []func() (unsub bool)
}

// grow will grow the bytes
func (b *Bytes) grow() (grew bool) {
	if b.tail < b.cap {
		return
	}

	for b.cap < b.tail {
		b.cap *= 2
	}

	bs := make([]byte, b.cap)
	copy(bs, b.bs)
	b.bs = bs
	return true
}

// Get will get bytes
func (b *Bytes) Get(offset, sz int64) []byte {
	return b.bs[offset : offset+sz]
}

// Allocate will allocate bytes
func (b *Bytes) Allocate(sz int64) (sp *Section, grew bool) {
	var s Section
	s.Offset = b.tail
	b.tail += sz
	grew = b.grow()

	s.Bytes = b.bs[s.Offset : s.Offset+sz]
	s.Size = sz
	sp = &s
	return
}

// Release will release a section
func (b *Bytes) Release(s *Section) {
	// Right now we just ignore it and let this grow
	return
}

// Close will close bytes
func (b *Bytes) Close() (err error) {
	return
}
