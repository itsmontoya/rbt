package rbTree

func newBlock(key []byte) (b Block) {
	// All new blocks start as red
	n.c = colorRed
	// Set parent and children to their zero values
	n.parent = -1
	n.children[0] = -1
	n.children[1] = -1
	// Set key length
	b.keyLen = len(key)
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

	parent   int
	children [2]int

	keyLen int
	valLen int
}

// Blob represents a Key/Value entry
// Blob is solid.
type Blob struct {
	Key []byte
	Val []byte
}
