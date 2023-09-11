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
		Boxs        map[string]*box

		creatingBox      sync.Mutex
		creatingSubScope sync.Mutex
	}
)

func newInnerScope(
	id string,
	parent *innerScope,
) *innerScope {

	return &innerScope{
		ID:               id,
		Parent:           parent,
		SucScopes:        map[string]*innerScope{},
		Descriptors:      []*ServiceDescriptor{},
		Boxs:             map[string]*box{},
		creatingBox:      sync.Mutex{},
		creatingSubScope: sync.Mutex{},
	}
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
	for _, option := range options {
		option(subScope)
	}

	scope.SucScopes[id] = subScope

	return nil, nil
}

func (scope *innerScope) GetSubScope(id string) (Scope, bool) {
	id = scope.fmtSubScopeID(id)

	subScope, exist := scope.SucScopes[id]
	return subScope, exist
}

func (scope *innerScope) AddService(descriptor *ServiceDescriptor) error {
	scope.Descriptors = append(scope.Descriptors, descriptor)
	return nil
}

func (scope *innerScope) GetService(serviceType reflect.Type) (reflect.Value, error) {

	id := exts.Reflect_GetTypeKey(serviceType)
	box := scope.findOrCreateBox(id)

	return box.GetInstance([]string{})
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
	if box != nil {
		scopeBox.Creators = append(scopeBox.Creators, box.Creators...)
	}
	scopeBox.Creators = append(scopeBox.Creators, scopeDescriptors...)
	scopeBox.LifeTime = scopeDescriptors[scopeCount-1].LifeTime

	scope.Boxs[id] = scopeBox

	return scopeBox
}

func (scope *innerScope) fillBox(dependPath []string, box *box) error {

	if box.CanntFill() {
		return fmt.Errorf("[%s] is not inject", box.ID)
	}

	lastCreater := box.Creators[len(box.Creators)-1]

	if lastCreater.hasInstance {
		box.Instance = lastCreater.Instance
		return nil
	}

	dependKeys := exts.Reflect_GetFuncParamKeys(lastCreater.Creator.Type())

	for _, dependKey := range dependKeys {
		for _, path := range dependPath {
			if path == dependKey {
				return fmt.Errorf("cycle depend for %s", box.ID)
			}
		}
	}

	inValues := []reflect.Value{}
	for _, dependKey := range dependKeys {
		dependBox := scope.findOrCreateBox(dependKey)
		dependValue, err := dependBox.GetInstance(dependPath)
		if err != nil {
			return err
		}
		inValues = append(inValues, dependValue)
	}

	outValues := lastCreater.Creator.Call(inValues)
	box.Instance = outValues[0]

	return nil
}

func (scope *innerScope) fmtSubScopeID(id string) string {
	return fmt.Sprintf("%s.%s", scope.ID, id)
}
