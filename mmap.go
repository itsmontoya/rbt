package whiskey

import (
	"os"
	"path"

	"github.com/edsrzf/mmap-go"
	"github.com/missionMeteora/journaler"
	"github.com/missionMeteora/toolkit/errors"
)

// newMMap will return a new Mmap
func newMMap(dir, name string) (mp *MMap, err error) {
	var m MMap
	if m.f, err = os.OpenFile(path.Join(dir, name), os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return
	}

	mp = &m
	return
}

// MMap manages the memory mapped file
type MMap struct {
	f   *os.File
	mm  mmap.MMap
	cap int64
}

func (m *MMap) unmap() (err error) {
	if m.mm == nil {
		return
	}

	return m.mm.Unmap()
}

func (m *MMap) grow(sz int64) (bs []byte) {
	var err error
	if m.cap == 0 {
		var fi os.FileInfo
		if fi, err = m.f.Stat(); err != nil {
			journaler.Error("Stat error: %v", err)
			return
		}

		if m.cap = fi.Size(); m.cap == 0 {
			m.cap = sz
		}
	}

	for m.cap < sz {
		m.cap *= 2
	}

	if err = m.unmap(); err != nil {
		journaler.Error("Unmap error: %v", err)
		return
	}

	if err = m.f.Truncate(m.cap); err != nil {
		journaler.Error("Truncate error: %v", err)
		return
	}

	if m.mm, err = mmap.Map(m.f, os.O_RDWR, 0); err != nil {
		journaler.Error("Map error: %v", err)
		return
	}

	return m.mm
}

// Close will close an MMap
func (m *MMap) Close() (err error) {
	if m.f == nil {
		return errors.ErrIsClosed
	}

	var errs errors.ErrorList
	errs.Push(m.mm.Flush())
	errs.Push(m.mm.Unmap())
	errs.Push(m.f.Close())
	m.f = nil
	return
}
