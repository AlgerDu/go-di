package di

import "reflect"

type box struct {
	ID       string
	LifeTime ServiceLifetime
	Instance reflect.Value
	Creators []*ServiceDescriptor
	Scope    *innerScope
}

func newBox(id string, scope *innerScope) *box {
	return &box{
		ID:       id,
		LifeTime: 0,
		Instance: reflect.Value{},
		Creators: []*ServiceDescriptor{},
		Scope:    scope,
	}
}
