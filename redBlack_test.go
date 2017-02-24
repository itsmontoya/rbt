package redBlack

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/OneOfOne/skiplist"
	"github.com/itsmontoya/harmonic"
)

var (
	testSortedList  = getStrSlice(getSorted(10000))
	testReverseList = getStrSlice(getReverse(10000))
	testRandomList  = getStrSlice(getRand(10000))

	testVal interface{}
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

func BenchmarkGet(b *testing.B) {
	benchGet(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkSortedGetPut(b *testing.B) {
	benchGetPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkSortedPut(b *testing.B) {
	benchPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkReversePut(b *testing.B) {
	benchPut(b, testReverseList)
	b.ReportAllocs()
}

func BenchmarkRandomPut(b *testing.B) {
	benchPut(b, testRandomList)
	b.ReportAllocs()
}

func BenchmarkForEach(b *testing.B) {
	benchForEach(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkMapGet(b *testing.B) {
	benchMapGet(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkMapSortedGetPut(b *testing.B) {
	benchMapGetPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkMapSortedPut(b *testing.B) {
	benchMapPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkMapReversePut(b *testing.B) {
	benchMapPut(b, testReverseList)
	b.ReportAllocs()
}

func BenchmarkMapRandomPut(b *testing.B) {
	benchMapPut(b, testRandomList)
	b.ReportAllocs()
}

func BenchmarkMapForEach(b *testing.B) {
	benchMapForEach(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkHarmonicSortedGetPut(b *testing.B) {
	benchHarmonicGetPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkHarmonicSortedPut(b *testing.B) {
	benchHarmonicPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkHarmonicReversePut(b *testing.B) {
	benchHarmonicPut(b, testReverseList)
	b.ReportAllocs()
}

func BenchmarkHarmonicRandomPut(b *testing.B) {
	benchHarmonicPut(b, testRandomList)
	b.ReportAllocs()
}

func BenchmarkHarmonicForEach(b *testing.B) {
	benchHarmonicForEach(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkSkiplistGet(b *testing.B) {
	benchSkiplistGet(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkSkiplistSortedGetPut(b *testing.B) {
	benchSkiplistGetPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkSkiplistSortedPut(b *testing.B) {
	benchSkiplistPut(b, testSortedList)
	b.ReportAllocs()
}

func BenchmarkSkiplistReversePut(b *testing.B) {
	benchSkiplistPut(b, testReverseList)
	b.ReportAllocs()
}

func BenchmarkSkiplistRandomPut(b *testing.B) {
	benchSkiplistPut(b, testRandomList)
	b.ReportAllocs()
}

func BenchmarkSkiplistForEach(b *testing.B) {
	benchSkiplistForEach(b, testSortedList)
	b.ReportAllocs()
}

func testPut(t *testing.T, s []int) {
	tr := New()
	cnt := len(s)
	tm := make(map[string]interface{}, cnt)

	for _, v := range s {
		key := fmt.Sprintf("%012d", v)
		tr.Put(key, v)
		tm[key] = v
	}

	for key, mv := range tm {
		val := tr.Get(key)
		if val != mv {
			t.Fatalf("invalid value:\nKey: %s\nExpected: %v\nReturned: %v\n", key, mv, val)
		}
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
}

func benchGet(b *testing.B, s []string) {
	tr := New()
	for i, key := range s {
		tr.Put(key, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key := range s {
			testVal = tr.Get(key)
		}
	}
}

func benchPut(b *testing.B, s []string) {
	b.ResetTimer()
	tr := New()

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			tr.Put(key, i)
		}
	}
}

func benchGetPut(b *testing.B, s []string) {
	b.ResetTimer()
	tr := New()

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			tr.Put(key, i)
			testVal = tr.Get(key)
		}
	}
}

func benchForEach(b *testing.B, s []string) {
	tr := New()
	for i, key := range s {
		tr.Put(key, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tr.ForEach(func(_ string, val interface{}) {
			testVal = val
		})
	}
}

func benchMapGet(b *testing.B, s []string) {
	m := make(map[string]interface{})
	for i, key := range s {
		m[key] = i
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key := range s {
			testVal = m[key]
		}
	}
}

func benchMapPut(b *testing.B, s []string) {
	b.ResetTimer()
	m := make(map[string]interface{})

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			m[key] = i
		}
	}
}

func benchMapGetPut(b *testing.B, s []string) {
	b.ResetTimer()
	m := make(map[string]interface{})

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			m[key] = i
			testVal = m[key]
		}
	}
}

func benchMapForEach(b *testing.B, s []string) {
	m := make(map[string]interface{})
	for i, key := range s {
		m[key] = i
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, val := range m {
			testVal = val
		}
	}
}

func benchHarmonicPut(b *testing.B, s []string) {
	b.ResetTimer()
	h := harmonic.New(0)

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			h.Put(key, i)
		}
	}
}

func benchHarmonicGetPut(b *testing.B, s []string) {
	b.ResetTimer()
	h := harmonic.New(0)

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			h.Put(key, i)
			testVal, _ = h.Get(key)
		}
	}
}

func benchHarmonicForEach(b *testing.B, s []string) {
	h := harmonic.New(0)
	for i, key := range s {
		h.Put(key, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.ForEach(func(_ string, val interface{}) {
			testVal = val
		})
	}
}

func benchSkiplistGet(b *testing.B, s []string) {
	sl := skiplist.New(32, skiplistCompare)
	for i, key := range s {
		sl.Set(key, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key := range s {
			testVal = sl.Get(key)
		}
	}
}

func benchSkiplistPut(b *testing.B, s []string) {
	b.ResetTimer()
	sl := skiplist.New(32, skiplistCompare)
	for i := 0; i < b.N; i++ {
		for i, key := range s {
			sl.Set(key, i)
		}
	}
}

func benchSkiplistGetPut(b *testing.B, s []string) {
	b.ResetTimer()
	sl := skiplist.New(32, skiplistCompare)

	for i := 0; i < b.N; i++ {
		for i, key := range s {
			sl.Set(key, i)
			testVal = sl.Get(key)
		}
	}
}

func benchSkiplistForEach(b *testing.B, s []string) {
	sl := skiplist.New(32, skiplistCompare)
	for i, key := range s {
		sl.Set(key, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.ForEach(func(_, val interface{}) bool {
			testVal = val
			return false
		})
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
	sort.Reverse(sort.IntSlice(s))
	return
}

func getRand(n int) (s []int) {
	return rand.Perm(n)
}

func skiplistCompare(a, b interface{}) int {
	astr := a.(string)
	bstr := b.(string)

	if astr > bstr {
		return 1
	} else if astr < bstr {
		return -1
	} else {
		return 0
	}
}
