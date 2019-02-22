package stringsx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCut(t *testing.T) {
	assert := require.New(t)

	data := []struct {
		str    string
		minLen int
		want   string
	}{
		{str: "12345", minLen: 2, want: "12"},
		{str: "12345", minLen: 0, want: ""},
		{str: "12345", minLen: 5, want: "12345"},
		{str: "12345", minLen: 6, want: "12345"},
		{str: "一二三四五", minLen: 2, want: "一二"},
		{str: "一二三四五", minLen: 0, want: ""},
		{str: "一二三四五", minLen: 5, want: "一二三四五"},
		{str: "一二三四五", minLen: 60, want: "一二三四五"},
		{str: "", minLen: 10, want: ""},
	}
	for i, x := range data {
		assert.Equal(x.want, CutRune(x.str, x.minLen), "cut:%d", i)
	}

}
