package di

import (
	"errors"
	"reflect"
)

var (
	ErrEmptyDescriptor = errors.New("empty descrptor")
)

type (
	ServiceCollector interface {
		AddService(descriptor *ServiceDescriptor) error
	}

	ServiceProvider interface {
		GetService(serviceType reflect.Type) (reflect.Value, error)
	}

	Scope interface {
		ServiceProvider

		CreateSubScope(id string, options ...func(ServiceCollector)) (Scope, error)
		GetSubScope(id string) (Scope, bool)
	}
)

type (
	boxState int

	box interface {
		GetID() uint64
		GetInstance(toGetTypeID string, dependPath ...string) (reflect.Value, error)
	}
)

const (
	bs_Empty = iota
	bs_Filling
	bs_OK
)
