package kast

import (
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
		return uint64(val.Uint()), nil
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
		return float64(val.Float()), nil
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

func reflectToString(src any) (string, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 64), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 64), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat((val.Float()), 'G', -1, 64), nil
	case reflect.String:
		return val.String(), nil
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return string(val.Bytes()), nil
		}
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	default:
		// time.Time
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
	default:
	}
	return nil, fmt.Errorf("不支持的类型转换")
}

func reflectToTime(src any) (time.Time, error) {
	val := reflect.ValueOf(src)
	switch val.Kind() {
	case reflect.String:
	default:
	}
	panic("")
}

func toInt64(src any) (int64, error) {
	switch v := src.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		vv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return vv, nil
	}
	return 0, fmt.Errorf("不支持的类型`%T`", src)
}

func toUint64(src any) (uint64, error) {
	switch v := src.(type) {
	case int:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		vv, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return vv, nil
	}
	return 0, fmt.Errorf("不支持的类型`%T`", src)
}

func toFloat64(src any) (float64, error) {
	switch v := src.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		vv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return vv, nil
	}
	return 0, fmt.Errorf("不支持的类型`%T`", src)
}

func toString(src any) (vv string, err error) {
	switch v := src.(type) {
	case int:
		vv = strconv.Itoa(v)
	case int8:
		vv = strconv.FormatInt(int64(v), 10)
	case int16:
		vv = strconv.FormatInt(int64(v), 10)
	case int32:
		vv = strconv.FormatInt(int64(v), 10)
	case int64:
		vv = strconv.FormatInt(int64(v), 10)
	case uint8:
		vv = strconv.FormatUint(uint64(v), 10)
	case uint16:
		vv = strconv.FormatUint(uint64(v), 10)
	case uint32:
		vv = strconv.FormatUint(uint64(v), 10)
	case uint64:
		vv = strconv.FormatUint(uint64(v), 10)
	case uint:
		vv = strconv.FormatUint(uint64(v), 10)
	case float32:
		vv = strconv.FormatFloat(float64(v), 'G', -1, 32)
	case float64:
		vv = strconv.FormatFloat(float64(v), 'G', -1, 64)
	case string:
		vv = v
	case []byte:
		vv = string(v)
	case bool:
		vv = strconv.FormatBool(v)
	default:

	}
	return vv, nil
}

func toBool(src any) (vv bool, err error) {
	switch v := src.(type) {
	case int:
		vv = v != 0
	case int8:
		vv = v != 0
	case int16:
		vv = v != 0
	case int32:
		vv = v != 0
	case int64:
		vv = v != 0
	case uint8:
		vv = v != 0
	case uint16:
		vv = v != 0
	case uint32:
		vv = v != 0
	case uint64:
		vv = v != 0
	case uint:
		vv = v != 0
	case float32:
		vv = v != 0
	case float64:
		vv = v != 0
	case string:
		vv, err = strconv.ParseBool(v)

	case bool:
		vv = v
	default:
	}
	return vv, err
}
