package kast

import (
	"reflect"
)

var (
	structInfoMap = map[reflect.Type]*structInfo{}
)

func addStructInfoToMap(info *structInfo) {
	structInfoMap[info.structType] = info
}

func getStructInfoFromMap(typ reflect.Type) *structInfo {
	info := structInfoMap[typ]
	return info
}

type structInfo struct {
	fields map[string]*fieldInfo
	//
	listFields []*fieldInfo
	fieldListIndex map[string]int
	// *struct  => struct
	// **struct => struct
	structType reflect.Type
	//
	fuzzyMatchFieldFunc func(srcMap map[string]any, fieldName string,usedMapKeys map[string]struct{})(mapKey string,mapValue any)
	//
	fieldTagsFunc func(structField reflect.StructField)[]string
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
	if convFieldInfo.otherFieldIndex == nil {
		convFieldInfo.otherFieldIndex = make([][]int, 0, 2)
	}
	convFieldInfo.otherFieldIndex = append(convFieldInfo.otherFieldIndex, fieldIndex)
}

func(s *structInfo) getFieldTags(field reflect.StructField) []string {
	tags:=s.fieldTagsFunc(field)
	if tags == nil || len(tags) == 0 {
		return []string{field.Name}
	}
	return tags
}


