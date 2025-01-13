# go-di

## 为什么

## 功能

### 三种生命周期

`SL_Singleton` 单实例，在整个 `Container` 内（包含 root 及所有的子作用域）是只会生成构造一个实例。
`SL_Scoped` 作用域内单实例，在每个 scope 内只会构造一个。
`SL_Transient` 临时，每次获取都会获取一个新的实列。

### 直接注入一个实例

```
func AddInstance(services ServiceCollector,ins any) error
```

### 作用域支持

可以将 `di.Scope` 作为依赖，例如：

```
func New(scope di.Scope) (*Struct,error) {}
```

生成的实列将会获取其所在的作用域，然后进行其他的操作。比如通过 `CreateSubScope` 创建并且获取一个新的作用域。

### 可以将 slice 作为依赖

1. 注入接口的不同实现：

```
di.AddSingletonFor[Controller](container, NewAuth)
di.AddSingletonFor[Controller](container, NewBook)
```

2. 注入依赖的服务：

```
func NewHttp(controllers []Controller) (*HttpServer, error) {
	return &HttpServer{
		controllers: controllers,
	}, nil
}


di.AddSingleton(container, NewHttp)
```

3. 在需要的地方使用：

```
func (server *HttpServer) Start() {
	for _, controller := range server.controllers {
		controller.Actions()
	}
}
```

详情见 example/slice

### 不注入，直接通过 `New` 函数，获取新的实列

```
printer, err := ResloveService[func() bool](scope, useNamePrinter)
```

### 通过使用泛型的扩展方法来简化使用

```
func AddSingleton()
func AddSingletonFor[forT any]()
func AddScope()
func AddScopeFor[forT any]()
...
func GetService[ServiceType any]()
```
