package rbt

import (
	"bytes"
	"unsafe"

	"github.com/itsmontoya/rbt/allocator"
	"github.com/itsmontoya/rbt/backend"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrCannotAllocate is returned when Tree cannot allocate the bytes it needs
	ErrCannotAllocate = errors.Error("cannot allocate needed bytes")
)

const (
	colorBlack color = iota
	colorRed
	colorDoubleBlack
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
	bs := backend.NewBytes()
	// The only error that can return is ErrCannotAllocate which will not occur for a simple Bytes backend
	t, _ = NewRaw(sz, bs)
	return
}

// NewMMAP will return a new MMAP Tree
// sz is the size (in bytes) to initially allocate for this db
func NewMMAP(dir, name string, sz int64) (t *Tree, err error) {
	var mm *backend.MMap
	if mm, err = backend.NewMMap(dir, name); err != nil {
		return
	}

	return NewRaw(sz, mm)
}

// NewRaw will return a new Tree with the provided size, grow func, and close func
// sz is the size (in bytes) to initially allocate for this db
// gfn is the function to call on grows
// cfn is the function to call on close (optional)
func NewRaw(sz int64, b backend.Backend) (tp *Tree, err error) {
	var t Tree
	t.b = b
	t.a = allocator.NewSimple(b, sz)
	if sz < trunkSize {
		sz = trunkSize
	}

	t.a.Grow(sz)
	if t.a.Len() == 0 {
		t.a.Allocate(trunkSize)
	}

	t.setTrunk()
	// Check if trunk has been initialized
	if t.t.root == 0 {
		// trunk has not been set, set inital values
		t.t.root = -1
		t.t.cnt = 0
	}

	t.a.OnPostGrow(t.setTrunk)
	tp = &t
	return
}

// Tree is a red-black tree data structure
type Tree struct {
	t *trunk
	b backend.Backend
	a allocator.Allocator
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

	return (*Block)(unsafe.Pointer(t.a.Byte(offset)))
}

func (t *Tree) getKey(b *Block) (key []byte) {
	return t.a.Bytes(b.blobOffset, b.keyLen)
}

func (t *Tree) getValue(b *Block) (value []byte) {
	return t.a.Bytes(b.blobOffset+b.keyLen, b.valLen)
}

func (t *Tree) setTrunk() {
	t.t = (*trunk)(unsafe.Pointer(t.a.First()))
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
		copy(t.a.Bytes(valueIndex, b.valLen), value)
		return
	}

	var offset, boffset int64
	offset = b.offset
	if boffset, grew = t.newBlob(key, value); grew {
		b = t.getBlock(offset)
	}

	b.blobOffset = boffset
	b.valLen = valLen
	return
}

func (t *Tree) growBlob(b *Block, key []byte, sz int64) (grew bool) {
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
	blobLen := int64(len(key)) + vlen
	offset := b.offset

	var boffset int64
	if boffset, grew = t.a.Allocate(blobLen); grew {
		b = t.getBlock(offset)
	}

	value := t.getValue(b)
	bs := t.a.Bytes(boffset, blobLen)
	copy(bs, key)
	copy(bs[b.keyLen:], value)

	for i := len(bs) - int(delta); i < len(bs); i++ {
		bs[i] = 0
	}

	b.blobOffset = boffset
	b.valLen = vlen
	return
}

