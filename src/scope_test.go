package di

import (
	"reflect"
	"testing"
)

type (
	student struct {
		Name string
	}
)

func newStudent() *student {
	return &student{
		Name: "123",
	}
}

func TestScope_Simple(t *testing.T) {

	studentType := reflect.TypeOf(&student{})

	scope := newInnerScope(nil)
	scope.AddService(&ServiceDescriptor{
		LifeTime:    SL_Scoped,
		Type:        studentType,
		DstType:     nil,
		Creator:     reflect.ValueOf(newStudent),
		hasInstance: false,
	})

	v, err := scope.GetService(studentType)
	if err != nil {
		t.Logf(err.Error())
	}

	s := v.Interface().(*student)
	t.Log(s.Name)
}
