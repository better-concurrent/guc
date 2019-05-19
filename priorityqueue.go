package guc

import (
	"container/heap"
	"unsafe"
)

var _ Queue = new(PriorityQueue)

// we don't want to expose golang heap api to users
// so here we create a struct that implements heap api
type priorityData struct {
	data       []interface{}
	queue      *PriorityQueue
	comparator Comparator
}

type PriorityQueue struct {
	data priorityData
}

type priorityQueueIter struct {
	idx   int
	data  []interface{}
	queue *PriorityQueue
}

func (this *priorityQueueIter) HasNext() bool {
	this.idx++
	return this.idx < len(this.data)
}

func (this *priorityQueueIter) Next() interface{} {
	r := this.data[this.idx]
	return r
}

func (this *priorityQueueIter) Remove() {
	heap.Remove(&this.queue.data, this.idx)
	this.data = this.queue.data.data
	this.idx--
}

func (this *priorityQueueIter) ForEachRemaining(consumer func(i interface{})) {
	for this.HasNext() {
		consumer(this.Next())
	}
}

func NewPriority() *PriorityQueue {
	queue := &PriorityQueue{
		data: priorityData{},
	}
	queue.data.queue = queue
	return queue
}

func NewPriorityWithComparator(comparator Comparator) *PriorityQueue {
	queue := &PriorityQueue{
		data: priorityData{
			comparator: comparator,
		},
	}
	queue.data.queue = queue
	return queue
}

func (this priorityData) Len() int {
	return len(this.data)
}

func (this priorityData) Less(i, j int) bool {
	c := this.comparator
	data := this.data
	if c != nil {
		if c.Compare(data[i], data[j]) < 0 {
			return true
		} else {
			return false
		}
	} else {
		if data[i].(Comparable).CompareTo(data[j]) < 0 {
			return true
		} else {
			return false
		}
	}
}

func (this priorityData) Swap(i, j int) {
	data := this.data
	data[j], data[i] = data[i], data[j]
}

func (this *priorityData) Push(x interface{}) {
	this.data = append(this.data, x)
}

func (this *priorityData) Pop() interface{} {
	old := this.data
	n := len(old)
	i := old[n-1]
	old[n-1] = nil //clear index, in order to avoid memory leak
	this.data = old[:n-1]
	return i
}

func (this *PriorityQueue) Iterator() Iterator {
	iter := new(priorityQueueIter)
	iter.idx = -1
	iter.data = this.data.data
	iter.queue = this
	return iter
}

func (this *PriorityQueue) ForEach(consumer func(i interface{})) {
	for _, v := range this.data.data {
		consumer(v)
	}
}

func (this *PriorityQueue) Size() int {
	return len(this.data.data)
}

func (this *PriorityQueue) IsEmpty() bool {
	return len(this.data.data) == 0
}

func (this *PriorityQueue) Contains(i interface{}) bool {
	data := this.data.data
	for _, v := range data {
		if v.(Object).Equals(i) {
			return true
		}
	}
	return false
}

func (this *PriorityQueue) ToArray() []interface{} {
	data := this.data.data
	result := make([]interface{}, 0, len(data))
	for _, v := range data {
		result = append(result, v)
	}
	return result
}

func (this *PriorityQueue) FillArray(arr []interface{}) []interface{} {
	data := this.data.data
	if len(arr) >= len(data) {
		for i, v := range data {
			arr[i] = v
		}
		return arr[:len(data)]
	} else {
		return this.ToArray()
	}
}

func (this *PriorityQueue) Add(i interface{}) bool {
	heap.Push(&this.data, i)
	return true
}

func (this *PriorityQueue) Remove(item interface{}) bool {
	for i, v := range this.data.data {
		if v.(Object).Equals(item) {
			heap.Remove(&this.data, i)
			return true
		}
	}
	return false
}

func (this *PriorityQueue) ContainsAll(coll Collection) bool {
	iter := coll.Iterator()
	for iter.HasNext() {
		if !this.Contains(iter.Next()) {
			return false
		}
	}
	return true
}

func (this *PriorityQueue) AddAll(coll Collection) bool {
	changed := false
	iter := coll.Iterator()
	for iter.HasNext() {
		changed = true
		heap.Push(&this.data, iter.Next())
	}
	return changed
}

func (this *PriorityQueue) RemoveAll(coll Collection) bool {
	removed := false
	iter := coll.Iterator()
	for iter.HasNext() {
		r := this.Remove(iter.Next())
		if r {
			removed = true
		}
	}
	return removed
}

func (this *PriorityQueue) RemoveIf(predicate func(i interface{}) bool) bool {
	idx := -1
	for i, v := range this.data.data {
		if predicate(v) {
			idx = i
			break
		}
	}
	if idx >= 0 {
		heap.Remove(&this.data, idx)
		return true
	} else {
		return false
	}
}

func (this *PriorityQueue) RetainAll(coll Collection) bool {
	iter := this.Iterator()
	changed := false
	for iter.HasNext() {
		if !coll.Contains(iter.Next()) {
			iter.Remove()
			changed = true
		}
	}
	return changed
}

func (this *PriorityQueue) Clear() {
	this.data.data = make([]interface{}, 0)
}

func (this *PriorityQueue) Equals(i interface{}) bool {
	p, ok := i.(*PriorityQueue)
	if ok {
		return p == this
	}
	return false
}

func (this *PriorityQueue) HashCode() int {
	return int(uintptr(unsafe.Pointer(this)))
}

func (this *PriorityQueue) Offer(i interface{}) bool {
	heap.Push(&this.data, i)
	return true
}

func (this *PriorityQueue) RemoveHead() interface{} {
	var i interface{} = nil
	if !this.IsEmpty() {
		i = heap.Pop(&this.data)
	}
	if i == nil {
		panic("queue is empty")
	}
	return i
}

func (this *PriorityQueue) Poll() interface{} {
	if this.IsEmpty() {
		return nil
	}
	i := heap.Pop(&this.data)
	return i
}

func (this *PriorityQueue) Element() interface{} {
	if this.IsEmpty() {
		panic("queue is empty")
	} else {
		return this.data.data[0]
	}
}

func (this *PriorityQueue) Peek() interface{} {
	if this.IsEmpty() {
		return nil
	} else {
		return this.data.data[0]
	}
}
