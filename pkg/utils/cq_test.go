package utils_test

import (
	"testing"

	"github.com/konrads/go-rate-limiter/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPushAndPop(t *testing.T) {
	cq := utils.NewCircularQueue(3)
	// test pushes
	assert.True(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 0, cq.Size())
	assert.True(t, cq.Push(1))
	assert.False(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 1, cq.Size())
	assert.True(t, cq.Push(2))
	assert.False(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 2, cq.Size())
	assert.True(t, cq.Push(3))
	assert.False(t, cq.IsEmpty())
	assert.True(t, cq.IsFull())
	assert.False(t, cq.Push(4))
	assert.False(t, cq.IsEmpty())
	assert.True(t, cq.IsFull())
	// test pops
	assert.Equal(t, 3, *cq.Pop())
	assert.False(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 2, *cq.Pop())
	assert.False(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 1, *cq.Pop())
	assert.True(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Nil(t, cq.Pop())
	assert.True(t, cq.IsEmpty())
}

func TestPushAndDequeue(t *testing.T) {
	cq := utils.NewCircularQueue(3)
	assert.True(t, cq.Push(1))
	assert.True(t, cq.Push(2))
	assert.True(t, cq.Push(3))
	// test dequeues
	assert.Equal(t, 1, *cq.Dequeue())
	assert.False(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 2, *cq.Dequeue())
	assert.False(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Equal(t, 3, *cq.Dequeue())
	assert.True(t, cq.IsEmpty())
	assert.False(t, cq.IsFull())
	assert.Nil(t, cq.Dequeue())
	assert.True(t, cq.IsEmpty())
}

func TestDropWhile2(t *testing.T) {
	cq := utils.NewCircularQueue(3)
	assert.True(t, cq.Push(1))
	assert.True(t, cq.Push(2))
	assert.True(t, cq.Push(3))
	assert.Equal(t, []interface{}{1, 2, 3}, cq.Elements())
	cq.DropWhile(func(x interface{}) bool { return x.(int) <= 0 })
	assert.Equal(t, []interface{}{1, 2, 3}, cq.Elements())
	cq.DropWhile(func(x interface{}) bool { return x.(int) <= 1 })
	assert.Equal(t, []interface{}{2, 3}, cq.Elements())
}

func TestRandom(t *testing.T) {
	cq := utils.NewCircularQueue(3)
	assert.True(t, cq.Push(1))
	assert.True(t, cq.Push(2))
	assert.True(t, cq.Push(3))
	assert.Equal(t, []interface{}{1, 2, 3}, cq.Elements())
	cq.DropWhile(func(x interface{}) bool { return x.(int) <= 1 })
	assert.Equal(t, []interface{}{2, 3}, cq.Elements())
	assert.True(t, cq.Push(4))
	cq.DropWhile(func(x interface{}) bool { return x.(int) <= 2 })
	assert.Equal(t, []interface{}{3, 4}, cq.Elements())
	assert.True(t, cq.Push(5))
	cq.DropWhile(func(x interface{}) bool { return x.(int) <= 3 })
	assert.Equal(t, []interface{}{4, 5}, cq.Elements())
}
