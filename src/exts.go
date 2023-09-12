package di

import (
	"reflect"
)

func Collector_AddInstance[insType any, dstType any](services ServiceCollector, ins insType) error {

	tmpDst := new(dstType)

	return services.AddService(&ServiceDescriptor{
		LifeTime:    SL_Singleton,
		Type:        reflect.TypeOf(ins),
		DstType:     reflect.TypeOf(tmpDst).Elem(),
		Instance:    reflect.ValueOf(ins),
		hasInstance: true,
	})
}

func Collector_AddSingleton(services ServiceCollector, creator any) error {
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	return services.AddService(&ServiceDescriptor{
		LifeTime:    SL_Singleton,
		Type:        insType,
		DstType:     insType,
		Instance:    reflect.Value{},
		Creator:     reflect.ValueOf(creator),
		hasInstance: false,
	})
}

func Collector_AddSingletonFor[forT any](services ServiceCollector, creator any) error {

	forType := reflect.TypeOf(new(forT)).Elem()
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	return services.AddService(&ServiceDescriptor{
		LifeTime:    SL_Scoped,
		Type:        insType,
		DstType:     forType,
		Instance:    reflect.Value{},
		Creator:     reflect.ValueOf(creator),
		hasInstance: false,
	})
}

func Collector_AddScopeFor[forT any](services ServiceCollector, creator any) error {

	forType := reflect.TypeOf(new(forT)).Elem()
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	return services.AddService(&ServiceDescriptor{
		LifeTime:    SL_Scoped,
		Type:        insType,
		DstType:     forType,
		Instance:    reflect.Value{},
		Creator:     reflect.ValueOf(creator),
		hasInstance: false,
	})
}

func Collector_AddScope(services ServiceCollector, creator any) error {
	creatorType := reflect.TypeOf(creator)
	insType := creatorType.Out(0)

	return services.AddService(&ServiceDescriptor{
		LifeTime:    SL_Scoped,
		Type:        insType,
		DstType:     insType,
		Instance:    reflect.Value{},
		Creator:     reflect.ValueOf(creator),
		hasInstance: false,
	})
}

func Provider_GetService[ServiceType any](provider ServiceProvider) (ServiceType, error) {
	newService := new(ServiceType)
	serviceType := reflect.TypeOf(newService).Elem()

	serviceValue, err := provider.GetService(serviceType)
	if err != nil {
		return *newService, err
	}

	service := serviceValue.Interface().(ServiceType)
	return service, nil
}
