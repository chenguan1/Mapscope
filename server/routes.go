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
		ds.Put("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesInsert)    // Insert or update a feature
		ds.Get("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesRetrive)   // Retrieve a feature
		ds.Delete("/{username}/{dataset_id}/features/{feature_id}", routes.DatasetFeaturesDelete) // Retrieve a feature
	}

	// fonts
	fs := app.Party("/fonts/v1")
	{
		fs.Get("/{username}/{font}/{rangepbf:string regexp(^[0-9]+-[0-9]+.pbf)}", routes.FontGlypRange) // Retrieve font glyph ranges
		fs.Get("/{username}", routes.FontList)                                       // List fonts
		fs.Post("/{username}", iris.LimitRequestBodySize(maxFontSize), routes.FontAdd)                                       // Add a font
	}

	// tilesets
	ts := app.Party("/tilesets/v1")
	{
		// tileset source
		tss := ts.Party("/sources")
		{

			tss.Post("/{username}/{id}", routes.TilesetSourceCreate) // Create a tileset source, id: tileset source ID
			tss.Get("/{username}/{id}", routes.TilesetSourceRetrieve) // Retrieve tileset source information
			tss.Get("/{username}",routes.TilesetSourceList) // List tileset sources
			tss.Delete("/{username}/{id}", routes.TilesetSourceDelete) // Beta
		}

		ts.Post("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}", routes.TilesetCreate) // Create a tileset
		ts.Post("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/publish", routes.TilesetPublish) // Publish a tileset
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/status",routes.TilesetStatus) // Retrieve the status of a tileset
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs/{job_id}",routes.TilesetJobInfo) // Retrieve information about a single tileset job
		ts.Get("/{tileset:string regexp(^[a-zA-Z_-]+.[a-zA-Z_-]+)}/jobs",routes.TilesetJobList) // Retrieve information about a single tileset job

	}

	return app
}
