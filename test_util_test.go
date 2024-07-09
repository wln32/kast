package kast_test

import (
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

func testRun(t *testing.T, fn func(assert *asserts.Assertions)) {
	fn(asserts.New(t))
}
