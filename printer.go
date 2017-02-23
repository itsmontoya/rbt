package redBlack

import (
	"fmt"

	fcolor "github.com/fatih/color"
)

func newPrinter(tr *Tree) *printer {
	var p printer
	p.tr = tr
	return &p
}

type printer struct {
	tr *Tree
	m  []printerPair
}

func (p *printer) Print() {
	printNodes(p.tr.root, 0)
}

type printerPair struct {
	key string
	c   color
}

/*
	fmt.Print(writeStrN(" ", indent))
	fmt.Printf("%s (%s)\n", n.key, n.getStrColor())
	if child := n.children[0]; child != nil {
		child.print(indent + 4)
	}

	if child := n.children[1]; child != nil {
		child.print(indent + 4)
	}
*/

func printNodes(n *node, row int) {
	if n == nil {
		return
	}

	printNode(n, row)
	printNodes(n.children[0], row+1)
	printNodes(n.children[1], row+1)
}

func printNode(n *node, row int) {
	if n == nil {
		return
	}

	key := n.key
	if n.c == colorRed {
		key = fcolor.RedString(key)
	}

	println(row, fmt.Sprintf("%s [%d]", key, row))
}

func print(row int, str string) {
	fmt.Print(writeStrN("    ", row) + str)
}

func println(row int, str string) {
	fmt.Println(writeStrN("    ", row) + str)
}

func writeStrN(str string, n int) string {
	b := make([]byte, 0, len(str)*n)
	for i := 0; i < n; i++ {
		b = append(b, str...)
	}

	return string(b)
}
