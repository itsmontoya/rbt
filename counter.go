package rbt

import (
	"sync/atomic"
	"unsafe"

	"github.com/missionMeteora/toolkit/errors"
)

func newCounter(bs []byte) (c counter) {
	c.v = (*int64)(unsafe.Pointer(&bs))
	return
}

type counter struct {
	v *int64
}

func (c *counter) Get() (n int64) {
	return atomic.LoadInt64(c.v)
}

func (c *counter) Set(n int64) {
	atomic.StoreInt64(c.v, n)
}

func (c *counter) Increment() (new int64) {
	return atomic.AddInt64(c.v, 1)
}

func (c *counter) Decrement() (new int64) {
	return atomic.AddInt64(c.v, -1)
}

func (c *counter) Close() (err error) {
	if c.v == nil {
		return errors.ErrIsClosed
	}

	c.v = nil
	return
}
