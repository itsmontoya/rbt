package allocator

import "unsafe"

var (
	headerSize = getHeaderSize()
)

type header struct {
	tail int64
	cap  int64
	hmm  int64
}

func getHeaderSize() int64 {
	ptr := unsafe.Sizeof(header{})
	sz := *(*int64)(unsafe.Pointer(&ptr))
	return sz
}