func (t *Tree) newBlock(key []byte) (b *Block, offset int64, grew bool) {
	offset, grew = t.a.Allocate(blockSize)
	b = t.getBlock(offset)

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

func (t *Tree) newBlob(key, value []byte) (offset int64, grew bool) {
	blobLen := int64(len(key) + len(value))
	offset, grew = t.a.Allocate(blobLen)
	bs := t.a.Bytes(offset, blobLen)
	copy(bs, key)
	copy(bs[int64(len(key)):], value)
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

// detachFromParent will detach a block from it's parent reference
// Note: This is never called on root node, parent will always exist
func (t *Tree) detachFromParent(b *Block) {
	// Get the parent of block
	parent := t.getBlock(b.parent)
	if parent == nil {
		return
	}

	// Set parent's child value for -1 where the block resided
	if b.ct == childLeft {
		parent.children[0] = -1
	} else if b.ct == childRight {
		parent.children[1] = -1
	}
}

// adoptChildren will set in-block children as out-block children
func (t *Tree) adoptChildren(in, out *Block) {
	var child *Block
	//	parent := t.getBlock(in.parent)
	// Set children of in-block to match the children of the out-block
	// Note: The in-block will always be a leaf. As a result, we know
	// that our next block does not have children.
	if out.children[0] != in.offset {
		if child = t.getBlock(in.children[0]); child != nil {
			child.parent = in.parent
			in.parent = -1
		}

		in.children[0] = out.children[0]
		if child = t.getBlock(in.children[0]); child != nil {
			child.parent = in.offset
		}
	}

	if out.children[1] != in.offset {
		if child = t.getBlock(in.children[1]); child != nil {
			child.parent = in.parent
			in.parent = -1
		}

		in.children[1] = out.children[1]
		if child = t.getBlock(in.children[1]); child != nil {
			child.parent = in.offset
		}
	}

	// Set children of out to nil values
	// Note: This technically isn't needed as the out block will be destoyed after this call
	// We could technically sqeeze out some more performance by avoiding two unnecessary write calls.
	// Consider removing this when going through the hyper-optimization phase
	out.children[0] = -1
	out.children[1] = -1
}

func (t *Tree) hasBlackChildpair(b *Block) bool {
	if b == nil {
		return false
	}

	var c *Block
	if c = t.getBlock(b.children[0]); c != nil && c.c == colorRed {
		return false
	}

	if c = t.getBlock(b.children[1]); c != nil && c.c == colorRed {
		return false
	}

	return true
}

func (t *Tree) hasRedChildpair(b *Block) bool {
	if b == nil {
		return false
	}

	var c *Block
	if c = t.getBlock(b.children[0]); c != nil && c.c == colorBlack {
		return false
	}

	if c = t.getBlock(b.children[1]); c != nil && c.c == colorBlack {
		return false
	}

	return true
}

func (t *Tree) zeroChildrenDelete(b, parent *Block) {
	t.replace(b, nil, parent)
}

func (t *Tree) oneChildDelete(b, parent *Block) (next *Block) {
	if b.children[1] != -1 {
		next = t.getBlock(b.children[1])
	} else {
		next = t.getBlock(b.children[0])
	}

	t.replace(b, next, parent)
	return
}

func (t *Tree) twoChildDelete(b, parent *Block) (next *Block) {
	var child *Block
	// Get the very next element following block
	// Note: Selecting the second child will ensure we move forward.
	// Calling getHead from this location will land us at the item directly
	// following the target block.
	next = t.getBlock(t.getHead(b.children[1]))

	if next.offset != b.children[1] {
		// Our next item is not the direct child of our target block, let's have our next block's parent adopt the orphan
		// Note: If the child doesn't exist, the offset of -1 will be applied as the orphan reference
		var coffset int64 = -1
		if child = t.getBlock(next.children[1]); child != nil {
			// Child exists, set child offset
			coffset = child.offset
			child.parent = next.parent
		}

		nextParent := t.getBlock(next.parent)
		// Have next parent adopt the orphan
		nextParent.children[0] = coffset
		// Detach next from it's parent
		next.parent = -1

		next.children[1] = b.children[1]
		if child = t.getBlock(next.children[1]); child != nil {
			child.parent = next.offset
		}
	}

	next.children[0] = b.children[0]
	if child = t.getBlock(next.children[0]); child != nil {
		child.parent = next.offset
	}

	t.replace(b, next, parent)
	return
}

func (t *Tree) replace(old, new, parent *Block) {
	var noffset int64 = -1
	if new != nil {
		t.detachFromParent(new)
		// Set next-block childtype as the block childtype
		new.ct = old.ct
		// Set the next-block parent as the block parent
		new.parent = old.parent
		noffset = new.offset
	}

	// Set the parent's child value as the offset to our next block
	switch old.ct {
	case childRoot:
		// If block is root, we need to update the trunk's reference to root
		t.t.root = noffset
		t.t.cnt = 0
	case childLeft:
		parent.children[0] = noffset
	case childRight:
		parent.children[1] = noffset
	}
}

func (t *Tree) deleteBalance(b, parent *Block) {
	if b.c != colorDoubleBlack {
		return
	}

	if b.ct == childRoot {
		b.c = colorBlack
		return
	}

	var sibling *Block
	// Acquire sibling
	switch {
	case b.ct == childLeft:
		sibling = t.getBlock(parent.children[1])
	case b.ct == childRight:
		sibling = t.getBlock(parent.children[0])
	case b.ct == childRoot:
		b.c = colorBlack
		return
	}

	// Set sibling black state
	// Note: We can eventually remove this for performance reasons once this function
	// is completely fleshed out. We need to ensure that we are not dealing with a nil sibling
	// This is just to avoid running into panic land
	siblingIsBlack := isBlack(sibling)
	var leftNephew, rightNephew *Block
	if sibling != nil {
		leftNephew = t.getBlock(sibling.children[0])
		rightNephew = t.getBlock(sibling.children[1])
	}

	// Sibling rotate bonanza
	switch {
	// Sibling is black and has both black children
	case siblingIsBlack && (isBlack(leftNephew) && isBlack(rightNephew)):
		if sibling != nil {
			sibling.c = colorRed
		}

		if parent.c == colorRed {
			parent.c = colorBlack
			return
		} else if parent.c == colorBlack {
			parent.c = colorDoubleBlack
			t.deleteBalance(parent, t.getBlock(parent.parent))
		}

	// Sibling is black and has at least one red child
	case siblingIsBlack:
		// Rotation cases:
		switch {
		// 1. Left Left Case (s is left child of its parent and r is left child of s or both children of s are red).
		case sibling.ct == childLeft && isRed(leftNephew):
			// Right rotate sibling
			t.rightRotate(sibling)
			// Recolor left nephew to black
			leftNephew.c = colorBlack

		// 2. Left Right Case (s is left child of its parent and r is right child).
		case sibling.ct == childLeft && isRed(rightNephew):
			// Left rotate right nephew
			t.leftRotate(rightNephew)
			// Right rotate right nephew
			t.rightRotate(rightNephew)

			// Note: Sibling is now left nephew and right nephew is now sibling

		// 3. Right Right Case (s is right child of its parent and r is right child of s or both children of s are red)
		case sibling.ct == childRight && isRed(rightNephew):
			// Left rotate sibling
			t.leftRotate(sibling)
			// Recolor right nephew to black
			rightNephew.c = colorBlack

		// 4. Right Left Case (s is right child of its parent and r is left child of s)
		case sibling.ct == childRight && isRed(leftNephew):
			// Right rotate left nephew
			t.rightRotate(leftNephew)
			// Left rotate left nephew
			t.leftRotate(leftNephew)

			// Note: Sibling is now right nephew and left nephew is now sibling
		}

	// Sibling is red
	default:
		if sibling.ct == childLeft {
			t.rightRotate(sibling)
		} else if sibling.ct == childRight {
			t.leftRotate(sibling)
		}
	}

	b.c = colorBlack
}

// Delete will remove an item from the tree
func (t *Tree) Delete(key []byte) {
	var (
		b      *Block
		next   *Block
		offset int64
	)

	if offset, _ = t.seekBlock(t.t.root, key, false); offset == -1 {
		return
	}

	b = t.getBlock(offset)

	parent := t.getBlock(b.parent)
	hasLeft := b.children[0] != -1
	hasRight := b.children[1] != -1

	// BST Delete switch
	switch {
	case hasLeft && hasRight:
		next = t.twoChildDelete(b, parent)

	case !hasLeft && !hasRight:
		t.zeroChildrenDelete(b, parent)
		// We are just using this as a placeholder for next
		next = b
	default:
		// Technically this is out of order, but it seems much more clean to check to see
		// if we have ALL or NONE. If neither cases exist, we know we have one child
		next = t.oneChildDelete(b, parent)

	}

	// Balancing cases
	if b.c == colorRed || next.c == colorRed {
		// Simple Case: If either u or v is red
		// Note: Because we are not disrupting the black-level, no rotation is needed
		next.c = colorBlack
	} else {
		next.c = colorDoubleBlack
	}

	t.deleteBalance(next, parent)

	root := t.getBlock(t.t.root)
	if root != nil && root.ct != childRoot {
		// Root has changed, update root reference to the new root
		t.t.root = root.parent
	}

	t.t.cnt--
	return
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

// Grow will grow a blob value to a given size
func (t *Tree) Grow(key []byte, sz int64) (bs []byte) {
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

	if grew = t.growBlob(b, key, sz); grew {
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

	bs = t.getValue(b)
	return
}

// Reset will clear the tree and keep the backend. Can be used as a fresh store
func (t *Tree) Reset() {
	t.a.Reset()
	t.t.root = -1
}

// Len will return the length of the data-store
func (t *Tree) Len() (n int) {
	return int(t.t.cnt)
}

// Close will close a tree
func (t *Tree) Close() (err error) {
	return t.b.Close()
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
