package allocator

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
	fl freelist

	bs   []byte
	tail int64
	cap  int64

	onGrow []OnGrowFn
}

func (b *Bytes) clean(s Section) {
	bs := b.bs[s.Offset : s.Offset+s.Size]

	for i := range bs {
		bs[i] = 0
	}
}

// Grow will grow the bytes
func (b *Bytes) Grow(sz int64) (grew bool) {
	if sz < b.cap {
		return
	}

	for b.cap <= sz {
		b.cap *= 2
	}

	bs := make([]byte, b.cap)
	copy(bs, b.bs)
	b.bs = bs

	for _, fn := range b.onGrow {
		fn()
	}

	return true
}

// EnsureSize will ensure the tail is at least at the requested size or greater
func (b *Bytes) EnsureSize(sz int64) (grew bool) {
	if b.tail >= sz {
		return
	}

	b.tail = sz
	return b.Grow(sz)
}

// Get will get bytes
func (b *Bytes) Get(offset, sz int64) []byte {
	return b.bs[offset : offset+sz]
}

// Allocate will allocate bytes
func (b *Bytes) Allocate(sz int64) (s Section, grew bool) {
	s.Size = sz
	if s.Offset = b.fl.acquire(sz); s.Offset != -1 {
		return
	}

	s.Offset = b.tail
	b.tail += sz
	grew = b.Grow(b.tail)
	return
}

// Release will release a section
func (b *Bytes) Release(s Section) {
	b.clean(s)
	b.fl.release(s)
	return
}

// OnGrow will append a function to be called on grows
func (b *Bytes) OnGrow(fn OnGrowFn) {
	b.onGrow = append(b.onGrow, fn)
}

// Close will close bytes
func (b *Bytes) Close() (err error) {
	return
}
