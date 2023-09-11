package di

import "sync"

type (
	book struct {
		Content string
	}

	student struct {
		ID   int
		Name string
	}

	reader interface {
		Read(b *book) error
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

func (student *student) Read(b *book) error {
	return nil
}
