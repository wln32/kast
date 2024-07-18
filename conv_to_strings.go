package kast

import (
	"reflect"
	"strconv"
)

func reflectConvToStrings(src reflect.Type) (fn convFunc) {
	srcElemKind := src.Elem().Kind()
	switch srcElemKind {
	case reflect.Int:
		fn = func(dest, src reflect.Value) error {
			v := intsToStrings(*(*[]int)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Int8:
		fn = func(dest, src reflect.Value) error {
			v := int8sToStrings(*(*[]int8)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Int16:
		fn = func(dest, src reflect.Value) error {
			v := int16sToStrings(*(*[]int16)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Int32:
		fn = func(dest, src reflect.Value) error {
			v := int32sToStrings(*(*[]int32)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Int64:
		fn = func(dest, src reflect.Value) error {
			v := int64sToStrings(*(*[]int64)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Uint:
		fn = func(dest, src reflect.Value) error {
			v := uintsToStrings(*(*[]uint)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Uint16:
		fn = func(dest, src reflect.Value) error {
			v := uint16sToStrings(*(*[]uint16)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Uint32:
		fn = func(dest, src reflect.Value) error {
			v := uint32sToStrings(*(*[]uint32)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Uint64:
		fn = func(dest, src reflect.Value) error {
			v := uint64sToStrings(*(*[]uint64)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Float32:
		fn = func(dest, src reflect.Value) error {
			v := float32sToStrings(*(*[]float32)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Float64:
		fn = func(dest, src reflect.Value) error {
			v := float64sToStrings(*(*[]float64)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.String:
		fn = func(dest, src reflect.Value) error {
			v := *(*[]string)(getEface(src.Interface()).data)
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	case reflect.Bool:
		fn = func(dest, src reflect.Value) error {
			v := boolsToStrings(*(*[]bool)(getEface(src.Interface()).data))
			dest.Set(reflect.ValueOf(v))
			return nil
		}
	}
	return nil
}

func intsToStrings(arr []int) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.Itoa(v))
	}
	return
}

func int8sToStrings(arr []int8) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.Itoa(int(v)))
	}
	return
}

func int16sToStrings(arr []int16) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.Itoa(int(v)))
	}
	return
}

func int32sToStrings(arr []int32) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.Itoa(int(v)))
	}
	return
}

func int64sToStrings(arr []int64) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatInt(v, 10))
	}
	return
}

func uintsToStrings(arr []uint) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatUint(uint64(v), 10))
	}
	return
}

func uint16sToStrings(arr []uint16) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatUint(uint64(v), 10))
	}
	return
}

func uint32sToStrings(arr []uint32) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatUint(uint64(v), 10))
	}
	return
}

func uint64sToStrings(arr []uint64) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatUint(uint64(v), 10))
	}
	return
}

func float32sToStrings(arr []float32) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatFloat(float64(v), 'G', -1, 32))
	}
	return
}

func float64sToStrings(arr []float64) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatFloat(v, 'G', -1, 64))
	}
	return
}

func boolsToStrings(arr []bool) (res []string) {
	for _, v := range arr {
		res = append(res, strconv.FormatBool(v))
	}
	return
}
