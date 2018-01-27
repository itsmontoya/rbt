package main

import (
	"fmt"
	"runtime"
	"strconv"
)

var val interface{}

func main() {
	var start, end runtime.MemStats
	runtime.ReadMemStats(&start)
	m := populateN(1000000)
	runtime.GC()
	val = m["1"]
	runtime.ReadMemStats(&end)
	fmt.Println(end.Alloc - start.Alloc)
}

func populateN(n int) (m map[string][]byte) {
	m = make(map[string][]byte, n)

	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		m[key] = []byte(key)
	}

	return
}
