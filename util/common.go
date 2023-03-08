package util

import (
	"reflect"
)

// If 三目
func If[T any](ok bool, v1, v2 T) T {
	if ok {
		return v1
	} else {
		return v2
	}
}

// In if in
func In[T comparable](arr []T, v T) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

// NewInstanceByType new
func NewInstanceByType(typ reflect.Type) interface{} {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		dst := reflect.New(typ).Elem()
		return dst.Addr().Interface()
	} else {
		dst := reflect.New(typ).Elem()
		return dst.Interface()
	}
}
