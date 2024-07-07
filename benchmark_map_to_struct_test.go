package kast_test

import (
	"testing"

	"github.com/wln32/kast"
)

func Benchmark_map_to_struct8(b *testing.B) {
	b.ReportAllocs()
	structPointer8 := &testStructType8{}
	for i := 0; i < b.N; i++ {
		kast.MapToStruct(mapFields8, structPointer8)
	}
}
