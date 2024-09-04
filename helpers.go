package rbt

type color uint8

type childType uint8

type trunk struct {
	root int64
	cnt  int64
}

// ForEachFn is used when calling ForEach from a Tree
type ForEachFn func(key, val []byte) (end bool)
