package utils

// based on https://github.com/Focinfi/gcircularqueue/blob/master/gcircularqueue.go
type CircularQueue struct {
	capacity int
	size     int
	elements []interface{}
	i        int
}

func NewCircularQueue(capacity int) *CircularQueue {
	return &CircularQueue{
		capacity: capacity,
		size:     0,
		elements: make([]interface{}, capacity),
		i:        0,
	}
}

func (cq *CircularQueue) IsEmpty() bool {
	return cq.size == 0
}

func (cq *CircularQueue) IsFull() bool {
	return cq.size == cq.capacity
}

// return false if failed to add due to queue full
func (cq *CircularQueue) Push(e interface{}) bool {
	if cq.IsFull() {
		return false
	}
	cq.elements[cq.i] = e
	cq.i = (cq.i + 1) % cq.capacity
	cq.size += 1
	return true
}

// return false if added even though queue full, hence overwriting
func (cq *CircularQueue) PushOverwrite(e interface{}) bool {
	needOverwrite := cq.IsFull()
	cq.elements[cq.i] = e
	cq.i = (cq.i + 1) % cq.capacity
	if !needOverwrite {
		cq.size += 1
	}
	return !needOverwrite
}

func (cq *CircularQueue) Pop() *interface{} {
	if cq.IsEmpty() {
		return nil
	}
	cq.i = (cq.i - 1 + cq.capacity) % cq.capacity
	res := cq.elements[cq.i]
	cq.size -= 1
	return &res
}

func (cq *CircularQueue) Dequeue() *interface{} {
	if cq.IsEmpty() {
		return nil
	}
	first := (cq.i - cq.size + cq.capacity) % cq.capacity
	res := cq.elements[first]
	cq.size -= 1
	return &res
}

func (cq *CircularQueue) Size() int {
	return cq.size
}

func (cq *CircularQueue) Elements() []interface{} {
	if cq.IsEmpty() {
		return []interface{}{}
	}
	end := (cq.i - 1 + cq.capacity) % cq.capacity
	start := (end - cq.size + 1 + cq.capacity) % cq.capacity
	if start < end {
		return cq.elements[start : end+1]
	} else {
		return append(cq.elements[start:], cq.elements[:end+1]...)
	}
}

func (cq *CircularQueue) DropWhile(predicate func(x interface{}) bool) {
	elems := cq.Elements()
	for _, x := range elems {
		if predicate(x) {
			cq.Dequeue()
		} else {
			break
		}
	}
}
