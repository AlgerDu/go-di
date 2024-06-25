package di

import (
	"reflect"
	"testing"
)

func TestScope_Base(t *testing.T) {

	scope := newInnerScope("root", nil)
	Collector_AddScope(scope, newStudent)

	s, err := Provider_GetService[*student](scope)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(s.ID, s.Name)
}

func TestScope_Singleton(t *testing.T) {

	scope := newInnerScope("root", nil)
	Collector_AddSingleton(scope, newStudent)

	s1, err := Provider_GetService[*student](scope)
	if err != nil {
		t.Error(err.Error())
	}

	s2, err := Provider_GetService[*student](scope)
	if err != nil {
		t.Error(err.Error())
	}

	if s1 != s2 {
		t.Errorf("addr not equal for Singleton")
	}
}

func TestScope_Singleton2(t *testing.T) {

	scope := newInnerScope("root", nil)
	Collector_AddSingletonFor[reader](scope, newStudent)

	r, err := Provider_GetService[reader](scope)
	if err != nil {
		t.Error(err.Error())
	}

	r.Read(&book{})
}

func TestScope_CreateSub(t *testing.T) {

	root := newInnerScope("root", nil)
	Collector_AddSingleton(root, newStudent)

	subScope, err := root.CreateSubScope("t")
	if err != nil {
		t.Error(err)
	}

	s1, err := Provider_GetService[*student](root)
	if err != nil {
		t.Error(err.Error())
	}

	s2, err := Provider_GetService[*student](subScope)
	if err != nil {
		t.Error(err.Error())
	}

	if s1.ID == s2.ID {
		t.Errorf("err")
	}

	t.Log(s1.ID)
	t.Log(s2.ID)
}

func TestScope_Slice(t *testing.T) {
	container := New()
	Collector_AddSingletonFor[bookStore](container, newABookStore)
	Collector_AddSingletonFor[bookStore](container, newBBookStore)

	stores, err := Provider_GetService[[]bookStore](container)
	if err != nil {
		t.Error(err)
	}

	t.Log(len(stores))
	t.Log(reflect.TypeOf(stores[0]).Elem().Name())
	t.Log(reflect.TypeOf(stores[1]).Elem().Name())
}

func TestScope_Slice2(t *testing.T) {
	container := New()
	AddSingletonFor[bookStore](container, newABookStore)
	AddSingletonFor[bookStore](container, newBBookStore)

	stores, err := Provider_GetService[[]bookStore](container)
	if err != nil {
		t.Error(err)
	}

	t.Log(len(stores))

	store, err := Provider_GetService[*aBookStore](container)
	if err != nil {
		t.Error(err)
	}

	t.Log(store)
}

func TestScope_GetScope(t *testing.T) {
	container := New()

	scope, err := GetService[Scope](container)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(reflect.ValueOf(scope).Type().Name())
	innerScope, ok := scope.(*innerScope)
	if ok {
		t.Log(innerScope.ID)
	}

}
