# rbTree [![GoDoc](https://godoc.org/github.com/itsmontoya/rbTree?status.svg)](https://godoc.org/github.com/itsmontoya/rbTree) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)
rbTree is a simple red-black tree for storing data in a sorted manner

## Benchmarks
```bash
## go --version
## > go version go1.9.3 linux/amd64

# rbTree
BenchmarkGet-16                             5000            211060 ns/op               0 B/op          0 allocs/op
BenchmarkSortedGetPut-16                    3000            387460 ns/op             240 B/op          0 allocs/op
BenchmarkSortedPut-16                      10000            217051 ns/op              72 B/op          0 allocs/op
BenchmarkReversePut-16                      5000            214899 ns/op             144 B/op          0 allocs/op
BenchmarkRandomPut-16                      10000            217789 ns/op              72 B/op          0 allocs/op
BenchmarkForEach-16                     100000000               12.7 ns/op             0 B/op          0 allocs/op

# Skiplist (github.com/OneOfOne/skiplist)
BenchmarkSkiplistGet-16                    10000            151220 ns/op               0 B/op          0 allocs/op
BenchmarkSkiplistSortedGetPut-16            5000            299129 ns/op               1 B/op          0 allocs/op
BenchmarkSkiplistSortedPut-16              10000            183248 ns/op               0 B/op          0 allocs/op
BenchmarkSkiplistReversePut-16              5000            230720 ns/op               1 B/op          0 allocs/op
BenchmarkSkiplistRandomPut-16              10000            190457 ns/op               0 B/op          0 allocs/op
BenchmarkSkiplistForEach-16             200000000                6.81 ns/op            0 B/op          0 allocs/op

# Standard library map (Used as a maximum speed measurement, not sorted like the others)
BenchmarkMapGet-16                         10000            121128 ns/op               0 B/op          0 allocs/op
BenchmarkMapSortedGetPut-16                10000            233717 ns/op               0 B/op          0 allocs/op
BenchmarkMapSortedPut-16                   10000            172372 ns/op               0 B/op          0 allocs/op
BenchmarkMapReversePut-16                  10000            170415 ns/op               0 B/op          0 allocs/op
BenchmarkMapRandomPut-16                   10000            170233 ns/op               0 B/op          0 allocs/op
BenchmarkMapForEach-16                  30000000                43.1 ns/op             0 B/op          0 allocs/op

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