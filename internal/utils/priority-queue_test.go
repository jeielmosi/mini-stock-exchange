package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertHeap[T any](t *testing.T, pq *PriorityQueue[T]) {
	for i := 0; i < len(pq.heap); i++ {
		assert.True(t, 0 <= pq.mn[i])
		assert.True(t, pq.mn[i] < len(pq.heap))

		left := 2*i + 1
		right := 2*i + 2

		if left < len(pq.heap) {
			assert.False(t, pq.cmp(pq.heap[left], pq.heap[i]))
		}
		if right < len(pq.heap) {
			assert.False(t, pq.cmp(pq.heap[right], pq.heap[i]))
		}

		c := pq.mn[i]

		if left < len(pq.heap) {
			l := pq.mn[left]
			assert.False(t, pq.cmp(pq.heap[c], pq.heap[l]))
		}
		if right < len(pq.heap) {
			r := pq.mn[right]
			assert.False(t, pq.cmp(pq.heap[c], pq.heap[r]))
		}
	}
}

func TestNewPriorityQueue_InvalidCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "capacity must be greater than 0", r)
		}
	}()

	NewPriorityQueue(func(a, b int) bool { return a > b }, 0)
	t.Fatal("should have panicked")
}

func TestNewPriorityQueue_NegativeCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "capacity must be greater than 0", r)
		}
	}()

	NewPriorityQueue(func(a, b int) bool { return a > b }, -5)
	t.Fatal("should have panicked")
}

func TestPriorityQueue_Peek_Empty(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	item, ok := pq.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, item)
}

func TestPriorityQueue_PushPeekBasic(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 5, item)
}

func TestPriorityQueue_PushMultiplePeek(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(10)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)

	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_DropSingle(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Drop()
	assertHeap(t, pq)

	item, ok := pq.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, item)
}

func TestPriorityQueue_DropMultiple(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(10)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)

	pq.Drop()
	assertHeap(t, pq)
	item, _ := pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 3, item)

	pq.Drop()
	assertHeap(t, pq)
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_DropEmpty(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Drop()
	assertHeap(t, pq)
	pq.Drop()
	assertHeap(t, pq)

	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_MinHeap(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a < b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(10)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)
	pq.Push(7)
	assertHeap(t, pq)

	item, _ := pq.Peek()
	assert.Equal(t, 3, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 7, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_PushBeyondCapacity(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(10)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)
}

type tieItem struct {
	priority int
	id       string
}

func TestPriorityQueue_TieElements(t *testing.T) {

	cmp := func(a, b tieItem) bool {
		return a.priority > b.priority
	}
	pq := NewPriorityQueue(cmp, 10)

	// Elements with ties in priority
	items := []tieItem{
		{priority: 10, id: "A"},
		{priority: 5, id: "D"},
		{priority: 2, id: "E"},
		{priority: 5, id: "C"},
		{priority: 10, id: "B"},
	}

	for _, it := range items {
		pq.Push(it)
		assertHeap(t, pq)
	}

	// Initial peek
	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 10, item.priority)
	pq.Drop()
	assertHeap(t, pq)

	item, ok = pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 10, item.priority)
	pq.Drop()
	assertHeap(t, pq)

	item, ok = pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 5, item.priority)
	pq.Drop()
	assertHeap(t, pq)

	item, ok = pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 5, item.priority)
	pq.Drop()
	assertHeap(t, pq)

	item, ok = pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 2, item.priority)
	pq.Drop()
	assertHeap(t, pq)

	_, ok = pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_TieElementsCapacity(t *testing.T) {
	cmp := func(a, b tieItem) bool {
		return a.priority > b.priority
	}
	// Capacity 3
	pq := NewPriorityQueue(cmp, 3)

	// Fill to capacity with some ties
	items := []tieItem{
		{priority: 10, id: "A"},
		{priority: 10, id: "B"},
		{priority: 10, id: "C"},
		{priority: 10, id: "D"},
		{priority: 10, id: "E"},
		{priority: 7, id: "F"},
		{priority: 5, id: "G"},
		{priority: 3, id: "H"},
	}

	for _, it := range items {
		pq.Push(it)
		assertHeap(t, pq)
	}

	for 0 < len(pq.heap) {
		item, ok := pq.Peek()
		assert.True(t, ok)
		assert.True(t, item.id == items[0].id || item.id == items[1].id || item.id == items[2].id)
		pq.Drop()
		assertHeap(t, pq)
	}
}

