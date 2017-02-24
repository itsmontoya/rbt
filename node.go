package rbTree

import "fmt"

func newNode(key string) *node {
	var n node
	// Set node key
	n.key = key
	// All new nodes start as red
	n.c = colorRed
	return &n
}

// node is a data node
type node struct {
	c  color
	ct childType

	key string
	val interface{}

	parent   *node
	children [2]*node
}

// getNode will return a node matching the provided key. It create is set to true, a new node will be created if no match is found
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

	return
}

// getHead will get the very first item starting from a given node
// Note: If called from root, will return the first item in the tree
func (n *node) getHead() *node {
	if child := n.children[0]; child != nil {
		return child.getHead()
	}

	return n
}

// getTail will get the very last item starting from a given node
// Note: If called from root, will return the last item in the tree
func (n *node) getTail() *node {
	if child := n.children[1]; child != nil {
		return child.getTail()
	}

	return n
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
	case n.c == colorBlack:
		return
	case n.ct == childRoot:
		if n.c == colorRed {
			n.c = colorBlack
			return
		}

	case n.getUncle().isRed():
		n.parent.c = colorBlack
		n.getUncle().c = colorBlack
		n.parent.parent.c = colorRed
		n.parent.parent.balance()

	case n.parent.isRed():
		parent := n.parent
		grandparent := parent.parent

		if n.isTriangle() {
			n.rotateParent()
			parent.balance()
		} else {
			// Is a line
			n.rotateGrandparent()
			grandparent.balance()
		}
	}
}

func (n *node) leftRotate() {
	parent := n.parent

	// Swap  children
	swapChild := n.children[0]
	parent.children[1] = swapChild
	n.children[0] = parent

	if swapChild != nil {
		swapChild.parent = parent
		swapChild.ct = childRight
	}

	// Update grandparent so that n is it's new child
	parent.updateParent(n)
	// Set n's grandparent as parent
	n.parent = parent.parent
	// Set n as parent to the original parent
	parent.parent = n

	// Set child types
	n.ct = parent.ct
	parent.ct = childLeft
}

func (n *node) rightRotate() {
	parent := n.parent

	// Swap  children
	swapChild := n.children[1]
	parent.children[0] = swapChild
	n.children[1] = parent

	if swapChild != nil {
		swapChild.parent = parent
		swapChild.ct = childLeft
	}

	// Update grandparent so that n is it's new child
	parent.updateParent(n)
	// Set n's grandparent as parent
	n.parent = parent.parent
	// Set n as parent to the original parent
	parent.parent = n

	// Set child types
	n.ct = parent.ct
	parent.ct = childRight
}

func (n *node) rotateParent() {
	switch n.ct {
	case childLeft:
		n.rightRotate()

	case childRight:
		n.leftRotate()

	default:
		panic("invalid child type for parent rotation")
	}
}

func (n *node) rotateGrandparent() {
	parent := n.parent
	grandparent := parent.parent

	switch n.parent.ct {
	case childLeft:
		n.parent.rightRotate()

	case childRight:
		n.parent.leftRotate()

	default:
		panic("invalid child type for grandparent rotation")
	}

	parent.c = colorBlack
	grandparent.c = colorRed
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

func (n *node) swapColor() {
	if n.c == colorBlack {
		n.c = colorRed
	} else {
		n.c = colorBlack
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

func (n *node) numBlack() (nb int) {
	if n.c == colorBlack {
		nb = 1
	}

	if child := n.children[0]; child != nil {
		nb += child.numBlack()
	}

	if child := n.children[1]; child != nil {
		nb += child.numBlack()
	}

	return
}

func (n *node) iterate(fn ForEachFn) (ended bool) {
	if child := n.children[0]; child != nil {
		if ended = child.iterate(fn); ended {
			return
		}
	}

	if ended = fn(n.key, n.val); ended {
		return
	}

	if child := n.children[1]; child != nil {
		if ended = child.iterate(fn); ended {
			return
		}
	}

	return
}
