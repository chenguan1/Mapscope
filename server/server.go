package server

import "github.com/kataras/iris/v12"

func Run()  {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

