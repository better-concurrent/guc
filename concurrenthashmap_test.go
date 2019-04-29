package main

import (
	"fmt"
	"testing"
)

type testKey struct {
	i int32
	s string
	k *testKey
}

func TestBase(t *testing.T) {
	testMap := make(map[testKey]bool)
	k1 := testKey{i: 1, s: ""}
	k2 := testKey{i: 1, s: ""}
	testMap[k1] = true
	fmt.Println(testMap[k1])
	fmt.Println(testMap[k2])
}

func TestTableSizeAlign(t *testing.T) {
	if tableSizeFor(1) != 1 {
		t.Error("resize error")
	}
	if tableSizeFor(14) != 16 {
		t.Error("resize error")
	}
}
