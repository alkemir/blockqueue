package blockqueue

import (
	"sync"
	"time"
)

// OrderedBlock holds a block contents and its position on the file for queuing writer
type OrderedBlock struct {
	Position int
	Hash     string
	FileID   uint64
	Content  []byte
}

// A BlocksQueue is a min-heap of orderedBlocks.
type BlocksQueue struct {
	queue  []*OrderedBlock
	lock   *sync.Mutex
	waiter *sync.Cond
}

// New returns a pointer to an empty BlocksQueue ready to be used.
func New() *BlocksQueue {
	l := &sync.Mutex{}
	ret := &BlocksQueue{queue: []*OrderedBlock{&OrderedBlock{Position: -1}}, lock: l, waiter: sync.NewCond(l)}
	go ret.watchdog()
	return ret
}

// Len returns the number of elements in the queue
func (q BlocksQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.queue) - 1
}

// Push add an element to the queue
func (q *BlocksQueue) Push(b *OrderedBlock) {
	q.lock.Lock()
	defer q.lock.Unlock()

	idx := len(q.queue)
	parent := idx / 2
	q.queue = append(q.queue, b)

	for q.queue[idx].Position < q.queue[parent].Position {
		q.queue[idx], q.queue[parent] = q.queue[parent], q.queue[idx]
		idx = parent
		parent = idx / 2
	}

	q.waiter.Signal()
}

// Pop removes the lowest priority element from the queue and returns it
func (q *BlocksQueue) Pop() *OrderedBlock {
	q.lock.Lock()
	defer q.lock.Unlock()

	ret := q.queue[1]
	n := len(q.queue)
	q.queue[1] = q.queue[n-1]
	q.queue = q.queue[:n-1] // this keeps the reference to the original slice, but it's ok for my use case
	// I don't need to nil the left out element, as I will be checking for len() != 0

	parent := 1

	for {
		lChild := parent * 2
		if lChild > n-2 || lChild < 0 {
			break
		}

		smallestChild := lChild

		rChild := lChild + 1
		if rChild <= n-2 && q.queue[rChild].Position < q.queue[lChild].Position {
			smallestChild = rChild
		}

		if q.queue[parent].Position < q.queue[smallestChild].Position {
			break
		}

		q.queue[parent], q.queue[smallestChild] = q.queue[smallestChild], q.queue[parent]
		parent = smallestChild
	}

	return ret
}

// Peek returns the lowest priority element from the queue
func (q BlocksQueue) Peek() *OrderedBlock {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue[1]
}

// Wait until a new block is pushed
func (q BlocksQueue) Wait() {
	q.lock.Lock()
	q.waiter.Wait()
	q.lock.Unlock()
}

func (q BlocksQueue) watchdog() {
	for {
		time.Sleep(2)
		q.waiter.Signal()
	}
}
