package jsync

import (
	"sync/atomic"
	"unsafe"
)

type Queue interface {
	Enqueue(value interface{})
	Dequeue() (value interface{}, ok bool)
}

type _node struct {
	value interface{}
	next  unsafe.Pointer // _node 포인터를 저장하기 위해 unsafe.Pointer 사용
}

type _queue struct {
	head unsafe.Pointer // _node 포인터를 저장하기 위해 unsafe.Pointer 사용
	tail unsafe.Pointer // _node 포인터를 저장하기 위해 unsafe.Pointer 사용
}

// NewQueue creates a new Queue
func NewQueue() *_queue {
	dummy := &_node{} // 더미 노드 생성
	return &_queue{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

// Enqueue adds an element to the end of the queue
func (q *_queue) Enqueue(value interface{}) {
	newNode := &_node{value: value}
	for {
		tail := (*_node)(atomic.LoadPointer(&q.tail))
		next := (*_node)(atomic.LoadPointer(&tail.next))
		if tail == (*_node)(atomic.LoadPointer(&q.tail)) { // tail이 여전히 유효한지 확인
			if next == nil { // 마지막 노드인지 확인
				if atomic.CompareAndSwapPointer(&tail.next, nil, unsafe.Pointer(newNode)) {
					//atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(newNode))
					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(newNode))
					return
				}
			}
			//  else { // 이미 다른 노드가 추가된 경우
			// 	atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), tail.next)
			// }
		}
	}
}

// Dequeue removes and returns the element at the front of the queue
func (q *_queue) Dequeue() (interface{}, bool) {
	for {
		head := (*_node)(atomic.LoadPointer(&q.head))
		tail := (*_node)(atomic.LoadPointer(&q.tail))
		next := (*_node)(atomic.LoadPointer(&head.next))
		if head == (*_node)(atomic.LoadPointer(&q.head)) { // head가 여전히 유효한지 확인
			if head == tail { // 큐가 비어 있는지 확인
				if next == nil {
					return nil, false
				}
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			} else {
				if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
					return next.value, true
				}
			}
		}
	}
}

////////////////////////////////////////////////////////////////////
