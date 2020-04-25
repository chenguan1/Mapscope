package server

import (
	"Mapscope/routes"
	"github.com/kataras/iris/v12"
)

// 30 MB
const maxFontSize = 30 << 20

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
		ds.Post("/{username}/{dataset_id}/features", routes.DatasetFeaturesPut)                   // List features
		ds.Put("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesInsert)    // Insert or update a feature
		ds.Get("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesRetrive)   // Retrieve a feature
		ds.Delete("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesDelete) // Retrieve a feature

		// 自己定义的接口，非mapbox定义的接口
		// 上传数据集，支持zip包，geojson，json，shp(zip)
		ds.Post("/{username}", routes.DatasetUpload)
	}

	// fonts
	fs := app.Party("/fonts/v1")
	{
		fs.Get("/{username}/{font}/{rangepbf:string regexp(^[0-9]+-[0-9]+.pbf)}", routes.FontGlypRange) // Retrieve font glyph ranges
		fs.Get("/{username}", routes.FontList)                                                          // List fonts
		fs.Post("/{username}", iris.LimitRequestBodySize(maxFontSize), routes.FontAdd)                  // Add a font
	}

	// tilesets
	ts := app.Party("/tilesets/v1")
	{
		// tileset source
		tss := ts.Party("/sources")
		{

			tss.Post("/{username}/{id}", routes.TilesetSourceCreate)   // Create a tileset source, id: tileset source ID
			tss.Get("/{username}/{id}", routes.TilesetSourceRetrieve)  // Retrieve tileset source information
			tss.Get("/{username}", routes.TilesetSourceList)           // List tileset sources
			tss.Delete("/{username}/{id}", routes.TilesetSourceDelete) // Beta
		}

		ts.Post("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}", routes.TilesetCreate)               // Create a tileset
		ts.Post("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/publish", routes.TilesetPublish)      // Publish a tileset
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/status", routes.TilesetStatus)         // Retrieve the status of a tileset
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs/{job_id}", routes.TilesetJobInfo) // Retrieve information about a single tileset job
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs", routes.TilesetJobList)          // List information about all jobs for a tileset
		ts.Put("/queue", routes.TilesetJobQueue)                                                          // View the Tilesets API global queue
		ts.Put("/validateRecipe", routes.TilesetRecipeValidate)                                           // Validate a recipe
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/recipe", routes.TilesetRecipe)         // Retrieve a tileset's recipe
		ts.Patch("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/recipe", routes.TilesetRecipeUpdate) // Update a tileset's recipe

		ts.Get("/{username}", routes.TilesetList)                                                        // List tilesets
		ts.Delete("/{username}/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}", routes.TilesetDelete) // Delete tileset
		ts.Get("/{tilesetjson:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+.json)}", routes.TilesetMetadata)    // Retrieve TileJSON metadata

	}

	// Retrieve vector tiles
	// app.Get("/v4/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.{[a-z]+})}")
	// Retrieve raster tiles
	// app.Get("/v4/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.{[a-z]+})}")

	return app
}
