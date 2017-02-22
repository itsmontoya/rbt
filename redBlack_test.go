package redBlack

import (
	//"fmt"

	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func TestSortedPut(t *testing.T) {
	testPut(t, getSorted(10000))
}

func TestReversePut(t *testing.T) {
	testPut(t, getReverse(10000))
}

func TestRandomPut(t *testing.T) {
	testPut(t, getRand(10000))
}

func testPut(t *testing.T, s []int) {
	tr := New()
	cnt := len(s)
	tm := make(map[string]interface{}, cnt)

	for _, v := range s {
		key := strconv.Itoa(v)
		tr.Put(key, v)
		tm[key] = v
	}

	var fecnt int
	tr.ForEach(func(key string, val interface{}) {
		if tm[key] != val {
			t.Fatalf("invalid value:\nKey: %s\nExpected: %v\nReturned: %v\n", key, tm[key], val)
		}

		fecnt++
	})

	if fecnt != cnt {
		t.Fatalf("invalid ForEach iterations:\nExpected: %v\nActual: %v\n", cnt, fecnt)
	}

	for key, mv := range tm {
		val := tr.Get(key)
		if val != mv {
			t.Fatalf("invalid value:\nKey: %s\nExpected: %v\nReturned: %v\n", key, mv, val)
		}
	}
}

func getSorted(n int) (s []int) {
	s = make([]int, n)

	for i := 0; i < n; i++ {
		s[i] = i
	}

	return
}

func getReverse(n int) (s []int) {
	s = getSorted(n)
	sort.Reverse(sort.IntSlice(s))
	return
}

func getRand(n int) (s []int) {
	return rand.Perm(n)
}
