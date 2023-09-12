package exts

import (
	"reflect"
	"testing"
)

type (
	student struct {
	}
)

func SetStudentName(s *student, name string) {

}

func TestReflectExts_Key(t *testing.T) {

	tp := reflect.TypeOf(student{})
	key := Reflect_GetTypeKey(tp)

	desire := "github.com/AlgerDu/go-di/src/exts/student"
	if key != desire {
		t.Errorf("key is [%s], not [%s]", key, desire)
	}
}

func TestReflectExts_Param(t *testing.T) {
	tt := reflect.TypeOf(SetStudentName)
	keys := Reflect_GetFuncParamKeys(tt)
	t.Log(keys)
}
