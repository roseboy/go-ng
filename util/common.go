package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

// SHA256Hex sha256
func SHA256Hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

// HMacSHA256 hmacsha256
func HMacSHA256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}
