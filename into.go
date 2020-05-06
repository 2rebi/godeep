package into

import (
	"errors"
	"reflect"
)

const (
	tagNameFrom = "into"
)

func Into(src, dst interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	} else if dst == nil {
		return errors.New("dst is nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return errors.New("destination must be pointer of struct")
	}
	into(reflect.ValueOf(src), dstVal.Elem())
	return nil
}


func into(src, dst reflect.Value)  {
	if !dst.CanSet() {
		return
	}

	if src.Kind() == reflect.Ptr {
		into(src.Elem(), dst)
	}

	dstKind := dst.Kind()
	dstType := dst.Type()

	if dstType == src.Type() {
		switch src.Kind() {
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
		//TODO case reflect.Interface:
		}
	} else {
		switch dstKind {
		case reflect.Ptr:
			elem := dst.Elem()
			if elem.Kind() == reflect.Invalid {
				dst.Set(reflect.New(dstType.Elem()))
				elem = dst.Elem()
			}
			into(src, elem)
		case reflect.Struct:
			for i, cnt := 0, dstType.NumField(); i < cnt; i++ {
				field := dstType.Field(i)
				key := ""
				if tag, ok := field.Tag.Lookup(tagNameFrom); ok {
					if tag == "-" {
						continue
					}
					key = tag
				} else {
					key = field.Name
				}

				srcField := src.FieldByName(key)

				if !srcField.IsValid() {
					continue
				}
				into(srcField, dst.Field(i))
			}
		}
		//TODO case reflect.Map:
	}

	return
}