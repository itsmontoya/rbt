package redBlack

import (
	"fmt"
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

// New will return a new Tree
func New() *Tree {
	var t Tree
	return &t
}

// Tree is a data structure
// Red-black tree specifications:
// - Root of tree is black
// - There are no two adjacent red nodes (A red node cannot have a red parent or red child).
// - Every path from root to a NULL node has same number of black nodes.
type Tree struct {
	root *node
}

func (t *Tree) balance() {
	if t.isRootRed() {
		t.root.c = colorBlack
	}

	if t.hasAdjacentReds() {
		// Fix adjacent reds
	}

	if !t.hasSamePathLength() {
		// Fix same path
	}
}

func (t *Tree) isRootRed() bool {
	return t.root.c == colorRed
}

func (t *Tree) hasAdjacentReds() bool {
	return t.root.hasAdjacentReds()
}

func (t *Tree) hasSamePathLength() bool {
	return true
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
		fmt.Println("Node is root")
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

// Print will print the current tree
func (t *Tree) Print() {
	t.root.print(0)
}

func newNode(key string) *node {
	var n node
	n.key = key
	n.c = colorRed
	return &n
}

// node is a data node
type node struct {
	c  color
	ct childType

	key string
	val interface{}

	parent *node

	children [2]*node

	//	right *node
}

func (n *node) getLeftmost() (ln *node) {
	//	if child := n.children[0]
	return
}

func (n *node) getNode(key string, create bool) (tn *node) {
	switch {
	case key > n.key:
		if n.children[1] == nil {
			if !create {
				return
			}

			tn = newNode(key)
			tn.ct = childRight
			tn.parent = n
			n.children[1] = tn
			return
		}

		return n.children[1].getNode(key, create)

	case key < n.key:
		if n.children[0] == nil {
			if !create {
				return
			}

			tn = newNode(key)
			tn.ct = childLeft
			tn.parent = n
			n.children[0] = tn
			return
		}

		return n.children[0].getNode(key, create)

	case key == n.key:
		return n
	}
	//t.root.key
	return
}

func (n *node) getUncle() (un *node) {
	grandparent := n.parent.parent
	if grandparent == nil {
		return
	}

	switch n.parent.ct {
	case childLeft:
		return grandparent.children[1]
	case childRight:
		return grandparent.children[0]

	}

	return
}

func (n *node) balance() {
	switch {
	case n.ct == childRoot:
		if n.c == colorRed {
			n.c = colorBlack
			return
		}

	case n.getUncle().isRed():
		n.parent.swapColor()
		n.getUncle().swapColor()
		n.parent.parent.swapColor()

	case n.parent.isRed():
		if n.isTriangle() {
			n.rotateParent()
		} else {
			// Is a line
			n.rotateGrandparent()
		}
	}
}

func (n *node) swapColor() {
	if n.c == colorBlack {
		n.c = colorRed
	} else {
		n.c = colorBlack
	}
}

func (n *node) rotateParent() {
	parent := n.parent
	fmt.Println("Rotating parent", n.key, parent.key, parent.parent.key)
	switch n.ct {
	case childLeft:
		fmt.Println("Child left")
		return
		n.ct = parent.ct

		parent.children[0] = nil
		parent.updateParent(n)
		parent.ct = childLeft
		parent.parent = n

		n.children[1] = parent

	case childRight:
		fmt.Println("Child right!")

		n.ct = parent.ct

		parent.children[1] = nil
		parent.updateParent(n)
		parent.ct = childRight
		parent.parent = n

		n.children[0] = parent

	default:
		panic("invalid child type for parent rotation")
	}
}

func (n *node) rotateGrandparent() {
	grandparent := n.parent.parent
	fmt.Println("Rotating grandparent", n.key, n.parent.key, grandparent.key)
	n.parent.parent = grandparent.parent

	// Swap colors
	pc := n.parent.c
	n.parent.c = grandparent.c
	grandparent.c = pc

	switch n.parent.ct {
	case childLeft:
		n.parent.ct = grandparent.ct
		grandparent.children[0] = n.parent.children[1]
		grandparent.updateParent(n.parent)
		grandparent.ct = childRight
		grandparent.parent = n.parent
		n.parent.children[1] = grandparent

	case childRight:
		n.parent.ct = grandparent.ct
		grandparent.children[1] = n.parent.children[0]
		grandparent.updateParent(n.parent)
		grandparent.ct = childLeft
		grandparent.parent = n.parent
		n.parent.children[0] = grandparent

	default:
		panic("invalid child type for grandparent rotation")
	}
}

func (n *node) updateParent(nc *node) {
	switch n.ct {
	case childLeft:
		n.parent.children[0] = nc
	case childRight:
		n.parent.children[1] = nc
	case childRoot:
		// No action is taken, tree will handle this at the end of put
	}
}

func (n *node) isRed() bool {
	if n != nil && n.c == colorRed {
		return true
	}

	return false
}

func (n *node) isTriangle() bool {
	if n.ct == childLeft && n.parent.ct == childRight {
		return true
	}

	if n.ct == childRight && n.parent.ct == childLeft {
		return true
	}

	return false
}

func (n *node) hasAdjacentReds() bool {
	if n == nil {
		return false
	}

	if n.c == colorRed && n.hasRedChildren() {
		return true
	}

	if n.children[0].hasAdjacentReds() {
		return true
	} else if n.children[1].hasAdjacentReds() {
		return true
	}

	return false
}

func (n *node) hasRedChildren() bool {
	return n.children[0].isRed() || n.children[1].isRed()
}

func (n *node) getStrColor() string {
	switch n.c {
	case colorBlack:
		return "black"
	case colorRed:
		return "red"
	default:
		panic("invalid color")
	}
}

// print will print the current node (and it's decendents)
func (n *node) print(indent int) {
	fmt.Print(writeStrN(" ", indent))
	fmt.Printf("%s (%s)\n", n.key, n.getStrColor())
	if child := n.children[0]; child != nil {
		child.print(indent + 4)
	}

	if child := n.children[1]; child != nil {
		child.print(indent + 4)
	}
}

func writeStrN(str string, n int) string {
	b := make([]byte, 0, len(str)*n)
	for i := 0; i < n; i++ {
		b = append(b, str...)
	}

	return string(b)
}

func (n *node) iterate(fn func(key string, val interface{})) {
	if child := n.children[0]; child != nil {
		child.iterate(fn)
	}

	fn(n.key, n.val)

	if child := n.children[1]; child != nil {
		child.iterate(fn)
	}

	return
}

type color uint8

type childType uint8
