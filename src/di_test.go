package di

import "sync"

type (
	student struct {
		ID   int
		Name string
	}
)

var (
	studentCount      int        = 0
	studentCreateLock sync.Mutex = sync.Mutex{}
)

func newStudent() *student {

	studentCreateLock.Lock()
	defer studentCreateLock.Unlock()

	studentCount++

	return &student{
		ID:   studentCount,
		Name: "123",
	}
}
