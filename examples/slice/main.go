package main

import (
	"fmt"

	di "github.com/AlgerDu/go-di/src"
)

type Controller interface {
	Actions()
}

type AuthController struct {
}

func NewAuth() (*AuthController, error) {
	return &AuthController{}, nil
}

func (controller *AuthController) Actions() {
	fmt.Println("/api/v1/auth")
}

type BookController struct {
}

func NewBook() (*BookController, error) {
	return &BookController{}, nil
}

func (controller *BookController) Actions() {
	fmt.Println("/api/v1/boog/add")
	fmt.Println("/api/v1/boog/get")
	fmt.Println("/api/v1/boog/list")
	fmt.Println("/api/v1/boog/delete")
}

type HttpServer struct {
	controllers []Controller
}

func NewHttp(controllers []Controller) (*HttpServer, error) {
	return &HttpServer{
		controllers: controllers,
	}, nil
}

func (server *HttpServer) Start() {
	fmt.Println("http start")
	fmt.Println("support urls:")

	for _, controller := range server.controllers {
		controller.Actions()
	}
}

func main() {
	container := di.New()

	di.AddSingletonFor[Controller](container, NewAuth)
	di.AddSingletonFor[Controller](container, NewBook)

	di.AddSingleton(container, NewHttp)

	http, err := di.GetService[*HttpServer](container)
	if err != nil {
		panic(err)
	}

	http.Start()
}
