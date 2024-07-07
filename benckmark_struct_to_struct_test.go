package kast_test

import (
	"testing"

	"github.com/wln32/kast"
)

var benchmarkSrcStructFields8 = testStructType8{
	Name: "gf",

	CategoryId:  "1",
	Price:       198.09,
	Code:        "1",
	Image:       "https://goframe.org",
	Description: "This is the data for testing eight fields",
	Status:      1,
	IdType:      2,

	Score: 100,
	Age:   98,
	ID:    199,
}

func Benchmark_struct_to_struct8(b *testing.B) {
	var destStructMapFields8 = &testStructType82{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		kast.StructToStruct(benchmarkSrcStructFields8, destStructMapFields8)
	}
}
