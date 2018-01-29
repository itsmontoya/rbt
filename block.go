package whiskey

func newBlock(key []byte) (b Block) {
	// All new blocks start as red
	b.c = colorRed
	// Set parent and children to their zero values
	b.parent = -1
	b.children[0] = -1
	b.children[1] = -1
	// Set key length
	b.keyLen = int64(len(key))
	return
}

// BlockAndBlob are friends
type BlockAndBlob struct {
	Block
	Blob
}

// Block is a reference to a data block
type Block struct {
	c  color
	ct childType

	offset     int64
	blobOffset int64
	parent     int64
	children   [2]int64

	keyLen int64
	valLen int64

	derp byte
}

// Blob represents a Key/Value entry
// Blob is solid.
type Blob struct {
	Key []byte
	Val []byte
}
