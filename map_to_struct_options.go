package kast

import (
	"reflect"
	"strings"
)

var (
	defaultMapToStructOptions = MapToStructOptions{}
)

type MapToStructOptions struct {
	// 返回值,默认返回字段的json tag和字段名，按照顺序
	// 例如:
	//	type Dest struct{
	//		LastName string `json:"lastName"`
	//		FirstName string `json:"firstName"`
	//		Omitempty string `json:",omitempty"`  // 会忽略空白的tag
	//		Name string `json:"-"` 				 // 同时也会忽略 -
	//	}
	// 默认返回：
	//		LastName : []string{lastName, LastName}
	// 		FirstName: []string{firstName, FirstName}
	//		Omitempty: []string{Omitempty}
	//		Name     : []string{Name}
	//		...
	FieldTagsFunc func(structField reflect.StructField) []string
	// 在srcMap中查找fieldName, usedMapKeys是已经被赋值过的key的集合，
	// 例如：
	//type Dest struct{
	//	LastName string `json:"lastName"`
	//}
	// srcMap = map{
	//	"last-name":"lastname",
	//  ...,
	//}
	// 在默认情况下，不会匹配到map里面的last-name
	// 因为默认的模糊匹配只会忽略大小写
	// 当然你可以把json tag指定为last-name，这样也可以匹配到
	// 但是如果你的json tag改变不了的话，改变了可能会影响其他更多的地方
	// 你就可以重写这个模糊的函数
	// 模糊匹配时会记住最后一次匹配到的key
	// 以此来加速匹配，前提是当你的key不会一直发生变化时
	// 如果你的key来回发生变化，则每次都需要模糊匹配，非常耗时
	FuzzyMatchFieldFunc func(srcMap map[string]any, fieldName string, usedMapKeys map[string]struct{}) (mapKey string, mapValue any)
}

func defaultMapToStructOptionsFieldTagsFunc(field reflect.StructField) []string {
	fieldName := field.Name
	fieldTag := field.Tag
	if fieldTag == "" {
		return []string{fieldName}
	}
	jsonTag := fieldTag.Get("json")
	// json:"test,omitempty"
	// json:"-"
	jsonTag = strings.Split(jsonTag, ",")[0]
	jsonTag = strings.TrimSpace(jsonTag)
	if jsonTag == "" || jsonTag == "-" {
		return []string{fieldName}
	}
	return []string{jsonTag, fieldName}
}

func defaultMapToOptionsFuzzyMatchField(src map[string]any, key string, usedKey map[string]struct{}) (string, any) {
	for k, v := range src {
		_, ok := usedKey[k]
		if ok {
			continue
		}
		if equalFold(key, k) {
			return k, v
		}
	}
	return "", nil
}

func equalFold(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}
