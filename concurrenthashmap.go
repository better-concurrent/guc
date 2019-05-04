package guc

import (
	"math/bits"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	defaultCapacity         = 16
	defaultContendCellCount = 8
	loadFactor              = 0.75
	maxCapacity             = 1 << 30
	treeifyThreshold        = 8
	resizeStampBits         = 16
	maxResizers             = (1 << (32 - resizeStampBits)) - 1
	resizeStampShift        = 32 - resizeStampBits
)

var hashSeed = generateHashSeed()

func generateHashSeed() uint32 {
	return Fastrand()
}

type CounterCell struct {
	// Volatile
	value   int64
	padding [CacheLineSize - 4]byte
}

// TODO list
// 1. LongAdder like total count
// 2. bucket tree degenerate, ps: golang has no build-in comparable interface
// 3. iterator
// 4. multi-goroutine cooperate resize
type ConcurrentHashMap struct {
	// The array of bins. Lazily initialized upon first insertion.
	// Volatile, type is []*node
	table unsafe.Pointer
	// The next table to use; non-nil only while resizing.
	// Volatile, type is []*node
	nextTable unsafe.Pointer
	// Table initialization and resizing control
	// When negative, the table is being initialized or resized: -1 for initialization,
	// else -(1 + the number of active resizing threads).  Otherwise,
	// when table is null, holds the initial table size to use upon
	// creation, or 0 for default. After initialization, holds the
	// next element count value upon which to resize the table.
	// Volatile
	sizeCtl int32
	// The next table index (plus one) to split while resizing.
	// Volatile
	transferIndex int32
	// Base counter value, used mainly when there is no contention,
	// but also as a fallback during table initialization
	// races. Updated via CAS.
	// Volatile
	baseCount int64
	// FIXME! j.u.c implementation is too complex, this is a simple version
	// Volatile, type is []CounterCell
	counterCells unsafe.Pointer
	// Spinlock (locked via CAS) used when resizing and/or creating CounterCells.
	// Volatile
	cellsBusy int32
}

// node const
const (
	moved    = -1
	treebin  = -2
	reserved = -3
	hashBits = 0x7fffffff
)

type externNode interface {
	find(n *node, h int32, k *interface{}) (node *node, ok bool)
	isTreeNode() bool
}

// base node
type node struct {
	hash int32
	// FIXME! move to head node, not each node
	m sync.Mutex
	// type is *interface
	key unsafe.Pointer
	// volatile, type is *interface
	val unsafe.Pointer
	// volatile, type is *node
	next unsafe.Pointer

	// FIXME! better design?
	extern externNode
}

func (n *node) getKey() interface{} {
	return *(*interface{})(n.key)
}

func (n *node) getValue() interface{} {
	return *(*interface{})(atomic.LoadPointer(&n.val))
}

func (n *node) getNext() *node {
	return (*node)(atomic.LoadPointer(&n.next))
}

type baseNode struct {
}

func (en *baseNode) find(n *node, h int32, k *interface{}) (*node, bool) {
	panic("NYI")
}

func (en *baseNode) isTreeNode() bool {
	return false
}

type forwardingNode struct {
}

func (en *forwardingNode) find(n *node, h int32, k *interface{}) (*node, bool) {
	panic("NYI")
}

func (en *forwardingNode) isTreeNode() bool {
	return false
}

type treeNode struct {
}

func (en *treeNode) find(n *node, h int32, k *interface{}) (*node, bool) {
	panic("NYI")
}

func (en *treeNode) isTreeNode() bool {
	return true
}

