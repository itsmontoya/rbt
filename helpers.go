package rbTree

import (
	"fmt"

	fcolor "github.com/fatih/color"
)

type color uint8

type childType uint8

func writeStrN(str string, n int) string {
	b := make([]byte, 0, len(str)*n)
	for i := 0; i < n; i++ {
		b = append(b, str...)
	}

	return string(b)
}

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
