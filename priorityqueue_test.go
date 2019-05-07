package guc

import (
	"container/list"
	"testing"
)

func TestNewPriority(t *testing.T) {
	p := NewPriority()
	if p != p.data.queue {
		t.Fatal("underlying queue reference should be equals to queue itself")
	}
}

type comparatorSample struct {
}

func (comparatorSample) Compare(o1, o2 interface{}) int {
	panic("implement me")
}

func TestNewPriorityWithComparator(t *testing.T) {
	p := NewPriorityWithComparator(comparatorSample{})
	if p.data.comparator == nil {
		t.Fatal("underlying comparator should not be nil")
	}
}

type sampleItem struct {
	Value int
}

func (this *sampleItem) Equals(i interface{}) bool {
	dst, ok := i.(*sampleItem)
	if !ok {
		return false
	}
	return this.Value == dst.Value
}

func (this *sampleItem) CompareTo(i interface{}) int {
	return this.Value - i.(*sampleItem).Value
}

func newSampleItem(value int) *sampleItem {
	return &sampleItem{
		Value: value,
	}
}

func newPreparedPriorityQueue() *PriorityQueue {
	p := NewPriority()
	p.Add(newSampleItem(6))
	p.Add(newSampleItem(8))
	p.Add(newSampleItem(3))
	p.Add(newSampleItem(6))
	p.Add(newSampleItem(33))
	p.Add(newSampleItem(7))
	p.Add(newSampleItem(2))
	return p
}

func TestPriorityQueue_Add(t *testing.T) {
	p := NewPriority()
	if !p.IsEmpty() {
		t.Fatal("queue should be empty")
	}
	b := p.Add(newSampleItem(10))
	if !b {
		t.Fatal("add result not true")
	}
	if len(p.data.data) != 1 {
		t.Fatal("data len should be 1")
	}
	if p.data.data[0].(*sampleItem).Value != 10 {
		t.Fatal("item value is not the origin value")
	}
	if p.IsEmpty() {
		t.Fatal("queue should not be empty")
	}
	if p.Size() != 1 {
		t.Fatal("queue size should be 1")
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	p := NewPriority()
	p.Add(newSampleItem(10))
	b := p.Remove(newSampleItem(10))
	if !b {
		t.Fatal("remove result should be true")
	}
	if !p.IsEmpty() {
		t.Fatal("queue should be empty")
	}
	if p.Size() != 0 {
		t.Fatal("queue size should be zero")
	}

	p.Add(newSampleItem(10))
	if p.Remove(newSampleItem(20)) {
		t.Fatal("remove result should be false")
	}
	if p.Size() != 1 {
		t.Fatal("queue size should be 1")
	}
}

func TestPriorityQueue_Add_Poll_Multi(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	prev := -1
	for i := 0; i < 7; i++ {
		s := p.Poll().(*sampleItem)
		if prev == -1 {
			prev = s.Value
			continue
		}
		if prev > s.Value {
			t.Fatal("prev value:", prev, "should be <= current value:", s.Value)
		}
		prev = s.Value
	}
}

func TestPriorityQueue_ForEach(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	l := list.New()
	p.ForEach(func(i interface{}) {
		l.PushBack(i)
	})
	if l.Len() != 7 {
		t.Fatal("total number of foreach items are not 7")
	}
}

func TestPriorityQueue_Contains(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	if !p.Contains(newSampleItem(6)) {
		t.Fatal("should have value 6 item")
	}
	if p.Contains(newSampleItem(100)) {
		t.Fatal("should not have value 100 item")
	}
}

func TestPriorityQueue_ToArray(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	arr := p.ToArray()
	if len(arr) != 7 {
		t.Fatal("array size must be 7")
	}
	matched := false
	for _, v := range arr {
		if v.(*sampleItem).Value == 6 {
			matched = true
			break
		}
	}
	if !matched {
		t.Fatal("must have value 6 in array")
	}
}

func TestPriorityQueue_FillArray(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	arr := p.FillArray(make([]interface{}, 0))
	if len(arr) != 7 {
		t.Fatal("array size must be 7")
	}
	matched := false
	for _, v := range arr {
		if v.(*sampleItem).Value == 6 {
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
		if v.(*sampleItem).Value == 6 {
			matched = true
			break
		}
	}
	if !matched {
		t.Fatal("must have value 6 in array")
	}
}

func TestPriorityQueue_RemoveIf(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	b := p.RemoveIf(func(i interface{}) bool {
		return i.(*sampleItem).Value == 6
	})
	if !b {
		t.Fatal("remove result should be true")
	}
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}

	b = p.RemoveIf(func(i interface{}) bool {
		return i.(*sampleItem).Value == 100
	})
	if b {
		t.Fatal("remove result should be false")
	}
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}
}