// TODO need test
func spread(hash uintptr) int32 {
	h := int32(hash)
	return (h ^ (h >> 16)) & hashBits
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

func hash(v *interface{}) uintptr {
	return Nilinterhash(unsafe.Pointer(unpackEFace(*v)), uintptr(hashSeed))
}

func equals(v1, v2 *interface{}) bool {
	return v1 == v2 || *v1 == *v2
}

func tabAt(tab *[]unsafe.Pointer, i int32) *node {
	return (*node)(atomic.LoadPointer(&(*tab)[i]))
}

func casTabAt(tab *[]unsafe.Pointer, i int32, c, v *node) bool {
	return atomic.CompareAndSwapPointer(&(*tab)[i], unsafe.Pointer(c), unsafe.Pointer(v))
}

func (m *ConcurrentHashMap) sumCount() int64 {
	cells := m.getCountCells()
	sum := atomic.LoadInt64(&m.baseCount)
	if cells != nil {
		for i := 0; i < len(*cells); i++ {
			c := (*cells)[i]
			sum += atomic.LoadInt64(&c.value)
		}
	}
	return sum
}

func (m *ConcurrentHashMap) getCountCells() *[]CounterCell {
	return (*[]CounterCell)(atomic.LoadPointer(&m.counterCells))
}

func (m *ConcurrentHashMap) getTable() *[]unsafe.Pointer {
	return (*[]unsafe.Pointer)(atomic.LoadPointer(&m.table))
}

func (m *ConcurrentHashMap) getNextTable() *[]unsafe.Pointer {
	return (*[]unsafe.Pointer)(atomic.LoadPointer(&m.nextTable))
}

func (m *ConcurrentHashMap) init(initialCapacity, concurrencyLevel int32) {
	if initialCapacity < 0 {
		panic("initialCapacity should > 0")
	}
	var capacity int32 = 0
	if initialCapacity < concurrencyLevel {
		initialCapacity = concurrencyLevel
	}
	if initialCapacity >= (maxCapacity >> 1) {
		capacity = maxCapacity
	} else {
		capacity = tableSizeFor(initialCapacity + (initialCapacity >> 1) + 1)
	}
	m.sizeCtl = capacity
}

func (m *ConcurrentHashMap) initTable() *[]unsafe.Pointer {
	for {
		tab := m.getTable()
		if tab != nil && len(*tab) > 0 {
			break
		}
		sc := atomic.LoadInt32(&m.sizeCtl)
		if sc < 0 {
			// lost initialization race; just spin
			runtime.Gosched()
		} else if atomic.CompareAndSwapInt32(&m.sizeCtl, sc, -1) {
			tab = m.getTable()
			if tab == nil || len(*tab) == 0 {
				var n int32
				if sc > 0 {
					n = sc
				} else {
					n = defaultCapacity
				}
				arr := make([]unsafe.Pointer, n)
				atomic.StorePointer(&m.table, unsafe.Pointer(&arr))
			}
			atomic.StoreInt32(&m.sizeCtl, sc)
		}
	}
	return m.getTable()
}

func (m *ConcurrentHashMap) Size() int32 {
	sum := m.sumCount()
	if sum < 0 {
		return 0
	} else {
		return int32(sum)
	}
}

func (m *ConcurrentHashMap) IsEmpty() bool {
	return m.sumCount() <= 0
}

func (m *ConcurrentHashMap) Load(key interface{}) (interface{}, bool) {
	if key == nil {
		panic("key is nil!")
	}
	h := spread(hash(&key))
	tab := m.getTable()
	// not initialized
	if tab == nil {
		return nil, false
	}
	// empty table
	n := int32(len(*tab))
	if n == 0 {
		return nil, false
	}
	// bin is empty
	e := tabAt(tab, (n-1)&h)
	if e == nil {
		return nil, false
	}
	eh := e.hash
	if h == eh {
		ek := e.getKey()
		if key == ek {
			return e.getValue(), true
		}
	} else if eh < 0 {
		p, ok := e.extern.find(e, h, &key)
		if ok {
			return p.getValue(), true
		} else {
			return nil, false
		}
	}
	for {
		next := e.getNext()
		if next == nil {
			break
		}
		if h == next.hash && key == next.getKey() {
			return next.getValue(), true
		}
	}
	return nil, false
}

func (m *ConcurrentHashMap) Contains(key interface{}) bool {
	_, ok := m.Load(key)
	return ok
}

func (m *ConcurrentHashMap) Store(key, value interface{}) interface{} {
	return m.storeVal(key, value, false)
}

func (m *ConcurrentHashMap) storeVal(key, value interface{}, onlyIfAbsent bool) interface{} {
	if key == nil || value == nil {
		panic("key or value is null")
	}
	var binCount int32 = 0
	h := spread(hash(&key))
	for {
		tab := m.getTable()
		var n int32
		var f *node
		if tab == nil || len(*tab) == 0 {
			tab = m.initTable()
		} else {
			n = int32(len(*tab)) // length
			i := (n - 1) & h
			f = tabAt(tab, i)
			if f == nil {
				// cas node
				newNode := &node{hash: h, key: unsafe.Pointer(&key),
					val: unsafe.Pointer(&value), next: nil, extern: &baseNode{}}
				if casTabAt(tab, i, nil, newNode) {
					// no lock when adding to empty bin
					break
				}
			} else {
				fh := f.hash
				if fh == moved {
					m.helpTransfer(tab, f)
				} else {
					var oldVal *interface{} = nil
					// slow path
					f.m.Lock()
					// re-check
					if tabAt(tab, i) != f {
						f.m.Unlock()
						continue
					}
					if fh >= 0 {
						binCount = 1
						for e := f; ; binCount++ {
							var ek *interface{}
							if e.hash == h {
								ek = (*interface{})(e.key)
								if key == *ek {
									oldVal = (*interface{})(e.val)
									if !onlyIfAbsent {
										e.val = unsafe.Pointer(&value)
										break
									}
								}
							}
							pred := e
							e = (*node)(e.next)
							if e == nil {
								pred.next = unsafe.Pointer(&node{hash: h, key: unsafe.Pointer(&key),
									val: unsafe.Pointer(&value), next: nil, extern: &baseNode{}})
								break
							}
						}
					} else if f.extern.isTreeNode() {
						panic("NYI")
					}
					f.m.Unlock()
					// treeify
					if binCount != 0 {
						if binCount > treeifyThreshold {
							m.treeifyBin(tab, i)
						}
						if oldVal != nil {
							return *oldVal
						}
					}
				}
			}
		}
	}
	m.addCount(1, binCount)
	return nil
}

func (m *ConcurrentHashMap) helpTransfer(tab *[]unsafe.Pointer, f *node) {
	panic("NYI")
}

// x: the count to add
// check: if <0, don't check resize, if <= 1 only check if uncontended
// FIXME! simple implementation
func (m *ConcurrentHashMap) addCount(x int64, check int32) {
	as := m.getCountCells()
	b := atomic.LoadInt64(&m.baseCount)
	s := b + x
	if as != nil || !atomic.CompareAndSwapInt64(&m.baseCount, b, s) {
		if as == nil {
			m.fullAddCount(x, false)
		} else {
			a := getRandomCountCell(as)
			incrementCountCell(a, x)
		}
		if check <= 1 {
			return
		}
		s = m.sumCount()
	}
	if check >= 0 {
		for {
			sc := atomic.LoadInt32(&m.sizeCtl)
			var tab, nt *[]unsafe.Pointer
			tab = m.getTable()
			if s >= int64(sc) && tab != nil {
				n := len(*tab)
				if n < maxCapacity {
					break
				}
				rs := resizeStamp(int32(n))
				if sc < 0 {
					nt = m.getNextTable()
					if (sc>>resizeStampBits) != rs || sc == rs+1 ||
						sc == rs+maxResizers || nt == nil {
						break
					} else {
						ti := atomic.LoadInt32(&m.transferIndex)
						if ti <= 0 {
							break
						}
					}
					if atomic.CompareAndSwapInt32(&m.sizeCtl, sc, sc+1) {
						m.transfer(tab, nt)
					}
				} else {
					if atomic.CompareAndSwapInt32(&m.sizeCtl, sc, (rs<<resizeStampShift)+2) {
						m.transfer(tab, nil)
					}
				}
				s = m.sumCount()
			} else {
				break
			}
		}
	}
}

func resizeStamp(n int32) int32 {
	return int32(bits.LeadingZeros(uint(n)) | (1 << (resizeStampBits - 1)))
}

// FIXME
func (m *ConcurrentHashMap) fullAddCount(x int64, wasUncontended bool) {
	// TODO hard code
	as := make([]CounterCell, defaultContendCellCount)
	asp := &as
	for {
		if !atomic.CompareAndSwapPointer(&m.counterCells, nil, unsafe.Pointer(asp)) {
			asp = m.getCountCells()
			if asp != nil {
				break
			}
		}
	}
	incrementCountCell(&(*asp)[0], x)
}

func incrementCountCell(a *CounterCell, x int64) {
	for i := 0; ; i++ {
		old := atomic.LoadInt64(&a.value)
		if !atomic.CompareAndSwapInt64(&a.value, old, old+x) {
			if !SyncRuntimeCanSpin(i) {
				runtime.Gosched()
			} else {
				// or sync.runtime_doSpin? FIXME
				continue
			}
		} else {
			break
		}
	}
}

// FIXME just need a random probe in G.m, no need re-rand
func getRandomCountCell(as *[]CounterCell) *CounterCell {
	i := int(Fastrand()) & 0xffffffff
	n := len(*as)
	return &(*as)[i%n]
}

// TODO
func (m *ConcurrentHashMap) treeifyBin(tab *[]unsafe.Pointer, i int32) {
	// NYI, golang has no build-in comparable interface
	return
}

// Moves and/or copies the nodes in each bin to new table.
func (m *ConcurrentHashMap) transfer(tab, nt *[]unsafe.Pointer) {
	// TODO
}
