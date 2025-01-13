package di

import (
	"reflect"
	"sync"
)

type (
	structBox struct {
		Descriptor *ServiceDescriptor
		Scope      *innerScope

		State    boxState
		Instance reflect.Value

		lock sync.Mutex
	}
)

func newStructBox(
	descriptor *ServiceDescriptor,
	scope *innerScope,
) *structBox {
	return &structBox{
		Descriptor: descriptor,
		Scope:      scope,
		State:      bs_Empty,
		Instance:   reflect.Value{},
		lock:       sync.Mutex{},
	}
}

func (box *structBox) GetID() uint64 {
	return box.Descriptor.id
}

func (box *structBox) GetInstance(toGetTypeID string, dependPath ...string) (reflect.Value, error) {

	if box.Descriptor.LifeTime == SL_Transient {
		return createInstance(toGetTypeID, box.Scope, box.Descriptor.Creator, dependPath)
	}

	if box.State == bs_OK {
		return box.Instance, nil
	}

	box.lock.Lock()
	defer box.lock.Unlock()

	if box.State == bs_OK {
		return box.Instance, nil
	}

	box.State = bs_Filling
	dependPath = append(dependPath, toGetTypeID)

	if box.Descriptor.LifeTime == SL_Singleton &&
		box.Scope != box.Descriptor.belongScope {

		parentBox := box.Descriptor.belongScope.FindOrCreateBoxByDescriptor(box.Descriptor)
		instance, err := parentBox.GetInstance(toGetTypeID, dependPath...)
		if err != nil {
			return reflect.Value{}, err
		}

		box.Instance = instance
	} else if box.Descriptor.hasInstance {
		box.Instance = box.Descriptor.Instance
	} else {

		instance, err := createInstance(toGetTypeID, box.Scope, box.Descriptor.Creator, dependPath)
		if err != nil {
			return reflect.Value{}, err
		}
		box.Instance = instance
	}

	box.State = bs_OK
	return box.Instance, nil
}
