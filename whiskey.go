package whiskey

import (
	"bytes"
	"unsafe"
)

const (
	colorBlack color = iota
	colorRed
)

const (
	childRoot childType = iota
	childLeft
	childRight
)

// Note: These have to be variables in order to be referenced
var (
	trunkSizePtr = unsafe.Sizeof(trunk{})
	blockSizePtr = unsafe.Sizeof(Block{})

	trunkSize = *(*int64)(unsafe.Pointer(&trunkSizePtr))
	blockSize = *(*int64)(unsafe.Pointer(&blockSizePtr))
)

// New will return a new Tree
// sz is the size (in bytes) to initially allocate for this db
func New(sz int64) (t *Tree) {
	bs := newBytes()
	t = newTree(sz, bs.grow, nil)
	return
}

// NewMMAP will return a new MMAP Tree
// sz is the size (in bytes) to initially allocate for this db
func NewMMAP(dir, name string, sz int64) (t *Tree, err error) {
	var mm *MMap
	if mm, err = newMMap(dir, name); err != nil {
		return
	}

	t = newTree(sz, mm.grow, mm.Close)
	return
}

// newTree will return a new Tree with the provided size, grow func, and close func
// sz is the size (in bytes) to initially allocate for this db
// gfn is the function to call on grows
// cfn is the function to call on close (optional)
func newTree(sz int64, gfn GrowFn, cfn CloseFn) *Tree {
	var t Tree
	t.gfn = gfn
	t.cfn = cfn
	t.bs = t.gfn(sz)
	t.setTrunk()
	// Check if trunk has been initialized
	if t.t.tail == 0 {
		// trunk has not been set, set inital values
		t.t.root = -1
		t.t.tail = trunkSize
		t.t.cap = sz
	}

	return &t
}

// Tree is a red-black tree data structure
type Tree struct {
	bs []byte
	t  *trunk

	gfn GrowFn
	cfn CloseFn
}

// getHead will get the very first item starting from a given node
// Note: If called from root, will return the first item in the tree
func (t *Tree) getHead(startOffset int64) (offset int64) {
	offset = -1

	if startOffset == -1 {
		return
	}

	b := t.getBlock(startOffset)
	if child := b.children[0]; child != -1 {
		return t.getHead(child)
	}

	return startOffset
}

// getTail will get the very last item starting from a given node
// Note: If called from root, will return the last item in the tree
func (t *Tree) getTail(startOffset int64) (offset int64) {
	offset = -1

	if startOffset == -1 {
		return
	}

	b := t.getBlock(startOffset)
	if child := b.children[1]; child != -1 {
		return t.getTail(child)
	}

	return startOffset
}

func (t *Tree) getUncle(startOffset int64) (offset int64) {
	offset = -1
	block := t.getBlock(startOffset)
	parent := t.getBlock(block.parent)
	if parent == nil {
		return
	}

	grandparent := t.getBlock(parent.parent)
	if grandparent == nil {
		return
	}

	switch parent.ct {
	case childLeft:
		return grandparent.children[1]
	case childRight:
		return grandparent.children[0]

	}

	return
}

func (t *Tree) getBlock(offset int64) (b *Block) {
	if offset == -1 {
		return
	}

	return (*Block)(unsafe.Pointer(&t.bs[offset]))
}

func (t *Tree) getKey(b *Block) (key []byte) {
	blobIndex := b.offset + blockSize
	return t.bs[blobIndex : blobIndex+b.keyLen]
}

func (t *Tree) getValue(b *Block) (value []byte) {
	blobIndex := b.offset + blockSize
	valueIndex := blobIndex + b.keyLen
	return t.bs[valueIndex : valueIndex+b.valLen]
}

func (t *Tree) setTrunk() {
	t.t = (*trunk)(unsafe.Pointer(&t.bs[0]))
	t.t.cap = int64(cap(t.bs))
}

func (t *Tree) setParentChild(b, parent, child *Block) {
	switch b.ct {
	case childLeft:
		parent.children[0] = child.offset
	case childRight:
		parent.children[1] = child.offset
	case childRoot:
		// No action is taken, tree will handle this at the end of put
	}
}

