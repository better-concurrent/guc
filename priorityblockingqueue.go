package guc

import (
	"sync"
	"time"
	"unsafe"
)

var _ BlockingQueue = new(PriorityBlockingQueue)

type PriorityBlockingQueue struct {
	lock          sync.Mutex
	priorityQueue PriorityQueue
	cond          *sync.Cond
}

type priorityBlockingQueueIter struct {
	idx   int
	data  []interface{}
	queue *PriorityBlockingQueue
}

func (this priorityBlockingQueueIter) HasNext() bool {
	this.idx++
	return this.idx < len(this.data)
}

func (this priorityBlockingQueueIter) Next() interface{} {
	r := this.data[this.idx]
	return r
}

func (this priorityBlockingQueueIter) Remove() {
	this.queue.Remove(this.data[this.idx])
}

func (this priorityBlockingQueueIter) ForEachRemaining(consumer func(i interface{})) {
	for this.HasNext() {
		consumer(this.Next())
	}
}

func (this *PriorityBlockingQueue) Iterator() Iterator {
	arr := this.ToArray()
	return priorityBlockingQueueIter{
		data: arr,
	}
}

func (this *PriorityBlockingQueue) ForEach(consumer func(i interface{})) {
	iter := this.Iterator()
	for iter.HasNext() {
		consumer(iter.Next())
	}
}

func (this *PriorityBlockingQueue) Size() int {
	this.lock.Lock()
	l := this.priorityQueue.Size()
	this.lock.Unlock()
	return l
}

func (this *PriorityBlockingQueue) IsEmpty() bool {
	return this.Size() == 0
}

func (this *PriorityBlockingQueue) Contains(i interface{}) bool {
	this.lock.Lock()
	r := this.priorityQueue.Contains(i)
	this.lock.Unlock()
	return r
}

func (this *PriorityBlockingQueue) ToArray() []interface{} {
	this.lock.Lock()
	data := this.priorityQueue.data.data
	result := make([]interface{}, 0, len(data))
	for _, v := range data {
		result = append(result, v)
	}
	this.lock.Unlock()
	return result
}

func (this *PriorityBlockingQueue) FillArray(arr []interface{}) []interface{} {
	this.lock.Lock()
	data := this.priorityQueue.data.data
	if len(arr) >= len(data) {
		for i, v := range data {
			arr[i] = v
		}
		this.lock.Unlock()
		return arr[:len(data)]
	} else {
		result := make([]interface{}, 0, len(data))
		for _, v := range data {
			result = append(result, v)
		}
		this.lock.Unlock()
		return result
	}
}

func (this *PriorityBlockingQueue) Add(i interface{}) bool {
	return this.Offer(i)
}

func (this *PriorityBlockingQueue) Remove(i interface{}) bool {
	this.lock.Lock()
	r := this.priorityQueue.Remove(i)
	this.lock.Unlock()
	return r
}

func (this *PriorityBlockingQueue) ContainsAll(coll Collection) bool {
	iter := coll.Iterator()
	for iter.HasNext() {
		if !this.Contains(iter.Next()) {
			return false
		}
	}
	return true
}

func (this *PriorityBlockingQueue) AddAll(coll Collection) bool {
	changed := false
	iter := coll.Iterator()
	for iter.HasNext() {
		changed = true
		this.Add(iter.Next())
	}
	return changed
}

func (this *PriorityBlockingQueue) RemoveAll(coll Collection) bool {
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

func (this *PriorityBlockingQueue) RemoveIf(predicate func(i interface{}) bool) bool {
	this.lock.Lock()
	r := this.priorityQueue.RemoveIf(predicate)
	this.lock.Unlock()
	return r
}

func (this *PriorityBlockingQueue) RetainAll(coll Collection) bool {
	this.lock.Lock()
	r := this.priorityQueue.RetainAll(coll)
	this.lock.Unlock()
	return r
}

func (this *PriorityBlockingQueue) Clear() {
	this.lock.Lock()
	this.priorityQueue.Clear()
	this.lock.Unlock()
}

func (this *PriorityBlockingQueue) Equals(i interface{}) bool {
	p, ok := i.(*PriorityBlockingQueue)
	if ok {
		return p == this
	}
	return false
}

func (this *PriorityBlockingQueue) HashCode() int {
	return int(uintptr(unsafe.Pointer(this)))
}

func (this *PriorityBlockingQueue) Offer(i interface{}) bool {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) RemoveHead() interface{} {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) Poll() interface{} {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) Element() interface{} {
	p := this.Peek()
	if p != nil {
		return p
	} else {
		panic("queue is empty")
	}
}

func (this *PriorityBlockingQueue) Peek() interface{} {
	this.lock.Lock()
	p := this.priorityQueue.Peek()
	this.lock.Unlock()
	return p
}

func (this *PriorityBlockingQueue) Put(i interface{}) {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) OfferWithTimeout(i interface{}, t time.Duration) bool {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) Take() interface{} {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) PollWithTimeout(t time.Duration) interface{} {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) RemainingCapacity() int {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) DrainTo(coll interface{}) int {
	//TODO
	panic("implement me")
}

func (this *PriorityBlockingQueue) DrainToWithLimit(coll interface{}, max int) int {
	//TODO
	panic("implement me")
}
