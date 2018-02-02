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

func (t *Txn) getRoot(key []byte, sz int64) (bs []byte) {
	return t.r.Get(key)
}

func (t *Txn) truncateScratch(key []byte, sz int64) (bs []byte) {
	t.w.Grow(key, sz)
	return t.w.Get(key)
}

func (t *Txn) truncateRoot(key []byte, sz int64) (bs []byte) {
	t.r.Grow(key, sz)
	return t.r.Get(key)
}

// Bucket will return a bucket for a provided key
func (t *Txn) Bucket(key []byte) (bp *Bucket) {
	t.setKeyBuffer(key)

	var rgfn, sgfn GrowFn
	if t.r != nil {
		rgfn = t.getRoot
	}

	if t.w != nil {
		sgfn = t.truncateScratch
	}

	return newBucket(t.kbuf, rgfn, sgfn)
}

// CreateBucket will create a bucket for a provided key
func (t *Txn) CreateBucket(key []byte) (bp *Bucket, err error) {
	if t.w == nil {
		err = ErrCannotWrite
		return
	}

	t.setKeyBuffer(key)

	var rgfn, sgfn GrowFn
	if t.r != nil {
		rgfn = t.getRoot
	}

	if t.w != nil {
		sgfn = t.truncateScratch
	}

	bp = newBucket(t.kbuf, rgfn, sgfn)
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

	bkt := newBucket(key, t.truncateRoot, t.truncateScratch)
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
