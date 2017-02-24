# redBlack [![GoDoc](https://godoc.org/github.com/itsmontoya/redBlack?status.svg)](https://godoc.org/github.com/itsmontoya/redBlack) ![Status](https://img.shields.io/badge/status-alpha-red.svg)
redBlack is a simple red-black tree for storing data in a sorted manner

## Benchmarks
```bash
## go --version
## > go version go1.8 linux/amd64

# redBlack
BenchmarkGet-4                    1000     1764255 ns/op        0 B/op        0 allocs/op # Fastest gets
BenchmarkSortedGetPut-4            500     3942287 ns/op    81280 B/op    10020 allocs/op # Fastest get/put
BenchmarkSortedPut-4              1000     2456036 ns/op    80640 B/op    10010 allocs/op # Fastest sorted put
BenchmarkReversePut-4             1000     2275650 ns/op    80640 B/op    10010 allocs/op # Fastest reverse put
BenchmarkRandomPut-4               500     3331017 ns/op    81287 B/op    10020 allocs/op # Fastest random put
BenchmarkForEach-4               20000       78819 ns/op        0 B/op        0 allocs/op

# Harmonic (github.com/itsmontoya/harmonic)
BenchmarkHarmonicSortedGetPut-4    100    10869235 ns/op    85211 B/op    10103 allocs/op
BenchmarkHarmonicSortedPut-4       200     6376283 ns/op    82605 B/op    10051 allocs/op
BenchmarkHarmonicReversePut-4      200     6094193 ns/op    82606 B/op    10051 allocs/op
BenchmarkHarmonicRandomPut-4       100    19305654 ns/op    85810 B/op    10104 allocs/op
BenchmarkHarmonicForEach-4       20000       69101 ns/op        0 B/op        0 allocs/op

# Skiplist (github.com/OneOfOne/skiplist)
BenchmarkSkiplistGet-4             300     6066402 ns/op   160000 B/op    10000 allocs/op
BenchmarkSkiplistSortedGetPut-4    100    14057207 ns/op   407568 B/op    30200 allocs/op
BenchmarkSkiplistSortedPut-4       300     5907818 ns/op   242525 B/op    20066 allocs/op
BenchmarkSkiplistReversePut-4      200     5559227 ns/op   243776 B/op    20100 allocs/op
BenchmarkSkiplistRandomPut-4       200     9626887 ns/op   243785 B/op    20100 allocs/op
BenchmarkSkiplistForEach-4       30000      57052 ns/op         0 B/op        0 allocs/op # Fastest iteration

# Standard library map (Used as a maximum speed measurement, not sorted like the others)
BenchmarkMapGet-4                 5000      339508 ns/op        0 B/op        0 allocs/op
BenchmarkMapSortedGetPut-4        1000     1348527 ns/op    81292 B/op    10000 allocs/op
BenchmarkMapSortedPut-4           2000     1009477 ns/op    80643 B/op    10000 allocs/op
BenchmarkMapReversePut-4          2000      889043 ns/op    80642 B/op    10000 allocs/op
BenchmarkMapRandomPut-4           2000     1041122 ns/op    80643 B/op    10000 allocs/op
BenchmarkMapForEach-4            10000      159575 ns/op        0 B/op        0 allocs/op

```
