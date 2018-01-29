package rbTree

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/OneOfOne/skiplist"
	cznic "github.com/cznic/b"
)

var (
	testSortedList  = getSorted(10000)
	testReverseList = getReverse(10000)
	testRandomList  = getRand(10000)

	testSortedListStr  = getStrSlice(testSortedList)
	testReverseListStr = getStrSlice(testReverseList)
	testRandomListStr  = getStrSlice(testRandomList)

	testVal      []byte
	testCznicVal interface{}
)

func TestBasic(t *testing.T) {
	tr, err := NewMMAP("data", "mmap.db", 64)
	if err != nil {
		t.Fatal(err)
	}

	tr.Put([]byte("1"), []byte("1"))
	tr.Put([]byte("2"), []byte("2"))
	tr.Put([]byte("3"), []byte("3"))
	tr.Put([]byte("4"), []byte("4"))
	tr.Put([]byte("5"), []byte("5"))
	tr.Put([]byte("6"), []byte("6"))
	tr.Put([]byte("7"), []byte("7"))
	tr.Put([]byte("8"), []byte("8"))
	tr.Put([]byte("9"), []byte("9"))
	tr.Put([]byte("10"), []byte("10"))

	if val := string(tr.Get([]byte("1"))); val != "1" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("2"))); val != "2" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("3"))); val != "3" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("4"))); val != "4" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("5"))); val != "5" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("6"))); val != "6" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("7"))); val != "7" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("8"))); val != "8" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("9"))); val != "9" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("10"))); val != "10" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	tr.Close()

	if tr, err = NewMMAP("data", "mmap.db", 64); err != nil {
		t.Fatal(err)
	}

	if val := string(tr.Get([]byte("1"))); val != "1" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("2"))); val != "2" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("3"))); val != "3" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("4"))); val != "4" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("5"))); val != "5" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("6"))); val != "6" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("7"))); val != "7" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("8"))); val != "8" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("9"))); val != "9" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

	if val := string(tr.Get([]byte("10"))); val != "10" {
		t.Fatalf("invalid value, expected \"%v\" and received \"%v\"", "1", val)
	}

}

func TestSortedPut(t *testing.T) {
	testPut(t, getSorted(10))
}

func TestReversePut(t *testing.T) {
	testPut(t, getReverse(10))
}

func TestRandomPut(t *testing.T) {
	testPut(t, getRand(10))
}

