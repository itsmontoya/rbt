package rbTree

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
// Count is the initial backend entries capacity
// keySize is the approximate key size (in bytes)
// valSize is the approximate value size (in bytes)
func New(cnt, keySize, valSize int64) *Tree {
	sz := keySize + valSize + blockSize
	sz *= cnt
	sz += trunkSize
	return newTree(NewBytes(sz))
}

/*
// NewMMAP will return a new MMAP
func NewMMAP(dir string, cnt, keySize, valSize int64) (t *Tree, err error) {
	sz := keySize + valSize + blockSize
	sz *= cnt
	sz += trunkSize

	var f *os.File
	if f, err = os.Create(path.Join(dir, "mmap.db")); err != nil {
		return
	}

	f.Truncate(sz)

	var mm mmap.MMap
	if mm, err = mmap.Map(f, os.O_RDWR, 0); err != nil {
		return
	}

	t = newTree(mm)
	return
}
*/

func newTree(be Backend) *Tree {
	var t Tree
	t.be = be
	t.setTrunk()
	t.t.root = -1
	return &t
}

// Tree is a red-black tree data structure
type Tree struct {
	be Backend
	t  *trunk
}

type trunk struct {
	root int64
	cnt  int64
}

func (t *Tree) setTrunk() {
	t.t = t.be.getTrunk()
}

// seekBlock will return a Block matching the provided key. It create is set to true, a new Block will be created if no match is found
func (t *Tree) seekBlock(startOffset int64, key []byte, create bool) (offset int64, grew bool) {
	offset = -1
	if startOffset == -1 {
		return
	}

	block := t.be.getBlock(startOffset)
	blockKey := t.be.getKey(block)

	switch bytes.Compare(key, blockKey) {
	case 1:
		child := block.children[1]
		if child == -1 {
			if !create {
				return
			}

			var nb *Block
			if nb, offset, grew = t.be.newBlock(key); grew {
				block = t.be.getBlock(startOffset)
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
			if nb, offset, grew = t.be.newBlock(key); grew {
				block = t.be.getBlock(startOffset)
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
	if grew = t.be.grow(sz); !grew {
		return
	}

	t.setTrunk()
	return
}

// getHead will get the very first item starting from a given node
// Note: If called from root, will return the first item in the tree
func (t *Tree) getHead(startOffset int64) (offset int64) {
	offset = -1

	if startOffset == -1 {
		return
	}

	b := t.be.getBlock(startOffset)
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

	b := t.be.getBlock(startOffset)
	if child := b.children[1]; child != -1 {
		return t.getTail(child)
	}

	return startOffset
}

func (t *Tree) getUncle(startOffset int64) (offset int64) {
	offset = -1
	block := t.be.getBlock(startOffset)
	parent := t.be.getBlock(block.parent)
	if parent == nil {
		return
	}

	grandparent := t.be.getBlock(parent.parent)
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

func (t *Tree) balance(b *Block) {
	parent := t.be.getBlock(b.parent)
	uncle := t.be.getBlock(t.getUncle(b.offset))

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

		grandparent := t.be.getBlock(parent.parent)
		grandparent.c = colorRed
		// Balance grandparent
		t.balance(grandparent)

	case parent.c == colorRed:
		grandparent := t.be.getBlock(parent.parent)

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
	parent := t.be.getBlock(b.parent)
	grandparent := t.be.getBlock(parent.parent)

	// Swap  children
	swapChild := t.be.getBlock(b.children[0])
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
	parent := t.be.getBlock(b.parent)
	grandparent := t.be.getBlock(parent.parent)

	// Swap  children
	swapChild := t.be.getBlock(b.children[1])
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
	parent := t.be.getBlock(b.parent)
	grandparent := t.be.getBlock(parent.parent)

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
		nb += t.numBlack(t.be.getBlock(childOffset))
	}

	if childOffset := b.children[1]; childOffset != -1 {
		nb += t.numBlack(t.be.getBlock(childOffset))
	}

	return
}

func (t *Tree) iterate(b *Block, fn ForEachFn) (ended bool) {
	if child := b.children[0]; child != -1 {
		if ended = t.iterate(t.be.getBlock(child), fn); ended {
			return
		}
	}

	if ended = fn(t.be.getKey(b), t.be.getValue(b)); ended {
		return
	}

	if child := b.children[1]; child != -1 {
		if ended = t.iterate(t.be.getBlock(child), fn); ended {
			return
		}

	}

	return
}

// Get will retrieve an item from a tree
func (t *Tree) Get(key []byte) (val []byte) {
	if offset, _ := t.seekBlock(t.t.root, key, false); offset != -1 {
		// Node was found, set value as the node's value
		val = t.be.getValue(t.be.getBlock(offset))
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
		b, offset, _ = t.be.newBlock(key)
		t.t.root = offset
	} else {
		// Find node whose key matches our provided key, if node does not exist - create it.
		offset, grew = t.seekBlock(t.t.root, key, true)
		b = t.be.getBlock(offset)
	}

	if grew = t.be.setBlob(b, key, val); grew {
		b = t.be.getBlock(offset)
	}

	// Balance tree after insert
	// TODO: This can be moved into the node-creation portion
	t.balance(b)

	root := t.be.getBlock(t.t.root)
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
	return t.iterate(t.be.getBlock(t.t.root), fn)
}

// Len will return the length of the data-store
func (t *Tree) Len() (n int) {
	return int(t.t.cnt)
}
