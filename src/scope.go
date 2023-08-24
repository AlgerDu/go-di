package di

import (
	"reflect"
	"sync"
)

type (
	innerScope struct {
		Parent *innerScope

		Descriptors []*ServiceDescriptor
		Boxs        map[string]*box

		creatingBox sync.Mutex
	}
)

func newInnerScope(
	parent *innerScope,
) *innerScope {

	return &innerScope{
		Parent:      parent,
		Descriptors: []*ServiceDescriptor{},
		Boxs:        map[string]*box{},
		creatingBox: sync.Mutex{},
	}
}

func (scope *innerScope) AddService(descriptor ServiceDescriptor) error {
	panic("not implemented") // TODO: Implement
}

func (scope *innerScope) GetService(serviceType reflect.Type) (reflect.Value, error) {

	return reflect.Value{}, nil
}

func (scope *innerScope) CreateScope(options ...func(ServiceCollector)) (Scope, error) {
	return nil, nil
}

func (scope *innerScope) findOrCreateBox(id string) *box {

	scope.creatingBox.Lock()
	defer scope.creatingBox.Unlock()

	box, exist := scope.Boxs[id]
	if exist {
		return box
	}

	if scope.Parent != nil {
		box = scope.Parent.findOrCreateBox(id)
	}

	scopeDescriptors := []*ServiceDescriptor{}
	for _, descriptor := range scope.Descriptors {
		if descriptor.IsSuport(id) {
			scopeDescriptors = append(scopeDescriptors, descriptor)
		}
	}

	scopeCount := len(scopeDescriptors)
	if scopeCount == 0 && box.LifeTime == SL_Singleton {
		scope.Boxs[id] = box
		return box
	}

	scopeBox := newBox(id, scope)
	scopeBox.Creators = append(box.Creators, scopeDescriptors...)
	scopeBox.LifeTime = scopeDescriptors[scopeCount-1].LifeTime

	scope.Boxs[id] = scopeBox

	return scopeBox
}

func (scope *innerScope) createBoxInstance(dependPath []string, box *box) error {
	return nil
}
