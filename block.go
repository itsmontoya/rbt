package rbt

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
