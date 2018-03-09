package backend

// Backend is a backend interface
type Backend interface {
	Get(offset, sz int64) []byte
	Allocate(sz int64) (_ *Section, grew bool)
	Release(*Section)
	Close() (err error)
}

// Section represents an allocated section of bytes
type Section struct {
	Offset int64
	Size   int64
	Bytes  []byte
}
