package di

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type (
	emptyBox struct {
		id string
	}
)

func newEmptyBox(id string) *emptyBox {
	return &emptyBox{
		id: id,
	}
}

func (box *emptyBox) GetID() uint64 { return 0 }

func (box *emptyBox) GetInstance(toGetTypeID string, dependPath ...string) (reflect.Value, error) {

	slices.Reverse(dependPath)

	return reflect.Value{}, fmt.Errorf("[%s] is not inject.\n%s", toGetTypeID, strings.Join(dependPath, "\n"))
}
