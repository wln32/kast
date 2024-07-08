package kast

import (
	"fmt"
	"reflect"
)

func getStructInfo(structType reflect.Type, convOptions ...MapToStructOptions) *structInfo {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("必须为struct类型"))
	}
	info := getStructInfoFromMap(structType)
	if info != nil {
		return info
	}
	var options = defaultMapToStructOptions
	if convOptions != nil && len(convOptions) > 0 {
		options = convOptions[0]
	}

	if options.FuzzyMatchFieldFunc == nil {
		options.FuzzyMatchFieldFunc = defaultMapToOptionsFuzzyMatchField
	}
	if options.FieldTagsFunc == nil {
		options.FieldTagsFunc = defaultMapToStructOptionsFieldTagsFunc
	}

	info = &structInfo{
		structType:          structType,
		fuzzyMatchFieldFunc: options.FuzzyMatchFieldFunc,
		fieldTagsFunc:       options.FieldTagsFunc,
	}
	parseStruct(structType, []int{}, info)
	addStructInfoToMap(info)
	return info
}

func parseStruct(structType reflect.Type, parentIndex []int, structInfo *structInfo) {
	var (
		structField reflect.StructField
		fieldType   reflect.Type
	)
	// TODO:
	//  Check if the structure has already been cached in the cache.
	//  If it has been cached, some information can be reused,
	//  but the [FieldIndex] needs to be reset.
	//  We will not implement it temporarily because it is somewhat complex
	for i := 0; i < structType.NumField(); i++ {
		structField = structType.Field(i)
		fieldType = structField.Type
		// Only do converting to public attributes.
		if structField.IsExported() == false {
			continue
		}
		if fieldType.Kind() == reflect.Interface {
			// TODO 支持空接口
			continue
		}
		if structField.Anonymous == false {
			structInfo.AddField(structField, append(parentIndex, i))
			continue
		}

		// Maybe it's struct/*struct embedded.
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}

		// TODO:
		if structField.Tag != "" {
			structInfo.AddField(structField, append(parentIndex, i))
		}
		// TODO：支持匿名的不带tag的结构体字段作为整个字段
		parseStruct(fieldType, append(parentIndex, i), structInfo)
	}
}
