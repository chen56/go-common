package stringsx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func ExampleReverse() {
	fmt.Println(Reverse("Hello, world"))
	fmt.Println(Reverse("Hello, 世界"))
	// Unordered Output:
	// 界世 ,olleH
	// dlrow ,olleH
}
func TestStartWith(t *testing.T) {
	require.True(t, StartWith("/chen56.admin4.XXXX", "/chen56.admin4"))
	require.True(t, StartWith("/chen56.admin4.XXXX", "/chen56"))

	require.False(t, StartWith("/chen56.admin4.XXXX", "/chen56.inner4"))

}
