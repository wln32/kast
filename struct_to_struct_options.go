package kast

import (
	"reflect"
)

var (
	defaultStructToStructOptions = StructToStructOptions{}
)

type StructToStructOptionFieldMappingFunc = func(srcStructType reflect.Type, destFieldName string) reflect.StructField

type StructToStructOptions struct {
	FieldMappingFunc func(srcStructType reflect.Type, destFieldName string) reflect.StructField

	DeepCopyPtrStruct bool
	DeepCopySlice     bool
	DeepCopyMap       bool
	// DeepCopy = DeepCopyPtrStruct + DeepCopySlice + DeepCopyMap
	DeepCopy bool
}
