package rbt

// Iterator is used when calling ForEach from a Tree
type Iterator func(key, val []byte) (end bool)
