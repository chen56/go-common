package slicex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var slice1_5 = []int{1, 2, 3, 4, 5}

//go test -timeout 30s  -run TestInterfaceToSlice
func TestSplitInt(t *testing.T) {
	assert := require.New(t)
	assert.Equal([][]int{[]int{1, 2}, []int{3, 4}, []int{5}}, SplitInt(slice1_5, 2))
	assert.Equal([][]int{[]int{1, 2, 3}, []int{4, 5}}, SplitInt(slice1_5, 3))
}

//看看如何使用分页结果
func TestSplitInt_for(t *testing.T) {
	assert := require.New(t)

	var x []int
	for _, page := range SplitInt(slice1_5, 5000) {
		x = append(x, page...)
	}
	assert.Equal(slice1_5, x)
}

func TestSplitIntError(t *testing.T) {
	assert := require.New(t)
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(r, "limit should ge 0")
		} else {
			assert.Fail("should panic , but not")
		}
	}()
	SplitInt(slice1_5, 0)
}
