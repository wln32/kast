package kast

import (
	"reflect"
	"sync"
)

var (
	// s2sInfosMap = map[reflect.Type]map[reflect.Type]*s2sInfo{}
	s2sInfosMap = sync.Map{}
)

func addS2SInfoToMap(src, dest reflect.Type, info *s2sInfo) {
	destMap, ok := s2sInfosMap.Load(src)
	if !ok {
		destMap = &sync.Map{}
		s2sInfosMap.Store(src, destMap)
	}
	destMap.(*sync.Map).Store(dest, info)
}

func getS2SInfoFromMap(src, dest reflect.Type) *s2sInfo {
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
		src: srcField,
		dst: dstField,
		// 需要检测字段类型是否和当前结构体类型一致，
		// 如果相同，不需要递归，否则会无限递归
		// TODO： 循环引用
		// convFunc: s2s.getS2SFieldConvFunc(dstField, srcField),
	}
	s2s.fields[dstField.fieldName] = fieldInfo
	s2s.listFields = append(s2s.listFields, fieldInfo)
}

func (s2s *s2sInfo) SetFieldConv(info *s2sFieldInfo) {
	info.convFunc = s2s.getS2SFieldConvFunc(info.dst, info.src)
}

func (s2s *s2sInfo) getS2SFieldConvFunc(dstField, srcField *fieldInfo) convFunc {
	dstFieldType, dstDeref := typeElem(dstField.fieldType, 0)
	srcFieldType, srcDeref := typeElem(srcField.fieldType, 0)

	conv := getCustomConvertFromMap(srcFieldType, dstFieldType)
	if conv != nil {
		return getDerefConvFunc(dstDeref, srcDeref, conv)
	}
	if srcFieldType == dstFieldType {
		return getDerefConvFunc(dstDeref, srcDeref, sameTypeFieldConv())
	}
	switch {
	case srcFieldType.Kind() == reflect.Struct && dstFieldType.Kind() == reflect.Struct:
		info := getOrSetS2SInfo(srcFieldType, dstFieldType, defaultStructToStructOptions)
		conv := toStruct(info)
		return getDerefConvFunc(dstDeref, srcDeref, conv)
	default:
		return getConvFunc(dstField.fieldType, srcField.fieldType)
	}
}

func getDerefConvFunc(dstDeref, srcDeref int, fn convFunc) convFunc {
	if dstDeref > 0 {
		return getDerefConvFunc(dstDeref-1, srcDeref, destPtrConvFunc(fn))
	}
	if srcDeref > 0 {
		return getDerefConvFunc(dstDeref, srcDeref-1, srcPtrConvFunc(fn))
	}
	return fn
}

func getOrSetS2SInfo(src, dest reflect.Type, options StructToStructOptions) *s2sInfo {
	s2s := getS2SInfoFromMap(src, dest)
	if s2s != nil {
		return s2s
	}
	srcInfo := getStructInfo(src)
	dstInfo := getStructInfo(dest)
	return setS2SInfo(src, dest, srcInfo, dstInfo, options)
}

func setS2SInfo(src, dest reflect.Type, srcInfo, dstInfo *structInfo, options StructToStructOptions) *s2sInfo {
	s2s := &s2sInfo{}
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
	for _, field := range s2s.fields {
		s2s.SetFieldConv(field)
	}
	return s2s
}

func StructToStruct(src, dest any) error {
	return structToStruct(src, dest, defaultStructToStructOptions)
}

func StructToStructWithOptions(src, dest any, options StructToStructOptions) error {
	return structToStruct(src, dest, options)
}

func structToStruct(src, dest any, options StructToStructOptions) error {
	if src == nil || dest == nil {
		panic(parameterCannotNilError("[src] or [dest]"))
	}
	destStructValue := reflect.ValueOf(dest)
	if destStructValue.Kind() != reflect.Ptr {
		panic(destTypeMismatchError(dest))
	}
	if destStructValue.IsNil() {
		panic(parameterCannotNilError("dest"))
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
