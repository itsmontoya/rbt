package rbTree

type color uint8

type childType uint8

// ForEachFn are used when calling ForEach from a Tree
type ForEachFn func(key, val []byte) (end bool)
