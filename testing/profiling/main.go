package main

import (
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/itsmontoya/rbt"
	"github.com/missionMeteora/journaler"

	"github.com/pkg/profile"
)

func main() {
	list := getSorted(10000)
	strs := getStrSlice(list)
	tr := rbt.New(10000)
	runtime.GC()
	time.Sleep(time.Second * 3)
	journaler.Notification("Values initialized, test starting")
	p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	defer p.Stop()

	for i := 0; i < 1000; i++ {
		for _, kv := range strs {
			tr.Put(kv.val, kv.val)
		}
	}
}

func getStrSlice(in []int) (out []kv) {
	out = make([]kv, len(in))

	for _, v := range in {
		var kv kv
		kv.key = strconv.Itoa(v)
		kv.val = []byte(kv.key)
		out = append(out, kv)
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

type kv struct {
	key string
	val []byte
}
