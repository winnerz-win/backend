package jsync

import (
	"sync/atomic"
	"unsafe"
)

type QueueT[T any] interface {
	Enqueue(value T)
	Dequeue() (value T, ok bool)
}

type nodeT[T any] struct {
	value T
	next  unsafe.Pointer // *nodeT[T]
}

type queueT[T any] struct {
	head unsafe.Pointer // *nodeT[T]
	tail unsafe.Pointer // *nodeT[T]
}

func newNodeT[T any](value T) *nodeT[T] {
	node := &nodeT[T]{value: value}
	atomic.StorePointer(&node.next, nil) // Initialize next pointer to nil
	return node
}

func NewQueueT[T any]() *queueT[T] {
	var zeroValue T
	node := newNodeT(zeroValue) // Create a new node with zero value
	q := &queueT[T]{}
	atomic.StorePointer(&q.head, unsafe.Pointer(node)) // Set head to new node
	atomic.StorePointer(&q.tail, unsafe.Pointer(node)) // Set tail to new node
	return q
}

// Enqueue adds an element to the end of the queue
func (q *queueT[T]) Enqueue(value T) {
	newNode := newNodeT(value)
	for {
		tail := (*nodeT[T])(atomic.LoadPointer(&q.tail))
		next := (*nodeT[T])(atomic.LoadPointer(&tail.next))
		if tail == (*nodeT[T])(atomic.LoadPointer(&q.tail)) { // Ensure tail is still valid
			if next == nil {
				// Try to link the new node to the end of the list
				if atomic.CompareAndSwapPointer(&tail.next, nil, unsafe.Pointer(newNode)) {
					// Successfully linked the new node, now update the tail to the new node directly
					//atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(newNode))
					atomic.StorePointer(&q.tail, unsafe.Pointer(newNode))
					return
				}
			}
			// else {
			// 	// Tail was not at the end, try to move it forward
			// 	atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			// }
		}
	}
}

// Dequeue removes and returns the element at the front of the queue
func (q *queueT[T]) Dequeue() (value T, ok bool) {
	for {
		head := (*nodeT[T])(atomic.LoadPointer(&q.head))
		tail := (*nodeT[T])(atomic.LoadPointer(&q.tail))
		next := (*nodeT[T])(atomic.LoadPointer(&head.next))
		if head == (*nodeT[T])(atomic.LoadPointer(&q.head)) { // Ensure head is still valid
			if head == tail {
				if next == nil {
					var zeroValue T
					return zeroValue, false // Queue is empty
				}
				// Tail is lagging behind, push it forward
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			} else {
				value := next.value
				// Attempt to move head to the next node
				if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
					return value, true
				}
			}
		}
	}
}

// import (
// 	"sync/atomic"
// )

// type QueueT[T any] interface {
// 	Enqueue(value T)
// 	Dequeue() (value T, ok bool)
// }

// type nodeT[T any] struct {
// 	value T
// 	next  atomic.Value // *nodeT[T]
// }

// type queueT[T any] struct {
// 	head atomic.Value // *nodeT[T]
// 	tail atomic.Value // *nodeT[T]
// }

// func newNodeT[T any](value T) *nodeT[T] {
// 	node := &nodeT[T]{value: value}
// 	node.next.Store((*nodeT[T])(nil))
// 	return node
// }

// func NewQueueT[T any]() *queueT[T] {
// 	var zeroValue T
// 	node := newNodeT(zeroValue)
// 	q := &queueT[T]{}
// 	q.head.Store(node)
// 	q.tail.Store(node)
// 	return q
// }

// // Enqueue adds an element to the end of the queue
// func (q *queueT[T]) Enqueue(value T) {
// 	newNode := newNodeT(value)
// 	for {
// 		tail := q.tail.Load().(*nodeT[T])
// 		next := tail.next.Load().(*nodeT[T])
// 		if tail == q.tail.Load().(*nodeT[T]) { // Ensure tail is still valid
// 			if next == nil {
// 				if tail.next.CompareAndSwap(next, newNode) {
// 					q.tail.CompareAndSwap(tail, newNode) // Move tail to the new node
// 					return
// 				}
// 			} else {
// 				q.tail.CompareAndSwap(tail, next) // Tail was not at the end, try to move it forward
// 			}
// 		}
// 	}
// }

// // Dequeue removes and returns the element at the front of the queue
// func (q *queueT[T]) Dequeue() (value T, ok bool) {
// 	for {
// 		head := q.head.Load().(*nodeT[T])
// 		tail := q.tail.Load().(*nodeT[T])
// 		next := head.next.Load().(*nodeT[T])
// 		if head == q.head.Load().(*nodeT[T]) { // Ensure head is still valid
// 			if head == tail {
// 				if next == nil {
// 					return *new(T), false // Queue is empty
// 				}
// 				q.tail.CompareAndSwap(tail, next) // Tail is lagging behind, push it forward
// 			} else {
// 				value := next.value
// 				if q.head.CompareAndSwap(head, next) {
// 					return value, true
// 				}
// 			}
// 		}
// 	}
// }
