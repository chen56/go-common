package reflectx

import (
	"reflect"
)

// IsNil checks if a specified object is nil or not, without Failing.
func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}
func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}

//IsPrimitive 判断是否是原始简单类型
func IsPrimitive(kind reflect.Kind) bool {
	// This switch parallels valueSortLess, except for the default case.
	switch kind {
	case reflect.Bool:
		return true
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return true
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	}
	return false
}

func IsZero(object interface{}) bool {

	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return IsZero(deref)
	// for all other types, compare against the zero value
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func VisitFields(x interface{}, vistor func(field reflect.Type) error) error {
	return nil
	//return visitFields(reflect.TypeOf(x), vistor)
}

//func visitFields(t reflect.Type, vistor func(value reflect.StructField) error) error {
//	//v := reflect.ValueOf(x)
//	fmt.Println(t)
//	//if t.Kind() != reflect.Struct {
//	//	return nil
//	//}
//	for i := 0; i < t.NumField(); i++ {
//
//		field := t.Field(i)
//		err := vistor(field)
//		if err != nil {
//			return err
//		}
//		//err = visitFields(field, vistor)
//		//if err != nil {
//		//	return err
//
//		//}
//	}
//	return nil
//}
