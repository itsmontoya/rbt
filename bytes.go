package rbt

// newBytes will return a new byteslice backend
func newBytes() *Bytes {
	var b Bytes
	return &b
}

// Bytes manages a byteslice backend
type Bytes []byte

func (b *Bytes) grow(sz int64) (bs []byte) {
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
