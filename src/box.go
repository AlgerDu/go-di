package di

import (
	"reflect"
	"sync"
)

type (
	boxState int

	box struct {
		ID       string
		LifeTime ServiceLifetime
		Instance reflect.Value
		Creators []*ServiceDescriptor
		Scope    *innerScope
		State    boxState
		DstType  reflect.Type

		isSlice bool
		lock    sync.Mutex
	}
)

const (
	bs_Empty = iota
	bs_Filling
	bs_OK
)

func newBox(id string, scope *innerScope, dstType reflect.Type) *box {

	isSlice := dstType.Kind() == reflect.Slice

	return &box{
		ID:       id,
		LifeTime: SL_Unknown,
		Instance: reflect.Value{},
		Creators: []*ServiceDescriptor{},
		Scope:    scope,
		State:    bs_Empty,
		DstType:  dstType,
		isSlice:  isSlice,
	}
}

func (box *box) GetInstance(dependPath ...string) (reflect.Value, error) {

	if box.State == bs_OK {
		return box.Instance, nil
	}

	box.lock.Lock()
	defer box.lock.Unlock()

	if box.State == bs_OK {
		return box.Instance, nil
	}

	box.State = bs_Filling

	dependPath = append(dependPath, box.ID)

	var err error

	if box.isSlice {
		err = box.Scope.fillSliceBox(dependPath, box)
	} else {
		err = box.Scope.fillBox(dependPath, box)
	}

	if err != nil {
		box.State = bs_Empty
		return reflect.Value{}, err
	}

	box.State = bs_OK
	return box.Instance, nil
}

func (box *box) CanntFill() bool {
	return len(box.Creators) <= 0
}
