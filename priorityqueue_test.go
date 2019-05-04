package guc

import "testing"

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
	dst, ok := i.(sampleItem)
	if !ok {
		return false
	}
	return this.Value == dst.Value
}

func (this *sampleItem) CompareTo(i interface{}) int {
	return this.Value - i.(sampleItem).Value
}

func newSampleItem(value int) *sampleItem {
	return &sampleItem{
		Value: value,
	}
}

func TestPriorityQueue_Add(t *testing.T) {
	p := NewPriority()
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
}
