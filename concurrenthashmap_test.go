package guc

import (
	"fmt"
	"testing"
	"unsafe"
)

type testKey struct {
	i int32
	s string
	k *testKey
}

func TestBase(t *testing.T) {
	testMap := make(map[interface{}]bool)
	k1 := testKey{i: 1, s: ""}
	k2 := testKey{i: 1, s: "1"}
	k2.s = ""
	testMap[k1] = true
	fmt.Println(testMap[k1])
	fmt.Println(testMap[k2])

	iMap := make(map[interface{}]interface{})
	iMap[k1] = 1
	fmt.Println(iMap[k2])

	piMap := make(map[*int32]bool)
	var i1 int32 = 1
	var i2 int32 = 1
	piMap[&i1] = true
	fmt.Println(piMap[&i2])

	pMap := make(map[unsafe.Pointer]bool)
	pMap[unsafe.Pointer(&k1)] = true
	fmt.Println(pMap[unsafe.Pointer(&k1)])
}

func TestTableSizeAlign(t *testing.T) {
	if tableSizeFor(1) != 1 {
		t.Error("resize error")
	}
	if tableSizeFor(14) != 16 {
		t.Error("resize error")
	}
}
