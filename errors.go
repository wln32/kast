package kast

import (
	"fmt"
	"reflect"
)

func unsupportedTypesError(from, to reflect.Type) error {
	return fmt.Errorf("kast: 不支持从`%v`转换到`%v`", from, to)
}

func duplicateRegisterError(from, to reflect.Type) error {
	return fmt.Errorf("kast: 重复注册自定义转换类型`%v`到`%v`", from, to)
}
func mapToStructUnsupportedTypesError(fieldType reflect.Type) error {
	return fmt.Errorf("kast: MapToStruct: 不支持的类型转换%v", fieldType)
}

func destTypeMismatchError(dest any) error {
	return fmt.Errorf("kast: dest必须是*struct或者**struct类型,不能是: %T", dest)
}

func typeMustBeStructError(typ reflect.Type) error {
	return fmt.Errorf("kast: 必须是*struct或者**struct类型,不能是: %v", typ)
}

func parameterCannotNilError(parameterName string) error {
	return fmt.Errorf("参数:`%s`不能为nil", parameterName)
}
