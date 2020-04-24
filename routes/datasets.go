package routes

import (
	"Mapscope/config"
	"Mapscope/controls"
	"Mapscope/models"
	"Mapscope/utils"
	"github.com/kataras/iris/v12/context"
	"github.com/teris-io/shortid"
	"net/http"
	"os"
	"path/filepath"
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
	limit := ctx.URLParam("limit")
	start := ctx.URLParam("start")

	//
	// limit [1,100]
	if limit == "" {
		limit = "10"
	}

	// start: The ID of the feature after which to start the listing.

	ctx.Application().Logger().Printf("get features %v.%v,%v,%v", user, dtid, limit, start)

	ctx.WriteString("return the json list of features.")
}

func DatasetFeaturesPut(ctx context.Context) {
	//
}

func DatasetFeaturesInsert(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("insert feature %v.%v,%v", user, dtid, ftid)

	geojson := make(map[string]interface{})
	ctx.ReadJSON(&geojson)
	ctx.JSON(geojson)
}

func DatasetFeaturesRetrive(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("retrive a feature %v.%v,%v", user, dtid, ftid)

	geojson := `
{
  "id": "{feature_id}",
  "type": "Feature",
  "geometry": {
    "type": "Polygon",
    "coordinates": [
      [
        [ 100, 0 ],
        [ 101, 0 ],
        [ 101, 1 ],
        [ 100, 1 ],
        [ 100, 0 ]
      ]
    ]
  },
  "properties": {
    "prop0": "value0"
  }
}
`
	ctx.WriteString(geojson)
}

func DatasetFeaturesDelete(ctx context.Context) {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")
	ftid := ctx.Params().Get("feature_id")

	ctx.Application().Logger().Printf("delete a feature %v.%v,%v", user, dtid, ftid)

	ctx.StatusCode(http.StatusNoContent)
}

// 数据上传，支持shp压缩包，geojson，csv等
// 其他格式的都转换成geojson格式的datasource，
// 并且入库成为dataset
func DatasetUpload(ctx context.Context)  {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	var err error

	// 用户权限验证
	// to do ...

	folder := filepath.Join(config.Get().DataFolder, "uploads", user)
	utils.EnsurePathExist(folder)

	dtfolder := filepath.Join(config.Get().DataFolder, "datasets", user)
	utils.EnsurePathExist(dtfolder)

	res := utils.NewRes(ctx)

	// 保存上传的文件
	files := utils.SaveFormFiles(ctx,folder)
	if len(files) == 0{
		res.FailMsg("upload files failed.")
		return
	}

	// 得到所有文件
	allfs := make([]string,0)
	for _, ff := range files{
		switch ff.Ext {
		case models.ZIPEXT:
			fs,err := utils.UnzipFile(ff.Path, "")
			if err != nil{
				ctx.Application().Logger().Error("unzip file failed, file:",ff.Path)
				res.FailMsg("unzip the file failed.")
			}
			allfs = append(allfs, fs...)
		}
	}

	// 过滤出支持的后缀文件格式，暂时只支持
	exts := models.GEOJSONEXT + models.SHPEXT + models.CSVEXT// + models.KMLEXT + models.GPXEXT
	vfs := utils.FilterByExt(allfs, exts)
	if len(vfs) == 0{
		ctx.Application().Logger().Error("file uploaded is not supported.")
		res.FailMsg("file type not supported.")
		return
	}

	// 将支持的格式转成geojson格式的datasorce
	dtsrcs := make([]models.Datasource,0) // datasorce 准备入库，记录必要信息
	for _, f := range vfs{
		gp := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f)) + ".geojson"
		if utils.PathExist(filepath.Join(dtfolder,gp)){
			sid,_ := shortid.Generate()
			gp = sid + "." + gp
		}

		err = controls.ToGeojson(f, gp)
		if err != nil{
			ctx.Application().Logger().Error("DatasetUpload, convert to geojson failed:",err)
			res.FailMsg("parse file failed.")
			return
		}

		info,err := os.Stat(gp)
		if err != nil{
			ctx.Application().Logger().Error("DatasetUpload error,",err)
			res.FailMsg("unexpacted error.")
			return
		}

		// to datasource
		st := models.Datasource{}
		st.Id,_ = shortid.Generate()
		st.Name = strings.TrimPrefix(filepath.Base(gp),filepath.Ext(gp))
		st.Owner = user
		st.Src = f
		st.Path = gp
		st.CreatedAt = time.Now()
		st.UpdatedAt = time.Now()
		st.Size = info.Size()
		st.Tag = "dataset." + dtid

		// append
		dtsrcs = append(dtsrcs, st)
	}


	// 入库 dtsrc，并转换成dataset
	dts := make([]*models.Dataset,0)
	for _,ds := range dtsrcs{
		// 入库
		err := ds.Save()
		if err != nil{
			ctx.Application().Logger().Error("Datasorce save failed,",ds.Path,err)
			res.FailMsg("Datasorce save failed.")
			return
		}
		// to dataset
		dt,err := ds.ToDataset()
		if err != nil{
			ctx.Application().Logger().Error("Dataset save failed,",ds.Path,err)
			res.FailMsg("Dataset save failed.")
			return
		}

		// dt 入库
		err = dt.Save()
		if err != nil{
			ctx.Application().Logger().Error("Dataset insert failed,",ds.Path,err)
			res.FailMsg("Dataset insert failed.")
			return
		}
		dts = append(dts, dt)
	}

	res.DoneData(dts)
}