package kast_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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
	"age":         98,
	"id":          199,
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
		"name": "lxq",
		"age":  98,
	}
	type Dest struct {
		Name string
		Age  int
	}
	testRun(t, func(assert *assert.Assertions) {
		var dest Dest
		err := kast.MapToStruct(m, &dest)
		assert.Nil(err)
		assert.EqualValues(m["name"], dest.Name)
		assert.EqualValues(m["age"], dest.Age)
	})
}

func Test_map_to_struct8(t *testing.T) {
	testRun(t, func(assert *assert.Assertions) {
		structPointer8 := &testStructType8{}
		err := kast.MapToStruct(mapFields8, structPointer8)
		assert.Nil(err)
		assert.EqualValues(mapFields8["name"], structPointer8.Name)
		assert.EqualValues(mapFields8["score"], structPointer8.Score)
		assert.EqualValues(mapFields8["age"], structPointer8.Age)
		assert.EqualValues(mapFields8["id"], structPointer8.ID)
		assert.EqualValues(mapFields8["categoryId"], structPointer8.CategoryId)
		assert.EqualValues(mapFields8["price"], structPointer8.Price)
		assert.EqualValues(mapFields8["code"], structPointer8.Code)
		assert.EqualValues(mapFields8["image"], structPointer8.Image)
		assert.EqualValues(mapFields8["description"], structPointer8.Description)
		assert.EqualValues(mapFields8["status"], structPointer8.Status)
		assert.EqualValues(mapFields8["idType"], structPointer8.IdType)
	})
}

func Test_map_to_struct_Duplicate_field(t *testing.T) {
	m := map[string]any{
		"ID":   100,
		"Name": "lxq",
	}
	type Nested1 struct {
		ID string
	}
	type Nested2 struct {
		ID uint
	}
	type Dest struct {
		ID   int
		Name string
		Nested1
		Nested2
	}
	testRun(t, func(assert *assert.Assertions) {
		var dest Dest
		err := kast.MapToStruct(m, &dest)
		assert.Nil(err)
		assert.Equal(m["ID"], dest.ID)
		assert.Equal(m["Name"], dest.Name)
		assert.EqualValues(strconv.Itoa(m["ID"].(int)), dest.Nested1.ID)
		assert.EqualValues(m["ID"], dest.Nested2.ID)
	})
}
