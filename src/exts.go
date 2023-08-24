package di

import "reflect"

func Collector_AddInstance[insType any, dstType any](services ServiceCollector, ins insType) error {

	tmpDst := new(dstType)

	return services.AddService(ServiceDescriptor{
		LifeTime: SL_Singleton,
		Type:     reflect.TypeOf(ins),
		DstType:  reflect.TypeOf(tmpDst).Elem(),
	})
}
