package di

type ServiceLifetime int

const (
	SL_Singleton ServiceLifetime = iota
	SL_Scoped
	SL_Transient
)
