package kast_test

import (
	"testing"

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

func Test_StructToStruct(t *testing.T) {
	type Src struct {
		Name string
		Age  int
	}

	type Dest struct {
		Name []byte
		Age  uint
	}

	var dest Dest
	err := kast.StructToStruct(Src{
		Name: "wln", Age: 98,
	}, &dest)
	t.Log(err)
	t.Log(dest)
}

func Test_struct_to_struct8(t *testing.T) {

	dest := &testStructType82{}
	err := kast.StructToStruct(srcStructFields8, dest)
	t.Log(err)
	t.Log(dest)
}
