package kast

import (
	"reflect"
	"sync"
)

var (
	//  map[reflect.Type]map[reflect.Type]convFunc
	customConvertMap = sync.Map{}
)

func RegisterCustomConvertFunc(src, dest reflect.Type, fn convFunc) {
	if dest == nil || src == nil || fn == nil {
		panic(parameterCannotNilError("[dest] or [src] or [fn]"))
	}
	dstType, _ := typeElem(dest, 0)
	srcType, _ := typeElem(src, 0)

	conv := getCustomConvertFromMap(srcType, dstType)
	if conv != nil {
		panic(duplicateRegisterError(src, dest))
	}
	addCustomConvertToMap(srcType, dstType, fn)
}

func addCustomConvertToMap(src, dest reflect.Type, fn convFunc) {
	destMap, ok := customConvertMap.Load(src)
	if !ok {
		destMap = &sync.Map{}
		customConvertMap.Store(src, destMap)
	}
	destMap.(*sync.Map).Store(dest, fn)
}

func getCustomConvertFromMap(src, dest reflect.Type) convFunc {
	destMap, ok := customConvertMap.Load(src)
	if !ok {
		return nil
	}
	info, ok := destMap.(*sync.Map).Load(dest)
	if !ok {
		return nil
	}
	return info.(convFunc)
}
