package routes

import (
	"Mapscope/internal/config"
	"Mapscope/internal/database"
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
// 并且入库成为dataset-list
func DatasetUpload(ctx context.Context)  {
	user := ctx.Params().Get("username")

	var err error

	// 用户权限验证
	// to do ...

	// 处理业务
	upfolder := config.PathUploads(user)
	utils.EnsurePathExist(upfolder)

	res := utils.NewRes(ctx)

	// 保存上传的文件
	files := utils.SaveFormFiles(ctx,upfolder)
	if len(files) == 0{
		ctx.Application().Logger().Errorf("DatasetUpload SaveFormFiles, err: %v", err)
		res.FailMsg("upload files failed.")
		return
	}

	// 上传的数据入库成为dataset
	dts, err := services.DatasetsFromUpload(files, user)
	if err != nil{
		ctx.Application().Logger().Errorf("DatasetUpload DatasetsFromUpload, err: %v", err)
		res.FailMsg("create dataset failed.")
		return
	}

	res.DoneData(dts)
}

//
func DatasetTile(ctx context.Context)  {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	zoom, _ := strconv.Atoi(ctx.Params().Get("zoom"))
	x, _ := strconv.Atoi(ctx.Params().Get("x"))
	yformat := strings.Split(ctx.Params().Get("yformat"),".")
	y,_ := strconv.Atoi(yformat[0])
	format := "."+yformat[1]

	ctx.Application().Logger().Info(user,"  ",dtid,"  ",zoom,"  ",x,"  ",y,"  ",format)
	ctx.WriteString("..")

}

func getMvt(tableName, colName string, z, x, y int) ([]byte,error) {
	queryFmt := `
	WITH mvtgeom AS
	(SELECT ST_AsMVTGeom(%v, ST_TileEnvelope(%v,%v,%v)) AS geom,"id","code","name","fclass"
	FROM %v 
	WHERE ST_Intersects(%v, ST_TileEnvelope(%v,%v,%v)))
	SELECT ST_AsMvt(mvtgeom.*,'%v') AS mvt
	FROM mvtgeom
`

	type geoItem struct {
		Mvt []byte `gorm:"column:mvt"`
	}


	queryStr := fmt.Sprintf(queryFmt,colName,z,x,y,tableName,colName,z,x,y,tableName)

	//fmt.Println(queryStr)

	var gi geoItem
	err := database.Get().Raw(queryStr).Scan(&gi).Error
	if err != nil{
		return nil, err
	}
	return gi.Mvt, nil
}