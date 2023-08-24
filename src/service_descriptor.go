package di

import "reflect"

type ServiceDescriptor struct {
	LifeTime ServiceLifetime

	Type    reflect.Type
	DstType reflect.Type

	Instance reflect.Value
	Creator  reflect.Value
}

func (descriptor *ServiceDescriptor) IsSuport(id string) bool {
	return descriptor.Type.Name() == id || descriptor.DstType.Name() == id
}
