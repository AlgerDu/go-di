package di

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/AlgerDu/go-di/src/exts"
)

func createInstance(
	toCreateTypeID string,
	scope *innerScope,
	creater *ServiceDescriptor,
	dependPath []string,
) (reflect.Value, error) {
	dependTypes := exts.Reflect_GetFuncParam(creater.Creator.Type())

	for _, dependType := range dependTypes {
		for _, path := range dependPath {
			id := exts.Reflect_GetTypeKey(dependType)
			if path == id {
				return reflect.Value{}, fmt.Errorf("cycle depend for [%s].\n%s", id, strings.Join(dependPath, "\n"))
			}
		}
	}

	inValues := []reflect.Value{}
	for _, dependType := range dependTypes {
		dependBox := scope.FindOrCreateBox(dependType)
		dependValue, err := dependBox.GetInstance(exts.Reflect_GetTypeKey(dependType), dependPath...)
		if err != nil {
			return reflect.Value{}, err
		}
		inValues = append(inValues, dependValue)
	}

	outValues := creater.Creator.Call(inValues)
	var err error
	if len(outValues) == 2 {

		errReturn, ok := outValues[1].Interface().(error)
		if !exts.Reflect_IsNil(outValues[1]) && !ok {
			return reflect.Value{}, fmt.Errorf("func creator sencond return value only support error, current is [%s]", outValues[1].Type().Name())
		}

		if errReturn != nil {
			err = fmt.Errorf("creator returns err for [%s].\n%s", toCreateTypeID, errReturn)
		}
	}

	return outValues[0], err
}
