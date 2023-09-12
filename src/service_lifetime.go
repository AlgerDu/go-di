package di

type ServiceLifetime int

const (
	SL_Unknown ServiceLifetime = iota
	SL_Singleton
	SL_Scoped
	SL_Transient
)
