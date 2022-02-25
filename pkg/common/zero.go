package common

import "reflect"

func IsZero(x interface{}) bool {
	if x == nil {
		return true
	}
	value := reflect.ValueOf(x)
	return IsValueZero(&value)
}

// s is one struct ptr
func IsZeroByFieldName(s interface{}, fieldName string) bool {
	if s == nil {
		return true
	}

	value_struct := reflect.ValueOf(s).Elem() // .Elem is useful if h_opt is a ptr
	value := value_struct.FieldByName(fieldName)
	return IsValueZero(&value)
}

func IsValueZero(v *reflect.Value) bool {
	value := *v
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