func TestPriorityQueue_Clear(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}
	p.Clear()
	if p.Size() != 0 {
		t.Fatal("queue size should be 0")
	}
	if !p.IsEmpty() {
		t.Fatal("queue should be empty")
	}
}

func TestPriorityQueue_Equals(t *testing.T) {
	p := newPreparedPriorityQueue()
	if !p.Equals(p) {
		t.Fatal("queue should be equals to itself")
	}
	if p.Equals(newPreparedPriorityQueue()) {
		t.Fatal("queue should not be equals to a new one")
	}
	if p.Equals(struct{}{}) {
		t.Fatal("queue should not be equals to another object of other type")
	}
}

func TestPriorityQueue_HashCode(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.HashCode() != p.HashCode() {
		t.Fatal("hashcode must same")
	}
}

func TestPriorityQueue_Offer(t *testing.T) {
	p := newPreparedPriorityQueue()
	p.Offer(newSampleItem(100))
	if p.Size() != 8 {
		t.Fatal("queue size must be 8")
	}
	if !p.Contains(newSampleItem(100)) {
		t.Fatal("queue must contains value 100")
	}
}

func TestPriorityQueue_RemoveHead(t *testing.T) {
	p := newPreparedPriorityQueue()
	h := p.RemoveHead().(*sampleItem)
	if h.Value != 2 {
		t.Fatal("remove head must value 2")
	}
	if p.Size() != 6 {
		t.Fatal("queue size must be 6")
	}

	emptyQueue := NewPriority()
	r := func(p *PriorityQueue) (result bool) {
		result = false
		defer func() {
			err := recover()
			if err != nil {
				result = true
			}
		}()
		p.RemoveHead()
		return
	}(emptyQueue)
	if !r {
		t.Fatal("should panic when RemoveHead of an empty queue")
	}
}

func TestPriorityQueue_Element(t *testing.T) {
	p := newPreparedPriorityQueue()
	h := p.Element().(*sampleItem)
	if h.Value != 2 {
		t.Fatal("element() must value 2")
	}

	emptyQueue := NewPriority()
	r := func(p *PriorityQueue) (result bool) {
		result = false
		defer func() {
			err := recover()
			if err != nil {
				result = true
			}
		}()
		p.Element()
		return
	}(emptyQueue)
	if !r {
		t.Fatal("should panic when element() of an empty queue")
	}
}

func TestPriorityQueue_Peek(t *testing.T) {
	p := newPreparedPriorityQueue()
	h := p.Peek().(*sampleItem)
	if h.Value != 2 {
		t.Fatal("remove head must value 2")
	}

	emptyQueue := NewPriority()
	r := emptyQueue.Peek()
	if r != nil {
		t.Fatal("peek of empty queue should be nil")
	}
}

func TestPriorityQueue_Iterator(t *testing.T) {
	p := newPreparedPriorityQueue()
	cnt := 0
	matched := false
	iter := p.Iterator()
	if iter == nil {
		t.Fatal("iter should not be nil")
	}
	for iter.HasNext() {
		cnt++
		if iter.Next().(*sampleItem).Value == 6 {
			matched = true
		}
	}
	if cnt != 7 {
		t.Fatal("iter count should be 7")
	}
	if !matched {
		t.Fatal("should contains value 6")
	}
}

