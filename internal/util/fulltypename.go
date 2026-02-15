package util

import (
	"reflect"
	"strings"
)

func FullTypeName(t reflect.Type) string {

	if t == nil {
		return "nil"
	}

	switch t.Kind() {
	case reflect.Array:
		return "[]" + FullTypeName(t.Elem())
	case reflect.Slice:
		return "[]" + FullTypeName(t.Elem())
	case reflect.Func:
		if t.PkgPath() != "" {
			path := strings.Split(t.PkgPath(), "/")
			return path[len(path)-1] + "." + t.Name()
		} else {
			return "func(...)"
		}

	case reflect.Interface:
		if t.Name() == "" {
			return "interface{...}"
		} else {
			path := strings.Split(t.PkgPath(), "/")
			return path[len(path)-1] + "." + t.Name()
		}
	case reflect.Struct:
		if t.Name() == "" {
			return "struct{...}"
		} else {
			path := strings.Split(t.PkgPath(), "/")
			return path[len(path)-1] + "." + t.Name()
		}
	case reflect.Map:
		return "map[" + FullTypeName(t.Key()) + "]" + FullTypeName(t.Elem())
	case reflect.Pointer:
		return "*" + FullTypeName(t.Elem())
	default:
		return t.Name()
	}
}
