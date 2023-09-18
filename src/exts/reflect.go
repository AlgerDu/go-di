package exts

import (
	"fmt"
	"reflect"
)

func Reflect_GetTypeKey(t reflect.Type) string {

	if t == nil {
		return ""
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
		return fmt.Sprintf("[]%s", Reflect_GetTypeKey(t))
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		return fmt.Sprintf("%s/*%s", t.PkgPath(), t.Name())
	}

	return fmt.Sprintf("%s/%s", t.PkgPath(), t.Name())
}

func Reflect_GetFuncParam(t reflect.Type) []reflect.Type {
	keys := []reflect.Type{}

	inCount := t.NumIn()
	for i := 0; i < inCount; i++ {
		inType := t.In(i)
		keys = append(keys, inType)
	}

	return keys
}

func Reflect_GetFuncParamKeys(t reflect.Type) []string {
	keys := []string{}

	inCount := t.NumIn()
	for i := 0; i < inCount; i++ {
		inType := t.In(i)
		keys = append(keys, Reflect_GetTypeKey(inType))
	}

	return keys
}
