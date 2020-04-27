package appsrv

import (
	"Mapscope/internal/config"
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/routes"
	"github.com/kataras/iris/v12"
)

func Run() {
	var err error

	if err = config.Initialize(); err != nil {
		panic(err)
	}
	if err = database.Initialize(); err != nil {
		panic(err)
	}

	database.Get().AutoMigrate(&models.Datasource{})
	database.Get().AutoMigrate(&models.Dataset{})

	app := iris.Default()
	app.Logger().SetLevel("debug")

	// set up routes
	routes.SetRoutes(app)

	app.Run(iris.Addr(":8080"))
}
