package kast

import (
	"reflect"
	"sync"
)

var (
	//  map[reflect.Type]*structInfo{}
	structInfoMap = sync.Map{}
)

func addStructInfoToMap(info *structInfo) {
	structInfoMap.Store(info.structType, info)
}

func getStructInfoFromMap(typ reflect.Type) *structInfo {
	info, ok := structInfoMap.Load(typ)
	if ok {
		return info.(*structInfo)
	}
	return nil
}

type structInfo struct {
	fields map[string]*fieldInfo
	//
	listFields     []*fieldInfo
	fieldListIndex map[string]int
	// *struct  => struct
	// **struct => struct
	structType reflect.Type
	//
	fuzzyMatchFieldFunc func(srcMap map[string]any, fieldName string, usedMapKeys map[string]struct{}) (mapKey string, mapValue any)
	//
	fieldTagsFunc func(structField reflect.StructField) []string
}

func (s *structInfo) GetField(fieldName string) *fieldInfo {
	return s.fields[fieldName]
}

func (s *structInfo) AddField(field reflect.StructField, fieldIndex []int) {
	if s.fields == nil {
		s.fields = make(map[string]*fieldInfo)
	}
	convFieldInfo, ok := s.fields[field.Name]
	if !ok {
		info := &fieldInfo{
			fieldName:  field.Name,
			fieldType:  field.Type,
			fieldIndex: fieldIndex,
			tags:       s.getFieldTags(field),
			convFunc:   getMapToStructFieldConvertFunc(field.Type),
		}
		info.lastFuzzKey.Store(field.Name)
		s.fields[field.Name] = info
		s.listFields = append(s.listFields, info)
		return
	}
	if convFieldInfo.otherField == nil {
		convFieldInfo.otherField = make([]*fieldInfo, 0, 2)
	}
	if convFieldInfo.fieldType == field.Type {
		convFieldInfo.otherField = append(convFieldInfo.otherField, convFieldInfo)
		return
	}
	info := &fieldInfo{
		fieldName:  field.Name,
		tags:       s.getFieldTags(field),
		fieldType:  field.Type,
		fieldIndex: fieldIndex,
		convFunc:   getMapToStructFieldConvertFunc(field.Type),
	}
	info.lastFuzzKey.Store(field.Name)
	convFieldInfo.otherField = append(convFieldInfo.otherField, info)
}

func (s *structInfo) getFieldTags(field reflect.StructField) []string {
	tags := s.fieldTagsFunc(field)
	if tags == nil || len(tags) == 0 {
		return []string{field.Name}
	}
	return tags
}
