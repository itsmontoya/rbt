package whiskey

import (
	"bytes"
	"unsafe"

	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrCannotAllocate is returned when Whiskey cannot allocate the bytes it needs
	ErrCannotAllocate = errors.Error("cannot allocate needed bytes")
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
	// The only error that can return is ErrCannotAllocate which will not occur for a simple Bytes backend
	w, _ = NewRaw(sz, bs.grow, nil)
	return
}

// NewMMAP will return a new MMAP Whiskey
// sz is the size (in bytes) to initially allocate for this db
func NewMMAP(dir, name string, sz int64) (w *Whiskey, err error) {
	var mm *MMap
	if mm, err = newMMap(dir, name); err != nil {
		return
	}

	return NewRaw(sz, mm.grow, mm.Close)
}

// NewRaw will return a new Whiskey with the provided size, grow func, and close func
// sz is the size (in bytes) to initially allocate for this db
// gfn is the function to call on grows
// cfn is the function to call on close (optional)
func NewRaw(sz int64, gfn GrowFn, cfn CloseFn) (wp *Whiskey, err error) {
	var w Whiskey
	w.gfn = gfn
	w.cfn = cfn

	if sz < labelSize {
		sz = labelSize
	}

	if w.bs = w.gfn(sz); int64(len(w.bs)) < sz {
		err = ErrCannotAllocate
		return
	}

	w.setLabel()
	// Check if trunk has been initialized
	if w.l.tail == 0 {
		// trunk has not been set, set inital values
		w.l.root = -1
		w.l.tail = labelSize
		w.l.cap = sz
	}

	wp = &w
	return
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
	return w.bs[b.blobOffset : b.blobOffset+b.keyLen]
}

func (w *Whiskey) getValue(b *Block) (value []byte) {
	valueIndex := b.blobOffset + b.keyLen
	return w.bs[valueIndex : valueIndex+b.valLen]
}

func (w *Whiskey) setLabel() {
	w.l = (*label)(unsafe.Pointer(&w.bs[0]))
	w.l.cap = int64(len(w.bs))
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

	var offset, boffset int64
	offset = b.offset
	if boffset, grew = w.newBlob(key, value); grew {
		b = w.getBlock(offset)
	}

	b.blobOffset = boffset
	b.valLen = valLen
	return
}

