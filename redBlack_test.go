package redBlack

import (
	//	"fmt"
	"fmt"
	"strconv"
	"testing"
)

func TestBasic(t *testing.T) {
	tr := New()

	for i := 10; i > -1; i-- {
		tr.Put(strconv.Itoa(i), i)
		tr.Print()
		fmt.Println("")
	}

	tr.Print()
	//	for i := 0; i < 10000; i++ {
	//		tr.Put(strconv.Itoa(i), i)
	//	}

	tr.ForEach(func(key string, val interface{}) {
		fmt.Println(key, val)
	})

	fmt.Println("Manual test")
	fmt.Println("0", tr.Get("0"))
	fmt.Println("1", tr.Get("1"))
	fmt.Println("2", tr.Get("2"))
	fmt.Println("3", tr.Get("3"))
	fmt.Println("4", tr.Get("4"))
	fmt.Println("5", tr.Get("5"))
	fmt.Println("5", tr.Get("6"))
	fmt.Println("5", tr.Get("7"))
	fmt.Println("5", tr.Get("8"))
	fmt.Println("5", tr.Get("9"))
	fmt.Println("5", tr.Get("10"))
	return
	//	fmt.Println("6", tr.Get("6"))
	//	fmt.Println("Child 0", tr.root.children[0])
	//	fmt.Println("Child 1", tr.root.children[1])
}
