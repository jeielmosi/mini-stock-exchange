package utils

type PriorityQueue[T any] struct {
	heap []T
	mn   []int             //keep the index of the lower priority leaf of sub-tree
	cmp  func(a, b T) bool //returns if a have higher priority than b
}

func NewPriorityQueue[T any](cmp func(a, b T) bool, capacity int) *PriorityQueue[T] {
	if capacity < 1 {
		panic("capacity must be greater than 0")
	}
	return &PriorityQueue[T]{
		heap: make([]T, 0, capacity),
		mn:   make([]int, 0, capacity),
		cmp:  cmp,
	}
}

func (pq *PriorityQueue[T]) Cap() int {
	return cap(pq.heap) - len(pq.heap)
}

func (pq *PriorityQueue[T]) Push(item T) {
	if len(pq.heap) == cap(pq.heap) {
		mn := pq.mn[0]
		if pq.cmp(item, pq.heap[mn]) {
			pq.heap[mn] = item
			pq.toRoot(mn)
		}
		return
	}
	curr := len(pq.heap)
	pq.heap = append(pq.heap, item)
	pq.mn = append(pq.mn, curr)
	pq.toRoot(curr)
}

func (pq *PriorityQueue[T]) Drop() {
	if len(pq.heap) == 0 {
		return
	}
	rm := len(pq.heap) - 1
	item := pq.heap[rm]
	pq.heap = pq.heap[:rm]
	pq.mn = pq.mn[:rm]
	if rm == 0 {
		return
	}

	curr := rm
	//update all the way to the root
	for 0 < curr {
		parent := (curr - 1) / 2
		if pq.mn[parent] != rm {
			break
		}
		pq.updateMn(parent)
		curr = parent
	}

	//insert the last element on root
	pq.heap[0] = item
	pq.toLeaf(0)
}

func (pq *PriorityQueue[T]) Peek() (T, bool) {
	if len(pq.heap) == 0 {
		var empty T
		return empty, false
	}
	return pq.heap[0], true
}

func (pq *PriorityQueue[T]) updateMn(root int) {
	size := len(pq.heap)
	left := root*2 + 1
	right := root*2 + 2

	pq.mn[root] = root
	if size <= left {
		return
	}
	left = pq.mn[left]

	pq.mn[root] = left
	if size <= right {
		return
	}
	right = pq.mn[right]

	if pq.cmp(pq.heap[left], pq.heap[right]) {
		pq.mn[root] = right
	}
}

// update the path from leaf to root
func (pq *PriorityQueue[T]) toRoot(child int) {
	for 0 < child {
		parent := (child - 1) / 2
		if pq.cmp(pq.heap[parent], pq.heap[child]) {
			break
		}
		pq.heap[parent], pq.heap[child] = pq.heap[child], pq.heap[parent]
		pq.updateMn(parent)
		child = parent
	}
	for 0 < child {
		parent := (child - 1) / 2
		pq.updateMn(parent)
		child = parent
	}
}

// update the leaf from curr until the leaf
func (pq *PriorityQueue[T]) toLeaf(root int) {
	size := len(pq.heap)

	defer pq.updateMn(root)

	left := root*2 + 1
	if size <= left {
		return
	}

	child := left

	right := root*2 + 2
	if right < size && pq.cmp(pq.heap[right], pq.heap[left]) {
		child = right
	}
	if pq.cmp(pq.heap[child], pq.heap[root]) {
		pq.heap[root], pq.heap[child] = pq.heap[child], pq.heap[root]
	}
	pq.toLeaf(child)
}
