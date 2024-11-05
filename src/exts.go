package di

import (
	"reflect"
)

func handleToSelef(
	descriptor *ServiceDescriptor,
	toSelf ...bool,
) *ServiceDescriptor {

	supportTypes := []reflect.Type{descriptor.DstType}
	if len(toSelf) > 0 && toSelf[0] {
		supportTypes = append(supportTypes, descriptor.Type)
	}

	descriptor.SupportTypes = supportTypes

	return descriptor
}

func AddSingleton(
	services ServiceCollector,
	creator any,
) error {
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	return services.AddService(&ServiceDescriptor{
		LifeTime:     SL_Singleton,
		Type:         insType,
		DstType:      insType,
		SupportTypes: []reflect.Type{insType},
		Instance:     reflect.Value{},
		Creator:      reflect.ValueOf(creator),
		hasInstance:  false,
	})
}

func AddSingletonFor[forT any](
	services ServiceCollector,
	creator any,
	toSelf ...bool,
) error {

	forType := reflect.TypeOf(new(forT)).Elem()
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	descriptor := &ServiceDescriptor{
		LifeTime:    SL_Singleton,
		Type:        insType,
		DstType:     forType,
		Instance:    reflect.Value{},
		Creator:     reflect.ValueOf(creator),
		hasInstance: false,
	}

	return services.AddService(handleToSelef(descriptor, toSelf...))
}

func AddScope(
	services ServiceCollector,
	creator any,
) error {
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	return services.AddService(&ServiceDescriptor{
		LifeTime:     SL_Scoped,
		Type:         insType,
		DstType:      insType,
		SupportTypes: []reflect.Type{insType},
		Instance:     reflect.Value{},
		Creator:      reflect.ValueOf(creator),
		hasInstance:  false,
	})
}

func AddScopeFor[forT any](
	services ServiceCollector,
	creator any,
	toSelf ...bool,
) error {

	forType := reflect.TypeOf(new(forT)).Elem()
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	descriptor := &ServiceDescriptor{
		LifeTime:    SL_Scoped,
		Type:        insType,
		DstType:     forType,
		Instance:    reflect.Value{},
		Creator:     reflect.ValueOf(creator),
		hasInstance: false,
	}
	return services.AddService(handleToSelef(descriptor, toSelf...))
}

func AddInstance(
	services ServiceCollector,
	ins any,
) error {

	insType := reflect.TypeOf(ins)

	return services.AddService(&ServiceDescriptor{
		LifeTime:     SL_Singleton,
		Type:         insType,
		DstType:      insType,
		SupportTypes: []reflect.Type{insType},
		Instance:     reflect.ValueOf(ins),
		hasInstance:  true,
	})
}

func AddInstanceFor[insType any, dstType any](
	services ServiceCollector,
	ins insType,
	toSelf ...bool,
) error {

	tmpDst := new(dstType)

	descriptor := &ServiceDescriptor{
		LifeTime:    SL_Singleton,
		Type:        reflect.TypeOf(ins),
		DstType:     reflect.TypeOf(tmpDst).Elem(),
		Instance:    reflect.ValueOf(ins),
		hasInstance: true,
	}

	return services.AddService(handleToSelef(descriptor, toSelf...))
}

func GetService[ServiceType any](provider ServiceProvider) (ServiceType, error) {
	newService := new(ServiceType)
	serviceType := reflect.TypeOf(newService).Elem()

	serviceValue, err := provider.GetService(serviceType)
	if err != nil {
		return *newService, err
	}

	service := serviceValue.Interface().(ServiceType)
	return service, nil
}
