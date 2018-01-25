# rbTree [![GoDoc](https://godoc.org/github.com/itsmontoya/rbTree?status.svg)](https://godoc.org/github.com/itsmontoya/rbTree) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)
rbTree is a simple red-black tree for storing data in a sorted manner

## Benchmarks
```bash
## go --version
## > go version go1.9.3 linux/amd64

# rbTree
BenchmarkGet-16                     500    3672168 ns/op         0 B/op        0 allocs/op # Fastest
BenchmarkSortedGetPut-16            200    7257202 ns/op      7208 B/op        0 allocs/op # Fastest
BenchmarkSortedPut-16               300    3897592 ns/op      4805 B/op        0 allocs/op # Fastest
BenchmarkReversePut-16              300    3831036 ns/op      4805 B/op        0 allocs/op # Fastest
BenchmarkRandomPut-16               300    4314028 ns/op      4805 B/op        0 allocs/op # Fastest
BenchmarkForEach-16               20000      94837 ns/op         0 B/op        0 allocs/op

# Skiplist (github.com/OneOfOne/skiplist)
BenchmarkSkiplistGet-16             500    2610549 ns/op         0 B/op        0 allocs/op
BenchmarkSkiplistSortedGetPut-16    200    8072147 ns/op    323778 B/op    10100 allocs/op
BenchmarkSkiplistSortedPut-16       300    5324434 ns/op    322520 B/op    10066 allocs/op
BenchmarkSkiplistReversePut-16      300    6091060 ns/op    322521 B/op    10066 allocs/op
BenchmarkSkiplistRandomPut-16       200    6747146 ns/op    323783 B/op    10100 allocs/op
BenchmarkSkiplistForEach-16       30000      51041 ns/op         0 B/op        0 allocs/op # Fastest

# Standard library map (Used as a maximum speed measurement, not sorted like the others)
BenchmarkMapGet-16                 3000     450179 ns/op         0 B/op        0 allocs/op
BenchmarkMapSortedGetPut-16        2000     717369 ns/op       791 B/op        0 allocs/op
BenchmarkMapSortedPut-16           3000     518363 ns/op       522 B/op        0 allocs/op
BenchmarkMapReversePut-16          2000     550347 ns/op       786 B/op        0 allocs/op
BenchmarkMapRandomPut-16           2000     514930 ns/op       789 B/op        0 allocs/op
BenchmarkMapForEach-16            10000     151997 ns/op         0 B/op        0 allocs/op
```

## Memory usage
```bash
# Memory usage test involves setting 1 million keys and checking the total allocations with 
# the populated data-store (See github.com/itsmontoya/rbTree/testing/allocs for source)

# rbTree
80005120 
# map (with pre-set length)
88565128 
# skiplist
91062448
```