func TestPriorityQueueIter_Remove(t *testing.T) {
	p := newPreparedPriorityQueue()
	iter := p.Iterator()
	iter.HasNext()
	iter.Remove()
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}
}

func TestPriorityQueueIter_ForEachRemaining(t *testing.T) {
	p := newPreparedPriorityQueue()
	iter := p.Iterator()

	l := list.New()
	iter.ForEachRemaining(func(i interface{}) {
		l.PushBack(i)
	})
	if l.Len() != 7 {
		t.Fatal("total number of foreach items are not 7")
	}
}

func TestPriorityQueue_ContainsAll(t *testing.T) {
	p := newPreparedPriorityQueue()

	c := NewPriority()
	c.Add(newSampleItem(6))
	c.Add(newSampleItem(2))
	if !p.ContainsAll(c) {
		t.Fatal("should contains all")
	}

	c = NewPriority()
	c.Add(newSampleItem(6))
	c.Add(newSampleItem(100))
	if p.ContainsAll(c) {
		t.Fatal("should not contains all - 1")
	}

	c = NewPriority()
	c.Add(newSampleItem(200))
	c.Add(newSampleItem(100))
	if p.ContainsAll(c) {
		t.Fatal("should not contains all - 2")
	}
}

func TestPriorityQueue_AddAll(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size must be 7")
	}

	c := NewPriority()
	c.Add(newSampleItem(200))
	c.Add(newSampleItem(100))
	if !p.AddAll(c) {
		t.Fatal("add result should be true")
	}

	if !p.Contains(newSampleItem(100)) {
		t.Fatal("queue should contains value 100")
	}
	if !p.Contains(newSampleItem(200)) {
		t.Fatal("queue should contains value 200")
	}
	if p.Size() != 9 {
		t.Fatal("queue size should be 9")
	}
}

func TestPriorityQueue_RemoveAll(t *testing.T) {
	p := newPreparedPriorityQueue()
	if p.Size() != 7 {
		t.Fatal("queue size must be 7")
	}

	c := NewPriority()
	c.Add(newSampleItem(200))
	c.Add(newSampleItem(100))
	if p.RemoveAll(c) {
		t.Fatal("remove all should return false")
	}
	if p.Size() != 7 {
		t.Fatal("queue size should be 7")
	}

	c = NewPriority()
	c.Add(newSampleItem(3))
	c.Add(newSampleItem(100))
	if !p.RemoveAll(c) {
		t.Fatal("remove all should return true")
	}
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}

	c = NewPriority()
	c.Add(newSampleItem(8))
	c.Add(newSampleItem(2))
	if !p.RemoveAll(c) {
		t.Fatal("remove all should return true")
	}
	if p.Size() != 4 {
		t.Fatal("queue size should be 4")
	}
}

func TestPriorityQueue_RetainAll(t *testing.T) {
	{
		p := newPreparedPriorityQueue()
		c := NewPriority()
		c.Add(newSampleItem(200))
		c.Add(newSampleItem(100))
		if !p.RetainAll(c) {
			t.Fatal("remove all should return true")
		}
		if p.Size() != 0 {
			t.Fatal("queue size should be 0")
		}
	}
	{
		p := newPreparedPriorityQueue()
		c := NewPriority()
		c.Add(newSampleItem(3))
		c.Add(newSampleItem(100))
		if !p.RetainAll(c) {
			t.Fatal("remove all should return true")
		}
		if p.Size() != 1 {
			t.Fatal("queue size should be 1")
		}
	}
	{
		p := newPreparedPriorityQueue()
		c := NewPriority()
		c.Add(newSampleItem(8))
		c.Add(newSampleItem(2))
		if !p.RetainAll(c) {
			t.Fatal("remove all should return true")
		}
		if p.Size() != 2 {
			t.Fatal("queue size should be 2")
		}
	}
}
