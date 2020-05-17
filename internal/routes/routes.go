package routes

import (
	"github.com/kataras/iris/v12"
)

// 50 MB
const maxFontSize = 50 << 20

// 设置路由
func SetRoutes(app *iris.Application) {
	app.Use(Cors)
	// 静态文件
	app.HandleDir("/", "./public")

	api := app.Party("/api")

	// datasets ok
	ds := api.Party("/datasets")
	{
		// dataset
		ds.Get("/list/{username}", DatasetList)      // list datasets ok
		ds.Post("/upload/{username}", DatasetUpload) // upload dataset, create a dataset to one user，geojson，json，shp(zip) ok

		ds.Get("/{dataset_id}", DatasetRetrive)   // Retrieve a dataset ok
		ds.Patch("/{dataset_id}", DatasetUpdate)  // Update a dataset name,public,description 可以修改 ok
		ds.Delete("/{dataset_id}", DatasetDelete) // Delete a dataset ok

		ds.Get(`/{dataset_ids}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.[a-z]+)}`, DatasetTile) // get mvt tile of dataset ok
		ds.Get("/{dataset_ids}/tile.json", DatasetTilejson)                                              // tilejson ok
		ds.Head("/{dataset_ids}/tile.json", DatasetTilejson)                                             // tilejson ok

		// backup
		ds.Put("/{dataset_id}/backup", DatasetBackup) // ok

		// commit changes
		ds.Put("/{dataset_id}/commit", DatasetCommit)
		// get backup list
		ds.Get("/{dataset_id}/backup", DatasetBackupList)
		// revert to version
		ds.Put("/{dataset_id}/revert", DatasetBackupRevert)

		// feature
		ds.Get("/{dataset_id}/features", DatasetFeatures)                       // List features      ok
		ds.Post("/{dataset_id}/features", DatasetFeaturesInsert)                // features insert    ok
		ds.Get("/{dataset_id}/features/{feature_id}", DatasetFeaturesRetrive)   // Retrieve a feature ok
		ds.Patch("/{dataset_id}/features/{feature_id}", DatasetFeaturesUpdate)  // features update    ok
		ds.Delete("/{dataset_id}/features/{feature_id}", DatasetFeaturesDelete) // Delete a feature   ok
	}

	// fonts ok
	fs := api.Party("/fonts")
	{
		fs.Get("/{fonts}/{rangepbf:string regexp(^[0-9]+-[0-9]+.pbf)}", FontGlypRange) // Retrieve font glyph ranges ok
		fs.Get("/list", FontList)                                                      // List fonts ok
		fs.Post("/upload", iris.LimitRequestBodySize(maxFontSize), FontAdd)            // upload a font ok
		fs.Delete("/{fontname}", FontDelete)                                           // delete ok
	}

	// tilesets
	ts := api.Party("/tilesets")
	{
		ts.Get("/list/{username}", TilesetList)      // List tilesets ok
		ts.Post("/upload/{username}", TilesetUpload) // upload mbtiles file ok
		ts.Delete("/{tileset_id}", TilesetDelete)    // Delete tileset ok

		ts.Post("/publish/{dataset_id}", TilesetPublish) // Publish a dataset to tileset
		ts.Get("/{tileset_id}/status", TilesetStatus)    // Retrieve the status of a tileset

		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs/{job_id}", TilesetJobInfo) // Retrieve information about a single tileset job
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs", TilesetJobList)          // List information about all jobs for a tileset
		ts.Put("/queue", TilesetJobQueue)                                                          // View the Tilesets API global queue
		ts.Put("/validateRecipe", TilesetRecipeValidate)                                           // Validate a recipe
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/recipe", TilesetRecipe)         // Retrieve a tileset's recipe
		ts.Patch("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/recipe", TilesetRecipeUpdate) // Update a tileset's recipe
	}
}