func BenchmarkRBTGet(b *testing.B) {
	benchGet(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTSortedGetPut(b *testing.B) {
	benchGetPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTSortedPut(b *testing.B) {
	benchPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTReversePut(b *testing.B) {
	benchPut(b, testReverseListStr)
	b.ReportAllocs()
}

func BenchmarkRBTRandomPut(b *testing.B) {
	benchPut(b, testRandomListStr)
	b.ReportAllocs()
}

func BenchmarkRBTForEach(b *testing.B) {
	benchForEach(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTMMapGet(b *testing.B) {
	benchMMAPGet(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTMMapSortedGetPut(b *testing.B) {
	benchMMAPGetPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTMMapSortedPut(b *testing.B) {
	benchMMAPPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkRBTMMapReversePut(b *testing.B) {
	benchMMAPPut(b, testReverseListStr)
	b.ReportAllocs()
}

func BenchmarkRBTMMapRandomPut(b *testing.B) {
	benchMMAPPut(b, testRandomListStr)
	b.ReportAllocs()
}

func BenchmarkRBTMMapForEach(b *testing.B) {
	benchMMAPForEach(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkMapGet(b *testing.B) {
	benchMapGet(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkMapSortedGetPut(b *testing.B) {
	benchMapGetPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkMapSortedPut(b *testing.B) {
	benchMapPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkMapReversePut(b *testing.B) {
	benchMapPut(b, testReverseListStr)
	b.ReportAllocs()
}

func BenchmarkMapRandomPut(b *testing.B) {
	benchMapPut(b, testRandomListStr)
	b.ReportAllocs()
}

func BenchmarkMapForEach(b *testing.B) {
	benchMapForEach(b, testSortedListStr)
	b.ReportAllocs()
}
func BenchmarkSkiplistGet(b *testing.B) {
	benchSkiplistGet(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkSkiplistSortedGetPut(b *testing.B) {
	benchSkiplistGetPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkSkiplistSortedPut(b *testing.B) {
	benchSkiplistPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkSkiplistReversePut(b *testing.B) {
	benchSkiplistPut(b, testReverseListStr)
	b.ReportAllocs()
}

func BenchmarkSkiplistRandomPut(b *testing.B) {
	benchSkiplistPut(b, testRandomListStr)
	b.ReportAllocs()
}

func BenchmarkSkiplistForEach(b *testing.B) {
	benchSkiplistForEach(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkCznicGet(b *testing.B) {
	benchCznicGet(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkCznicSortedGetPut(b *testing.B) {
	benchCznicGetPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkCznicSortedPut(b *testing.B) {
	benchCznicPut(b, testSortedListStr)
	b.ReportAllocs()
}

func BenchmarkCznicReversePut(b *testing.B) {
	benchCznicPut(b, testReverseListStr)
	b.ReportAllocs()
}

func BenchmarkCznicRandomPut(b *testing.B) {
	benchCznicPut(b, testRandomListStr)
	b.ReportAllocs()
}

func BenchmarkCznicForEach(b *testing.B) {
	benchCznicForEach(b, testSortedListStr)
	b.ReportAllocs()
}

func testPut(t *testing.T, s []int) {
	cnt := len(s)
	tr := New(1024 * 1024)
	tm := make(map[string][]byte, cnt)

	for _, v := range s {
		key := fmt.Sprintf("%012d", v)
		val := []byte(strconv.Itoa(v))
		tr.Put([]byte(key), val)
		tm[key] = val
	}

	for key, mv := range tm {
		val := tr.Get([]byte(key))
		if !bytes.Equal(val, mv) {
			t.Fatalf("invalid value:\nKey: %s\nExpected: %v\nReturned: %v\n", key, mv, val)
		}
	}

	var fecnt int
	tr.ForEach(func(key, val []byte) (end bool) {
		if !bytes.Equal(val, tm[string(key)]) {
			t.Fatalf("invalid value:\nKey: %s\nExpected: %v\nReturned: %v\n", key, tm[string(key)], val)
		}

		fecnt++
		return
	})

	if fecnt != cnt {
		t.Fatalf("invalid ForEach iterations:\nExpected: %v\nActual: %v\n", cnt, fecnt)
	}
}

func benchGet(b *testing.B, s []kv) {
	tr := New(1024 * 1024)
	for _, kv := range s {
		tr.Put(kv.val, kv.val)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			testVal = tr.Get(kv.val)
		}
	}
}

func benchPut(b *testing.B, s []kv) {
	tr := New(1024 * 1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			tr.Put(kv.val, kv.val)
		}
	}
}

func benchGetPut(b *testing.B, s []kv) {
	tr := New(1024 * 1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			tr.Put(kv.val, kv.val)
			testVal = tr.Get(kv.val)
		}
	}
}

func benchForEach(b *testing.B, s []kv) {
	tr := New(1024 * 1024)

	for _, kv := range s {
		tr.Put(kv.val, kv.val)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tr.ForEach(func(_, val []byte) (end bool) {
			testVal = val
			return
		})
	}
}

func benchMMAPGet(b *testing.B, s []kv) {
	tr, err := NewMMAP("data", "test.db", 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove("/data/test.db")

	for _, kv := range s {
		tr.Put(kv.val, kv.val)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			testVal = tr.Get(kv.val)
		}
	}
}

func benchMMAPPut(b *testing.B, s []kv) {
	tr, err := NewMMAP("data", "test.db", 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	//	defer os.Remove("/data/test.db")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			tr.Put(kv.val, kv.val)
		}
	}
}

func benchMMAPGetPut(b *testing.B, s []kv) {
	tr, err := NewMMAP("data", "test.db", 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove("/data/test.db")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			tr.Put(kv.val, kv.val)
			testVal = tr.Get(kv.val)
		}
	}
}

func benchMMAPForEach(b *testing.B, s []kv) {
	tr, err := NewMMAP("data", "test.db", 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove("/data/test.db")

	for _, kv := range s {
		tr.Put(kv.val, kv.val)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tr.ForEach(func(_, val []byte) (end bool) {
			testVal = val
			return
		})
	}
}

func benchMapGet(b *testing.B, s []kv) {
	m := make(map[string][]byte)
	for _, kv := range s {
		m[kv.key] = kv.val
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			testVal = m[kv.key]
		}
	}
}

func benchMapPut(b *testing.B, s []kv) {
	b.ResetTimer()
	m := make(map[string][]byte)

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			m[kv.key] = kv.val
		}
	}
}

func benchMapGetPut(b *testing.B, s []kv) {
	b.ResetTimer()
	m := make(map[string][]byte)

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			testVal = m[kv.key]
			m[kv.key] = kv.val
		}
	}
}

func benchMapForEach(b *testing.B, s []kv) {
	m := make(map[string][]byte)
	for _, kv := range s {
		m[kv.key] = kv.val
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, val := range m {
			testVal = val
		}
	}
}

func benchSkiplistGet(b *testing.B, s []kv) {
	sl := skiplist.New(32)
	for _, kv := range s {
		sl.Set(kv.key, kv.val)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			testVal = sl.Get(kv.key).([]byte)
		}
	}
}

func benchSkiplistPut(b *testing.B, s []kv) {
	b.ResetTimer()
	sl := skiplist.New(32)
	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			sl.Set(kv.key, kv.val)
		}
	}
}

func benchSkiplistGetPut(b *testing.B, s []kv) {
	b.ResetTimer()
	sl := skiplist.New(32)

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			sl.Set(kv.key, kv.val)
			testVal = sl.Get(kv.key).([]byte)
		}
	}
}

func benchSkiplistForEach(b *testing.B, s []kv) {
	sl := skiplist.New(32)
	for _, kv := range s {
		sl.Set(kv.key, kv.val)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.ForEach(func(_ string, val interface{}) bool {
			testVal = val.([]byte)
			return false
		})
	}
}

func benchCznicGet(b *testing.B, s []kv) {
	tr := cznic.TreeNew(byteLess)
	for _, kv := range s {
		tr.Put(kv.val, func(_ interface{}, exists bool) (interface{}, bool) {
			return kv.val, true
		})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			testCznicVal, _ = tr.Get(kv.val)
		}
	}
}

func byteLess(a, b interface{}) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

func benchCznicPut(b *testing.B, s []kv) {
	tr := cznic.TreeNew(byteLess)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			tr.Put(kv.val, func(_ interface{}, exists bool) (interface{}, bool) {
				return kv.val, true
			})
		}
	}
}

func benchCznicGetPut(b *testing.B, s []kv) {
	tr := cznic.TreeNew(byteLess)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range s {
			tr.Put(kv.val, func(_ interface{}, exists bool) (interface{}, bool) {
				return kv.val, true
			})
			testCznicVal, _ = tr.Get(kv.val)
		}
	}
}

func benchCznicForEach(b *testing.B, s []kv) {
	tr := cznic.TreeNew(byteLess)

	for _, kv := range s {
		tr.Put(kv.val, func(_ interface{}, exists bool) (interface{}, bool) {
			return kv.val, true
		})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e, ok := tr.Seek([]byte(""))
		if !ok {
			b.Fatal("error calling iterator")
		}

		var (
			v   interface{}
			err error
			cnt int
		)

		for _, v, err = e.Next(); err == nil; _, v, err = e.Next() {
			testCznicVal = v
			cnt++
		}

		if cnt != len(s) {
			b.Fatalf("invalid count, expected %v and received %v (%v)", len(s), cnt, err)
		}
	}
}

func getStrSlice(in []int) (out []kv) {
	out = make([]kv, 0, len(in))

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
