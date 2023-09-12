package di

type (
	Container struct {
		Scope
	}
)

func New() *Container {
	rootScope := newInnerScope("root", nil)

	return &Container{
		Scope: rootScope,
	}
}
