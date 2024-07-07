package kast

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	// s2sInfosMap = map[reflect.Type]map[reflect.Type]*s2sInfo{}
	s2sInfosMap = sync.Map{}
)

func addS2SInfoToMap(src, dest reflect.Type, info *s2sInfo) {
	//destMap := s2sInfosMap[src]
	//if destMap == nil {
	//	destMap = make(map[reflect.Type]*s2sInfo)
	//	s2sInfosMap[src] = destMap
	//}
	//destMap[dest] = info

	// destMap = make(map[reflect.Type]*s2sInfo)
	destMap, ok := s2sInfosMap.Load(src)
	if !ok {
		destMap = &sync.Map{}
		s2sInfosMap.Store(src, destMap)
	}
	destMap.(*sync.Map).Store(dest, info)
}

func getS2SInfoFromMap(src, dest reflect.Type) *s2sInfo {
	//destMap := s2sInfosMap[src]
	//if destMap == nil {
	//	return nil
	//}
	//return destMap[dest]
	destMap, ok := s2sInfosMap.Load(src)
	if !ok {
		return nil
	}
	info, ok := destMap.(*sync.Map).Load(dest)
	if !ok {
		return nil
	}
	return info.(*s2sInfo)
}

type s2sFieldInfo struct {
	dst      *fieldInfo
	src      *fieldInfo
	convFunc convFunc
}

type s2sInfo struct {
	fields     map[string]*s2sFieldInfo
	listFields []*s2sFieldInfo
	// TODO
	sameTypeConv convFunc
}

func (s2s *s2sInfo) AddField(srcField, dstField *fieldInfo) {
	if s2s.fields == nil {
		s2s.fields = map[string]*s2sFieldInfo{}
	}
	fieldInfo := &s2sFieldInfo{
		src:      srcField,
		dst:      dstField,
		convFunc: getConvFunc(dstField.fieldType, srcField.fieldType),
	}
	s2s.fields[dstField.fieldName] = fieldInfo
	s2s.listFields = append(s2s.listFields, fieldInfo)
}

func getOrSetS2SInfo(src, dest reflect.Type, options StructToStructOptions) *s2sInfo {

	s2s := getS2SInfoFromMap(src, dest)
	if s2s != nil {
		return s2s
	}

	srcInfo := getStructInfo(src)
	dstInfo := getStructInfo(dest)
	s2s = &s2sInfo{}

	if srcInfo.structType == dstInfo.structType {
		s2s.sameTypeConv = sameStructTypeConv(src, dest)
		addS2SInfoToMap(src, dest, s2s)
		return s2s
	}

	if options.FieldMappingFunc != nil {
		for dstFieldName, dstFieldInfo := range dstInfo.fields {
			srcField := options.FieldMappingFunc(srcInfo.structType, dstFieldName)
			srcFieldInfo := srcInfo.GetField(srcField.Name)
			if srcFieldInfo == nil {
				srcFieldInfo = &fieldInfo{
					fieldName:  srcField.Name,
					fieldType:  srcField.Type,
					fieldIndex: srcField.Index,
				}
			}
			s2s.AddField(srcFieldInfo, dstFieldInfo)
		}
	} else {
		for dstFieldName, dstFieldInfo := range dstInfo.fields {
			srcFieldInfo := srcInfo.GetField(dstFieldName)
			if srcFieldInfo != nil {
				s2s.AddField(srcFieldInfo, dstFieldInfo)
			}
		}
	}
	addS2SInfoToMap(src, dest, s2s)
	return s2s
}

func StructToStruct(src, dest any) error {
	return structToStruct(src, dest, defaultStructToStructOptions)
}

func StructToStructWithOptions(src, dest any, options StructToStructOptions) error {
	return structToStruct(src, dest, options)
}

func structToStruct(src, dest any, options StructToStructOptions) error {
	destStructValue := reflect.ValueOf(dest)
	if destStructValue.Kind() != reflect.Ptr {
		panic(fmt.Errorf("dest必须为*struct or **struct"))
	}
	if destStructValue.IsNil() {
		panic(fmt.Errorf("dest不能为nil"))
	}

	srcType := reflect.TypeOf(src)
	destType := reflect.TypeOf(dest)

	srcStructValue := reflect.ValueOf(src)

	info := getOrSetS2SInfo(srcType, destType, options)

	if info.sameTypeConv != nil {
		return info.sameTypeConv(destStructValue, srcStructValue)
	}

	/////////////////////////////////////////////
	if srcStructValue.Kind() == reflect.Ptr {
		if srcStructValue.IsNil() {
			return nil
		}
		srcStructValue = srcStructValue.Elem()
	}
	/////////////////////////////////////////
	if destStructValue.Kind() == reflect.Ptr {
		if destStructValue.IsNil() {
			destStructValue.Set(reflect.New(destStructValue.Type().Elem()))
		}
		destStructValue = destStructValue.Elem()
	}
	if destStructValue.Kind() == reflect.Ptr {
		if destStructValue.IsNil() {
			destStructValue.Set(reflect.New(destStructValue.Type().Elem()))
		}
		destStructValue = destStructValue.Elem()
	}
	for _, fieldInfo := range info.listFields {
		srcFieldValue := fieldInfo.src.getFieldReflectValue(srcStructValue)
		destFieldValue := fieldInfo.dst.getFieldReflectValue(destStructValue)
		err := fieldInfo.convFunc(destFieldValue, srcFieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}
