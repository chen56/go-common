package slicex

import reflect "reflect"

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
