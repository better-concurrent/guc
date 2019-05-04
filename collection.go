package guc

import "time"

// this interface imposes basic operations on objects
type Object interface {
	Equals(i interface{}) bool
}

// this interface is used for ordering the objects of each struct
type Comparable interface {
	CompareTo(i interface{}) int
}

// this interface represents the function that compares two objects
type Comparator interface {
	Compare(o1, o2 interface{}) int
}

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Remove()
	ForEachRemaining(consumer func(i interface{}))
}

type Iterable interface {
	Iterator() Iterator
	ForEach(consumer func(i interface{}))
}

type Collection interface {
	Iterable
	Size() int
	IsEmpty() bool
	Contains(i interface{}) bool
	ToArray() []interface{}
	FillArray(arr []interface{}) []interface{}

	Add(i interface{}) bool
	Remove(i interface{}) bool
	ContainsAll(coll Collection) bool
	AddAll(coll Collection) bool
	RemoveAll(coll Collection) bool
	RemoveIf(predicate func(i interface{}) bool) bool
	RetainAll(coll Collection) bool
	Clear()
	Equals(i interface{}) bool
	HashCode() int
}

type Queue interface {
	Collection

	// default inherits
	// Add(i interface{}) bool

	Offer(i interface{}) bool
	// retrieve and remove head
	// panic if empty
	RemoveHead() interface{}
	// retrieve and remove head
	// return nil if empty
	Poll() interface{}
	// retrieve head of the queue
	// panic if empty
	Element() interface{}
	// retrieve head of the queue
	// return nil if empty
	Peek() interface{}
}

type BlockingQueue interface {
	Queue

	// default inherits
	// Add(i interface{}) bool
	// Offer(i interface{}) bool
	// Remove(i interface{}) bool
	// Contains(i interface{}) bool

	Put(i interface{})
	OfferWithTimeout(i interface{}, t time.Duration) bool
	Take() interface{}
	PollWithTimeout(t time.Duration) interface{}
	RemainingCapacity() int
	DrainTo(coll interface{}) int
	DrainToWithLimit(coll interface{}, max int) int
}
