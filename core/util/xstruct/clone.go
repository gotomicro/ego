package xstruct

import "reflect"

// CopyStruct ...
func CopyStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		value := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name

		dstValue := dstVal.FieldByName(name)
		if dstValue.IsValid() == false {
			continue
		}
		dstValue.Set(value)
	}
}
