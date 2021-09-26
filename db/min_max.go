package db

import (
	"reflect"
	"time"
)

func Min(x interface{}, f func(int) time.Time) time.Time {
	rv := reflect.ValueOf(x)
	length := rv.Len()
	v := time.Time{}
	if length > 0 {
		v = f(0)
	}
	for i := 1; i < length; i++ {
		tt := f(i)
		if tt.Before(v) {
			v = tt
		}
	}
	return v
}
func Max(x interface{}, f func(int) time.Time) time.Time {
	rv := reflect.ValueOf(x)
	length := rv.Len()
	v := time.Time{}
	if length > 0 {
		v = f(0)
	}
	for i := 1; i < length; i++ {
		tt := f(i)
		if tt.After(v) {
			v = tt
		}
	}
	return v
}
