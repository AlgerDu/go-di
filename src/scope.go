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
		return err
	}

	canInject, err = isTypeCanInject(descriptor.DstType)
	if !canInject {
		return err
	}

	if !descriptor.hasInstance {
		validCreator, err := isValidCreator(descriptor.Creator.Type())
		if !validCreator {
			return err
		}
	}

	scope.Descriptors = append(scope.Descriptors, descriptor)
	return nil
}

func (scope *innerScope) GetService(serviceType reflect.Type) (reflect.Value, error) {
	return scope.findOrCreateBox(serviceType).GetInstance()
}

func (scope *innerScope) findOrCreateBox(serviceType reflect.Type) *box {

	id := exts.Reflect_GetTypeKey(serviceType)

	scope.creatingBox.Lock()
	defer scope.creatingBox.Unlock()

	box, exist := scope.Boxs[id]
	if exist {
		return box
	}

	if scope.Parent != nil {
		box = scope.Parent.findOrCreateBox(serviceType)
	}

	scopeDescriptors := []*ServiceDescriptor{}
	for _, descriptor := range scope.Descriptors {
		if descriptor.IsSuport(id) {
			scopeDescriptors = append(scopeDescriptors, descriptor)
		}
	}

	scopeCount := len(scopeDescriptors)
	if scopeCount == 0 && box != nil && box.LifeTime == SL_Singleton {
		scope.Boxs[id] = box
		return box
	}

	scopeBox := newBox(id, scope, serviceType)
	if box != nil {
		scopeBox.Creators = append(scopeBox.Creators, box.Creators...)
	}
	scopeBox.Creators = append(scopeBox.Creators, scopeDescriptors...)
	if scopeCount > 0 {
		scopeBox.LifeTime = scopeDescriptors[scopeCount-1].LifeTime
	} else {
		scopeBox.LifeTime = SL_Scoped
	}

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

	dependTypes := exts.Reflect_GetFuncParam(lastCreater.Creator.Type())

	for _, dependType := range dependTypes {
		for _, path := range dependPath {
			id := exts.Reflect_GetTypeKey(dependType)
			if path == id {
				return fmt.Errorf("cycle depend for %s", box.ID)
			}
		}
	}

	inValues := []reflect.Value{}
	for _, dependType := range dependTypes {
		dependBox := scope.findOrCreateBox(dependType)
		dependValue, err := dependBox.GetInstance(dependPath...)
		if err != nil {
			return err
		}
		inValues = append(inValues, dependValue)
	}

	outValues := lastCreater.Creator.Call(inValues)
	box.Instance = outValues[0]

	return nil
}

func (scope *innerScope) fillSliceBox(dependPath []string, box *box) error {

	// TODO 这里先不处理支持了
	creators := box.Creators
	elemType := box.DstType.Elem()

	if len(creators) == 0 {
		elemBox := scope.findOrCreateBox(elemType)

		creators = elemBox.Creators
	}

	sliceValue := reflect.MakeSlice(box.DstType, 0, 1)

	for _, elemCreator := range creators {

		if elemCreator.hasInstance {
			reflect.AppendSlice(sliceValue, elemCreator.Instance)
			continue
		}

		dependTypes := exts.Reflect_GetFuncParam(elemCreator.Creator.Type())

		for _, dependType := range dependTypes {
			for _, path := range dependPath {
				id := exts.Reflect_GetTypeKey(dependType)
				if path == id {
					return fmt.Errorf("cycle depend for %s", box.ID)
				}
			}
		}

		inValues := []reflect.Value{}
		for _, dependType := range dependTypes {
			dependBox := scope.findOrCreateBox(dependType)
			dependValue, err := dependBox.GetInstance(dependPath...)
			if err != nil {
				return err
			}
			inValues = append(inValues, dependValue)
		}

		outValues := elemCreator.Creator.Call(inValues)
		sliceValue = reflect.Append(sliceValue, outValues[0])
	}

	box.Instance = sliceValue
	return nil
}

func (scope *innerScope) fmtSubScopeID(id string) string {
	return fmt.Sprintf("%s.%s", scope.ID, id)
}
