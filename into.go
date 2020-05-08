package godeep

import (
	"math"
	"reflect"
	"strings"
	"unsafe"
)

const (
	tagNameFrom = "from"
	tagValueDeepCopy = "deep"
)

func Copy(src, dst interface{}) {
	if src == nil {
		panic("src is nil")
	} else if dst == nil {
		panic("dst is nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.Type().Elem().Kind() != reflect.Struct {
		panic("destination must be pointer of struct")
	}
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Ptr || srcVal.Type().Elem().Kind() != reflect.Struct {
		panic("destination must be pointer of struct")
	}
	into(srcVal.Elem(), dstVal.Elem(), true,false)
}


func into(src, dst reflect.Value, isExport, isFieldDeepCopy bool)  {
	if !dst.CanSet() {
		return
	}

	dstKind := dst.Kind()
	dstType := dst.Type()
	srcType := src.Type()
	srcKind := src.Kind()

	switch dstKind {
	case reflect.Ptr:
		if !isExport {
			into(fieldClone(src), dst, true, isFieldDeepCopy)
		} else if !isFieldDeepCopy && srcType.AssignableTo(dstType) {
			dst.Set(src)
		} else if !isFieldDeepCopy && srcType.ConvertibleTo(dstType) {
			dst.Set(src.Convert(dstType))
		} else {
			dst.Set(reflect.New(dstType.Elem()))
			into(src, dst.Elem(), isExport, isFieldDeepCopy)
		}
	case reflect.Bool:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
			dst.SetBool(src.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
			dst.SetInt(src.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
			dst.SetUint(src.Uint())
		}
	case reflect.Float32, reflect.Float64:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
			dst.SetFloat(src.Float())
		}
	case reflect.Complex64, reflect.Complex128:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
			dst.SetComplex(src.Complex())
		}
	case reflect.String:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
			dst.SetString(src.String())
		}
	case reflect.Func:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else if !isExport {
			into(fieldClone(src.Addr()), dst, true, isFieldDeepCopy)
		} else {
			dst.Set(src)
		}
	case reflect.Chan:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else if !isExport {
			into(fieldClone(src), dst, true, isFieldDeepCopy)
		} else {
			dst.Set(src)
		}
	case reflect.Map:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else if !isExport {
			into(fieldClone(src), dst, true, isFieldDeepCopy)
		} else if !isFieldDeepCopy && srcType.AssignableTo(dstType) {
			dst.Set(src)
		} else if !isFieldDeepCopy && srcType.ConvertibleTo(dstType) {
			dst.Set(src.Convert(dstType))
		} else if sKeyT, sElemT, dKeyT, dElemT := srcType.Key(), srcType.Elem(), dstType.Key(), dstType.Elem();
		!isFieldDeepCopy && dst.CanAddr() &&
			(sKeyT.Size() == dKeyT.Size() && sKeyT.ConvertibleTo(dKeyT)) &&
			(sElemT.Size() == dElemT.Size() && sElemT.ConvertibleTo(dElemT)) {
			*(*unsafe.Pointer)(unsafe.Pointer(dst.Addr().Pointer())) = unsafe.Pointer(src.Pointer())
		} else {
			keys := src.MapKeys()
			dst.Set(reflect.MakeMapWithSize(dstType, len(keys)))
			keyPtr := reflect.New(dKeyT)
			valPtr := reflect.New(dElemT)
			for i := range keys {
				into(keys[i], keyPtr.Elem(), true, false)
				into(src.MapIndex(keys[i]), valPtr.Elem(), true, false)
				dst.SetMapIndex(keyPtr.Elem(), valPtr.Elem())
			}
		}
	case reflect.Array:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else if !isExport {
			into(fieldClone(src), dst, true, isFieldDeepCopy)
		} else {
			for i, cnt := 0, int(math.Min(float64(dst.Len()), float64(src.Len()))); i < cnt; i++ {
				into(src.Index(i), dst.Index(i), true, false)
			}
		}
	case reflect.Slice:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else if !isExport {
			into(fieldClone(src), dst, true, isFieldDeepCopy)
		} else if !isFieldDeepCopy && srcType.AssignableTo(dstType) {
			dst.Set(src)
		} else if !isFieldDeepCopy && srcType.ConvertibleTo(dstType) {
			dst.Set(src.Convert(dstType))
		} else {
			dst.Set(reflect.MakeSlice(dst.Type(), src.Len(), int(float64(src.Len()) * 1.5)))
			for i, cnt := 0, dst.Len(); i < cnt; i++ {
				into(src.Index(i), dst.Index(i), true, false)
			}
		}
	case reflect.Interface:
		if !srcType.AssignableTo(dstType) && srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else if !isExport {
			into(fieldClone(src), dst, true, isFieldDeepCopy)
		} else {
			dst.Set(src)
		}
	case reflect.Struct:
		if srcKind == reflect.Ptr {
			into(src.Elem(), dst, isExport, isFieldDeepCopy)
		} else {
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

				srcFieldType, ok := srcType.FieldByName(key)
				if !ok {
					continue
				}
				srcField := src.FieldByIndex(srcFieldType.Index)
				into(srcField, dst.Field(i), isExport && (key[0] <= 'Z' && key[0] >= 'A'), isDeepCopy)
			}
		}
	}

	return
}

func fieldClone(v reflect.Value) reflect.Value {
	bowl := reflect.New(v.Type())
	*(*unsafe.Pointer)(unsafe.Pointer(bowl.Pointer())) = unsafe.Pointer(v.Pointer())
	return bowl.Elem()
}