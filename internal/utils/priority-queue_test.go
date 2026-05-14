package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 5, item)
}

func TestPriorityQueue_PushMultiplePeek(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	pq.Push(10)
	pq.Push(3)

	item, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_DropSingle(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	pq.Drop()

	item, ok := pq.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, item)
}

func TestPriorityQueue_DropMultiple(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	pq.Push(10)
	pq.Push(3)

	pq.Drop()
	item, _ := pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 3, item)

	pq.Drop()
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_DropEmpty(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Drop()
	pq.Drop()

	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_MinHeap(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a < b }, 10)

	pq.Push(5)
	pq.Push(10)
	pq.Push(3)
	pq.Push(7)

	item, _ := pq.Peek()
	assert.Equal(t, 3, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 7, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_PushBeyondCapacity(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	pq.Push(5)
	pq.Push(10)
	pq.Push(3)
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)

	pq.Push(20)
	item, _ = pq.Peek()
	assert.Equal(t, 20, item)

	pq.Push(1)
	item, _ = pq.Peek()
	assert.Equal(t, 20, item)

	pq.Push(15)
	item, _ = pq.Peek()
	assert.Equal(t, 20, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 15, item)
}

func TestPriorityQueue_PushBeyondCapacity_LowerPriority(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	pq.Push(10)
	pq.Push(5)
	pq.Push(3)
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)

	pq.Push(2)
	item, _ = pq.Peek()
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_ComplexOrdering(t *testing.T) {
	cmp := func(a, b int) bool { return a > b }
	pq := NewPriorityQueue(cmp, 10)

	values := []int{5, 10, 3, 8, 1, 9, 2, 7, 4, 6}
	for _, v := range values {
		pq.Push(v)
	}
	last, ok := pq.Peek()
	assert.True(t, ok)
	pq.Drop()

	for true {
		item, ok := pq.Peek()
		if !ok {
			break
		}
		pq.Drop()
		assert.True(t, cmp(last, item))
		last = item
	}
}

func TestPriorityQueue_SingleElement(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(42)
	item, _ := pq.Peek()
	assert.Equal(t, 42, item)

	pq.Drop()
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_Cap(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)
	assert.Equal(t, 10, pq.Cap())

	pq.Push(1)
	assert.Equal(t, 9, pq.Cap())

	pq.Push(2)
	pq.Push(3)
	assert.Equal(t, 7, pq.Cap())
}

func TestPriorityQueue_DropLastElement(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(5)
	pq.Push(10)
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
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
	pq.Push(5)
	pq.Push(5)

	item, _ := pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 5, item)

	pq.Drop()
	_, ok := pq.Peek()
	assert.False(t, ok)
}

func TestPriorityQueue_LargeValues(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(1000000)
	pq.Push(999999)
	pq.Push(1000001)

	item, _ := pq.Peek()
	assert.Equal(t, 1000001, item)
}

func TestPriorityQueue_AscendingOrderPush(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 1; i <= 5; i++ {
		pq.Push(i)
	}

	item, _ := pq.Peek()
	assert.Equal(t, 5, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 4, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 3, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 2, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 1, item)
	pq.Drop()
}

func TestPriorityQueue_DescendingOrderPush(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 5; i >= 1; i-- {
		pq.Push(i)
	}

	item, _ := pq.Peek()
	assert.Equal(t, 5, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 4, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 3, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 2, item)
	pq.Drop()
	item, _ = pq.Peek()
	assert.Equal(t, 1, item)
	pq.Drop()
}

func TestPriorityQueue_AlternatingPush(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	pq.Push(1)
	pq.Push(10)
	pq.Push(2)
	pq.Push(9)
	pq.Push(3)
	pq.Push(8)
	pq.Push(4)
	pq.Push(7)
	pq.Push(5)
	pq.Push(6)

	item, _ := pq.Peek()
	assert.Equal(t, 10, item)
}

func TestPriorityQueue_CapacityExactlyFilled(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 3)

	pq.Push(1)
	pq.Push(2)
	pq.Push(3)
	item, _ := pq.Peek()
	assert.Equal(t, 3, item)
}

func TestPriorityQueue_PushAfterPartialDrain(t *testing.T) {
	pq := NewPriorityQueue(func(a, b int) bool { return a > b }, 10)

	for i := 1; i <= 5; i++ {
		pq.Push(i * 10)
	}

	pq.Drop()
	pq.Drop()

	pq.Push(100)
	item, _ := pq.Peek()
	assert.Equal(t, 100, item)
}

func assertHeap(t *testing.T, pq *PriorityQueue[int]) {
	for i := 0; i < len(pq.heap); i++ {
		assert.True(t, 0 <= pq.mn[i])
		assert.True(t, pq.mn[i] < len(pq.heap))

		left := 2*i + 1
		right := 2*i + 2

		if left < len(pq.heap) {
			assert.True(t, pq.cmp(pq.heap[i], pq.heap[left]))
		}
		if right < len(pq.heap) {
			assert.True(t, pq.cmp(pq.heap[i], pq.heap[right]))
		}

		if left < len(pq.heap) && right < len(pq.heap) {
			l := pq.mn[left]
			r := pq.mn[right]
			if pq.cmp(pq.heap[l], pq.heap[r]) {
				assert.True(t, pq.mn[i] == r)
			} else {
				assert.True(t, pq.mn[i] == l)
			}
			continue
		}
		if left < len(pq.heap) {
			l := pq.mn[left]
			assert.True(t, pq.mn[i] == l)
			continue
		}
		assert.True(t, pq.mn[i] == i)
	}
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
	pq.Push(2)
	pq.Push(3)

	// Verify initial state
	assert.Len(t, pq.heap, 3)
	assert.Len(t, pq.mn, 3)

	// Verify heap property
	for i := 0; i < len(pq.heap); i++ {
		left := 2*i + 1
		right := 2*i + 2

		if left < len(pq.heap) {
			assert.True(t, pq.cmp(pq.heap[i], pq.heap[left]),
				"Heap property violated at index %d", i)
		}
		if right < len(pq.heap) {
			assert.True(t, pq.cmp(pq.heap[i], pq.heap[right]),
				"Heap property violated at index %d", i)
		}
	}

	// Verify mn array validity
	for i := 0; i < len(pq.mn); i++ {
		assert.GreaterOrEqual(t, pq.mn[i], 0, "mn[%d] should be non-negative", i)
		assert.Less(t, pq.mn[i], len(pq.heap), "mn[%d] should be less than heap length", i)
	}

	// Push beyond capacity with higher priority item
	pq.Push(10) // Should replace the lowest priority item (1)

	// Verify heap property still holds
	for i := 0; i < len(pq.heap); i++ {
		left := 2*i + 1
		right := 2*i + 2

		if left < len(pq.heap) {
			assert.True(t, pq.cmp(pq.heap[i], pq.heap[left]),
				"Heap property violated after replacement at index %d", i)
		}
		if right < len(pq.heap) {
			assert.True(t, pq.cmp(pq.heap[i], pq.heap[right]),
				"Heap property violated after replacement at index %d", i)
		}
	}

	// Verify mn array validity
	for i := 0; i < len(pq.mn); i++ {
		assert.GreaterOrEqual(t, pq.mn[i], 0, "mn[%d] should be non-negative after replacement", i)
		assert.Less(t, pq.mn[i], len(pq.heap), "mn[%d] should be less than heap length after replacement", i)
	}

	// The highest priority item should be at the root
	item, _ := pq.Peek()
	assert.Equal(t, 10, item)
}
