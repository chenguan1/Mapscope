package server

import (
	"Mapscope/routes"
	"github.com/kataras/iris/v12"
)

// 设置路由
func newApp() *iris.Application {
	app := iris.Default()
	app.Logger().SetLevel("info")

	// 静态文件
	app.HandleDir("/", "./public")

	// datasets
	ds := app.Party("/datasets/v1")
	{
		ds.Get("/{username}", routes.DatasetList)                                                 // list datasets
		ds.Post("/{username}", routes.DatasetCreate)                                              // create a new empty dataset
		ds.Get("/{username}/{dataset_id}", routes.DatasetRetrive)                                 // Retrieve a dataset
		ds.Patch("/{username}/{dataset_id}", routes.DatasetUpdate)                                // Retrieve a dataset
		ds.Delete("/{username}/{dataset_id}", routes.DatasetDelete)                               // Delete a dataset
		ds.Get("/{username}/{dataset_id}/features", routes.DatasetFeatures)                       // List features
		ds.Put("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesInsert)    // Insert or update a feature
		ds.Get("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesRetrive)   // Retrieve a feature
		ds.Delete("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesDelete) // Retrieve a feature
	}

	fs := app.Party("/fonts/v1")
	{
		fs.Get("/{username}/{font}/{start}.ddd", routes.FontGlypRange) // Retrieve font glyph ranges
		fs.Get("/{username}", routes.FontList)                                       // List fonts
		fs.Post("/{username}", routes.FontAdd)                                       // Add a font
	}

	return app
}
