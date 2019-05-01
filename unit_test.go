package guc

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestRuntimeLink(t *testing.T) {
	fmt.Println("now test " + t.Name())
	if SyncRuntimeCanSpin(0) {
		t.Error("SyncRuntimeCanSpin failed!")
	}
	// nanoTime
	fmt.Println(SyncRuntimeNanoTime())
	// do spin
	SyncRuntimeDoSpin()

	// sema
	var sema uint32 = 0
	c := make(chan string, 1)
	go func() {
		SyncRuntimeSemacquire(&sema)
		c <- "ok"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		SyncRuntimeSemrelease(&sema, false)
	}()

	select {
	case ok := <-c:
		fmt.Println("sema " + ok)
	case <-time.After(2 * time.Second):
		t.Error("sema error!")
	}
}

type testS struct {
	i int32
	s string
}

type testS2 struct {
	i int32
}

func TestEface(t *testing.T) {
	test1 := testS{i: 1, s: "test"}
	fmt.Println(unpackEFace(test1))
}

func TestSizeof(t *testing.T) {
	test := testS{}
	fmt.Println(unsafe.Sizeof(test))
	test2 := testS2{}
	fmt.Println(unsafe.Sizeof(test2))
	var a interface{}
	fmt.Println(unsafe.Sizeof(a))
	var arr [2]interface{}
	fmt.Println(unsafe.Sizeof(arr))
}

func TestStructHash(t *testing.T) {
	test1 := testS{i: 1, s: "test"}
	test2 := testS{i: 1, s: "test"}
	h1 := Nilinterhash(unsafe.Pointer(unpackEFace(test1)), 0xffff)
	h2 := Nilinterhash(unsafe.Pointer(unpackEFace(test2)), 0xffff)
	if h1 != h2 {
		t.Error("hash error!")
	}
}

func TestStructEquals(t *testing.T) {
	test1 := testS{i: 1, s: "test"}
	test2 := testS{i: 1, s: "test"}
	if !Nilinterequal(unsafe.Pointer(unpackEFace(test1)), unsafe.Pointer(unpackEFace(test2))) {
		t.Error("equals error!")
	}
}
