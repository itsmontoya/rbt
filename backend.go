package rbTree

import (
	"unsafe"
)

// NewBytes will return a new bytes backend
func NewBytes(sz int64) Backend {
	var b Bytes
	b.bs = make([]byte, sz)
	b.tail = trunkSize
	b.cap = sz
	return &b
}

// Backend is a data backend
type Backend interface {
	getTrunk() (t *trunk)
	getBlock(offset int64) (b *Block)
	getKey(b *Block) (key []byte)
	getValue(b *Block) (value []byte)

	setBlob(b *Block, key, value []byte) (grew bool)

	newBlock(key []byte) (b *Block, offset int64, grew bool)
	newBlob(key, value []byte) (offset int64, grew bool)

	grow(sz int64) (grew bool)
}

// Bytes is a simple byteslice backend
type Bytes struct {
	bs   []byte
	tail int64
	cap  int64
}

func (bs *Bytes) setBlob(b *Block, key, value []byte) (grew bool) {
	valLen := int64(len(value))
	if valLen == b.valLen {
		blobIndex := b.offset + blockSize
		valueIndex := blobIndex + b.keyLen
		copy(bs.bs[valueIndex:], value)
		return
	}

	b.blobOffset, grew = bs.newBlob(key, value)
	b.valLen = valLen
	return
}

func (bs *Bytes) getTrunk() (t *trunk) {
	return (*trunk)(unsafe.Pointer(&bs.bs[0]))
}

func (bs *Bytes) getBlock(offset int64) (b *Block) {
	if offset == -1 {
		return
	}

	return (*Block)(unsafe.Pointer(&bs.bs[offset]))
}

func (bs *Bytes) getKey(b *Block) (key []byte) {
	blobIndex := b.offset + blockSize
	return bs.bs[blobIndex : blobIndex+b.keyLen]
}

func (bs *Bytes) getValue(b *Block) (value []byte) {
	blobIndex := b.offset + blockSize
	valueIndex := blobIndex + b.keyLen
	return bs.bs[valueIndex : valueIndex+b.valLen]
}

func (bs *Bytes) newBlock(key []byte) (b *Block, offset int64, grew bool) {
	offset = bs.tail
	grew = bs.grow(offset + blockSize)

	b = bs.getBlock(offset)
	bs.tail += blockSize

	// All new blocks start as red
	b.c = colorRed
	// Set offset and blob offset
	b.offset = offset
	b.blobOffset = -1
	// Set parent and children to their zero values
	b.parent = -1
	b.children[0] = -1
	b.children[1] = -1
	// Set key length
	b.keyLen = int64(len(key))
	// Debug
	b.derp = 67
	return
}

func (bs *Bytes) newBlob(key, value []byte) (offset int64, grew bool) {
	offset = bs.tail
	blobLen := int64(len(key) + len(value))
	grew = bs.grow(offset + blobLen)
	copy(bs.bs[offset:], key)
	copy(bs.bs[offset+int64(len(key)):], value)
	bs.tail += blobLen
	return
}

func (bs *Bytes) grow(sz int64) (grew bool) {
	for bs.cap < sz {
		bs.cap *= 2
		grew = true
	}

	if !grew {
		return
	}

	nbs := make([]byte, bs.cap)
	copy(nbs, bs.bs)
	bs.bs = nbs
	return
}
