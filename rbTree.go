package rbTree

import "strings"

const (
	colorBlack color = iota
	colorRed
)

const (
	childRoot childType = iota
	childLeft
	childRight
)

// New will return a new Tree
func New(sz int) *Tree {
	var t Tree
	t.root = -1
	t.nodes = make([]node, 0, sz)
	return &t
}

// Tree is a red-black tree data structure
type Tree struct {
	nodes []node
	root  int
	tail  int
	cnt   uint32
}

// getNode will return a node matching the provided key. It create is set to true, a new node will be created if no match is found
func (t *Tree) getNode(nidx int, key string, create bool) (idx int) {
	idx = -1
	if nidx == -1 {
		return
	}

	switch strings.Compare(t.nodes[nidx].key, key) {
	case 1:
		child := t.nodes[nidx].children[1]
		if child == -1 {
			if !create {
				return
			}

			tn := newNode(key)
			tn.ct = childRight
			tn.parent = nidx
			idx = len(t.nodes)
			t.nodes[nidx].children[1] = idx
			t.nodes = append(t.nodes, tn)
			return
		}

		return t.getNode(child, key, create)

	case -1:
		child := t.nodes[nidx].children[0]
		if child == -1 {
			if !create {
				return
			}

			tn := newNode(key)
			tn.ct = childLeft
			tn.parent = nidx
			idx = len(t.nodes)
			t.nodes[nidx].children[0] = idx
			t.nodes = append(t.nodes, tn)
			return
		}

		return t.getNode(child, key, create)

	case 0:
		return nidx
	}

	return
}

// getHead will get the very first item starting from a given node
// Note: If called from root, will return the first item in the tree
func (t *Tree) getHead(nidx int) (idx int) {
	idx = -1

	if nidx == -1 {
		return
	}

	for {
		if child := t.nodes[nidx].children[0]; child != -1 {
			nidx = child
			continue
		}

		return nidx
	}
}

// getTail will get the very last item starting from a given node
// Note: If called from root, will return the last item in the tree
func (t *Tree) getTail(nidx int) (idx int) {
	idx = -1

	if nidx == -1 {
		return
	}

	for {
		if child := t.nodes[nidx].children[1]; child != -1 {
			nidx = child
			continue
		}

		return nidx
	}
}

func (t *Tree) getUncle(nidx int) (idx int) {
	idx = -1
	parent := t.nodes[nidx].parent
	grandparent := t.nodes[parent].parent

	if grandparent == -1 {
		return
	}

	switch t.nodes[parent].ct {
	case childLeft:
		return t.nodes[grandparent].children[1]
	case childRight:
		return t.nodes[grandparent].children[0]

	}

	return
}

func (t *Tree) setColor(nidx int, c color) {
	t.nodes[nidx].c = c
}

func (t *Tree) setChildType(nidx int, ct childType) {
	t.nodes[nidx].ct = ct
}

func (t *Tree) setParent(nidx, pidx int) {
	t.nodes[nidx].parent = pidx
}

func (t *Tree) setParentChild(nidx, cidx int) {
	switch t.nodes[nidx].ct {
	case childLeft:
		t.nodes[t.getParent(nidx)].children[0] = cidx
	case childRight:
		t.nodes[t.getParent(nidx)].children[1] = cidx
	case childRoot:
		// No action is taken, tree will handle this at the end of put
	}
}

func (t *Tree) getParent(nidx int) (idx int) {
	return t.nodes[nidx].parent
}

func (t *Tree) getGrandparent(nidx int) (idx int) {
	return t.getParent(t.getParent(nidx))
}

func (t *Tree) balance(nidx int) {
	switch {
	case t.nodes[nidx].c == colorBlack:
		return
	case t.nodes[nidx].ct == childRoot:
		if t.nodes[nidx].c == colorRed {
			t.nodes[nidx].c = colorBlack
			return
		}

	case t.isRed(t.getUncle(nidx)):
		t.setColor(t.nodes[nidx].parent, colorBlack)
		t.setColor(t.getUncle(nidx), colorBlack)
		grandparent := t.getGrandparent(nidx)
		t.setColor(grandparent, colorRed)
		// Balance grandparent
		t.balance(grandparent)

	case t.isRed(t.nodes[nidx].parent):
		parent := t.nodes[nidx].parent
		grandparent := t.getParent(parent)

		if t.isTriangle(nidx) {
			t.rotateParent(nidx)
			// Balance parent
			t.balance(parent)
		} else {
			// Is a line
			t.rotateGrandparent(nidx)
			// Balance grandparent
			t.balance(grandparent)
		}
	}
}

