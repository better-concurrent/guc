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
