package glass

import (
	"github.com/itsmontoya/whiskey"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrCannotWrite is returned when a write action is attempted during a read transaction
	ErrCannotWrite = errors.Error("cannot write during a read transaction")
)

// Txn is a transaction type
type Txn struct {
	r *whiskey.Whiskey
	w *whiskey.Whiskey

	kbuf []byte
}

func (t *Txn) setKeyBuffer(key []byte) {
	// Reset before using
	t.kbuf = t.kbuf[:0]
	// Append bucket prefix
	t.kbuf = append(t.kbuf, bucketPrefix)
	// Append key
	t.kbuf = append(t.kbuf, key...)
}

func (t *Txn) getBucketBytes(key []byte) (rbs, wbs []byte) {
	t.setKeyBuffer(key)
	if t.r != nil {
		rbs = t.r.Get(t.kbuf)
	}

	if t.w != nil {
		// This is a write transaction, let's check if this value has been changed
		wbs = t.w.Get(t.kbuf)
	}

	return
}

func (t *Txn) truncate(key []byte, sz int64) (bs []byte) {
	t.w.Put(key, make([]byte, sz))
	return t.w.Get(key)
}

func (t *Txn) truncateRoot(key []byte, sz int64) (bs []byte) {
	t.r.Put(key, make([]byte, sz))
	return t.r.Get(key)
}

// Bucket will return a bucket for a provided key
func (t *Txn) Bucket(key []byte) (bp *Bucket) {
	t.setKeyBuffer(key)
	rbs, wbs := t.getBucketBytes(t.kbuf)
	if rbs == nil && wbs == nil {
		return
	}

	if wbs == nil && t.w != nil {
		t.w.Put(t.kbuf, make([]byte, bucketInitSize))
		wbs = t.w.Get(t.kbuf)
	}

	return newBucket(t.kbuf, rbs, wbs, t.truncate)
}

// CreateBucket will create a bucket for a provided key
func (t *Txn) CreateBucket(key []byte) (bp *Bucket, err error) {
	if t.w == nil {
		err = ErrCannotWrite
		return
	}

	t.setKeyBuffer(key)
	rbs, wbs := t.getBucketBytes(t.kbuf)
	if rbs == nil && wbs == nil {
		t.w.Put(t.kbuf, make([]byte, bucketInitSize))
		wbs = t.w.Get(t.kbuf)
	}

	bp = newBucket(t.kbuf, rbs, wbs, t.truncate)
	return
}

// Get will retrieve a value for a given key
func (t *Txn) Get(key []byte) (val []byte, err error) {
	if key[0] == bucketPrefix {
		return nil, ErrInvalidKey
	}

	if t.w != nil {
		if val = t.w.Get(key); val != nil {
			return
		}
	}

	if t.r != nil {
		val = t.r.Get(key)
	}

	return
}

// Put will put a value for a given key
func (t *Txn) Put(key []byte, val []byte) (err error) {
	if key[0] == bucketPrefix {
		return ErrInvalidKey
	}

	if t.w == nil {
		return ErrCannotWrite
	}

	t.w.Put(key, val)
	return
}

func (t *Txn) writeEntry(key, val []byte) (end bool) {
	if key[0] != bucketPrefix {
		// Flush value to main branch
		t.r.Put(key, val)
		return
	}

	rbs := t.r.Get(key)
	if rbs == nil {
		t.r.Put(t.kbuf, make([]byte, bucketInitSize))
		rbs = t.r.Get(key)

	}

	bkt := newBucket(key, rbs, val, t.truncate)
	bkt.rgfn = t.truncateRoot
	end = bkt.w.ForEach(func(key, val []byte) (end bool) {
		return bkt.writeEntry(key, val)
	})

	return
}

func (t *Txn) flush() (err error) {
	t.w.ForEach(t.writeEntry)
	return
}

// TxnFn is a transaction func
type TxnFn func(txn *Txn) error
