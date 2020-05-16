package appsrv

import (
	"Mapscope/internal/cache"
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
	defer database.Destroy()

	if err = cache.Initialize(); err != nil {
		panic(err)
	}
	defer cache.Destroy()

	database.Get().AutoMigrate(&models.Datasource{})
	database.Get().AutoMigrate(&models.Dataset{})
	database.Get().AutoMigrate(&models.Font{})
	database.Get().AutoMigrate(&models.DataBackup{})

	app := iris.Default()
	app.Logger().SetLevel("debug")

	// set up routes
	routes.SetRoutes(app)

	app.Run(iris.Addr(":8080"))
}
