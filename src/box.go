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
		lock     sync.Mutex
	}
)

const (
	bs_Empty = iota
	bs_Filling
	bs_OK
)

func newBox(id string, scope *innerScope) *box {
	return &box{
		ID:       id,
		LifeTime: 0,
		Instance: reflect.Value{},
		Creators: []*ServiceDescriptor{},
		Scope:    scope,
		State:    bs_Empty,
	}
}

func (box *box) GetInstance(dependPath []string) (reflect.Value, error) {

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
	err := box.Scope.fillBox(dependPath, box)
	if err != nil {
		box.State = bs_Empty
		return reflect.Value{}, err
	}

	box.State = bs_OK
	return box.Instance, nil
}

func (box *box) CanntFill() bool {
	if len(box.Creators) <= 0 {
		return true
	}
	return false
}
