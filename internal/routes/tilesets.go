package routes

import (
	"Mapscope/internal/models"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

/*
https://docs.mapbox.com/api/maps/#create-a-tileset
*/

func TilesetCreate(ctx context.Context) {
	ts := ctx.Params().Get("tileset")

	ctx.Application().Logger().Debug(ts)

	tform := models.TilesetCreateForm{}
	tform.Private = true

	err := ctx.ReadJSON(&tform)
	if err != nil {
		ctx.JSON(err)
		return
	}

	ctx.JSON(tform)
}

func TilesetPublish(ctx context.Context) {
	ts := ctx.Params().Get("tileset")

	ctx.Application().Logger().Debug(ts)

	res := map[string]string{
		"message": "Processing " + ts,
		"job_id":  "afaddfa5654",
	}

	ctx.JSON(res)
}

func TilesetStatus(ctx context.Context) {
	ts := ctx.Params().Get("tileset")

	ctx.Application().Logger().Debug(ts)

	res := map[string]string{
		"id":         ts,
		"latest_job": "afaddfa5654",
		"status":     "success",
	}

	ctx.JSON(res)
}

// https://docs.mapbox.com/api/maps/#retrieve-information-about-a-single-tileset-job
func TilesetJobInfo(ctx context.Context) {
	ts := ctx.Params().Get("tileset")
	job := ctx.Params().Get("job_id")

	ctx.Application().Logger().Debug(ts, job)

	res := map[string]interface{}{
		"id":           "unique_hash",
		"stage":        "success",
		"created":      1560981902377,
		"created_nice": "Wed Jun 19 2019 22:05:02 GMT+0000 (UTC)",
		"published":    1560982158721,
		"tilesetId":    "user.id",
		"errors":       make([]string, 0),
		"warnings":     make([]string, 0),
	}

	ctx.JSON(res)
}

func TilesetJobList(ctx context.Context) {
	ts := ctx.Params().Get("tileset")

	ctx.Application().Logger().Debug(ts)

	res := map[string]interface{}{
		"id":           "unique_hash",
		"stage":        "success",
		"created":      1560981902377,
		"created_nice": "Wed Jun 19 2019 22:05:02 GMT+0000 (UTC)",
		"published":    1560982158721,
		"tilesetId":    "user.id",
		"errors":       make([]string, 0),
		"warnings":     make([]string, 0),
	}

	joblist := make([]map[string]interface{}, 0)
	joblist = append(joblist, res)
	joblist = append(joblist, res)

	ctx.JSON(joblist)
}

func TilesetJobQueue(ctx context.Context) {
	res := map[string]interface{}{
		"total": 42,
	}

	ctx.JSON(res)
}

func TilesetRecipeValidate(ctx context.Context) {
	rcp := models.Recipe{}
	err := ctx.ReadJSON(&rcp)
	if err != nil {
		ctx.JSON(err)
		return
	}

	res := map[string]interface{}{
		"valid": true,
	}
	ctx.JSON(res)
}

func TilesetRecipe(ctx context.Context) {
	ts := ctx.Params().Get("tileset")

	ctx.Application().Logger().Debug(ts)

	tr := models.TilesetCreateForm{}
	ctx.JSON(tr)

}

func TilesetRecipeUpdate(ctx context.Context) {
	ts := ctx.Params().Get("tileset")

	ctx.Application().Logger().Debug(ts)

	tr := models.TilesetCreateForm{}
	ctx.JSON(tr)

}

func TilesetList(ctx context.Context)  {
	username := ctx.Params().Get("username")
	ctx.Application().Logger().Debug(username)
	ts := models.Tileset{}
	tslist := make([]models.Tileset,0)
	tslist = append(tslist, ts)
	ctx.JSON(tslist)
}

func TilesetDelete(ctx context.Context)  {
	ctx.StatusCode(http.StatusNoContent)
}

func TilesetMetadata(ctx context.Context)  {
	ctx.StatusCode(http.StatusNoContent)
}
