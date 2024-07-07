package kast_test

import (
	"testing"

	"github.com/wln32/kast"
)
type testStructType82 struct {
	Name        string  `json:"name"   `
	CategoryId  string  `json:"categoryId" `
	Price       float64 `json:"price"    `
	Code        string  `json:"code"       `
	Image       string  `json:"image"   `
	Description string  `json:"description" `
	Status      int     `json:"status"   `
	IdType      int     `json:"idType"`
	Score       int
	Age         int
	ID          int
}

type testStructType8 struct {
	Name        string  `json:"name"   `
	CategoryId  string  `json:"categoryId" `
	Price       float64 `json:"price"    `
	Code        string  `json:"code"       `
	Image       string  `json:"image"   `
	Description string  `json:"description" `
	Status      int     `json:"status"   `
	IdType      int     `json:"idType"`
	Score       int
	Age         int
	ID          int
}

var mapFields8 = map[string]interface{}{
	"name":        "github",
	"score":       100,
	"Age":         98,
	"ID":          199,
	"categoryId":  "1",
	"price":       198.09,
	"code":        "1",
	"image":       "https://github.com",
	"description": "This is the data for testing eight fields",
	"status":      1,
	"idType":      2,
}

func TestMapToStruct(t *testing.T) {
	m := map[string]any{
		"name": "wln",
		"age":  98,
	}

	type Dest struct {
		Name string
		Age  int
	}
	var dest Dest
	err := kast.MapToStruct(m, &dest)
	t.Log(err)
	t.Log(dest)
}
