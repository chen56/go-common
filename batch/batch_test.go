package batch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"time"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	var batch = NewBuilder().Limit(2).Callback(func(items []interface{}, b *Batch) {}).Build()
	batch.TickOk("a")
	batch.TickOk("b")
	batch.TickError()
	batch.TickIgnore()
	for i, v := range batch.items {
		fmt.Printf("TestBatch %v %v \n", i, v)
	}
	assert.Equal(4, batch.Total())
	assert.Equal(2, batch.TotalOk())
	assert.Equal(1, batch.TotalError())
	assert.Equal(1, batch.TotalIgnore())
}

func TestIsReachLimit(t *testing.T) {
	assert := assert.New(t)
	var history string
	var batch = NewBuilder().Limit(2).Callback(func(items []interface{}, b *Batch) {
		history = history + fmt.Sprintf("%v,", b.items)
	}).Build()
	batch.TickOk("a")
	assert.Equal("", history)

	batch.TickOk("b")
	assert.Equal("[a b],", history)

	batch.TickError()
	assert.Equal("[a b],", history)

	batch.TickIgnore()
	assert.Equal("[a b],", history)

	batch.TickOk("c")
	assert.Equal("[a b],", history)
}

func TestEmpty(t *testing.T) {
	assert := assert.New(t)
	var last time.Time
	assert.NotNil(last)
	var batch = NewBuilder().Limit(2).Callback(func(items []interface{}, b *Batch) {}).Build()
	for i, v := range batch.items {
		fmt.Printf("TestEmpty %v %v \n", i, v)
	}

}
func TestItems(t *testing.T) {
	assert := assert.New(t)
	var history string
	var batch = NewBuilder().Limit(2).Callback(func(items []interface{}, b *Batch) {
		history = history + fmt.Sprintf("%v,", b.items)
	}).Build()
	batch.TickOk("a")
	assert.Equal("", history)
	assert.Equal([]interface{}{"a"}, batch.items)

	batch.TickOk("b")
	assert.Equal("[a b],", history)

	batch.TickOk("c")
	assert.Equal("[a b],", history)
}

func TestLimitTime(t *testing.T) {
	assert := assert.New(t)
	var history string
	var batch = NewBuilder().Limit(10).
		LimitDuration(10 * time.Millisecond).
		Callback(func(items []interface{}, b *Batch) {
		history = history + fmt.Sprintf("%v,", b.items)
	}).Build()

	batch.TickOk("a")
	time.Sleep(11 * time.Millisecond)

	batch.TickOk("b")
	assert.Equal("[a b],", history)
}
