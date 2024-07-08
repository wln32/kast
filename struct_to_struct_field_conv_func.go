package kast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type convFunc func(dest, src reflect.Value) error

// 这个是用来给structToStruct用的
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

// 目前对于指针类型的会做深拷贝
// 对于slice和map的不会做深拷贝
func getConvFunc(dest, src reflect.Type) (fn convFunc) {

	switch dest.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fn = convToInt(src)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fn = convToUint(src)
	case reflect.Float32, reflect.Float64:
		fn = convToFloat(src)
	case reflect.String:
		fn = convToString(src)
	case reflect.Bool:
		fn = convToBool(src)
	case reflect.Slice /*,reflect.Array */ :
		fn = convToSlice(dest, src)
	case reflect.Struct:
		fn = convToStruct(dest, src)
	case reflect.Map:
		fn = convToMap(dest, src)
	case reflect.Pointer:
		dest = dest.Elem()
		fn = getConvFunc(dest, src)
		if fn == nil {
			break
		}
		return destPtrConvFunc(fn)
	}
	if fn == nil {
		panic(unsupportedTypesError(src, dest))
	}
	return fn
}

func destPtrConvFunc(fn convFunc) convFunc {
	return func(dest, src reflect.Value) error {
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		return fn(dest.Elem(), src)
	}
}

// 数字类型可以互相转换
// bool转int: true=1,false=0
// string转int: strconv.ParseInt
// 支持任意级指针解引用
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
		return srcPtrConvFunc(fn)
	}
	return nil
}

// 数字类型可以互相转换
// bool转int: true=1,false=0
// string转int: strconv.ParseUint
// 支持任意级指针解引用
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
		return srcPtrConvFunc(fn)
	}

	return nil
}

// 数字类型可以互相转换
// bool转int: true=1,false=0
// string转int: strconv.ParseFloat
// 支持任意级指针解引用
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
			dest.SetFloat(src.Float())
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
		return srcPtrConvFunc(fn)
	}
	return nil
}

// 数字类型可以转string
// bool类型可以转string
// []byte类型可以转string
// 如果是以下3种类型的话,会调用json.Marshal来序列化，然后在转成string
//
//	1: slice (除了[]byte类型之外或者是类似于 type MySlice []byte的)
//	2: struct
//	3: map
//
// 支持任意级指针解引用
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
		return srcPtrConvFunc(fn)
	}
	return nil
}

// 数字类型转bool: 非0值=true，0=false
// string转bool: strconv.ParseBool
// 支持任意级指针解引用
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
		return srcPtrConvFunc(fn)
	}
	return nil
}

func srcPtrConvFunc(fn convFunc) convFunc {
	return func(dest, src reflect.Value) error {
		if src.IsNil() {
			// TODO nil值覆盖，或者直接退出，不设置dest
			return nil
		}
		return fn(dest, src.Elem())
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

func convToSlice(dest, src reflect.Type) convFunc {
	dest, _ = typeElem(dest, 0)
	src, _ = typeElem(src, 0)
	switch dest.Kind() {
	case reflect.Uint8:
		return convToBytes(src)
	}
	if dest == src {
		// TODO 深拷贝
		// []int []int64 []float64 []bool
		return func(dest, src reflect.Value) error {
			switch src.Kind() {
			case reflect.Ptr:
				if src.IsNil() {
					return nil
				}
				src = src.Elem()
			case reflect.Slice:
			}
			dest.Set(src)
			return nil
		}
	}
	return nil
}

// []byte
// string
// 如果是以下3种类型的话,会调用json.Marshal来序列化，然后在转成string
//
//	1: slice (除了[]byte类型之外或者是类似于 type MySlice []byte的)
//	2: struct
//	3: map
//
// 支持任意级指针解引用
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
		return srcPtrConvFunc(fn)
	}
	return nil
}

// struct转map: 会先用json.Marshal序列化，然后在反序列 (暂时这样做)
// map转map: 目前只支持类型相同的map
// 支持任意级指针解引用
func convToMap(dest, src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Struct:
		// to map[string]any
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
	case reflect.Map:
		if dest == src {
			return func(dest, src reflect.Value) error {
				dest.Set(src)
				return nil
			}
		}
	case reflect.Ptr:
		fn := convToMap(dest, src.Elem())
		if fn == nil {
			break
		}
		return srcPtrConvFunc(fn)
	}
	return nil
}

// map转struct: 如果map的类型是map[string]any,会调用mapToStruct
// struct转struct:
// 支持任意级指针解引用
func convToStruct(dest, src reflect.Type) convFunc {
	switch src.Kind() {
	case reflect.Map:
		if src == strAnyMapType {
			return func(dest, src reflect.Value) error {
				return mapToStruct(src.Interface().(map[string]any), dest, defaultMapToStructOptions)
			}
		}
	case reflect.Struct:
		fmt.Println("convToStruct ==>", dest, src)
		if src == dest {
			return sameStructTypeFieldConv()
		}
		info := getOrSetS2SInfo(src, dest, defaultStructToStructOptions)
		return toStruct(info)
	case reflect.Ptr:
		fn := convToStruct(dest, src.Elem())
		if fn == nil {
			break
		}
		return srcPtrConvFunc(fn)
	}
	return nil
}

func sameStructTypeFieldConv() convFunc {
	return func(dest, src reflect.Value) error {
		dest.Set(src)
		return nil
	}
}

func toStruct(info *s2sInfo) func(destStructValue, srcStructValue reflect.Value) error {
	return func(destStructValue, srcStructValue reflect.Value) error {
		/////////////////////////////////////////////
		if srcStructValue.Kind() == reflect.Ptr {
			if srcStructValue.IsNil() {
				return nil
			}
			srcStructValue = srcStructValue.Elem()
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
