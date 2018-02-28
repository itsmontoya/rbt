package backend

// Backend is the backend interface
type Backend interface {
	Grow(sz int64) []byte
	Close() error
}
