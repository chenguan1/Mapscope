package routes

import (
	"Mapscope/internal/config"
	"Mapscope/internal/models"
	"Mapscope/internal/services"
	"Mapscope/internal/thirdparty/teris-io/shortid"
	"Mapscope/internal/utils"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"strings"
)

/*
https://docs.mapbox.com/api/maps/#create-a-tileset
*/

// 上传Mbtiles
func TilesetUpload(ctx context.Context) {
	user := ctx.Params().Get("username")

	uploadPath := config.PathUploads(user)
	utils.EnsurePathExist(uploadPath)
	files := utils.SaveFormFiles(ctx, uploadPath, false)

	ctx.Application().Logger().Debug(files)

	res := utils.NewRes(ctx)

	tss := make([]*models.Tileset, 0, len(files))
	for _, f := range files {
		switch strings.ToLower(f.Ext) {
		case ".mbtiles":
			ts, err := services.TilesetLoad(f.Path)
			if err != nil {
				ctx.Application().Logger().Errorf("TilesetLoad err: %v", err)
				res.FailMsg("load Tileset file failed.")
				return
			}
			ts.Id, _ = shortid.GenerateLower()
			ts.Owner = user

			if err = ts.Save(); err != nil {
				ctx.Application().Logger().Errorf("Tileset save err: %v", err)
				res.FailMsg("Tileset save failed.")
				return
			}
			// ctx.Application().Logger().Info(f)
			tss = append(tss, ts)
		default:
			continue
		}
	}

	res.DoneData(tss)
}

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


/*
发布一个dataset 为 tileset
要使用job异步处理
*/
func TilesetPublish(ctx context.Context) {
	dtid := ctx.Params().Get("dataset_id")
	res := utils.NewRes(ctx)
	ds,err := services.Dataset2Tileset(dtid)
	if err != nil{
		ctx.Application().Logger().Error(err)
		res.FailMsg("Tileset public failed, err:" + err.Error())
		return
	}

	res.DoneData(ds)

	/*res := map[string]string{
		"message": "Processing " + ts,
		"job_id":  "afaddfa5654",
	}

	ctx.JSON(res)*/
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

// get tileset list by username
func TilesetList(ctx context.Context) {
	username := ctx.Params().Get("username")
	res := utils.NewRes(ctx)
	tss, err := services.TilesetList(username)
	if err != nil {
		ctx.Application().Logger().Errorf("TilesetList err: %v", err)
		res.FailMsg("Get tileset list failed.")
		return
	}
	res.DoneData(tss)
}

// delete tileset
/*
1. remove record
2. remove mbtiles file
*/
func TilesetDelete(ctx context.Context) {
	tsid := ctx.Params().Get("tileset_id")

	res := utils.NewRes(ctx)

	// 删除
	if err := services.TilesetDelete(tsid); err != nil {
		ctx.Application().Logger().Errorf("tilesetDelete err: %v", err)
		res.FailMsg("delete tileset failed.")
		return
	}

	res.DoneMsg("tileset has been deleted.")
}

func TilesetMetadata(ctx context.Context) {
	ctx.StatusCode(http.StatusNoContent)
}
