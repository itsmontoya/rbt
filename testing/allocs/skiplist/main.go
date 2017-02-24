package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/OneOfOne/skiplist"
)

var val interface{}

func main() {
	var start, end runtime.MemStats
	runtime.ReadMemStats(&start)
	s := populateN(1000000)
	runtime.GC()
	val = s.Get("1")
	runtime.ReadMemStats(&end)
	fmt.Println(end.Alloc - start.Alloc)
}

func populateN(n int) (s *skiplist.List) {
	s = skiplist.New(32)

	for i := 0; i < n; i++ {
		s.Set(strconv.Itoa(i), i)
	}

	return
}
