package reflectx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go test -timeout 30s  -run TestInterfaceToSlice
func TestInterfaceToSlice(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		strSlice []string
		expected string
	}{
		{strSlice: []string{"a", "b", "c"}, expected: "[a b c]"},
		{strSlice: []string{}, expected: "[]"},
	}

	for i, x := range data {
		result := fmt.Sprintf("%v", InterfaceToSlice(x.strSlice))
		assert.Equal(x.expected, result, "cut:%d", i)
	}

	fmt.Printf("xxxxx:%v \n", InterfaceToSlice([]string{"a", "b", "c"}))
}
