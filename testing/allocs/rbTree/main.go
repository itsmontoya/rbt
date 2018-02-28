package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/itsmontoya/rbt"
)

var val interface{}

func main() {
	var start, end runtime.MemStats
	runtime.ReadMemStats(&start)
	t := populateN(1000000)
	runtime.GC()
	val = t.Get([]byte("1"))
	runtime.ReadMemStats(&end)
	fmt.Println(end.Alloc - start.Alloc)
}

func populateN(n int) (t *rbt.Tree) {
	t = rbt.New(int64(n) * 32)

	for i := 0; i < n; i++ {
		key := []byte(strconv.Itoa(i))
		t.Put(key, key)
	}

	return
}
