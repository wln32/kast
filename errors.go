package kast

import (
	"fmt"
	"reflect"
)

func unsupportedTypesError(from, to reflect.Type) error {
	return fmt.Errorf("不支持从`%v`转换到`%v`", from, to)
}
