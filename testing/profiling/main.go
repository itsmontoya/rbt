package main

import (
	"math/rand"
	"runtime"
	"sort"
	"strconv"

	"github.com/missionMeteora/journaler"

	"github.com/itsmontoya/rbTree"
	"github.com/pkg/profile"
)

func main() {
	list := getSorted(10000)
	strs := getStrSlice(list)
	tr := rbTree.New(10000)
	runtime.GC()
	journaler.Notification("Values initialized, test starting")
	p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	defer p.Stop()

	for i := 0; i < 10; i++ {
		for _, v := range list {
			tr.Put(strs[v], v)
		}
	}
}

func getStrSlice(in []int) (out []string) {
	out = make([]string, len(in))

	for i, v := range in {
		out[i] = strconv.Itoa(v)
	}

	return
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
	_ = sort.Reverse(sort.IntSlice(s))
	return
}

func getRand(n int) (s []int) {
	return rand.Perm(n)
}
