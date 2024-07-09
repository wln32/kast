package kast_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wln32/kast"
)

var srcStructFields8 = testStructType8{
	Name: "github",

	CategoryId:  "1",
	Price:       198.09,
	Code:        "1",
	Image:       "https://github.com",
	Description: "This is the data for testing eight fields",
	Status:      1,
	IdType:      2,

	Score: 100,
	Age:   98,
	ID:    199,
}

func Test_s2s(t *testing.T) {
	type Src struct {
		Name string
		Age  int
	}
	type Dest struct {
		Name []byte
		Age  uint
	}
	var dest Dest
	src := Src{
		Name: "lxq",
		Age:  98,
	}
	testRun(t, func(assert *assert.Assertions) {
		err := kast.StructToStruct(src, &dest)
		assert.NoError(err)
		assert.Equal(src.Age, int(dest.Age))
		assert.Equal(src.Name, string(dest.Name))
	})
}

func Test_s2s_field8(t *testing.T) {
	const N = 100
	dest := &testStructType82{}
	start := time.Now()
	for i := 0; i < N; i++ {
		kast.StructToStruct(srcStructFields8, dest)
	}
	cost := time.Since(start)
	t.Log("耗时：", cost)
}

func Test_s2s_recursion(t *testing.T) {
	type Inner1 struct {
		Name string
		Age  int
		Next **Inner1
	}

	type Src struct {
		Inner *Inner1
	}
	type Dest struct {
		Inner *Inner1
	}
	srcInner3 := &Inner1{
		Name: "lxq333",
		Age:  334,
	}
	srcInner2 := &Inner1{
		Name: "lxq111",
		Age:  1111,
		Next: &srcInner3,
	}
	srcInner := &Inner1{
		Name: "lxq",
		Age:  100,
		Next: &srcInner2,
	}

	var recursion func(dest, src *Inner1, assert *assert.Assertions)

	recursion = func(dest, src *Inner1, assert *assert.Assertions) {
		assert.Equal(src.Name, dest.Name)
		assert.Equal(src.Age, dest.Age)
		if src.Next == nil {
			assert.Nil(dest.Next)
		} else {
			recursion(*dest.Next, *src.Next, assert)
		}
	}

	testRun(t, func(assert *assert.Assertions) {
		var dest Dest
		src := Src{srcInner}
		err := kast.StructToStruct(src, &dest)
		assert.Nil(err)
		srcInner := src.Inner
		destInner := dest.Inner
		recursion(srcInner, destInner, assert)
	})
}

func Test_s2s_SameType(t *testing.T) {
	type Src struct {
		Name *string
	}
	type Dest struct {
		Name *string
	}
	name := "lxq"
	src := Src{
		Name: &name,
	}
	testRun(t, func(assert *assert.Assertions) {
		var dest Dest
		err := kast.StructToStruct(src, &dest)
		assert.Nil(err)
		assert.Equal(src.Name, dest.Name)
	})
}

func Test_s2s_CustomConvert(t *testing.T) {
	type customSrcType int
	type customDstType string
	type Src struct {
		Custom customSrcType
	}
	type Dst struct {
		Custom customDstType
	}

	kast.RegisterCustomConvertFunc(
		reflect.TypeOf(customSrcType(0)),
		reflect.TypeOf(customDstType("")), func(dest, src reflect.Value) error {
			if src.Int() > 100 {
				dest.SetString("custom 100")
			} else {
				dest.SetString("custom 200")
			}
			return nil
		})

	testRun(t, func(assert *assert.Assertions) {
		var dest Dst
		err := kast.StructToStruct(Src{
			Custom: 198,
		}, &dest)
		assert.Nil(err)
		assert.EqualValues("custom 100", dest.Custom)
		err = kast.StructToStruct(Src{
			Custom: 98,
		}, &dest)
		assert.Nil(err)
		assert.EqualValues("custom 200", dest.Custom)
	})
}
