package routes

import (
	"Mapscope/internal/config"
	"Mapscope/internal/models"
	"Mapscope/internal/services"
	"Mapscope/internal/utils"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
https://docs.mapbox.com/api/maps/#datasets
Datasets API errors
Unauthorized	401	Check the access token you used in the query.
No such user	404	Check the username you used in the query.
Not found	404	The access token used in the query needs the datasets:list scope.
Invalid start key	422	Check the start key used in the query.
No dataset	422	Check the dataset ID used in the query (or, if you are retrieving a feature, check the feature ID).
*/

func DatasetList(ctx context.Context) {
	user := ctx.Params().Get("username")
	dts := make([]models.Dataset, 0)
	dts = append(dts, models.Dataset{
		Owner: user,
	})
	ctx.JSON(dts)
}

func DatasetCreate(ctx context.Context) {
	user := ctx.Params().Get("username")
	type post struct {
		Name        string `json"name"`
		Description string `json:description`
	}

	bd := post{}
	if err := ctx.ReadJSON(&bd); err != nil {
		ctx.JSON(err)
		return
	}

	dt := models.Dataset{
		Id:          "123456",
		Owner:       user,
		Name:        bd.Name,
		Created:     time.Now(),
		Modified:    time.Now(),
		Description: bd.Description,
	}
	// create dataset

	ctx.JSON(dt)
}

func DatasetRetrive(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	dt := models.Dataset{
		Id:    dtid,
		Owner: user,
	}

	ctx.JSON(dt)
}

func DatasetUpdate(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	type post struct {
		Name        string `json"name"`
		Description string `json:description`
	}

	bd := post{}
	if err := ctx.ReadJSON(&bd); err != nil {
		ctx.JSON(err)
		return
	}

	dt := models.Dataset{
		Id:          dtid,
		Owner:       user,
		Name:        bd.Name,
		Created:     time.Now(),
		Modified:    time.Now(),
		Description: bd.Description,
	}

	ctx.JSON(dt)
}

func DatasetDelete(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	ctx.Application().Logger().Printf("delete %v.%v", user, dtid)

	ctx.StatusCode(http.StatusNoContent)
}

func DatasetFeatures(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	limit := ctx.URLParamDefault("limit", "10")
	start := ctx.URLParamDefault("start", "0")

	lmt, _ := strconv.Atoi(limit)
	srt, _ := strconv.Atoi(start)
	if lmt > 100 {
		lmt = 100
	}

	ctx.Application().Logger().Printf("get features %v.%v,%v,%v", user, dtid, limit, start)
	res := utils.NewRes(ctx)

	gj, err := services.FeatureListByGeojson(dtid, srt, lmt)
	if err != nil {
		ctx.Application().Logger().Errorf("retrive feature failed err: %v", err)
		res.FailMsg("retrive feature failed.")
		return
	}

	res.DoneData(gj)

}

func DatasetFeaturesUpdate(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("update feature %v.%v.%v", user, dtid, ftid)

	res := utils.NewRes(ctx)

	gjson, err := ctx.GetBody()
	if err != nil {
		res.FailMsg("failed to get body.")
		ctx.Application().Logger().Errorf("DatasetFeaturesUpdate failed get body. err: %v", err)
		return
	}

	err = services.FeatureUpdate(dtid, ftid, gjson)
	if err != nil {
		res.FailMsg("failed update feature.")
		ctx.Application().Logger().Errorf("DatasetFeaturesUpdate.FeatureUpdate failed get body. err: %v", err)
		return
	}

	res.Done()
}

func DatasetFeaturesInsert(ctx context.Context) {

	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("update feature %v.%v.%v", user, dtid, ftid)

	res := utils.NewRes(ctx)

	gjson, err := ctx.GetBody()
	if err != nil {
		res.FailMsg("failed to get body.")
		ctx.Application().Logger().Errorf("DatasetFeaturesInsert failed get body. err: %v", err)
		return
	}

	err = services.FeatureInsert(dtid, gjson)
	if err != nil {
		res.FailMsg("failed update feature.")
		ctx.Application().Logger().Errorf("DatasetFeaturesInsert.FeatureUpdate failed get body. err: %v", err)
		return
	}

	res.Done()
}

func DatasetFeaturesRetrive(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("retrive a feature %v.%v,%v", user, dtid, ftid)

	res := utils.NewRes(ctx)

	gj, err := services.FeatureGetByGeojson(dtid, ftid)
	if err != nil {
		ctx.Application().Logger().Errorf("retrive feature failed err: %v", err)
		res.FailMsg("retrive feature failed.")
		return
	}

	res.DoneData(gj)
}

func DatasetFeaturesDelete(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("delete a feature %v.%v,%v", user, dtid, ftid)

	res := utils.NewRes(ctx)

	err := services.FeatureDelete(dtid, ftid)
	if err != nil {
		ctx.Application().Logger().Errorf("DatasetFeaturesDelete failed err: %v", err)
		res.FailMsg("delete feature failed.")
		return
	}

	res.Done()
}

// 数据上传，支持shp压缩包，geojson，csv等
// 其他格式的都转换成geojson格式的datasource，
// 并且入库成为dataset-list
func DatasetUpload(ctx context.Context) {
	user := ctx.Params().Get("username")

	var err error

	// 用户权限验证
	// to do ...

	// 处理业务
	upfolder := config.PathUploads(user)
	utils.EnsurePathExist(upfolder)

	res := utils.NewRes(ctx)

	// 保存上传的文件
	files := utils.SaveFormFiles(ctx, upfolder)
	if len(files) == 0 {
		ctx.Application().Logger().Errorf("DatasetUpload SaveFormFiles, err: %v", err)
		res.FailMsg("upload files failed.")
		return
	}

	// 上传的数据入库成为dataset
	dts, err := services.DatasetsFromUpload(files, user)
	if err != nil {
		ctx.Application().Logger().Errorf("DatasetUpload DatasetsFromUpload, err: %v", err)
		res.FailMsg("create dataset failed.")
		return
	}

	res.DoneData(dts)
}

// 获取数据切片
func DatasetTile(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	zoom, _ := strconv.Atoi(ctx.Params().Get("zoom"))
	x, _ := strconv.Atoi(ctx.Params().Get("x"))
	yformat := strings.Split(ctx.Params().Get("yformat"), ".")
	y, _ := strconv.Atoi(yformat[0])
	format := "." + yformat[1]

	res := utils.NewRes(ctx)

	// 用户认证

	// 获取切片
	if format == ".mvt" || format == ".pbf" {
		mvtbuf, err := services.DatasetToMvtBuf(dtid, zoom, x, y)
		if err != nil {
			ctx.Application().Logger().Errorf("DatasetTile, get mvt failed, err: %v", err)
			res.FailMsg("get mvt tile failed.")
			return
		}
		ctx.ContentType("application/vnd.mapbox-vector-tile")
		ctx.Binary(mvtbuf)
	} else {
		ctx.Application().Logger().Errorf("DatasetTile, format %s not surpported.", format)
		res.FailMsg(fmt.Sprintf("format %s not surported.", format))
		return
	}
}

// 获取Tilejson
func DatasetTilejson(ctx context.Context)  {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	var tj *models.Tilejson
	var err error

	res := utils.NewRes(ctx)

	if tj,err = services.DatasetTilejson(dtid); err != nil{
		res.FailMsg("get tile json failed.")
		return
	}

	res.Json(tj)
}
