package kast

import "reflect"

func typeElem(typ reflect.Type, n int) (reflect.Type, int) {
	if typ.Kind() == reflect.Ptr {
		return typeElem(typ.Elem(), n+1)
	}
	return typ, n
}

func valueElem(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
