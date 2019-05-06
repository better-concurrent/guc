package guc

import (
	"fmt"
	"runtime"
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

func TestTabAt(t *testing.T) {
	n0 := node{hash: 1}
	tab := make([]unsafe.Pointer, 4)
	if tabAt(&tab, 0) != nil {
		t.Error("tab at error")
	}
	if casTabAt(&tab, 0, nil, &n0) {
		if tabAt(&tab, 0) != &n0 {
			t.Error("tab at error")
		}
	}
}

type innerStruct struct {
	j int32
}

type keyObject struct {
	i     int32
	s     string
	inner innerStruct
}

type keyObject2 struct {
	i int
}

type valueObject struct {
	v string
}

func TestBasicOperation(t *testing.T) {
	key0 := keyObject{i: 0, s: "a", inner: innerStruct{32}}
	value0 := valueObject{v: "v"}

	cmap := ConcurrentHashMap{}
	cmap.init(16, 4)
	fmt.Println(cmap)
	oldValue := cmap.Store(key0, value0)
	if oldValue != nil {
		t.Fatalf("Store error")
	}
	oldValue1 := cmap.Store(key0, value0)
	if oldValue1 == nil || value0 != oldValue1 {
		t.Fatalf("Store error")
	}
	value1, _ := cmap.Load(key0)
	if value0 != value1 {
		t.Fatalf("Load error")
	}
	fmt.Println(cmap)
}

func TestMapResize(t *testing.T) {
	cmap := NewConcurrentHashMap(4, 4)
	total := 32
	for i := 0; i < total; i++ {
		key := keyObject{i: int32(i), s: "a", inner: innerStruct{32}}
		value0 := valueObject{v: "v"}
		cmap.Store(key, value0)
		_, ok := cmap.Load(key)
		if !ok {
			t.Fatalf("get error")
		}
	}
	if cmap.Size() != total {
		t.Fatalf("cmap size is %d\n", cmap.Size())
	}
	cmap.printTableDetail()
}

func TestContendedCell(t *testing.T) {
	cc := CounterCell{}
	fmt.Println(unsafe.Sizeof(cc))
	if unsafe.Sizeof(cc) != CacheLineSize {
		t.Fatalf("padding error")
	}
}

func TestMultiGoroutine(t *testing.T) {
	runtime.GOMAXPROCS(4)
	gc := 4
	countPerG := 1024 * 32
	cmap := NewConcurrentHashMap(4, 4)
	endCh := make(chan int)
	value0 := valueObject{v: "v"}
	for i := 0; i < gc; i++ {
		go func() {
			begin := i * countPerG
			for n := 0; n < begin+countPerG; n++ {
				key := keyObject2{i: n}
				cmap.Store(key, value0)
			}
			endCh <- 1
		}()
	}
	for i := 0; i < gc; i++ {
		<-endCh
	}
	cmap.printTableDetail()
	cmap.printCountDetail()
	fmt.Println("end")
}
