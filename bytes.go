package rbTree

import (
	"github.com/missionMeteora/journaler"
)

// NewBytes will return a new byteslice backend
func NewBytes() *Bytes {
	var b Bytes
	return &b
}

// Bytes manages a byteslice backend
type Bytes []byte

func (b *Bytes) grow(sz int64) (bs []byte) {
	cap := int64(cap(*b))
	for cap < sz {
		cap *= 2
	}

	*b = make([]byte, cap)
	bs = *b
	journaler.Error("Grow")
	return
}
