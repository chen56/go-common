package slicex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//go test -timeout 30s  -run TestInterfaceToSlice
func TestSplitInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([][]int{[]int{1, 2}, []int{3, 4}, []int{5}}, SplitInt([]int{1, 2, 3, 4, 5}, 2))
	assert.Equal([][]int{[]int{1, 2, 3}, []int{4, 5}}, SplitInt([]int{1, 2, 3, 4, 5}, 3))
}

func TestSplitIntError(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(r, "limit should ge 0")
		} else {
			assert.Fail("should panic , but not")
		}
	}()
	SplitInt([]int{1, 2, 3, 4, 5}, 0)
}
