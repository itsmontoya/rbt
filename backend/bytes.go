package backend

// NewBytes will return a new byteslice backend
func NewBytes() *Bytes {
	var b Bytes
	return &b
}

// Bytes manages a byteslice backend
type Bytes []byte

// Grow will grow the byteslice to the requested size
func (b *Bytes) Grow(sz int64) (bs []byte) {
	cap := int64(cap(*b))
	if cap == 0 {
		cap = sz
	}

	for cap < sz {
		cap *= 2
	}

	bs = make([]byte, cap)
	copy(bs, *b)
	*b = bs
	return
}

// Close will close the bytes
func (b *Bytes) Close() (err error) {
	*b = nil
	return
}
