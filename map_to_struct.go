package kast

import (
	"fmt"
	"reflect"
)

// MapToStruct dest 必须是*struct或者**struct
// 字段匹配机制：例如以下结构体和map
//
//	var src = map[string]any{
//		"last-name": "lastName",
//		"firstName": "firstName",
//		"Name":      "name",
//	}
//	type Dest struct {
//		LastName  string `json:"lastName"`
//		FirstName string `json:"firstName"`
//		Name      string `json:"name"`
//	}
//
// 解析出来的fieldInfo可能是以下的样子，tags包含json和字段名本身
// fieldInfo{ fieldName=LastName  tags={lastName,LastName} }
// fieldInfo{ fieldName=FirstName tags={firstName,FirstName} }
// fieldInfo{ fieldName=Name      tags={name,Name} }
// 那么会按照tags里面的值来循环去map中查找对应的key是否存在
// firstName 和 name 都存在map中
// LastName的两个tag在map中都找不到，会进入模糊匹配机制
// 模糊匹配机制：
//
//	忽略LastName的大小写，循环遍历map所有的key(key也会忽略大小写)
//	直到匹配到为止，如果没有找到，则不进行任何操作
func MapToStruct(src map[string]any, dest any) (err error) {
	return mapToStruct(src, dest, defaultMapToStructOptions)
}

func MapToStructWithOptions(src map[string]any, dest any, options MapToStructOptions) (err error) {
	return mapToStruct(src, dest, options)
}

func mapToStruct(src map[string]any, dest any, options MapToStructOptions) (err error) {
	if src == nil || len(src) == 0 || dest == nil {
		return nil
	}
	structValue := reflect.ValueOf(dest)
	if structValue.Kind() != reflect.Ptr {
		panic(fmt.Errorf("dest必须是*struct **struct类型"))
	}
	if structValue.IsNil() {
		panic(fmt.Errorf("dest不能为nil"))
	}
	structValue = structValue.Elem()
	if structValue.Kind() != reflect.Struct {
		panic(fmt.Errorf("dest必须是*struct **struct类型"))
	}
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			structValue.Set(reflect.New(structValue.Type().Elem()))
		}
		structValue = structValue.Elem()
	}

	m2sInfo := getStructInfo(reflect.TypeOf(dest), options)
	if m2sInfo == nil {
		return nil
	}

	var (
		ok          bool
		matched     bool
		lastFuzzKey string
		usedKey     = getUsedKeyMap()
		key         string
		val         any
	)
	defer putUsedKeyMap(usedKey)
	// for fieldName, fieldInfo := range m2sInfo.fields {
	for _, fieldInfo := range m2sInfo.listFields {
		for _, tag := range fieldInfo.tags {
			v, ok := src[tag]
			if ok {
				fieldValue := fieldInfo.getFieldReflectValue(structValue)
				err = fieldInfo.convFunc(fieldValue, v)
				if err != nil {
					return err
				}
				for i := 0; i < len(fieldInfo.otherFieldIndex); i++ {
					fieldValue := fieldInfo.getOtherFieldReflectValue(structValue, i)
					err = fieldInfo.convFunc(fieldValue, v)
					if err != nil {
						return err
					}
				}
				matched = true
				// usedKey[tag] = struct{}{}
				break
			}
		}
		if matched {
			matched = false
			continue
		}
		lastFuzzKey = fieldInfo.lastFuzzKey.Load().(string)
		val, ok = src[lastFuzzKey]
		if !ok {

			key, val = m2sInfo.fuzzyMatchFieldFunc(src, fieldInfo.fieldName, usedKey)
		}
		if key != "" {
			fieldInfo.lastFuzzKey.Store(key)
			if val != nil {
				//////////////////////////////////////////////////////////
				fieldValue := fieldInfo.getFieldReflectValue(structValue)
				err = fieldInfo.convFunc(fieldValue, val)
				if err != nil {
					return err
				}
				for i := 0; i < len(fieldInfo.otherFieldIndex); i++ {
					fieldValue = fieldInfo.getOtherFieldReflectValue(structValue, i)
					err = fieldInfo.convFunc(fieldValue, val)
					if err != nil {
						return err
					}
				}
			}
			usedKey[key] = struct{}{}
		}
	}
	return nil
}
