package kast

import "reflect"

func typeElem(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
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
