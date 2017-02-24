package rbTree

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
func New() *Tree {
	return &Tree{}
}

// Tree is a red-black tree data structure
type Tree struct {
	root *node
	cnt  uint32
}

// Get will retrieve an item from a tree
func (t *Tree) Get(key string) (val interface{}) {
	if t.root == nil {
		// Root doesn't exist, return early
		return
	}

	if n := t.root.getNode(key, false); n != nil {
		// Node was found, set value as the node's value
		val = n.val
	}

	return
}

// Put will insert an item into the tree
func (t *Tree) Put(key string, val interface{}) {
	var n *node
	if t.root == nil {
		// Root doesn't exist, we can create one
		n = newNode(key)
		t.root = n
	} else {
		// Find node whose key matches our provided key, if node does not exist - create it.
		n = t.root.getNode(key, true)
	}

	n.val = val
	// Balance tree after insert
	// TODO: This can be moved into the node-creation portion
	n.balance()
	// TODO: Remove this, I don't believe it's actually needed now that the percolation is working properly
	t.root.balance()

	if t.root.ct != childRoot {
		// Root has changed, update root reference to the new root
		t.root = t.root.parent
	}

	t.cnt++
}

// ForEach will iterate through each tree item
func (t *Tree) ForEach(fn ForEachFn) (ended bool) {
	if t.root == nil {
		// Root doesn't exist, return early
		return
	}

	// Call iterate from root
	return t.root.iterate(fn)
}

// Len will return the length of the data-store
func (t *Tree) Len() (n int) {
	return int(t.cnt)
}
