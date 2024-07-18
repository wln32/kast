package kast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func reflectToInt64(src any) (int64, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(val.Float()), nil
	case reflect.String:
		return strconv.ParseInt(val.String(), 10, 64)
	case reflect.Bool:
		if val.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的类型转换")
	}
}

func reflectToUint64(src any) (uint64, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(val.Float()), nil
	case reflect.String:
		return strconv.ParseUint(val.String(), 10, 64)
	case reflect.Bool:
		if val.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的类型转换")
	}
}

func reflectToFloat64(src any) (float64, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	case reflect.String:
		return strconv.ParseFloat(val.String(), 64)
	case reflect.Bool:
		if val.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的类型转换")
	}
}

func sliceOrMapOrStructToString(val reflect.Value) (string, error) {
	// TODO 配置开关,注册类型转换
	s, err := json.Marshal(val.Interface())
	return string(s), err
}
func sliceOrMapOrStructToBytes(val reflect.Value) ([]byte, error) {
	// TODO 配置开关,注册类型转换
	s, err := json.Marshal(val.Interface())
	return s, err
}

func reflectToString(src any) (string, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'G', -1, 64), nil
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return string(val.Bytes()), nil
		}
		return sliceOrMapOrStructToString(val)
	case reflect.Struct:
		if val.Type().String() == "time.Time" {
			return val.Interface().(time.Time).String(), nil
		}
		return sliceOrMapOrStructToString(val)
	case reflect.Map:
		return sliceOrMapOrStructToString(val)
	default:
	}
	return "", fmt.Errorf("不支持的类型转换")
}

func reflectToBool(src any) (bool, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0, nil
	case reflect.String:
		return strconv.ParseBool(val.String())
	case reflect.Bool:
		return val.Bool(), nil
	default:
		return false, fmt.Errorf("不支持的类型转换")
	}
}

func reflectToBytes(src any) ([]byte, error) {
	switch v := src.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	}
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return val.Bytes(), nil
		}
		return sliceOrMapOrStructToBytes(val)
	case reflect.Map:
		return sliceOrMapOrStructToBytes(val)
	case reflect.Struct:
		if val.Type().String() == "time.Time" {
			return val.Interface().(time.Time).MarshalJSON()
		}
		return sliceOrMapOrStructToBytes(val)
	default:
	}
	return nil, fmt.Errorf("不支持的类型转换")
}

func reflectToStruct(dstType reflect.Type) func(dest reflect.Value, src any) error {
	switch {
	case dstType == timeType:
		return func(dest reflect.Value, src any) error {
			t, err := reflectToTime(src)
			if err != nil {
				return err
			}
			*dest.Addr().Interface().(*time.Time) = t
			return nil
		}
	default:
		return func(dest reflect.Value, src any) error {
			srcValue := reflect.ValueOf(src)
			srcValue = valueElem(srcValue)
			if srcValue.Kind() == reflect.Struct {
				return structToStruct(src, dest.Interface(), defaultStructToStructOptions)
			}
			return fmt.Errorf("不支持的类型`%T`转换", src)
		}
	}
}

func reflectToTime(src any) (time.Time, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Struct:
		if val.Type().String() == "time.Time" {
			return val.Interface().(time.Time), nil
		}
	// case reflect.String: 字符串转time比较复杂，因为格式很多种
	default:
	}
	return time.Time{}, fmt.Errorf("不支持的类型转换")
}

func reflectToMap(dstType reflect.Type) func(dest reflect.Value, src any) error {
	if dstType == strAnyMapType {
		return func(dest reflect.Value, src any) error {
			return mapToStruct(src.(map[string]any), dest, defaultMapToStructOptions)
		}
	}
	panic(fmt.Errorf("不支持的类型转换%v", dstType))
}

func reflectToAny(src any) (any, error) {
	return src, nil
}

func reflectToStringSlice(src any) (res []string, err error) {
	switch arr := src.(type) {
	case []int:
		res = intsToStrings(arr)
	case []int64:
		res = int64sToStrings(arr)
	case []uint:
		res = uintsToStrings(arr)
	case []uint64:
		res = uint64sToStrings(arr)
	case []float32:
		res = float32sToStrings(arr)
	case []float64:
		res = float64sToStrings(arr)
	case []bool:
		res = boolsToStrings(arr)
	case []string:
		return arr, nil
	default:
		switch arr := src.(type) {
		case []int8:
			res = int8sToStrings(arr)
		case []int16:
			res = int16sToStrings(arr)
		case []int32:
			res = int32sToStrings(arr)
		case []uint16:
			res = uint16sToStrings(arr)
		case []uint32:
			res = uint32sToStrings(arr)
		}
		if res != nil {
			return
		}
		typ := reflect.TypeOf(src)
		if typ.Kind() != reflect.Slice {
			return nil, fmt.Errorf("暂时不支持将`%T`类型转换到`[]string`", src)
		}
		v := reflect.ValueOf(src)
		switch typ.Elem().Kind() {
		case reflect.Int:
			res = intsToStrings(*(*[]int)(getEface(v.Interface()).data))
		case reflect.Int8:
			res = int8sToStrings(*(*[]int8)(getEface(v.Interface()).data))
		case reflect.Int16:
			res = int16sToStrings(*(*[]int16)(getEface(v.Interface()).data))
		case reflect.Int32:
			res = int32sToStrings(*(*[]int32)(getEface(v.Interface()).data))
		case reflect.Int64:
			res = int64sToStrings(*(*[]int64)(getEface(v.Interface()).data))
		case reflect.Uint:
			res = uintsToStrings(*(*[]uint)(getEface(v.Interface()).data))
		case reflect.Uint16:
			res = uint16sToStrings(*(*[]uint16)(getEface(v.Interface()).data))
		case reflect.Uint32:
			res = uint32sToStrings(*(*[]uint32)(getEface(v.Interface()).data))
		case reflect.Uint64:
			res = uint64sToStrings(*(*[]uint64)(getEface(v.Interface()).data))
		case reflect.Float32:
			res = float32sToStrings(*(*[]float32)(getEface(v.Interface()).data))
		case reflect.Float64:
			res = float64sToStrings(*(*[]float64)(getEface(v.Interface()).data))
		case reflect.String:
			res = *(*[]string)(getEface(v.Interface()).data)
		case reflect.Bool:
			res = boolsToStrings(*(*[]bool)(getEface(v.Interface()).data))
		}
		return nil, fmt.Errorf("暂时不支持将`%T`类型转换到`[]string`", src)
	}
	return
}
