package common

import (
	"fmt"
	"reflect"
)

// a good blog for golang-reflect: https://blog.csdn.net/chen1415886044/article/details/104929244
// To merge flags with config-file, flag_opt、file_opt and default_opt should be ptr.
func Merge(flag_opt interface{}, config_opt interface{}, default_opt interface{}, fieldName string) {
	flag_value := reflect.ValueOf(flag_opt).Elem() // .Elem is useful if h_opt is a ptr

	// check if value is able to be set
	if !flag_value.FieldByName(fieldName).CanSet() {
		fmt.Println("bad args when use common.args, structs should be ptr")
		return
	}

	v := flag_value.FieldByName(fieldName)
	if IsValueZero(&v) {
		// flag => config
		if !IsZeroByFieldName(config_opt, fieldName) {
			config_value := reflect.ValueOf(config_opt).Elem()
			v.Set(config_value.FieldByName(fieldName))
		} else {
			// flag => default value
			default_value := reflect.ValueOf(default_opt).Elem()
			v.Set(default_value.FieldByName(fieldName))
		}
	}
}
