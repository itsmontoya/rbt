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
	var t Tree
	return &t
}

// Tree is a red-black tree data structure
type Tree struct {
	root *node
}

// Get will retrieve an item from a tree
func (t *Tree) Get(key string) (val interface{}) {
	if t.root == nil {
		return
	}

	if n := t.root.getNode(key, false); n != nil {
		val = n.val
	}

	return
}

// Put will insert an item into the tree
func (t *Tree) Put(key string, val interface{}) {
	var n *node
	if t.root == nil {
		n = newNode(key)
		t.root = n
	} else {
		n = t.root.getNode(key, true)
	}

	n.val = val
	n.balance()
	t.root.balance()

	if t.root.ct != childRoot {
		t.root = t.root.parent
	}
}

// ForEach will iterate through each tree item
func (t *Tree) ForEach(fn func(key string, val interface{})) {
	if t.root == nil {
		return
	}

	t.root.iterate(fn)
}
