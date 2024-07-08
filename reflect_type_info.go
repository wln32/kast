package kast

import (
	"reflect"
	"time"
)

var (
	timeType      = reflect.TypeOf((*time.Time)(nil)).Elem()
	timePtrType   = reflect.PointerTo(timeType)
	strAnyMapType = reflect.TypeOf((*map[string]any)(nil)).Elem()

	stringType = reflect.TypeOf((*string)(nil)).Elem()
	anyType    = reflect.TypeOf((*any)(nil)).Elem()

	bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
	boolType  = reflect.TypeOf((*bool)(nil)).Elem()
)
