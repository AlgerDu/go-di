package di

import (
	"fmt"
	"reflect"
)

var (
	canInjectTypeKind map[reflect.Kind]bool = map[reflect.Kind]bool{
		reflect.Interface: true,
		reflect.Struct:    true,
	}
)

func isTypeCanInject(t reflect.Type) (bool, error) {

	k := t.Kind()
	if k == reflect.Interface {
		return true, nil
	}

	if k != reflect.Pointer {
		return false, fmt.Errorf("use [*%s] instead of [%s]", t.Name(), t.Name())
	}

	k = t.Elem().Kind()
	v, exist := canInjectTypeKind[k]
	if !exist {
		return false, fmt.Errorf("not support for [%s]", t.Name())
	}
	return v, nil
}

func isValidCreator(t reflect.Type) (bool, error) {

	if t.Kind() != reflect.Func {
		return false, fmt.Errorf("use func instead of [%s]", t.Name())
	}

	outCount := t.NumOut()
	if outCount != 1 {
		return false, fmt.Errorf("func creator should has 1 return, curr is %d", outCount)
	}

	outType := t.Out(0)
	canInject, err := isTypeCanInject(outType)
	if !canInject {
		return false, fmt.Errorf("func creator return can inject, %s", err)
	}

	return true, nil
}
