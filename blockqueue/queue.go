package blockqueue

// OrderedBlock holds a block contents and its position on the file for queuing writer
type OrderedBlock struct {
	Position int
	Hash     string
	Content  []byte
}

// A BlocksQueue is a min-heap of orderedBlocks.
type BlocksQueue []*OrderedBlock

// New returns a pointer to an empty BlocksQueue ready to be used.
func New() *BlocksQueue {
	return &BlocksQueue{&OrderedBlock{Position: -1}}
}

// Len returns the number of elements in the queue
func (q BlocksQueue) Len() int {
	return len(q) - 1
}

// Push add an element to the queue
func (q *BlocksQueue) Push(b *OrderedBlock) {
	idx := len(*q)
	parent := idx / 2
	*q = append(*q, b)

	for (*q)[idx].Position < (*q)[parent].Position {
		(*q)[idx], (*q)[parent] = (*q)[parent], (*q)[idx]
		idx = parent
		parent = idx / 2
	}
}

// Pop removes the lowest priority element from the queue and returns it
func (q *BlocksQueue) Pop() *OrderedBlock {
	ret := (*q)[1]
	n := len(*q)
	(*q)[1] = (*q)[n-1]
	*q = (*q)[:n-1] // this keeps the reference to the original slice, but it's ok for my use case
	// I don't need to nil the left out element, as I will be checking for len() != 0

	parent := 1

	for {
		lChild := parent * 2
		if lChild > n-2 || lChild < 0 {
			break
		}

		smallestChild := lChild

		rChild := lChild + 1
		if rChild < n-2 && (*q)[rChild].Position < (*q)[lChild].Position {
			smallestChild = rChild
		}

		if (*q)[parent].Position < (*q)[smallestChild].Position {
			break
		}

		(*q)[parent], (*q)[smallestChild] = (*q)[smallestChild], (*q)[parent]
		parent = smallestChild
	}

	return ret
}

// Peek returns the lowest priority element from the queue
func (q BlocksQueue) Peek() *OrderedBlock {
	return q[1]
}
