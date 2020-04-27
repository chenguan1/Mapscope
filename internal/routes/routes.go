package routes

import (
	"github.com/kataras/iris/v12"
)

// 30 MB
const maxFontSize = 30 << 20

// 设置路由
func SetRoutes(app *iris.Application) {
	// 静态文件
	app.HandleDir("/", "./public")

	// datasets
	ds := app.Party("/datasets/v1")
	{
		ds.Get("/{username}", DatasetList)                                                 // list datasets
		ds.Post("/{username}", DatasetCreate)                                              // create a new empty dataset
		ds.Get("/{username}/{dataset_id}", DatasetRetrive)                                 // Retrieve a dataset
		ds.Patch("/{username}/{dataset_id}", DatasetUpdate)                                // Retrieve a dataset
		ds.Delete("/{username}/{dataset_id}", DatasetDelete)                               // Delete a dataset
		ds.Get("/{username}/{dataset_id}/features", DatasetFeatures)                       // List features
		ds.Post("/{username}/{dataset_id}/features", DatasetFeaturesPut)                   // List features
		ds.Put("/{username}/{dataset_id}/features/{feature_id}", DatasetFeaturesInsert)    // Insert or update a feature
		ds.Get("/{username}/{dataset_id}/features/{feature_id}", DatasetFeaturesRetrive)   // Retrieve a feature
		ds.Delete("/{username}/{dataset_id}/features/{feature_id}", DatasetFeaturesDelete) // Retrieve a feature

		// 自己定义的接口，非mapbox定义的接口
		// 上传数据集，支持zip包，geojson，json，shp(zip)
		ds.Post("/{username}", DatasetUpload)
		ds.Get(`/{username}/{dataset_id}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.[a-z]+)}`, DatasetTile)
	}

	// fonts
	fs := app.Party("/fonts/v1")
	{
		fs.Get("/{username}/{font}/{rangepbf:string regexp(^[0-9]+-[0-9]+.pbf)}", FontGlypRange) // Retrieve font glyph ranges
		fs.Get("/{username}", FontList)                                                          // List fonts
		fs.Post("/{username}", iris.LimitRequestBodySize(maxFontSize), FontAdd)                  // Add a font
	}

	// tilesets
	ts := app.Party("/tilesets/v1")
	{
		// tileset source
		tss := ts.Party("/sources")
		{

			tss.Post("/{username}/{id}", TilesetSourceCreate)   // Create a tileset source, id: tileset source ID
			tss.Get("/{username}/{id}", TilesetSourceRetrieve)  // Retrieve tileset source information
			tss.Get("/{username}", TilesetSourceList)           // List tileset sources
			tss.Delete("/{username}/{id}", TilesetSourceDelete) // Beta
		}

		ts.Post("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}", TilesetCreate)               // Create a tileset
		ts.Post("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/publish", TilesetPublish)      // Publish a tileset
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/status", TilesetStatus)         // Retrieve the status of a tileset
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs/{job_id}", TilesetJobInfo) // Retrieve information about a single tileset job
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs", TilesetJobList)          // List information about all jobs for a tileset
		ts.Put("/queue", TilesetJobQueue)                                                          // View the Tilesets API global queue
		ts.Put("/validateRecipe", TilesetRecipeValidate)                                           // Validate a recipe
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/recipe", TilesetRecipe)         // Retrieve a tileset's recipe
		ts.Patch("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/recipe", TilesetRecipeUpdate) // Update a tileset's recipe

		ts.Get("/{username}", TilesetList)                                                        // List tilesets
		ts.Delete("/{username}/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}", TilesetDelete) // Delete tileset
		ts.Get("/{tilesetjson:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+.json)}", TilesetMetadata)    // Retrieve TileJSON metadata

	}

	// Retrieve vector tiles
	// app.Get("/v4/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.{[a-z]+})}")
	// Retrieve raster tiles
	// app.Get("/v4/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.{[a-z]+})}")
}
