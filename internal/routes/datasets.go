package routes

import (
	"Mapscope/internal/config"
	"Mapscope/internal/models"
	"Mapscope/internal/services"
	"Mapscope/internal/utils"
	"fmt"
	"github.com/kataras/iris/v12/context"
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

func DatasetRetrive(ctx context.Context) {
	dtid := ctx.Params().Get("dataset_id")

	res := utils.NewRes(ctx)

	dt, err := services.DatasetGet(dtid)
	if err != nil {
		res.FailMsg("Get dataset list failed.")
		return
	}
	res.DoneData(dt)
}

func DatasetList(ctx context.Context) {
	user := ctx.Params().Get("username")

	res := utils.NewRes(ctx)

	dts, err := services.DatasetList(user)
	if err != nil {
		res.FailMsg("Get dataset list failed.")
		return
	}
	res.DoneData(dts)
}

/*
更新dataset信息
只有以下属性可以修改：
name，public，description
*/
func DatasetUpdate(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	res := utils.NewRes(ctx)

	var dtForm models.Dataset
	if err := ctx.ReadJSON(&dtForm); err != nil {
		res.FailMsg("read dataset info from request failed.")
		return
	}

	dt, err := services.DatasetGet(dtid)
	if err != nil {
		res.FailMsg(fmt.Sprintf("server can not find dataset %v.", dtid))
		return
	}

	dt.Name = dtForm.Name
	dt.Public = dtForm.Public
	dt.Description = dtForm.Description
	dt.Modified = time.Now()

	if err := dt.Save(); err != nil {
		res.FailErr(fmt.Errorf("store dataset info failed. %v", dtid))
		return
	}

	res.DoneData(dt)
}

// 删除datast
func DatasetDelete(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	ctx.Application().Logger().Printf("delete %v.%v", user, dtid)

	res := utils.NewRes(ctx)

	// 删除
	if err := services.DatasetDelete(dtid); err != nil {
		ctx.Application().Logger().Errorf("DatasetDelete err: %v", err)
		res.FailMsg("delete dataset failed.")
		return
	}

	res.DoneMsg("dataset has been deleted.")
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
	files := utils.SaveFormFiles(ctx, upfolder, false)
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
// 要支持dtid1,dtid3,dtid3这种形式传入多个id
func DatasetTile(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtidstr := ctx.Params().Get("dataset_ids")

	dtids := strings.Split(dtidstr, ",")

	zoom, _ := strconv.Atoi(ctx.Params().Get("zoom"))
	x, _ := strconv.Atoi(ctx.Params().Get("x"))
	yformat := strings.Split(ctx.Params().Get("yformat"), ".")
	y, _ := strconv.Atoi(yformat[0])
	format := "." + yformat[1]

	res := utils.NewRes(ctx)

	// 用户认证

	if format != ".mvt" && format != ".pbf" {
		ctx.Application().Logger().Errorf("DatasetTile, format %s not surpported.", format)
		res.FailMsg(fmt.Sprintf("format %s not surported.", format))
		return
	}

	data := make([]byte, 0)

	for _, dtid := range dtids {
		// 获取切片
		mvtbuf, err := services.DatasetToMvtBuf(dtid, zoom, x, y)
		if err != nil {
			ctx.Application().Logger().Errorf("DatasetTile, get mvt failed, err: %v", err)
			res.FailMsg("get mvt tile failed.")
			return
		}
		data = append(data, mvtbuf...)
	}

	ctx.ContentType("application/vnd.mapbox-vector-tile")
	ctx.Binary(data)

}

// 获取Tilejson
func DatasetTilejson(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtids := strings.Split(ctx.Params().Get("dataset_ids"), ",")

	var err error

	res := utils.NewRes(ctx)

	var tj *models.Tilejson

	if tj, err = services.DatasetsTilejson(dtids); err != nil {
		res.FailMsg("get tile json failed.")
		return
	}

	res.Json(tj)
}

// 备份dataset
/*
策略：
1. 开始编辑，进入编辑状态
2. 提交编辑，version++
3. 备份，对当前version进行备份，如果同一个version已经存在，则备份失败，除非使用参数force强制备份
   编辑状态下不可以备份。
*/
func DatasetBackup(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	res := utils.NewRes(ctx)

	bk, err := services.DatasetBackup(dtid, true)
	if err != nil {
		ctx.Application().Logger().Error("DatasetBackup err: %v", err)
		res.FailMsg("databack up failed, err: " + err.Error())
		return
	}

	res.DoneData(bk)
}

// 提交dataset修改，版本加1
func DatasetCommit(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	services.DatasetCommit(dtid)
}

// 获取某一dataset的备份列表
func DatasetBackupList(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	res := utils.NewRes(ctx)

	bks, err := services.DatasetBackupList(dtid)
	if err != nil {
		ctx.Application().Logger().Error("DatasetBackup err: %v", err)
		res.FailMsg("get databack up err: " + err.Error())
		return
	}

	res.DoneData(bks)
}

// 恢复到某一备份版本，根据version号
func DatasetBackupRevert(ctx context.Context) {
	//user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	res := utils.NewRes(ctx)
	version := ctx.URLParam("version")
	if version == "" {
		ctx.Application().Logger().Error("DatasetBackupRevert err: version not spacifyed.")
		res.FailMsg("databack up failed, err: version not spacifyed")
		return
	}

	versionNum, _ := strconv.Atoi(version)

	dt, err := services.DatasetRevertTo(dtid, versionNum)
	if err != nil {
		ctx.Application().Logger().Error("DatasetBackupRevert err: %v", err)
		res.FailMsg("revert failed.")
		return
	}

	res.DoneData(dt)
}
