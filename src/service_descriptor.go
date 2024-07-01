package di

import (
	"reflect"
	"sync/atomic"

	"github.com/AlgerDu/go-di/src/exts"
)

var (
	serviceDescriptorStartID uint64 = 0
)

type ServiceDescriptor struct {
	LifeTime ServiceLifetime

	Type         reflect.Type
	DstType      reflect.Type
	SupportTypes []reflect.Type

	Instance reflect.Value
	Creator  reflect.Value

	belongScope    *innerScope
	hasInstance    bool
	supportTypeIDs []string
	id             uint64
}

func (descriptor *ServiceDescriptor) IsSuport(id string) bool {
	for _, supportID := range descriptor.supportTypeIDs {
		if supportID == id {
			return true
		}
	}
	return false
}

func copyDescriptor(
	descriptor *ServiceDescriptor,
	belongScope *innerScope,
) *ServiceDescriptor {

	supportTypeIDs := []string{}
	for _, supportType := range descriptor.SupportTypes {
		supportTypeIDs = append(supportTypeIDs, exts.Reflect_GetTypeKey(supportType))
	}

	return &ServiceDescriptor{
		LifeTime:     descriptor.LifeTime,
		Type:         descriptor.Type,
		DstType:      descriptor.DstType,
		SupportTypes: descriptor.SupportTypes,
		Instance:     descriptor.Instance,
		Creator:      descriptor.Creator,
		hasInstance:  descriptor.hasInstance,

		id:             getServiceDescriptorID(),
		belongScope:    belongScope,
		supportTypeIDs: supportTypeIDs,
	}
}

func getServiceDescriptorID() uint64 {
	return atomic.AddUint64(&serviceDescriptorStartID, 1)
}
