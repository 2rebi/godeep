package godeep

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	type empty struct {}
	// nil
	assert.Panics(t, func() {
		Copy(nil, &empty{})
	})
	assert.Panics(t, func() {
		Copy(&empty{}, nil)
	})
	// uint
	assert.Panics(t, func() {
		Copy(uint(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(uint8(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(uint16(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(uint32(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(uint64(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(&empty{}, uint(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, uint8(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, uint16(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, uint32(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, uint64(0))
	})
	// int
	assert.Panics(t, func() {
		Copy(int(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(int8(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(int16(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(int32(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(int64(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(&empty{}, int(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, int8(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, int16(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, int32(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, int64(0))
	})
	// float
	assert.Panics(t, func() {
		Copy(float32(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(float64(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(&empty{}, float32(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, float64(0))
	})
	//complex
	assert.Panics(t, func() {
		Copy(complex64(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(complex128(0), &empty{})
	})
	assert.Panics(t, func() {
		Copy(&empty{}, complex64(0))
	})
	assert.Panics(t, func() {
		Copy(&empty{}, complex128(0))
	})
	//not ptr
	assert.Panics(t, func() {
		Copy(empty{}, &empty{})
	})
	assert.Panics(t, func() {
		Copy(&empty{}, empty{})
	})

	//correct
	assert.NotPanics(t, func() {
		Copy(&empty{}, &empty{})
		Copy(&empty{}, map[string]interface{}{})
	})

}


func TestFieldClone(t *testing.T) {
	// panic
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(uint(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(uint8(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(uint16(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(uint32(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(uint64(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(int(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(int8(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(int16(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(int32(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(int64(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(float32(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(float64(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(complex64(0)))
	})
	assert.Panics(t, func() {
		fieldClone(reflect.ValueOf(complex128(0)))
	})

	//correct
	assert.NotPanics(t, func() {
		val := int(0)
		fieldClone(reflect.ValueOf(&val))
	})
	assert.NotPanics(t, func() {
		val := func() {}
		fieldClone(reflect.ValueOf(&val))
	})
}

func TestCpy(t *testing.T) {

	type dst struct {
		StringField string
		StringField1 string `from:"unexportField"`
		StringField2 string `from:"-"`
		Stringfield3 *string `from:"unexportField,deep"`
		Stringfield4 **string `from:"ptrStringField,deep"`

		NumberField int
		NumberField1 int8 `from:"NumberField"`
		NumberField2 int16 `from:"NumberField"`
		//NumberField3 *int32 `from:"NumberField"`
		//NumberField4 **int64 `from:"NumberField"`

		ActionField func() string `from:"action"`
	}

	type src struct {
		unexportField string
		StringField string
		ptrStringField *string

		NumberField int

		action func() string
	}
	str := "ptr string"

	d := dst{}
	s := src{
		unexportField:  "this is unexported field",
		StringField:    "string string field",
		ptrStringField: &str,
		NumberField:    1234,
		action: func() string {
			return "hello action world"
		},
	}

	assert.NotPanics(t, func() {
		cpy(reflect.ValueOf(&d).Elem(), reflect.ValueOf(&s).Elem(), true,false)
	})

	assert.Equal(t, d.StringField, s.StringField)
	assert.Equal(t, d.StringField1, s.unexportField)
	assert.Zero(t, d.StringField2)
	assert.Equal(t, *d.Stringfield3, s.unexportField)
	assert.Equal(t, **d.Stringfield4, *s.ptrStringField)

	type empty struct {}
	assert.NotPanics(t, func() {
		cpy(reflect.ValueOf(&empty{}), reflect.ValueOf(&empty{}), true,false)
	})
	assert.NotPanics(t, func() {
		cpy(reflect.ValueOf(&empty{}), reflect.ValueOf(&empty{}), true,true)
	})
	assert.NotPanics(t, func() {
		cpy(reflect.ValueOf(&empty{}), reflect.ValueOf(&empty{}), false,true)
	})
}