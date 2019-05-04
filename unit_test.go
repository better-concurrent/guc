package guc

import (
	"fmt"
	"sync"
	"sync/atomic"
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

type testS3 struct {
	i int32
	//arr []testS
	arr unsafe.Pointer
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

	if !Nilinterequal(unsafe.Pointer(unpackEFace(*&test1)), unsafe.Pointer(unpackEFace(test2))) {
		t.Error("equals error!")
	}
}

func TestUnsafeLoad(t *testing.T) {
	test := testS3{i: 1, arr: nil}
	fmt.Println(test)
	arr := []testS{{i: 1}}
	value := atomic.Value{}
	fmt.Println(value.Load())
	value.Store(arr)
	fmt.Println(value.Load())

	arrp := [2]*testS{}
	arrp[0] = &testS{i: 1}
	arrup := [2]unsafe.Pointer{}
	p := unsafe.Pointer(&testS{i: 3})
	arrup[0] = p
	atomic.StorePointer(&arrup[1], p)
	fmt.Println((*testS)(atomic.LoadPointer(&arrup[0])))
	fmt.Println((*testS)(atomic.LoadPointer(&arrup[1])))
}

type testS4 struct {
	i int32
	t *testS
}

func (test *testS4) getT() interface{} {
	return *test.t
}

func TestObjectCopy(t *testing.T) {
	test := testS4{i: 1, t: &testS{1, "s"}}
	fmt.Println(test)
	var local interface{}
	local = test.getT()
	fmt.Println(&local)
	fmt.Println(unpackEFace(local).rtype)
	fmt.Println(unpackEFace(local).data)
	fmt.Printf("%T\n", local)
	fmt.Println(local == *test.t)
}

type testS5 struct {
	i int32
	p unsafe.Pointer // is *[]*testS
}

func TestPointerLoad(t *testing.T) {
	test := testS5{i: 2, p: nil}
	fmt.Println(test)
	arr := make([]*testS, test.i)
	arr[0] = &testS{i: 1, s: "a"}
	test.p = unsafe.Pointer(&arr)
	pp := (*[]*testS)(test.p)
	fmt.Println((*pp)[0])

	newT := &testS{i: 2, s: "b"}
	//p := unsafe.Pointer((*pp)[0])
	p := (unsafe.Pointer)((*pp)[0])
	fmt.Println((*testS)(p))
	test.p = unsafe.Pointer(newT)
	atomic.StorePointer(&p, unsafe.Pointer(newT))
	fmt.Println((*testS)(p))
	fmt.Println((*pp)[0])
}

func TestMutex(t *testing.T) {
	fmt.Println(unsafe.Sizeof(sync.Mutex{}))
}

func TestGolangInt(t *testing.T) {
	for i := 0; i < 100000; i++ {
		n := int(Fastrand()) & 0xffffffff
		if n < 0 {
			fmt.Println(n)
		}
	}
}

func TestArrayPointRef(t *testing.T) {
	arr := make([]testS, 2)
	arr[0].s = "a"
	arr[1].s = "b"
	ap := &arr
	acp0 := &(*ap)[0]
	acp1 := &(*ap)[0]
	fmt.Printf("%p\n", acp0)
	fmt.Printf("%p\n", acp1)
}
