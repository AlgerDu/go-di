package di

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/AlgerDu/go-di/src/exts"
)

type (
	innerScope struct {
		ID string

		Parent    *innerScope
		SucScopes map[string]*innerScope

		Descriptors []*ServiceDescriptor

		boxs       map[uint64]box
		idToBoxIDs map[string]uint64

		creatingBox      sync.Mutex
		creatingSubScope sync.Mutex
	}
)

func newInnerScope(
	id string,
	parent *innerScope,
) *innerScope {

	scope := &innerScope{
		ID:               id,
		Parent:           parent,
		SucScopes:        map[string]*innerScope{},
		Descriptors:      []*ServiceDescriptor{},
		boxs:             map[uint64]box{},
		idToBoxIDs:       map[string]uint64{},
		creatingBox:      sync.Mutex{},
		creatingSubScope: sync.Mutex{},
	}

	AddInstanceFor[*innerScope, Scope](scope, scope)

	return scope
}

func (scope *innerScope) CreateSubScope(id string, options ...func(ServiceCollector)) (Scope, error) {
	id = scope.fmtSubScopeID(id)

	_, exist := scope.SucScopes[id]
	if exist {
		return nil, fmt.Errorf("scope %s exist", id)
	}

	scope.creatingSubScope.Lock()
	defer scope.creatingSubScope.Unlock()

	_, exist = scope.SucScopes[id]
	if exist {
		return nil, fmt.Errorf("scope %s exist", id)
	}

	subScope := newInnerScope(id, scope)
	subScope.Descriptors = append(subScope.Descriptors, scope.Descriptors...)

	for _, option := range options {
		option(subScope)
	}

	scope.SucScopes[id] = subScope

	return subScope, nil
}

func (scope *innerScope) GetSubScope(id string) (Scope, bool) {
	id = scope.fmtSubScopeID(id)

	subScope, exist := scope.SucScopes[id]
	return subScope, exist
}

func (scope *innerScope) AddService(descriptor *ServiceDescriptor) error {

	canInject, err := isTypeCanInject(descriptor.Type)
	if !canInject {
		panic(err)
	}

	canInject, err = isTypeCanInject(descriptor.DstType)
	if !canInject {
		panic(err)
	}

	if !descriptor.hasInstance {
		validCreator, err := isValidCreator(descriptor.Creator.Type())
		if !validCreator {
			panic(err)
		}
	}

	scope.Descriptors = append(scope.Descriptors, copyDescriptor(descriptor, scope))
	return nil
}

func (scope *innerScope) GetService(serviceType reflect.Type) (reflect.Value, error) {
	return scope.FindOrCreateBox(serviceType).GetInstance(exts.Reflect_GetTypeKey(serviceType))
}

func (scope *innerScope) FindSupportDescriptors(id string) []*ServiceDescriptor {
	descriptors := []*ServiceDescriptor{}

	if scope.Parent != nil {
		descriptors = append(descriptors, scope.Parent.FindSupportDescriptors(id)...)
	}

	for _, descriptor := range scope.Descriptors {
		if descriptor.IsSuport(id) {
			descriptors = append(descriptors, descriptor)
		}
	}

	return descriptors
}

func (scope *innerScope) FindOrCreateBoxByDescriptor(descriptor *ServiceDescriptor) box {

	if box, exist := scope.boxs[descriptor.id]; exist {
		return box
	}

	scope.creatingBox.Lock()
	defer scope.creatingBox.Unlock()

	if box, exist := scope.boxs[descriptor.id]; exist {
		return box
	}

	box := newStructBox(descriptor, scope)
	scope.boxs[box.GetID()] = box

	return box
}

func (scope *innerScope) FindOrCreateBox(serviceType reflect.Type) box {
	id := exts.Reflect_GetTypeKey(serviceType)
	if boxID, exist := scope.idToBoxIDs[id]; exist {
		return scope.boxs[boxID]
	}

	scope.creatingBox.Lock()
	defer scope.creatingBox.Unlock()

	if boxID, exist := scope.idToBoxIDs[id]; exist {
		return scope.boxs[boxID]
	}

	descriptors := scope.FindSupportDescriptors(id)
	if len(descriptors) == 0 {
		return newEmptyBox(id)
	}

	var box box
	if serviceType.Kind() == reflect.Slice {
		box = newSliceBox(descriptors, scope, serviceType)
	} else {
		descriptor := descriptors[0]
		existBox, exist := scope.boxs[descriptor.id]
		if exist {
			box = existBox
		} else {
			box = newStructBox(descriptors[0], scope)
		}
	}

	scope.boxs[box.GetID()] = box
	scope.idToBoxIDs[id] = box.GetID()

	return box
}

func (scope *innerScope) fmtSubScopeID(id string) string {
	return fmt.Sprintf("%s.%s", scope.ID, id)
}
