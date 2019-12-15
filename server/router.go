package server

import (
	"github.com/gin-gonic/gin"
)

// 设置路由
func setupRouter() *gin.Engine  {
	r := gin.Default()
	{
		r.GET("/ping", handlePing)

		// Fonts
		ft := r.Group("/fonts")
		{
			v1 := ft.Group("/v1")
			v1.GET("/:user/:font/:range.pbf")
		}

		// datasets
		ds := r.Group("/datasets")
		{
			v1 := ds.Group("/v1")
			v1.GET("/:user", handleDatasetList)
			v1.POST("/:user", handleDatesetCreate)
			v1.GET("/:user/:id", handleDatasetGet)
			v1.PATCH("/:user/:id",handleDatasetUpdate)
			v1.DELETE("/:user/:id", handleDatasetDelete)
			v1.GET("/:user/:id/features",handleDatasetFeatureList)
			v1.PUT("/:user/:id/features/:feature_id", handleDatasetFeaturePut)
			v1.GET("/:user/:id/features/:feature_id", handleDatasetFeatureGet)
			v1.DELETE("/:user/:id/features/:feature_id", handleDatasetFeatureDelete)
		}

		// tilesets
		ts := r.Group("/tilesets")
		{
			v1 := ts.Group("/v1")
			v1.POST("/sources/:user/:id", handleTilesetSourceCreate)
			v1.GET("/sources/:user/:id", handleTilesetSourceGet)
			v1.GET("/sources/:user", handleTilesetSourceList)
			v1.DELETE("/sources/:user/:id", handleTilesetSourceDelete)

			v1.POST("/:id", handleTilesetCreate)
			v1.POST("/:id/publish", handleTilesetPublish)
			v1.GET("/:id/status", handleTilesetStatus)

			v1.GET("/:id/jobs/:job_id", handleTilesetJobInfo)
			v1.GET("/:id/jobs", handleTilesetJobInfoList)

			v1.PUT("/squeue", handleTilesetQueue)

			v1.PUT("/validateRecipe", handleTilesetValidateRecipe)
			v1.GET("/:id/recipe", handleTilesetRecipeGet)
			v1.PATCH("/:id/recipe", handleTilesetRecipeUpdate)

			v1.GET("/:user", handleTilesetList)
			v1.DELETE("/:user.:id", handleTilesetDelete)

			v4 := ts.Group("/v4")
			v4.GET("/:id.json", handleTilesetMetadataGet)
		}

		// styles
		sl := r.Group("/styles")
		{
			v1 := sl.Group("/v1")
			v1.GET("/", handleStylesList)
		}
	}


	return r
}