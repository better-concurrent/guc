package main

import (
	"sync/atomic"
	"unsafe"
)

const (
	defaultCapacity = 16
	loadFactor      = 0.75
	maxCapacity     = 1 << 30
)

type ConcurrentHashMap struct {
}

type node struct {
	hash int32
	key  interface{}
	val  atomic.Value
	next unsafe.Pointer
}

func (*node) find(h int32, k interface{}) (node *node, ok bool) {
	return nil, false
}

func tableSizeFor(c int32) int32 {
	n := c - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return 1
	} else {
		if n >= maxCapacity {
			return maxCapacity
		} else {
			return n + 1
		}
	}
}