func (w *Whiskey) growBlob(b *Block, key []byte, sz int64) (grew bool) {
	if sz <= b.valLen {
		return
	}

	vlen := b.valLen
	if vlen == 0 {
		vlen = sz
	}

	for vlen < sz {
		vlen *= 2
	}

	delta := vlen - b.valLen

	offset := b.offset
	boffset := w.l.tail
	blobLen := int64(len(key)) + vlen
	if grew = w.grow(boffset + blobLen); grew {
		b = w.getBlock(offset)
	}

	value := w.getValue(b)
	copy(w.bs[boffset:], key)
	copy(w.bs[boffset+b.keyLen:], value)
	w.l.tail += blobLen

	for i := w.l.tail - delta; i < w.l.tail; i++ {
		w.bs[i] = 0
	}

	b.blobOffset = boffset
	b.valLen = vlen
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
	b.valLen = 0
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

// detachFromParent will detach a block from it's parent reference
// Note: This is never called on root node, parent will always exist
func (w *Whiskey) detachFromParent(b *Block) {
	// Get the parent of block
	parent := w.getBlock(b.parent)
	// Set parent's child value for -1 where the block resided
	if parent.ct == childLeft {
		parent.children[0] = -1
	} else if parent.ct == childRight {
		parent.children[1] = -1
	}
}

// adoptChildren will set in-block children as out-block children
func (w *Whiskey) adoptChildren(in, out *Block) {
	// Set children of in-block to match the children of the out-block
	// Note: The in-block will always be a leaf. As a result, we know
	// that our next block does not have children.
	in.children[0] = out.children[0]
	in.children[1] = out.children[1]

	// Set children of out to nil values
	// Note: This technically isn't needed as the out block will be destoyed after this call
	// We could technically sqeeze out some more performance by avoiding two unnecessary write calls.
	// Consider removing this when going through the hyper-optimization phase
	out.children[0] = -1
	out.children[1] = -1
}

func (w *Whiskey) hasBlackChildpair(b *Block) bool {
	if b == nil {
		return false
	}

	var c *Block
	if c = w.getBlock(b.children[0]); c != nil && c.c == colorRed {
		return false
	}

	if c = w.getBlock(b.children[1]); c != nil && c.c == colorRed {
		return false
	}

	return true
}

func (w *Whiskey) zeroChildrenDelete(b, parent *Block) {
	if parent == nil {
		return
	}
	// Set reference to block in parent to -1
	if b.ct == childLeft {
		parent.children[0] = -1
	} else {
		parent.children[1] = -1
	}
}

func (w *Whiskey) oneChildDelete(b, parent *Block) (next *Block) {
	if b.children[0] != -1 {
		next = w.getBlock(b.children[0])
	} else {
		next = w.getBlock(b.children[1])
	}

	w.replace(b, next, parent)
	return
}

func (w *Whiskey) twoChildDelete(b, parent *Block) (next *Block) {
	// Get the very next element following block
	// Note: Selecting the second child will ensure we move forward.
	// Calling getHead from this location will land us at the item directly
	// following the target block.
	next = w.getBlock(w.getHead(b.children[1]))
	w.replace(b, next, parent)
	return
}

func (w *Whiskey) replace(old, new, parent *Block) {
	w.adoptChildren(new, old)
	w.detachFromParent(new)
	// Set next-block childtype as the block childtype
	new.ct = old.ct
	// Set the next-block parent as the block parent
	new.parent = old.parent

	// Set the parent's child value as the offset to our next block
	if old.ct == childLeft {
		parent.children[0] = new.offset
	} else if old.ct == childRight {
		parent.children[1] = new.offset
	} else {
		// If block is root, we need to update the label's reference to root
		w.l.root = new.offset
	}

}

// Delete will remove an item from the tree
func (w *Whiskey) Delete(key []byte) {
	var (
		b      *Block
		next   *Block
		offset int64
	)

	if offset, _ = w.seekBlock(w.l.root, key, false); offset == -1 {
		return
	}

	b = w.getBlock(offset)
	parent := w.getBlock(b.parent)
	hasLeft := b.children[0] != -1
	hasRight := b.children[1] != -1
	// BST Delete switch
	switch {
	case hasLeft && hasRight:
		next = w.twoChildDelete(b, parent)

	case !hasLeft && !hasRight:
		w.zeroChildrenDelete(b, parent)

	default:
		// Technically this is out of order, but it seems much more clean to check to see
		// if we have ALL or NONE. If neither cases exist, we know we have one child
		next = w.oneChildDelete(b, parent)

	}

	// Balancing cases
	if b.c == colorRed || (next != nil && next.c == colorRed) {
		// Simple Case: If either u or v is red
		// Note: Because we are not disrupting the black-level, no rotation is needed
		next.c = colorBlack
		return
	}

	var sibling *Block
	// Acquire sibling
	switch {
	case b.ct == childRoot:

	case b.ct == childLeft:
		sibling = w.getBlock(parent.children[1])
	case b.ct == childRight:
		sibling = w.getBlock(parent.children[0])
	}

	// Set sibling black state
	// Note: We can eventually remove this for performance reasons once this function
	// is completely fleshed out. We need to ensure that we are not dealing with a nil sibling
	// This is just to avoid running into panic land
	sib := isBlack(sibling)

	// Sibling rotate bonanza
	switch {
	case sib && w.hasBlackChildpair(sibling):
		// Sibling is black and has both black children
	case sib:
		// Sibling is black and has at least one red child
	default:
		// Sibling is red
		//	case sib &&:
	}

	// Complex Case: If Both u and v are Black.
	// Here we go..
	// TODO: Balance
	return
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

// Grow will grow a blob value to a given size
func (w *Whiskey) Grow(key []byte, sz int64) (bs []byte) {
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

	if grew = w.growBlob(b, key, sz); grew {
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

	bs = w.getValue(b)
	return
}

// Reset will clear the tree and keep the backend. Can be used as a fresh store
func (w *Whiskey) Reset() {
	w.l.tail = labelSize
	w.l.root = -1
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

func isBlack(b *Block) bool {
	if b == nil {
		return true
	}

	return b.c == colorBlack
}

func isRed(b *Block) bool {
	if b == nil {
		return false
	}

	return b.c == colorRed
}