func (t *Tree) setBlob(b *Block, key, value []byte) (grew bool) {
	valLen := int64(len(value))
	if valLen == b.valLen {
		blobIndex := b.offset + blockSize
		valueIndex := blobIndex + b.keyLen
		copy(t.bs[valueIndex:], value)
		return
	}

	b.blobOffset, grew = t.newBlob(key, value)
	b.valLen = valLen
	return
}

func (t *Tree) newBlock(key []byte) (b *Block, offset int64, grew bool) {
	offset = t.t.tail
	grew = t.grow(offset + blockSize)

	b = t.getBlock(offset)
	t.t.tail += blockSize

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

func (t *Tree) newBlob(key, value []byte) (offset int64, grew bool) {
	offset = t.t.tail
	blobLen := int64(len(key) + len(value))
	grew = t.grow(offset + blobLen)
	copy(t.bs[offset:], key)
	copy(t.bs[offset+int64(len(key)):], value)
	t.t.tail += blobLen
	return
}

// seekBlock will return a Block matching the provided key. It create is set to true, a new Block will be created if no match is found
func (t *Tree) seekBlock(startOffset int64, key []byte, create bool) (offset int64, grew bool) {
	offset = -1
	if startOffset == -1 {
		return
	}

	block := t.getBlock(startOffset)
	blockKey := t.getKey(block)

	switch bytes.Compare(key, blockKey) {
	case 1:
		child := block.children[1]
		if child == -1 {
			if !create {
				return
			}

			var nb *Block
			if nb, offset, grew = t.newBlock(key); grew {
				block = t.getBlock(startOffset)
			}

			nb.ct = childRight
			nb.parent = startOffset
			block.children[1] = offset
			return
		}

		return t.seekBlock(child, key, create)

	case -1:
		child := block.children[0]
		if child == -1 {
			if !create {
				return
			}

			var nb *Block
			if nb, offset, grew = t.newBlock(key); grew {
				block = t.getBlock(startOffset)
			}

			nb.ct = childLeft
			nb.parent = startOffset
			block.children[0] = offset
			return
		}

		return t.seekBlock(child, key, create)

	case 0:
		offset = startOffset
		return
	}

	return
}

func (t *Tree) grow(sz int64) (grew bool) {
	if t.t.cap > sz {
		return
	}

	t.bs = t.gfn(sz)
	t.setTrunk()
	return true
}

func (t *Tree) balance(b *Block) {
	parent := t.getBlock(b.parent)
	uncle := t.getBlock(t.getUncle(b.offset))

	switch {
	case b.c == colorBlack:
		return
	case b.ct == childRoot:
		if b.c == colorRed {
			b.c = colorBlack
			return
		}

	case uncle != nil && uncle.c == colorRed:
		parent.c = colorBlack
		uncle.c = colorBlack

		grandparent := t.getBlock(parent.parent)
		grandparent.c = colorRed
		// Balance grandparent
		t.balance(grandparent)

	case parent.c == colorRed:
		grandparent := t.getBlock(parent.parent)

		if t.isTriangle(b, parent) {
			t.rotateParent(b)
			// Balance parent
			t.balance(parent)
		} else {
			// Is a line
			t.rotateGrandparent(b)
			// Balance grandparent
			t.balance(grandparent)
		}
	}
}

func (t *Tree) leftRotate(b *Block) {
	parent := t.getBlock(b.parent)
	grandparent := t.getBlock(parent.parent)

	// Swap  children
	swapChild := t.getBlock(b.children[0])
	b.children[0] = parent.offset

	if swapChild != nil {
		parent.children[1] = swapChild.offset
		swapChild.parent = parent.offset
		swapChild.ct = childRight
		// Set nidx as the child for our grandparent
	} else {
		parent.children[1] = -1
	}

	// Set block as the child for our grandparent
	t.setParentChild(parent, grandparent, b)

	// Set n's grandparent as parent
	if grandparent == nil {
		b.parent = -1
	} else {
		b.parent = grandparent.offset
	}
	// Set n as parent to the original parent
	parent.parent = b.offset

	// Set child types
	b.ct = parent.ct
	parent.ct = childLeft
}

func (t *Tree) rightRotate(b *Block) {
	parent := t.getBlock(b.parent)
	grandparent := t.getBlock(parent.parent)

	// Swap  children
	swapChild := t.getBlock(b.children[1])
	b.children[1] = parent.offset

	if swapChild != nil {
		parent.children[0] = swapChild.offset
		swapChild.parent = parent.offset
		swapChild.ct = childLeft
	} else {
		parent.children[0] = -1
		// Set nidx as the child for our grandparent
	}

	// Set block as the child for our grandparent
	t.setParentChild(parent, grandparent, b)

	// Set n's grandparent as parent
	if grandparent == nil {
		b.parent = -1
	} else {
		b.parent = grandparent.offset
	}
	// Set n as parent to the original parent
	parent.parent = b.offset

	// Set child types
	b.ct = parent.ct
	parent.ct = childRight
}

func (t *Tree) rotateParent(b *Block) {
	switch b.ct {
	case childLeft:
		t.rightRotate(b)

	case childRight:
		t.leftRotate(b)

	default:
		panic("invalid child type for parent rotation")
	}
}

func (t *Tree) rotateGrandparent(b *Block) {
	parent := t.getBlock(b.parent)
	grandparent := t.getBlock(parent.parent)

	switch parent.ct {
	case childLeft:
		t.rightRotate(parent)

	case childRight:
		t.leftRotate(parent)

	default:
		panic("invalid child type for grandparent rotation")
	}

	parent.c = colorBlack
	grandparent.c = colorRed
}

func (t *Tree) isTriangle(b, parent *Block) (isTriangle bool) {
	if b == nil {
		return
	}

	if b.ct == childLeft && parent.ct == childRight {
		return true
	}

	if b.ct == childRight && parent.ct == childLeft {
		return true
	}

	return
}

func (t *Tree) numBlack(b *Block) (nb int) {
	if b == nil {
		return
	}

	if b.c == colorBlack {
		nb = 1
	}

	if childOffset := b.children[0]; childOffset != -1 {
		nb += t.numBlack(t.getBlock(childOffset))
	}

	if childOffset := b.children[1]; childOffset != -1 {
		nb += t.numBlack(t.getBlock(childOffset))
	}

	return
}

func (t *Tree) iterate(b *Block, fn ForEachFn) (ended bool) {
	if child := b.children[0]; child != -1 {
		if ended = t.iterate(t.getBlock(child), fn); ended {
			return
		}
	}

	if ended = fn(t.getKey(b), t.getValue(b)); ended {
		return
	}

	if child := b.children[1]; child != -1 {
		if ended = t.iterate(t.getBlock(child), fn); ended {
			return
		}

	}

	return
}

// Get will retrieve an item from a tree
func (t *Tree) Get(key []byte) (val []byte) {
	if offset, _ := t.seekBlock(t.t.root, key, false); offset != -1 {
		// Node was found, set value as the node's value
		val = t.getValue(t.getBlock(offset))
	}

	return
}

// Put will insert an item into the tree
func (t *Tree) Put(key, val []byte) {
	var (
		b      *Block
		grew   bool
		offset int64
	)

	if t.t.root == -1 {
		// Root doesn't exist, we can create one
		b, offset, _ = t.newBlock(key)
		t.t.root = offset
	} else {
		// Find node whose key matches our provided key, if node does not exist - create it.
		offset, grew = t.seekBlock(t.t.root, key, true)
		b = t.getBlock(offset)
	}

	if grew = t.setBlob(b, key, val); grew {
		b = t.getBlock(offset)
	}

	// Balance tree after insert
	// TODO: This can be moved into the node-creation portion
	t.balance(b)

	root := t.getBlock(t.t.root)
	if root.ct != childRoot {
		// Root has changed, update root reference to the new root
		t.t.root = root.parent
	}

	t.t.cnt++
}

// ForEach will iterate through each tree item
func (t *Tree) ForEach(fn ForEachFn) (ended bool) {
	if t.t.root == -1 {
		// Root doesn't exist, return early
		return
	}

	// Call iterate from root
	return t.iterate(t.getBlock(t.t.root), fn)
}

// Len will return the length of the data-store
func (t *Tree) Len() (n int) {
	return int(t.t.cnt)
}

// Close will close a tree
func (t *Tree) Close() (err error) {
	if t.cfn == nil {
		return
	}

	return t.cfn()
}
