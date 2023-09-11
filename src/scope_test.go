package di

import (
	"testing"
)

func TestScope_Base(t *testing.T) {

	scope := newInnerScope(nil)
	Collector_AddScope(scope, newStudent)

	s, err := Provider_GetService[student](scope)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(s.ID, s.Name)
}

func TestScope_Singleton(t *testing.T) {

	scope := newInnerScope(nil)
	Collector_AddSingleton(scope, newStudent)

	s1, err := Provider_GetService[student](scope)
	if err != nil {
		t.Error(err.Error())
	}

	s2, err := Provider_GetService[student](scope)
	if err != nil {
		t.Error(err.Error())
	}

	if s1 != s2 {
		t.Errorf("addr not equal for Singleton")
	}
}
