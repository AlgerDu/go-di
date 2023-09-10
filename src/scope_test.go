package di

import (
	"testing"
)

func TestScope_Simple(t *testing.T) {

	scope := newInnerScope(nil)
	Collector_AddScope(scope, newStudent)

	s, err := Provider_GetService[student](scope)
	if err != nil {
		t.Logf(err.Error())
	}

	t.Log(s.Name)
}
