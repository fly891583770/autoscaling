package utils

import (
	"reflect"
)

func VerifyType(t reflect.Type) reflect.Type {

	switch t.Kind() {
	//指针类型
	case reflect.Ptr:
		return VerifyType(t.Elem())
	default:
		return t
	}
}
