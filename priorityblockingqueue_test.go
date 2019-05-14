package guc

import (
	"container/list"
	"testing"
)

func TestNewPriorityBlockingQueue(t *testing.T) {
	p := NewPriorityBlockingQueue()
	if !p.IsEmpty() {
		t.Fatal("queue should be empty")
	}
}

type sampleBlockingItem struct {
	Value int
}

func (this *sampleBlockingItem) Equals(i interface{}) bool {
	dst, ok := i.(*sampleBlockingItem)
	if !ok {
		return false
	}
	return this.Value == dst.Value
}

func (this *sampleBlockingItem) CompareTo(i interface{}) int {
	return this.Value - i.(*sampleBlockingItem).Value
}

func newSampleBlockingItem(value int) *sampleBlockingItem {
	return &sampleBlockingItem{
		Value: value,
	}
}

func newPreparedPriorityBlockingQueue() *PriorityBlockingQueue {
	p := NewPriorityBlockingQueue()
	p.Add(newSampleBlockingItem(6))
	p.Add(newSampleBlockingItem(8))
	p.Add(newSampleBlockingItem(3))
	p.Add(newSampleBlockingItem(6))
	p.Add(newSampleBlockingItem(33))
	p.Add(newSampleBlockingItem(7))
	p.Add(newSampleBlockingItem(2))
	return p
}

func TestPriorityBlockingQueue_Add(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.IsEmpty() {
		t.Fatal("queue should not be empty")
	}
	if p.priorityQueue.data.data[0].(*sampleBlockingItem).Value != 2 {
		t.Fatal("first should be value 2")
	}
	if p.Size() != 7 {
		t.Fatal("queue size should be 7")
	}
}

func TestPriorityBlockingQueue_Contains(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.Contains(newSampleBlockingItem(100)) {
		t.Fatal("queue should not contains value 100")
	}
	if !p.Contains(newSampleBlockingItem(6)) {
		t.Fatal("queue should contains value 6")
	}
	if !p.Contains(newSampleBlockingItem(3)) {
		t.Fatal("queue should contains value 3")
	}
}

func TestPriorityBlockingQueue_ForEach(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()

	l := list.New()
	p.ForEach(func(i interface{}) {
		l.PushBack(i)
	})
	if l.Len() != 7 {
		t.Fatal("total number of foreach items are not 7")
	}
}

func TestPriorityBlockingQueue_ToArray(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()

	arr := p.ToArray()
	if len(arr) != 7 {
		t.Fatal("array size must be 7")
	}
	matched := false
	for _, v := range arr {
		if v.(*sampleBlockingItem).Value == 6 {
			matched = true
			break
		}
	}
	if !matched {
		t.Fatal("must have value 6 in array")
	}
}

func TestPriorityBlockingQueue_FillArray(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()

	arr := p.FillArray(make([]interface{}, 0))
	if len(arr) != 7 {
		t.Fatal("array size must be 7")
	}
	matched := false
	for _, v := range arr {
		if v.(*sampleBlockingItem).Value == 6 {
			matched = true
			break
		}
	}
	if !matched {
		t.Fatal("must have value 6 in array")
	}

	newArr := make([]interface{}, 7)
	p.FillArray(newArr)
	matched = false
	for _, v := range newArr {
		if v.(*sampleBlockingItem).Value == 6 {
			matched = true
			break
		}
	}
	if !matched {
		t.Fatal("must have value 6 in array")
	}
}

//TODO need to test iterator methods
func TestPriorityBlockingQueue_Iterator(t *testing.T) {
}
