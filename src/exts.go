package di

import "reflect"

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

func Collector_AddScope[T any](services ServiceCollector, creator any) error {

	t := new(T)

	return services.AddService(&ServiceDescriptor{
		LifeTime:    SL_Scoped,
		Type:        reflect.TypeOf(t),
		DstType:     nil,
		Instance:    reflect.Value{},
		Creator:     reflect.Value{},
		hasInstance: false,
	})

}
