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
	labelSizePtr = unsafe.Sizeof(label{})
	blockSizePtr = unsafe.Sizeof(Block{})

	labelSize = *(*int64)(unsafe.Pointer(&labelSizePtr))
	blockSize = *(*int64)(unsafe.Pointer(&blockSizePtr))
)

// New will return a new Whiskey
// sz is the size (in bytes) to initially allocate for this db
func New(sz int64) (w *Whiskey) {
	bs := newBytes()
	w = NewRaw(sz, bs.grow, nil)
	return
}

// NewMMAP will return a new MMAP Whiskey
// sz is the size (in bytes) to initially allocate for this db
func NewMMAP(dir, name string, sz int64) (w *Whiskey, err error) {
	var mm *MMap
	if mm, err = newMMap(dir, name); err != nil {
		return
	}

	w = NewRaw(sz, mm.grow, mm.Close)
	return
}

// NewRaw will return a new Whiskey with the provided size, grow func, and close func
// sz is the size (in bytes) to initially allocate for this db
// gfn is the function to call on grows
// cfn is the function to call on close (optional)
func NewRaw(sz int64, gfn GrowFn, cfn CloseFn) *Whiskey {
	var w Whiskey
	w.gfn = gfn
	w.cfn = cfn
	w.bs = w.gfn(sz)
	w.setLabel()
	// Check if trunk has been initialized
	if w.l.tail == 0 {
		// trunk has not been set, set inital values
		w.l.root = -1
		w.l.tail = labelSize
		w.l.cap = sz
	}

	return &w
}

// Whiskey is a red-black tree data structure
type Whiskey struct {
	bs []byte
	l  *label

	gfn GrowFn
	cfn CloseFn
}

// getHead will get the very first item starting from a given node
// Note: If called from root, will return the first item in the tree
func (w *Whiskey) getHead(startOffset int64) (offset int64) {
	offset = -1

	if startOffset == -1 {
		return
	}

	b := w.getBlock(startOffset)
	if child := b.children[0]; child != -1 {
		return w.getHead(child)
	}

	return startOffset
}

// getTail will get the very last item starting from a given node
// Note: If called from root, will return the last item in the tree
func (w *Whiskey) getTail(startOffset int64) (offset int64) {
	offset = -1

	if startOffset == -1 {
		return
	}

	b := w.getBlock(startOffset)
	if child := b.children[1]; child != -1 {
		return w.getTail(child)
	}

	return startOffset
}

