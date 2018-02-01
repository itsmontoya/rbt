package glass

import (
	"github.com/itsmontoya/whiskey"
)

const (
	bucketInitSize = 256
)

func newBucket(key, rbs, wbs []byte, gfn GrowFn) *Bucket {
	var b Bucket
	// Because the provided key is a reference to the DB's key buffer, we will need to copy
	// the contents into a new slice so that we don't encounter any race conditions later

	// Make a new byteslice with the length of the provided key
	b.key = make([]byte, len(key))
	// Copy key buffer to key
	copy(b.key, key)
	if rbs != nil {
		b.rbs = rbs
		b.r = whiskey.NewRaw(bucketInitSize, b.growSlave, nil)
	}

	if wbs != nil {
		b.wbs = wbs
		b.w = whiskey.NewRaw(bucketInitSize, b.grow, nil)
		b.gfn = gfn
	}

	return &b
}

// Bucket represents a database bucket
type Bucket struct {
	Txn

	key []byte
	wbs []byte
	rbs []byte

	rgfn GrowFn
	gfn  GrowFn
}

func (b *Bucket) grow(sz int64) (bs []byte) {
	if sz <= int64(len(b.wbs)) {
		return b.wbs
	}

	n := int64(len(b.wbs))
	for n < sz {
		n *= 2
	}

	bs = b.gfn(b.key, n)
	copy(bs, b.wbs)
	b.wbs = bs
	return
}

func (b *Bucket) growSlave(sz int64) (bs []byte) {
	if sz <= int64(len(b.rbs)) {
		return b.rbs
	}

	if b.rgfn == nil {
		panic("slave is attempting to grow past it's intended size")
	}

	n := int64(len(b.rbs))
	for n < sz {
		n *= 2
	}

	// This saves segfault, we need to figure out why
	// Note: This is a huge performance regression
	bbs := make([]byte, len(b.rbs))
	copy(bbs, b.rbs)
	bs = b.rgfn(b.key, n)
	copy(bs, bbs)
	b.rbs = bs
	return
}

// Close will close a bucket
func (b *Bucket) Close() (err error) {
	if err = b.w.Close(); err != nil {
		return
	}

	// Technically we will panic if close is called twice.
	// Let's take some time to really plan this and see how we
	// want to approach this
	b.w = nil

	b.key = nil
	b.rbs = nil
	b.r = nil
	b.wbs = nil
	b.w = nil
	b.gfn = nil
	return
}

// GrowFn is called on grows
type GrowFn func(key []byte, sz int64) []byte
