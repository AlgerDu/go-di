package di

import "reflect"

type ServiceDescriptor struct {
	LifeTime ServiceLifetime
	Type     reflect.Type
	DstType  reflect.Type
	Value    reflect.Value
}
