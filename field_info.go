package kast

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type fieldInfo struct {
	fieldName string
	// tag = {json,fieldName}
	tags      []string
	fieldType reflect.Type

	fieldIndex []int
	// src = map.v
	// dest = struct field
	convFunc func(dest reflect.Value, src any) error
	// 主要用于重名的字段
	// 一般是嵌套结构体才会有重名字段
	// TODO 需要根据字段类型不同重新设置convFunc
	otherField []*fieldInfo
	// lastFuzzKey string
	lastFuzzKey atomic.Value
}

func (c *fieldInfo) getFieldReflectValue(structValue reflect.Value) reflect.Value {
	if len(c.fieldIndex) == 1 {
		return structValue.Field(c.fieldIndex[0])
	}
	return fieldReflectValue(structValue, c.fieldIndex)
}

func (c *fieldInfo) getOtherFieldReflectValue(structValue reflect.Value, fieldLevel int) reflect.Value {
	field := c.otherField[fieldLevel]
	if len(field.fieldIndex) == 1 {
		return structValue.Field(field.fieldIndex[0])
	}
	return fieldReflectValue(structValue, field.fieldIndex)
}

func fieldReflectValue(v reflect.Value, fieldIndex []int) reflect.Value {
	for i, x := range fieldIndex {
		if i > 0 {
			switch v.Kind() {
			case reflect.Pointer:
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
				//case reflect.Interface:
				//	// Compatible with previous code
				//	// Interface => struct
				//	v = v.Elem()
				//	if v.Kind() == reflect.Ptr {
				//		// maybe *struct or other types
				//		v = v.Elem()
				//	}
			}
		}
		v = v.Field(x)
	}
	return v
}

func getMapToStructFieldConvertFunc(fieldType reflect.Type) func(dest reflect.Value, src any) error {
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dest reflect.Value, src any) error {
			v, err := reflectToInt64(src)
			//v, err := toInt64(src)
			if err != nil {
				return err
			}
			dest.SetInt(v)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dest reflect.Value, src any) error {
			v, err := reflectToUint64(src)
			if err != nil {
				return err
			}
			dest.SetUint(v)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(dest reflect.Value, src any) error {
			v, err := reflectToFloat64(src)
			//v, err := toFloat64(src)
			if err != nil {
				return err
			}
			dest.SetFloat(v)
			return nil
		}
	case reflect.String:
		return func(dest reflect.Value, src any) error {
			v, err := reflectToString(src)
			//v, err := toString(src)
			if err != nil {
				return err
			}
			dest.SetString(v)
			return nil
		}
	case reflect.Bool:
		return func(dest reflect.Value, src any) error {
			v, err := reflectToBool(src)
			if err != nil {
				return err
			}
			dest.SetBool(v)
			return nil
		}
	case reflect.Slice:
		return func(dest reflect.Value, src any) error {
			b, err := reflectToBytes(src)
			if err != nil {
				return err
			}
			dest.SetBytes(b)
			return nil
		}
	case reflect.Struct:
		return reflectToStruct(fieldType)
	case reflect.Map:
		return reflectToMap
	case reflect.Ptr:
		conv := getMapToStructFieldConvertFunc(fieldType.Elem())
		if conv != nil {
			return mapToStructFieldPtrConvFunc(conv)
		}
	}
	panic(fmt.Errorf("不支持的类型转换%v", fieldType))
}

func mapToStructFieldPtrConvFunc(convFunc func(to reflect.Value, from any) error) func(to reflect.Value, from any) error {
	return func(to reflect.Value, from any) error {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		return convFunc(to.Elem(), from)
	}
}