func TestPriorityQueue_PushBeyondCapacity_LowerPriority(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	pq.Push(10)
	assertHeap(t, pq)
	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)

	pq.Push(2)
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_ComplexOrdering(t *testing.T) {
	cmp := func(a, b int) bool { return a > b }
	pq := NewPriorityQueue(cmp, 10)

	values := []int{5, 10, 3, 8, 1, 9, 2, 7, 4, 6}
	for _, v := range values {
		pq.Push(v)
		assertHeap(t, pq)
	}
	last, ok := pq.Peek()
	assert.True(t, ok)
	pq.Drop()
	assertHeap(t, pq)

	for true {
		item, ok := pq.Peek()
		if !ok {
			break
		}
		pq.Drop()
		assertHeap(t, pq)
		assert.True(t, cmp(last, item))
		last = item
	}
}

func TestPriorityQueue_SingleElement(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(42)
	assertHeap(t, pq)
	item, _ := pq.Peek()
	assert.Equal(t, 42, item)

	pq.Drop()
	assertHeap(t, pq)
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_Cap(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)
	assert.Equal(t, 10, pq.Cap())

	pq.Push(1)
	assertHeap(t, pq)
	assert.Equal(t, 9, pq.Cap())

	pq.Push(2)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)
	assert.Equal(t, 7, pq.Cap())
}

func TestPriorityQueue_DropLastElement(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(10)
	assertHeap(t, pq)
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	assertHeap(t, pq)
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_StringType(t *testing.T) {
	pq := NewPriorityQueue(func(a, b string) bool { return len(a) > len(b) }, 10)

	pq.Push("hello")
	pq.Push("hi")
	pq.Push("word")

	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, "hello", item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, "word", item)
}

func TestPriorityQueue_AllSamePriority(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(5)
	assertHeap(t, pq)
	pq.Push(5)
	assertHeap(t, pq)

	item, _ := pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	assertHeap(t, pq)
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	assertHeap(t, pq)
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_AscendingOrderPush(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 1; i <= 5; i++ {
		pq.Push(i)
		assertHeap(t, pq)
	}

	for i := 5; 0 < i; i-- {
		item, ok := pq.Peek()
		assert.True(t, ok)
		assert.Equal(t, i, item)
		pq.Drop()
		assertHeap(t, pq)
	}
}

func TestPriorityQueue_DescendingOrderPush(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 5; i >= 1; i-- {
		pq.Push(i)
		assertHeap(t, pq)
	}

	for i := 5; i >= 1; i-- {
		item, ok := pq.Peek()
		assert.True(t, ok)
		assert.Equal(t, i, item)
		pq.Drop()
		assertHeap(t, pq)
	}
}

func TestPriorityQueue_AlternatingPush(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 1; i <= 5; i++ {
		pq.Push(i)
		assertHeap(t, pq)
		pq.Push(i + 5)
		assertHeap(t, pq)
	}

	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_CapacityExactlyFilled(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	pq.Push(1)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)
	pq.Push(2)
	assertHeap(t, pq)
	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 3, item)
}

func TestPriorityQueue_PushAfterPartialDrain(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 1; i <= 5; i++ {
		pq.Push(i * 10)
		assertHeap(t, pq)
	}

	pq.Drop()
	assertHeap(t, pq)
	pq.Drop()
	assertHeap(t, pq)

	pq.Push(100)
	assertHeap(t, pq)
	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 100, item)
}

// Test to verify heap and mn array validity
func TestPriorityQueue_HeapAndMnValidity(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	// Verify heap property: parent should have higher priority than children
	// For equal values, neither has higher priority, which is acceptable

	// Test with various values
	values := []int{3, 1, 4, 7, 5, 9, 2, 6, 8, 10, 0, 11}
	for _, v := range values {
		pq.Push(v)
		assertHeap(t, pq)
	}

	// Test after drops
	for 0 < len(pq.heap) {
		pq.Drop()
		assertHeap(t, pq)
	}
}

// Test to verify heap and mn array validity when pushing beyond capacity
func TestPriorityQueue_HeapAndMnValidity_BeyondCapacity(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	// Fill to capacity
	pq.Push(1)
	assertHeap(t, pq)
	pq.Push(2)
	assertHeap(t, pq)
	pq.Push(3)
	assertHeap(t, pq)

	// Verify initial state
	assert.Len(t, pq.heap, 3)
	assert.Len(t, pq.mn, 3)

	// Push beyond capacity with higher priority item
	pq.Push(10) // Should replace the lowest priority item (1)
	assertHeap(t, pq)

	// The highest priority item should be at the root
	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 10, item)
}
