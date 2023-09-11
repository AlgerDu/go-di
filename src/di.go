package di

import "reflect"

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
