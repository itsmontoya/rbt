package rbTree

func newNode(key string) (n node) {
	// Set node key
	n.key = key
	// All new nodes start as red
	n.c = colorRed
	n.parent = -1
	n.children[0] = -1
	n.children[1] = -1
	return
}

// node is a data node
type node struct {
	c  color
	ct childType

	key string
	val []byte

	parent   int
	children [2]int
}
