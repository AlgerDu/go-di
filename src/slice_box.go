package di

import (
	"reflect"
	"sync"

	"github.com/AlgerDu/go-di/src/exts"
)

type (
	sliceBox struct {
		ID          uint64
		Descriptors []*ServiceDescriptor
		DstType     reflect.Type
		Scope       *innerScope

		State    boxState
		Instance reflect.Value

		lock sync.Mutex
	}
)

func newSliceBox(
	descriptors []*ServiceDescriptor,
	scope *innerScope,
	dstType reflect.Type,
) *sliceBox {

	id := getServiceDescriptorID()

	return &sliceBox{
		ID:          id,
		Scope:       scope,
		Descriptors: descriptors,
		DstType:     dstType,
		State:       bs_Empty,
		Instance:    reflect.Value{},
		lock:        sync.Mutex{},
	}
}

func (box *sliceBox) GetID() uint64 {
	return box.ID
}

func (box *sliceBox) GetInstance(toGetTypeID string, dependPath ...string) (reflect.Value, error) {

	if box.State == bs_OK {
		return box.Instance, nil
	}

	box.lock.Lock()
	defer box.lock.Unlock()

	if box.State == bs_OK {
		return box.Instance, nil
	}

	dependPath = append(dependPath, toGetTypeID)

	sliceValue := reflect.MakeSlice(box.DstType, 0, 1)
	for _, descriptor := range box.Descriptors {
		ins, err := box.Scope.
			FindOrCreateBoxByDescriptor(descriptor).
			GetInstance(exts.Reflect_GetTypeKey(descriptor.Type), dependPath...)
		if err != nil {
			return reflect.Value{}, err
		}
		sliceValue = reflect.Append(sliceValue, ins)
	}

	box.Instance = sliceValue
	box.State = bs_OK

	return sliceValue, nil
}
