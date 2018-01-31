package testUtils

import (
	"math/rand"
	"sort"
	"strconv"
)

// GetStrSlice will get a string slice
func GetStrSlice(in []int) (out []KV) {
	out = make([]KV, 0, len(in))

	for _, v := range in {
		var kv KV
		kv.Key = strconv.Itoa(v)
		kv.Val = []byte(kv.Key)
		out = append(out, kv)
	}

	return
}

// GetSorted will get a sorted int slice
func GetSorted(n int) (s []int) {
	s = make([]int, n)

	for i := 0; i < n; i++ {
		s[i] = i
	}

	return
}

// GetReverse will get a reversed int slice
func GetReverse(n int) (s []int) {
	s = GetSorted(n)
	_ = sort.Reverse(sort.IntSlice(s))
	return
}

// GetRand will get a random ized int slice
func GetRand(n int) (s []int) {
	return rand.Perm(n)
}

// KV pair
type KV struct {
	Key string
	Val []byte
}
