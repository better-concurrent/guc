package guc

import (
	"container/list"
	"testing"
	"time"
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

func TestPriorityBlockingQueue_Remove(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()

	b := p.Remove(newSampleBlockingItem(100))
	if b {
		t.Fatal("should return false")
	}
	if p.Size() != 7 {
		t.Fatal("size should be 7")
	}

	b = p.Remove(newSampleBlockingItem(3))
	if !b {
		t.Fatal("should return true")
	}
	if p.Size() != 6 {
		t.Fatal("size should be 6")
	}

	b = p.Remove(newSampleBlockingItem(6))
	if !b {
		t.Fatal("should return true")
	}
	if p.Size() != 5 {
		t.Fatal("size should be 5")
	}
}

func TestPriorityBlockingQueue_ContainsAll(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()

	c := NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(6))
	c.Add(newSampleBlockingItem(2))
	if !p.ContainsAll(c) {
		t.Fatal("should contains all")
	}

	c = NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(6))
	c.Add(newSampleBlockingItem(100))
	if p.ContainsAll(c) {
		t.Fatal("should not contains all - 1")
	}

	c = NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(200))
	c.Add(newSampleBlockingItem(100))
	if p.ContainsAll(c) {
		t.Fatal("should not contains all - 2")
	}
}

func TestPriorityBlockingQueue_AddAll(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.Size() != 7 {
		t.Fatal("queue size must be 7")
	}

	c := NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(200))
	c.Add(newSampleBlockingItem(100))
	if !p.AddAll(c) {
		t.Fatal("add result should be true")
	}

	if !p.Contains(newSampleBlockingItem(100)) {
		t.Fatal("queue should contains value 100")
	}
	if !p.Contains(newSampleBlockingItem(200)) {
		t.Fatal("queue should contains value 200")
	}
	if p.Size() != 9 {
		t.Fatal("queue size should be 9")
	}
}

func TestPriorityBlockingQueue_RemoveAll(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.Size() != 7 {
		t.Fatal("queue size must be 7")
	}

	c := NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(200))
	c.Add(newSampleBlockingItem(100))
	if p.RemoveAll(c) {
		t.Fatal("remove all should return false")
	}
	if p.Size() != 7 {
		t.Fatal("queue size should be 7")
	}

	c = NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(3))
	c.Add(newSampleBlockingItem(100))
	if !p.RemoveAll(c) {
		t.Fatal("remove all should return true")
	}
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}

	c = NewPriorityBlockingQueue()
	c.Add(newSampleBlockingItem(8))
	c.Add(newSampleBlockingItem(2))
	if !p.RemoveAll(c) {
		t.Fatal("remove all should return true")
	}
	if p.Size() != 4 {
		t.Fatal("queue size should be 4")
	}
}

func TestPriorityBlockingQueue_RemoveIf(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	b := p.RemoveIf(func(i interface{}) bool {
		return i.(*sampleBlockingItem).Value == 6
	})
	if !b {
		t.Fatal("remove result should be true")
	}
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}

	b = p.RemoveIf(func(i interface{}) bool {
		return i.(*sampleBlockingItem).Value == 100
	})
	if b {
		t.Fatal("remove result should be false")
	}
	if p.Size() != 6 {
		t.Fatal("queue size should be 6")
	}
}

func TestPriorityBlockingQueue_RetainAll(t *testing.T) {
	{
		p := newPreparedPriorityBlockingQueue()
		c := NewPriorityBlockingQueue()
		c.Add(newSampleBlockingItem(200))
		c.Add(newSampleBlockingItem(100))
		if !p.RetainAll(c) {
			t.Fatal("remove all should return true")
		}
		if p.Size() != 0 {
			t.Fatal("queue size should be 0")
		}
	}
	{
		p := newPreparedPriorityBlockingQueue()
		c := NewPriorityBlockingQueue()
		c.Add(newSampleBlockingItem(3))
		c.Add(newSampleBlockingItem(100))
		if !p.RetainAll(c) {
			t.Fatal("remove all should return true")
		}
		if p.Size() != 1 {
			t.Fatal("queue size should be 1")
		}
	}
	{
		p := newPreparedPriorityBlockingQueue()
		c := NewPriorityBlockingQueue()
		c.Add(newSampleBlockingItem(8))
		c.Add(newSampleBlockingItem(2))
		if !p.RetainAll(c) {
			t.Fatal("remove all should return true")
		}
		if p.Size() != 2 {
			t.Fatal("queue size should be 2")
		}
	}
}

