package allocator

// Allocator is an allocator interface
type Allocator interface {
	Allocate(sz int64) (offset int64, grew bool)
	First() *byte
	Byte(offset int64) *byte
	Grow(sz int64) (grew bool)
	Bytes(offset, sz int64) []byte
	Release(offset, sz int64)
	Len() int64
	Reset()

	OnPreGrow(func())
	OnPostGrow(func())
}

type pair struct {
	offset int64
	sz     int64
}
