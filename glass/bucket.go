package glass

import (
	"github.com/itsmontoya/whiskey"
)

func newBucket(key, bs []byte, gfn func(key []byte, sz int64) []byte) *Bucket {
	var b Bucket
	// Because the provided key is a reference to the DB's key buffer, we will need to copy
	// the contents into a new slice so that we don't encounter any race conditions later

	// Make a new byteslice with the length of the provided key
	b.key = make([]byte, len(key))
	// Copy key buffer to key
	copy(b.key, key)

	b.bs = bs
	b.gfn = gfn
	b.w = whiskey.NewRaw(256, b.grow, nil)
	return &b
}

// Bucket represents a database bucket
type Bucket struct {
	w   *whiskey.Whiskey
	key []byte
	bs  []byte
	gfn func(key []byte, sz int64) []byte
}

func (b *Bucket) grow(sz int64) (bs []byte) {
	if sz <= int64(len(b.bs)) {
		return b.bs
	}

	bs = b.gfn(b.key, sz)
	copy(bs, b.bs)
	b.bs = bs
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
	b.bs = nil
	b.gfn = nil
	return
}
