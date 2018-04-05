package allocator

import (
	"sort"
)

type freelist struct {
	s []*Section
}

func (f *freelist) acquire(sz int64) (offset int64) {
	for i, s := range f.s {
		switch {
		case s.Size < sz:
			// Pair size is smaller than requested size, continue
			continue
		case s.Size == sz:
			// Pair size is the same as the requested size
			// Set offset as pair
			offset = s.Offset
			// Remove pair from freelist
			f.remove(i)
		case s.Size > sz:
			// Pair size is bigger than requested size,
			// Set offset as pair
			offset = s.Offset
			// Move pair's offset up by the size amount
			s.Offset += sz
			// Reduce pair's size by the size amount
			s.Size -= sz
		}

		return
	}

	return -1
}

func (f *freelist) release(s Section) {
	f.s = append(f.s, &s)
	f.sort()
	f.merge()
}

func (f *freelist) sort() {
	sort.Slice(f.s, f.isLess)
}

func (f *freelist) merge() {
	var (
		last    *Section
		removed int
	)

	for i, s := range f.s {
		if last == nil || (s.Offset != last.Offset+last.Size) {
			last = s
			continue
		}

		last.Size += s.Size
		f.remove(i - removed)
		removed++
	}
}

func (f *freelist) remove(i int) {
	f.s = append(f.s[:i], f.s[i+1:]...)
}

func (f *freelist) isLess(i, j int) bool {
	return f.s[i].Offset < f.s[j].Offset
}
