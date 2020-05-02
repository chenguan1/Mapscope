package routes

import (
	"github.com/kataras/iris/v12"
)

// 50 MB
const maxFontSize = 50 << 20

// 设置路由
func SetRoutes(app *iris.Application) {
	// 静态文件
	app.HandleDir("/", "./public")

	// datasets ok
	ds := app.Party("/datasets/v1")
	{
		ds.Get("/{username}", DatasetList)                   // list datasets ok

		ds.Get("/{username}/{dataset_id}", DatasetRetrive)   // Retrieve a dataset ok
		ds.Patch("/{username}/{dataset_id}", DatasetUpdate)  // Update a dataset name,public,description 可以修改 ok
		ds.Delete("/{username}/{dataset_id}", DatasetDelete) // Delete a dataset ok

		// dataset
		ds.Post("/{username}", DatasetUpload)                                                                      // 上传数据集，支持zip包，geojson，json，shp(zip) ok
		ds.Get(`/{username}/{dataset_id}/{zoom:int}/{x:int}/{yformat:string regexp(^[0-9]+.[a-z]+)}`, DatasetTile) // format=mvt ok
		ds.Get("/{username}/{dataset_id}/tile.json", DatasetTilejson) // ok
		ds.Head("/{username}/{dataset_id}/tile.json", DatasetTilejson) // ok

		// feature
		ds.Get("/{username}/{dataset_id}/features", DatasetFeatures)                       // List features      ok
		ds.Post("/{username}/{dataset_id}/features", DatasetFeaturesInsert)                // features insert    ok
		ds.Get("/{username}/{dataset_id}/features/{feature_id}", DatasetFeaturesRetrive)   // Retrieve a feature ok
		ds.Patch("/{username}/{dataset_id}/features/{feature_id}", DatasetFeaturesUpdate)  // features update    ok
		ds.Delete("/{username}/{dataset_id}/features/{feature_id}", DatasetFeaturesDelete) // Delete a feature   ok
	}

	// fonts ok
	fs := app.Party("/fonts/v1")
	{
		fs.Get("/{username}/{font}/{rangepbf:string regexp(^[0-9]+-[0-9]+.pbf)}", FontGlypRange) // Retrieve font glyph ranges ok
		fs.Get("/{username}", FontList)                                                          // List fonts ok
		fs.Post("/{username}", iris.LimitRequestBodySize(maxFontSize), FontAdd)                  // Add a font ok
		fs.Delete("/{username}/{fontname}", FontDelete)                                          // delete ok
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
}
