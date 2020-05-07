package into

import (
	"errors"
	"reflect"
	"strings"
	"unsafe"
)

const (
	tagNameFrom = "into"
	tagValueDeepCopy = "deep"
)

func Into(src, dst interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	} else if dst == nil {
		return errors.New("dst is nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.Type().Elem().Kind() != reflect.Struct {
		return errors.New("destination must be pointer of struct")
	}

	srcVal := reflect.ValueOf(src)
	srcKind := srcVal.Kind()
	if srcKind == reflect.Ptr {
		srcVal = srcVal.Elem()
		srcKind = srcVal.Kind()
	}

	if srcKind != reflect.Struct {
		return errors.New("source must be pointer of struct")
	}

	into(srcVal, dstVal.Elem(), false)
	return nil
}


func into(src, dst reflect.Value, isFieldDeepCopy bool)  {
	if !dst.CanSet() {
		return
	}

	dstKind := dst.Kind()
	dstType := dst.Type()

	if dstType == src.Type() {
		switch src.Kind() {
		case reflect.Ptr:
			if isFieldDeepCopy {
				into(src.Elem(), dst, false)
			} else {
				dst.SetPointer(unsafe.Pointer(src.Pointer()))
			}
		case reflect.Func, //reflect.Interface,
			reflect.Chan:
			dst.SetPointer(unsafe.Pointer(src.Pointer()))
		case reflect.String:
			dst.SetString(src.String())
		case reflect.Bool:
			dst.SetBool(src.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			dst.SetInt(src.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			dst.SetUint(src.Uint())
		case reflect.Float32, reflect.Float64:
			dst.SetFloat(src.Float())
		case reflect.Complex64, reflect.Complex128:
			dst.SetComplex(src.Complex())
		case reflect.Map:
			keys := src.MapKeys()
			dst.Set(reflect.MakeMapWithSize(dstType, len(keys)))
			for i := range keys {
				dst.SetMapIndex(keys[i], src.MapIndex(keys[i]))
			}
		case reflect.Array:
		case reflect.Slice:
		case reflect.Struct:
		}
	} else if dstKind == src.Kind() {
		switch dstKind {
		case reflect.Ptr:
			elem := dst.Elem()
			if elem.Kind() == reflect.Invalid {
				dst.Set(reflect.New(dstType.Elem()))
				elem = dst.Elem()
			}
			into(src, elem, false)
		case reflect.Struct:
			for i, cnt := 0, dstType.NumField(); i < cnt; i++ {
				field := dstType.Field(i)
				key := ""
				isDeepCopy := false
				if tag, ok := field.Tag.Lookup(tagNameFrom); ok {
					if tag == "-" {
						continue
					}
					tagVal := strings.SplitN(tag, ",", 2)
					if len(tagVal) == 2 {
						isDeepCopy = tagVal[1] == tagValueDeepCopy
					}
					key = tagVal[0]
				} else {
					key = field.Name
				}

				srcField := src.FieldByName(key)

				if !srcField.IsValid() {
					continue
				}

				into(srcField, dst.Field(i), isDeepCopy)
			}
		case reflect.Interface:
		case reflect.Map:
		}

	} else {

	}

	return
}
