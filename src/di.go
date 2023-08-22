package di

import "reflect"

type (
	ServiceCollector interface {
		AddService(ServiceDescriptor) error
	}

	ServiceProvider interface {
		GetService(serviceType reflect.Type) (reflect.Value, error)
	}

	Scope interface {
		ServiceProvider

		CreateScope(options ...func(ServiceCollector)) (Scope, error)
	}
)