func (t *Tree) leftRotate(nidx int) {
	parent := t.nodes[nidx].parent
	grandparent := t.getParent(parent)

	// Swap  children
	swapChild := t.nodes[nidx].children[0]
	t.nodes[parent].children[1] = swapChild
	t.nodes[nidx].children[0] = parent

	if swapChild != -1 {
		t.setParent(swapChild, parent)
		t.setChildType(swapChild, childRight)
	}

	// Set nidx as the child for our grandparent
	t.setParentChild(parent, nidx)
	// Set n's grandparent as parent
	t.setParent(nidx, grandparent)
	// Set n as parent to the original parent
	t.setParent(parent, nidx)

	// Set child types
	t.setChildType(nidx, t.nodes[parent].ct)
	t.setChildType(parent, childLeft)
}

func (t *Tree) rightRotate(nidx int) {
	parent := t.nodes[nidx].parent
	grandparent := t.getParent(parent)

	// Swap  children
	swapChild := t.nodes[nidx].children[1]
	t.nodes[parent].children[0] = swapChild
	t.nodes[nidx].children[1] = parent

	if swapChild != -1 {
		t.setParent(swapChild, parent)
		t.setChildType(swapChild, childLeft)
	}

	// Set nidx as the child for our grandparent
	t.setParentChild(parent, nidx)
	// Set n's grandparent as parent
	t.setParent(nidx, grandparent)
	// Set n as parent to the original parent
	t.setParent(parent, nidx)

	// Set child types
	t.setChildType(nidx, t.nodes[parent].ct)
	t.setChildType(parent, childRight)
}

func (t *Tree) rotateParent(nidx int) {
	switch t.nodes[nidx].ct {
	case childLeft:
		t.rightRotate(nidx)

	case childRight:
		t.leftRotate(nidx)

	default:
		panic("invalid child type for parent rotation")
	}
}

func (t *Tree) rotateGrandparent(nidx int) {
	parent := t.getParent(nidx)
	grandparent := t.getParent(parent)

	switch t.nodes[parent].ct {
	case childLeft:
		t.rightRotate(parent)

	case childRight:
		t.leftRotate(parent)

	default:
		panic("invalid child type for grandparent rotation")
	}

	t.setColor(parent, colorBlack)
	t.setColor(grandparent, colorRed)
}

func (n *node) swapColor() {
	if n.c == colorBlack {
		n.c = colorRed
	} else {
		n.c = colorBlack
	}
}

func (t *Tree) isRed(nidx int) (isRed bool) {
	if nidx == -1 {
		return
	}

	return t.nodes[nidx].c == colorRed
}

func (t *Tree) isTriangle(nidx int) (isTriangle bool) {
	if nidx == -1 {
		return
	}

	parent := t.nodes[nidx].parent
	if t.nodes[nidx].ct == childLeft && t.nodes[parent].ct == childRight {
		return true
	}

	if t.nodes[nidx].ct == childRight && t.nodes[parent].ct == childLeft {
		return true
	}

	return
}

func (t *Tree) numBlack(nidx int) (nb int) {
	if nidx == -1 {
		return
	}

	if t.nodes[nidx].c == colorBlack {
		nb = 1
	}

	if child := t.nodes[nidx].children[0]; child != -1 {
		nb += t.numBlack(child)
	}

	if child := t.nodes[nidx].children[1]; child != -1 {
		nb += t.numBlack(child)
	}

	return
}

func (t *Tree) iterate(nidx int, fn ForEachFn) (ended bool) {
	if child := t.nodes[nidx].children[0]; child != -1 {
		if ended = t.iterate(child, fn); ended {
			return
		}
	}

	if ended = fn(t.nodes[nidx].key, t.nodes[nidx].val); ended {
		return
	}

	if child := t.nodes[nidx].children[1]; child != -1 {
		if ended = t.iterate(child, fn); ended {
			return
		}

	}

	return
}

// Get will retrieve an item from a tree
func (t *Tree) Get(key string) (val []byte) {
	if nidx := t.getNode(t.root, key, false); nidx != -1 {
		// Node was found, set value as the node's value
		val = t.nodes[nidx].val
	}

	return
}

// Put will insert an item into the tree
func (t *Tree) Put(key string, val []byte) {
	var nidx int
	if t.root == -1 {
		// Root doesn't exist, we can create one
		tn := newNode(key)
		nidx = len(t.nodes)
		t.nodes = append(t.nodes, tn)
		t.root = nidx
	} else {
		// Find node whose key matches our provided key, if node does not exist - create it.
		nidx = t.getNode(t.root, key, true)
	}

	t.nodes[nidx].val = val
	// Balance tree after insert
	// TODO: This can be moved into the node-creation portion
	t.balance(nidx)

	if t.nodes[t.root].ct != childRoot {
		// Root has changed, update root reference to the new root
		t.root = t.getParent(t.root)
	}

	t.cnt++
}

// ForEach will iterate through each tree item
func (t *Tree) ForEach(fn ForEachFn) (ended bool) {
	if t.root == -1 {
		// Root doesn't exist, return early
		return
	}

	// Call iterate from root
	return t.iterate(t.root, fn)
}

// Len will return the length of the data-store
func (t *Tree) Len() (n int) {
	return int(t.cnt)
}
