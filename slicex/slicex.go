package slicex

import (
	reflect "reflect"

	"github.com/chen56/go-common/must"
)

func ToSlice(slice interface{}) []interface{} {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		panic("param slice should be a slice")
	}

	length := value.Len()
	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		result[i] = value.Index(i).Interface()
	}
	return result
}

func RemoveIndexInt64(slice []int64, s int) []int64 {
	return append(slice[:s], slice[s+1:]...)
}

func SplitInt(slice []int, limit int) [][]int {
	var result [][]int
	must.Truef(limit > 0, "limit should ge 0")

	for i := 0; i < len(slice); i += limit {
		end := i + limit

		if end > len(slice) {
			end = len(slice)
		}

		result = append(result, slice[i:end])
	}
	return result
}
func SplitInt64(slice []int64, limit int) [][]int64 {
	var result [][]int64
	must.Truef(limit > 0, "limit should ge 0")

	for i := 0; i < len(slice); i += limit {
		end := i + limit

		if end > len(slice) {
			end = len(slice)
		}

		result = append(result, slice[i:end])
	}
	return result
}
