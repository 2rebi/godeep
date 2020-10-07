package godeep

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	type empty struct {}
	// nil
	assert.Error(t, Copy(nil, &empty{}))
	assert.Error(t, Copy(&empty{}, nil))
	// uint
	assert.Error(t, Copy(uint(0), &empty{}))
	assert.Error(t, Copy(uint8(0), &empty{}))
	assert.Error(t, Copy(uint16(0), &empty{}))
	assert.Error(t, Copy(uint32(0), &empty{}))
	assert.Error(t, Copy(uint64(0), &empty{}))
	assert.Error(t, Copy(&empty{}, uint(0)))
	assert.Error(t, Copy(&empty{}, uint8(0)))
	assert.Error(t, Copy(&empty{}, uint16(0)))
	assert.Error(t, Copy(&empty{}, uint32(0)))
	assert.Error(t, Copy(&empty{}, uint64(0)))
	// int
	assert.Error(t, Copy(int(0), &empty{}))
	assert.Error(t, Copy(int8(0), &empty{}))
	assert.Error(t, Copy(int16(0), &empty{}))
	assert.Error(t, Copy(int32(0), &empty{}))
	assert.Error(t, Copy(int64(0), &empty{}))
	assert.Error(t, Copy(&empty{}, int(0)))
	assert.Error(t, Copy(&empty{}, int8(0)))
	assert.Error(t, Copy(&empty{}, int16(0)))
	assert.Error(t, Copy(&empty{}, int32(0)))
	assert.Error(t, Copy(&empty{}, int64(0)))
	// float
	assert.Error(t, Copy(float32(0), &empty{}))
	assert.Error(t, Copy(float64(0), &empty{}))
	assert.Error(t, Copy(&empty{}, float32(0)))
	assert.Error(t, Copy(&empty{}, float64(0)))
	//complex
	assert.Error(t, Copy(complex64(0), &empty{}))
	assert.Error(t, Copy(complex128(0), &empty{}))
	assert.Error(t, Copy(&empty{}, complex64(0)))
	assert.Error(t, Copy(&empty{}, complex128(0)))
	//not ptr
	assert.Error(t, Copy(empty{}, &empty{}))
	assert.Error(t, Copy(&empty{}, empty{}))

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
		val := 0
		fieldClone(reflect.ValueOf(&val))
	})
	assert.NotPanics(t, func() {
		val := func() {}
		fieldClone(reflect.ValueOf(&val))
	})
}

func TestCpy(t *testing.T) {

	type dst struct {
		StringField  string
		StringField1 string   `from:"unexportField"`
		StringField2 string   `from:"-"`
		StringField3 *string  `from:"unexportField,deep"`
		StringField4 **string `from:"ptrStringField,deep"`

		NumberField int
		NumberField1 int8 `from:"NumberField"`
		NumberField2 int16 `from:"NumberField"`
		NumberField3 *int32 `from:"NumberField,deep"`
		NumberField4 **int64 `from:"NumberField,deep"`

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
	assert.Equal(t, *d.StringField3, s.unexportField)
	assert.Equal(t, **d.StringField4, *s.ptrStringField)

	assert.Equal(t, d.NumberField, s.NumberField)
	assert.Equal(t, d.NumberField1, int8(s.NumberField))
	assert.Equal(t, d.NumberField2, int16(s.NumberField))
	assert.Equal(t, *d.NumberField3, int32(s.NumberField))
	assert.Equal(t, **d.NumberField4, int64(s.NumberField))


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