package kast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type convFunc func(dest, src reflect.Value) error

func sameStructTypeConv(src, dest reflect.Type) convFunc {
	getPtrTypeDeref := func(typ reflect.Type) int {
		deref := 0
		for typ.Kind() == reflect.Ptr {
			deref++
			typ = typ.Elem()
		}
		return deref
	}
	deref := getPtrTypeDeref(dest) - getPtrTypeDeref(src)
	return func(dest, src reflect.Value) error {
		switch deref {
		case 0:
			// 如果是一个字段的话
			if dest.Kind() == reflect.Struct {
				dest.Set(src.Elem())
			}
			// *struct => *struct
			if src.IsNil() {
				return nil
			}
			dest.Elem().Set(src.Elem())
			return nil
		case 1:
			// *struct => **struct
			if src.Kind() == reflect.Pointer {
				if src.IsNil() {
					return nil
				}
				dest.Elem().Set(src)
				return nil
			}
			// struct => *struct
			dest.Elem().Set(src)
			return nil
		case 2:
			// struct => **struct
			dest = dest.Elem()
			if dest.IsNil() {
				dest.Set(reflect.New(dest.Type().Elem()))
			}
			dest.Elem().Set(src)
		}
		return nil
	}
}

func getConvFunc(dest, src reflect.Type) convFunc {
	switch dest.Kind() {
	case reflect.Pointer:
		dest = dest.Elem()
		destType := dest
		fn := getConvFunc(dest, src)
		return func(dest, src reflect.Value) error {
			if dest.IsNil() {
				dest.Set(reflect.New(destType))
			}
			return fn(dest.Elem(), src)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convToInt(src)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convToUint(src)
	case reflect.Float32, reflect.Float64:
		return convToFloat(src)
	case reflect.String:
		return convToString(src)
	case reflect.Bool:
		return convToBool(src)
	case reflect.Slice:
		return convToBytes(src)
	case reflect.Struct:
		return convToStruct(dest, src)
	case reflect.Map:
		return convToMap(dest, src)
	}

	panic(fmt.Errorf("不支持的dest类型转换`%v`", dest))
}

func convToInt(src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dest, src reflect.Value) error {
			dest.SetInt(src.Int())
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dest, src reflect.Value) error {
			dest.SetInt(int64(src.Uint()))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(dest, src reflect.Value) error {
			dest.SetInt(int64(src.Float()))
			return nil
		}
	case reflect.String:
		return func(dest, src reflect.Value) error {
			v, err := strconv.ParseInt(src.String(), 10, 64)
			if err != nil {
				return err
			}
			dest.SetInt(v)
			return nil
		}
	case reflect.Bool:
		return func(dest, src reflect.Value) error {
			b := src.Bool()
			if b {
				dest.SetInt(1)
			} else {
				dest.SetInt(0)
			}
			return nil
		}
	case reflect.Pointer:
		fn := convToInt(src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}

	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func convToUint(src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dest, src reflect.Value) error {
			dest.SetUint(uint64(src.Int()))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dest, src reflect.Value) error {
			dest.SetUint((src.Uint()))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(dest, src reflect.Value) error {
			dest.SetUint(uint64(src.Float()))
			return nil
		}
	case reflect.String:
		return func(dest, src reflect.Value) error {
			v, err := strconv.ParseUint(src.String(), 10, 64)
			if err != nil {
				return err
			}
			dest.SetUint(v)
			return nil
		}
	case reflect.Bool:
		return func(dest, src reflect.Value) error {
			b := src.Bool()
			if b {
				dest.SetUint(1)
			} else {
				dest.SetUint(0)
			}
			return nil
		}

	case reflect.Pointer:
		fn := convToUint(src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}

	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func convToFloat(src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dest, src reflect.Value) error {
			dest.SetFloat(float64(src.Int()))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dest, src reflect.Value) error {
			dest.SetFloat(float64(src.Uint()))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(dest, src reflect.Value) error {
			dest.SetFloat((src.Float()))
			return nil
		}
	case reflect.String:
		return func(dest, src reflect.Value) error {
			v, err := strconv.ParseFloat(src.String(), 64)
			if err != nil {
				return err
			}
			dest.SetFloat(v)
			return nil
		}
	case reflect.Bool:
		return func(dest, src reflect.Value) error {
			b := src.Bool()
			if b {
				dest.SetFloat(1)
			} else {
				dest.SetFloat(0)
			}
			return nil
		}

	case reflect.Pointer:
		fn := convToFloat(src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}

	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func convToString(src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dest, src reflect.Value) error {
			v := strconv.FormatInt(src.Int(), 10)
			dest.SetString(v)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dest, src reflect.Value) error {
			v := strconv.FormatUint(src.Uint(), 10)
			dest.SetString(v)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(dest, src reflect.Value) error {
			v := strconv.FormatFloat(src.Float(), 'G', -1, 64)
			dest.SetString(v)
			return nil
		}
	case reflect.String:
		return func(dest, src reflect.Value) error {
			dest.SetString(src.String())
			return nil
		}
	case reflect.Bool:
		return func(dest, src reflect.Value) error {
			v := strconv.FormatBool(src.Bool())
			dest.SetString(v)
			return nil
		}
	case reflect.Slice:
		if src.Elem().Kind() == reflect.Uint8 {
			return func(dest, src reflect.Value) error {
				dest.SetString(string(src.Bytes()))
				return nil
			}
		}
		return jsonMarshalToString
	case reflect.Struct, reflect.Map:
		return jsonMarshalToString
	case reflect.Pointer:
		fn := convToString(src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}
	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func convToBool(src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dest, src reflect.Value) error {
			dest.SetBool(src.Int() != 0)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dest, src reflect.Value) error {

			dest.SetBool(src.Uint() != 0)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(dest, src reflect.Value) error {

			dest.SetBool(src.Float() != 0)
			return nil
		}
	case reflect.String:
		return func(dest, src reflect.Value) error {
			v, err := strconv.ParseBool(src.String())
			if err != nil {
				return err
			}
			dest.SetBool(v)
			return nil
		}
	case reflect.Bool:
		return func(dest, src reflect.Value) error {
			dest.SetBool(src.Bool())
			return nil
		}

	case reflect.Pointer:
		fn := convToBool(src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}
	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func ptrConvFunc(fn convFunc) convFunc {
	return func(dest, src reflect.Value) error {
		if src.IsNil() {
			return fn(dest, reflect.New(src.Type().Elem()))
		} else {
			return fn(dest, src.Elem())
		}
	}
}

func jsonMarshalToString(dest, src reflect.Value) error {
	v := src.Interface()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	dest.SetString(string(b))
	return nil
}

func jsonMarshalToBytes(dest, src reflect.Value) error {
	v := src.Interface()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	dest.SetBytes(b)
	return nil
}

func convToBytes(src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Slice:
		if src.Elem().Kind() == reflect.Uint8 {
			return func(dest, src reflect.Value) error {
				dest.SetBytes(src.Bytes())
				return nil
			}
		}
		return jsonMarshalToBytes
	case reflect.String:
		return func(dest, src reflect.Value) error {
			dest.SetBytes([]byte(src.String()))
			return nil
		}
	case reflect.Struct:
		return jsonMarshalToBytes
	case reflect.Map:
		return jsonMarshalToBytes
	case reflect.Pointer:
		fn := convToBytes(src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}
	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func convToMap(dest, src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Struct:
		// to map[string]any
		strAnyMapType := reflect.TypeOf((*map[string]any)(nil)).Elem()
		if dest == strAnyMapType {
			return func(dest, src reflect.Value) error {
				v := src.Interface()
				m := make(map[string]any)
				bytes, err := json.Marshal(v)
				if err != nil {
					return err
				}
				err = json.Unmarshal(bytes, &m)
				if err != nil {
					return err
				}
				dest.Set(reflect.ValueOf(m))
				return nil
			}
		}
		panic(fmt.Errorf("暂时不支持从`%v`转换到`%v`", src, dest))
	case reflect.Map:
		if dest == src {
			return func(dest, src reflect.Value) error {
				dest.Set(src)
				return nil
			}
		}
		panic(fmt.Errorf("map的类型不一样，暂时不支持转换,dest=%v, src=%v", dest, src))
	case reflect.Ptr:
		fn := convToMap(dest, src.Elem())
		if fn == nil {
			break
		}
		return ptrConvFunc(fn)
	}
	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func convToStruct(dest, src reflect.Type) convFunc {
	src = typeElem(src)
	dest = typeElem(dest)

	switch src.Kind() {
	case reflect.Map:

	case reflect.Struct:
		info := getOrSetS2SInfo(src, dest, defaultStructToStructOptions)
		return toStruct(info)
	}
	panic(fmt.Errorf("不支持的src类型`%v`", src))
}

func toStruct(info *s2sInfo) func(destStructValue, srcStructValue reflect.Value) error {
	return func(destStructValue, srcStructValue reflect.Value) error {
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
}