func (w *Whiskey) getUncle(startOffset int64) (offset int64) {
	offset = -1
	block := w.getBlock(startOffset)
	parent := w.getBlock(block.parent)
	if parent == nil {
		return
	}

	grandparent := w.getBlock(parent.parent)
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

func (w *Whiskey) getBlock(offset int64) (b *Block) {
	if offset == -1 {
		return
	}

	return (*Block)(unsafe.Pointer(&w.bs[offset]))
}

func (w *Whiskey) getKey(b *Block) (key []byte) {
	blobIndex := b.offset + blockSize
	return w.bs[blobIndex : blobIndex+b.keyLen]
}

func (w *Whiskey) getValue(b *Block) (value []byte) {
	blobIndex := b.offset + blockSize
	valueIndex := blobIndex + b.keyLen
	return w.bs[valueIndex : valueIndex+b.valLen]
}

func (w *Whiskey) setLabel() {
	w.l = (*label)(unsafe.Pointer(&w.bs[0]))
	w.l.cap = int64(cap(w.bs))
}

func (w *Whiskey) setParentChild(b, parent, child *Block) {
	switch b.ct {
	case childLeft:
		parent.children[0] = child.offset
	case childRight:
		parent.children[1] = child.offset
	case childRoot:
		// No action is taken, tree will handle this at the end of put
	}
}

func (w *Whiskey) setBlob(b *Block, key, value []byte) (grew bool) {
	valLen := int64(len(value))
	if valLen == b.valLen {
		blobIndex := b.offset + blockSize
		valueIndex := blobIndex + b.keyLen
		copy(w.bs[valueIndex:], value)
		return
	}

	b.blobOffset, grew = w.newBlob(key, value)
	b.valLen = valLen
	return
}

func (w *Whiskey) newBlock(key []byte) (b *Block, offset int64, grew bool) {
	offset = w.l.tail
	grew = w.grow(offset + blockSize)

	b = w.getBlock(offset)
	w.l.tail += blockSize

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

func (w *Whiskey) newBlob(key, value []byte) (offset int64, grew bool) {
	offset = w.l.tail
	blobLen := int64(len(key) + len(value))
	grew = w.grow(offset + blobLen)
	copy(w.bs[offset:], key)
	copy(w.bs[offset+int64(len(key)):], value)
	w.l.tail += blobLen
	return
}

// seekBlock will return a Block matching the provided key. It create is set to true, a new Block will be created if no match is found
func (w *Whiskey) seekBlock(startOffset int64, key []byte, create bool) (offset int64, grew bool) {
	offset = -1
	if startOffset == -1 {
		return
	}

	block := w.getBlock(startOffset)
	blockKey := w.getKey(block)

	switch bytes.Compare(key, blockKey) {
	case 1:
		child := block.children[1]
		if child == -1 {
			if !create {
				return
			}

			var nb *Block
			if nb, offset, grew = w.newBlock(key); grew {
				block = w.getBlock(startOffset)
			}

			nb.ct = childRight
			nb.parent = startOffset
			block.children[1] = offset
			return
		}

		return w.seekBlock(child, key, create)

	case -1:
		child := block.children[0]
		if child == -1 {
			if !create {
				return
			}

			var nb *Block
			if nb, offset, grew = w.newBlock(key); grew {
				block = w.getBlock(startOffset)
			}

			nb.ct = childLeft
			nb.parent = startOffset
			block.children[0] = offset
			return
		}

		return w.seekBlock(child, key, create)

	case 0:
		offset = startOffset
		return
	}

	return
}

func (w *Whiskey) grow(sz int64) (grew bool) {
	if w.l.cap > sz {
		return
	}

	w.bs = w.gfn(sz)
	w.setLabel()
	return true
}

func (w *Whiskey) balance(b *Block) {
	parent := w.getBlock(b.parent)
	uncle := w.getBlock(w.getUncle(b.offset))

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

		grandparent := w.getBlock(parent.parent)
		grandparent.c = colorRed
		// Balance grandparent
		w.balance(grandparent)

	case parent.c == colorRed:
		grandparent := w.getBlock(parent.parent)

		if w.isTriangle(b, parent) {
			w.rotateParent(b)
			// Balance parent
			w.balance(parent)
		} else {
			// Is a line
			w.rotateGrandparent(b)
			// Balance grandparent
			w.balance(grandparent)
		}
	}
}

func (w *Whiskey) leftRotate(b *Block) {
	parent := w.getBlock(b.parent)
	grandparent := w.getBlock(parent.parent)

	// Swap  children
	swapChild := w.getBlock(b.children[0])
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
	w.setParentChild(parent, grandparent, b)

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

func (w *Whiskey) rightRotate(b *Block) {
	parent := w.getBlock(b.parent)
	grandparent := w.getBlock(parent.parent)

	// Swap  children
	swapChild := w.getBlock(b.children[1])
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
	w.setParentChild(parent, grandparent, b)

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

func (w *Whiskey) rotateParent(b *Block) {
	switch b.ct {
	case childLeft:
		w.rightRotate(b)

	case childRight:
		w.leftRotate(b)

	default:
		panic("invalid child type for parent rotation")
	}
}

func (w *Whiskey) rotateGrandparent(b *Block) {
	parent := w.getBlock(b.parent)
	grandparent := w.getBlock(parent.parent)

	switch parent.ct {
	case childLeft:
		w.rightRotate(parent)

	case childRight:
		w.leftRotate(parent)

	default:
		panic("invalid child type for grandparent rotation")
	}

	parent.c = colorBlack
	grandparent.c = colorRed
}

func (w *Whiskey) isTriangle(b, parent *Block) (isTriangle bool) {
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

func (w *Whiskey) numBlack(b *Block) (nb int) {
	if b == nil {
		return
	}

	if b.c == colorBlack {
		nb = 1
	}

	if childOffset := b.children[0]; childOffset != -1 {
		nb += w.numBlack(w.getBlock(childOffset))
	}

	if childOffset := b.children[1]; childOffset != -1 {
		nb += w.numBlack(w.getBlock(childOffset))
	}

	return
}

func (w *Whiskey) iterate(b *Block, fn ForEachFn) (ended bool) {
	if child := b.children[0]; child != -1 {
		if ended = w.iterate(w.getBlock(child), fn); ended {
			return
		}
	}

	if ended = fn(w.getKey(b), w.getValue(b)); ended {
		return
	}

	if child := b.children[1]; child != -1 {
		if ended = w.iterate(w.getBlock(child), fn); ended {
			return
		}

	}

	return
}

// Get will retrieve an item from a tree
func (w *Whiskey) Get(key []byte) (val []byte) {
	if offset, _ := w.seekBlock(w.l.root, key, false); offset != -1 {
		// Node was found, set value as the node's value
		val = w.getValue(w.getBlock(offset))
	}

	return
}

// Put will insert an item into the tree
func (w *Whiskey) Put(key, val []byte) {
	var (
		b      *Block
		grew   bool
		offset int64
	)

	if w.l.root == -1 {
		// Root doesn't exist, we can create one
		b, offset, _ = w.newBlock(key)
		w.l.root = offset
	} else {
		// Find node whose key matches our provided key, if node does not exist - create it.
		offset, grew = w.seekBlock(w.l.root, key, true)
		b = w.getBlock(offset)
	}

	if grew = w.setBlob(b, key, val); grew {
		b = w.getBlock(offset)
	}

	// Balance tree after insert
	// TODO: This can be moved into the node-creation portion
	w.balance(b)

	root := w.getBlock(w.l.root)
	if root.ct != childRoot {
		// Root has changed, update root reference to the new root
		w.l.root = root.parent
	}

	w.l.cnt++
}

// ForEach will iterate through each tree item
func (w *Whiskey) ForEach(fn ForEachFn) (ended bool) {
	if w.l.root == -1 {
		// Root doesn't exist, return early
		return
	}

	// Call iterate from root
	return w.iterate(w.getBlock(w.l.root), fn)
}

// Len will return the length of the data-store
func (w *Whiskey) Len() (n int) {
	return int(w.l.cnt)
}

// Close will close a tree
func (w *Whiskey) Close() (err error) {
	if w.cfn == nil {
		return
	}

	return w.cfn()
}