func TestPriorityBlockingQueue_Clear(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
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

func TestPriorityBlockingQueue_Equals(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if !p.Equals(p) {
		t.Fatal("queue should be equals to itself")
	}
	if p.Equals(newPreparedPriorityBlockingQueue()) {
		t.Fatal("queue should not be equals to a new one")
	}
	if p.Equals(struct{}{}) {
		t.Fatal("queue should not be equals to another object of other type")
	}
}

func TestPriorityBlockingQueue_HashCode(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.HashCode() != p.HashCode() {
		t.Fatal("hashcode must same")
	}
}

func TestPriorityBlockingQueue_RemoveHead(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	h := p.RemoveHead().(*sampleBlockingItem)
	if h.Value != 2 {
		t.Fatal("remove head must value 2")
	}
	if p.Size() != 6 {
		t.Fatal("queue size must be 6")
	}

	emptyQueue := NewPriorityBlockingQueue()
	r := func(p *PriorityBlockingQueue) (result bool) {
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

func TestPriorityBlockingQueue_Poll(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.Size() != 7 {
		t.Fatal("queue size not match")
	}

	prev := -1
	for i := 0; i < 7; i++ {
		s := p.Poll().(*sampleBlockingItem)
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

func TestPriorityBlockingQueue_Peek(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	h := p.Peek().(*sampleBlockingItem)
	if h.Value != 2 {
		t.Fatal("remove head must value 2")
	}

	emptyQueue := NewPriorityBlockingQueue()
	r := emptyQueue.Peek()
	if r != nil {
		t.Fatal("peek of empty queue should be nil")
	}
}

func TestPriorityBlockingQueue_Element(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	h := p.Element().(*sampleBlockingItem)
	if h.Value != 2 {
		t.Fatal("element() must value 2")
	}

	emptyQueue := NewPriorityBlockingQueue()
	r := func(p *PriorityBlockingQueue) (result bool) {
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

func TestPriorityBlockingQueue_Put(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	p.Put(newSampleBlockingItem(100))
	if p.Size() != 8 {
		t.Fatal("queue size must be 8")
	}
	if !p.Contains(newSampleBlockingItem(100)) {
		t.Fatal("queue must contains value 100")
	}
}

func TestPriorityBlockingQueue_Offer(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	p.Offer(newSampleBlockingItem(100))
	if p.Size() != 8 {
		t.Fatal("queue size must be 8")
	}
	if !p.Contains(newSampleBlockingItem(100)) {
		t.Fatal("queue must contains value 100")
	}
}

func TestPriorityBlockingQueue_OfferWithTimeout(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	p.OfferWithTimeout(newSampleBlockingItem(100), 1*time.Hour)
	if p.Size() != 8 {
		t.Fatal("queue size must be 8")
	}
	if !p.Contains(newSampleBlockingItem(100)) {
		t.Fatal("queue must contains value 100")
	}
}

func TestPriorityBlockingQueue_Take(t *testing.T) {
	{
		p := newPreparedPriorityBlockingQueue()
		i := p.Take()
		if i == nil {
			t.Fatal("take from queue should return a value")
		}
	}
	{
		p := NewPriorityBlockingQueue()
		ch := make(chan struct{}, 1)
		go func() {
			<-ch
			time.Sleep(50 * time.Millisecond)
			p.Offer(newSampleBlockingItem(1))
		}()
		ch <- struct{}{}
		i := p.Take()
		if i.(*sampleBlockingItem).Value != 1 {
			t.Fatal("shoudl have value 1")
		}
	}
}

func TestPriorityBlockingQueue_PollWithTimeout(t *testing.T) {
	{
		p := newPreparedPriorityBlockingQueue()
		i := p.PollWithTimeout(1 * time.Second)
		if i == nil {
			t.Fatal("take from queue should return a value")
		}
	}
	{
		p := NewPriorityBlockingQueue()
		ch := make(chan struct{}, 1)
		go func() {
			<-ch
			time.Sleep(50 * time.Millisecond)
			p.Offer(newSampleBlockingItem(1))
		}()
		ch <- struct{}{}
		i := p.PollWithTimeout(1 * time.Second)
		if i.(*sampleBlockingItem).Value != 1 {
			t.Fatal("shoudl have value 1")
		}
	}
}

func TestPriorityBlockingQueue_RemainingCapacity(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	if p.RemainingCapacity() <= 0 {
		t.Fatal("remaining capacity should greater than 0")
	}
	p = NewPriorityBlockingQueue()
	if p.RemainingCapacity() <= 0 {
		t.Fatal("remaining capacity should greater than 0")
	}
}

func TestPriorityBlockingQueue_DrainTo(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	d := NewPriorityBlockingQueue()
	if p.DrainTo(d) != 7 {
		t.Fatal("drain result should be 7")
	}
	if d.Size() != 7 {
		t.Fatal("dest size should be 7")
	}
	if !d.Contains(newSampleBlockingItem(6)) {
		t.Fatal("dest should contain value 7")
	}
	if !d.Contains(newSampleBlockingItem(33)) {
		t.Fatal("dest should contain value 33")
	}
	if !d.Contains(newSampleBlockingItem(3)) {
		t.Fatal("dest should contain value 3")
	}
	if d.Contains(newSampleBlockingItem(0)) {
		t.Fatal("dest should not contain value 0")
	}
}

func TestPriorityBlockingQueue_DrainToWithLimit(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	d := NewPriorityBlockingQueue()
	if p.DrainToWithLimit(d, 2) != 2 {
		t.Fatal("drain result should be 2")
	}
	if d.Size() != 2 {
		t.Fatal("dest size should be 2")
	}
}

func TestPriorityBlockingQueue_Iterator(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	cnt := 0
	matched := false
	iter := p.Iterator()
	if iter == nil {
		t.Fatal("iter should not be nil")
	}
	for iter.HasNext() {
		cnt++
		if iter.Next().(*sampleBlockingItem).Value == 6 {
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

func TestPriorityBlockingQueueIter_ForEachRemaining(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	iter := p.Iterator()

	l := list.New()
	iter.ForEachRemaining(func(i interface{}) {
		l.PushBack(i)
	})
	if l.Len() != 7 {
		t.Fatal("total number of foreach items are not 7")
	}
}

func TestPriorityBlockingQueueIter_Remove(t *testing.T) {
	p := newPreparedPriorityBlockingQueue()
	iter := p.Iterator()
	iter.HasNext()
	iter.Remove()
	iterImpl := iter.(*priorityBlockingQueueIter)
	if len(iterImpl.data) != 7 {
		t.Fatal("iter size should be 7")
	}
	if p.Size() != 6 {
		t.Fatal("iter size should be 6")
	}
}
