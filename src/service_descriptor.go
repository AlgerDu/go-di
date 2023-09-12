package di

import (
	"reflect"

	"github.com/AlgerDu/go-di/src/exts"
)

type ServiceDescriptor struct {
	LifeTime ServiceLifetime

	Type    reflect.Type
	DstType reflect.Type

	Instance reflect.Value
	Creator  reflect.Value

	hasInstance bool
}

func (descriptor *ServiceDescriptor) IsSuport(id string) bool {
	return exts.Reflect_GetTypeKey(descriptor.Type) == id || exts.Reflect_GetTypeKey(descriptor.DstType) == id
}
