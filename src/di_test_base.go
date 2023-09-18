package di

import "sync"

type (
	reader interface {
		Read(b *book) error
	}

	bookStore interface {
		Find(name string) *book
	}

	aBookStore struct {
	}

	bBookStore struct {
	}

	book struct {
		Content string
	}

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

func (student *student) Read(b *book) error {
	return nil
}

func newABookStore() *aBookStore {
	return &aBookStore{}
}

func (sotre *aBookStore) Find(name string) *book {
	return &book{}
}

func newBBookStore() *bBookStore {
	return &bBookStore{}
}

func (sotre *bBookStore) Find(name string) *book {
	return &book{}
}
