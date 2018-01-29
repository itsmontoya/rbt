# whiskey [![GoDoc](https://godoc.org/github.com/itsmontoya/whiskey?status.svg)](https://godoc.org/github.com/itsmontoya/whiskey) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)
whiskey is a simple red-black tree for storing data in a sorted manner

## Benchmarks
```bash
## go --version
## > go version go1.9.3 linux/amd64

# Whiskey (byteslice)
BenchmarkWhiskeyGet-16                  1000     1802499 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeySortedGetPut-16          500     3538649 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeySortedPut-16            1000     1941610 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyReversePut-16            500     2005640 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyRandomPut-16             500     2766245 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyForEach-16             20000       95022 ns/op         0 B/op        0 allocs/op

# Whiskey (MMap)
BenchmarkWhiskeyMMapGet-16              1000     1824530 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyMMapSortedGetPut-16      500     3576144 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyMMapSortedPut-16        1000     1989553 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyMMapReversePut-16       1000     1963757 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyMMapRandomPut-16         500     2728154 ns/op         0 B/op        0 allocs/op
BenchmarkWhiskeyMMapForEach-16         20000       93281 ns/op         0 B/op        0 allocs/op

# Skiplist (github.com/OneOfOne/skiplist)
BenchmarkSkiplistGet-16                  500     2505922 ns/op         0 B/op        0 allocs/op
BenchmarkSkiplistSortedGetPut-16         200     8779344 ns/op    323786 B/op    10100 allocs/op
BenchmarkSkiplistSortedPut-16            300     4912010 ns/op    322520 B/op    10066 allocs/op
BenchmarkSkiplistReversePut-16           300     5556532 ns/op    322521 B/op    10066 allocs/op
BenchmarkSkiplistRandomPut-16            200     7497191 ns/op    323779 B/op    10100 allocs/op
BenchmarkSkiplistForEach-16            30000       51236 ns/op         0 B/op        0 allocs/op

# B+Tree (github.com/cznic/b)
BenchmarkCznicGet-16                     300     5902947 ns/op    320000 B/op    10000 allocs/op
BenchmarkCznicSortedGetPut-16            100    14052786 ns/op    963912 B/op    30001 allocs/op
BenchmarkCznicSortedPut-16               200     7624983 ns/op    641956 B/op    20000 allocs/op
BenchmarkCznicReversePut-16              200     7599889 ns/op    641956 B/op    20000 allocs/op
BenchmarkCznicRandomPut-16               200     9212829 ns/op    642221 B/op    20000 allocs/op
# Note: Could not get iteration to work with B+Tree

# Standard library map (Used as a maximum speed measurement, not sorted like the others)
BenchmarkMapGet-16                      3000      471404 ns/op         0 B/op        0 allocs/op
BenchmarkMapSortedGetPut-16             2000      772270 ns/op       787 B/op        0 allocs/op
BenchmarkMapSortedPut-16                2000      518880 ns/op       785 B/op        0 allocs/op
BenchmarkMapReversePut-16               2000      558830 ns/op       790 B/op        0 allocs/op
BenchmarkMapRandomPut-16                2000      538387 ns/op       786 B/op        0 allocs/op
BenchmarkMapForEach-16                 10000      153389 ns/op         0 B/op        0 allocs/op

```

## Memory usage
```bash
# Memory usage test involves setting 1 million keys and checking the total allocations with 
# the populated data-store (See github.com/itsmontoya/whiskey/testing/allocs for source)

# whiskey
88025720 
# map (with pre-set length)
109605384 
# skiplist
123108784
```