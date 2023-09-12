package di

type (
	Container struct {
		Scope
		ServiceCollector
	}
)

func New() *Container {
	rootScope := newInnerScope("root", nil)

	return &Container{
		Scope:            rootScope,
		ServiceCollector: rootScope,
	}
}
