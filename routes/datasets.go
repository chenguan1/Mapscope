package routes

import (
	"Mapscope/config"
	"Mapscope/controls"
	"Mapscope/model"
	"Mapscope/utils"
	"github.com/kataras/iris/v12/context"
	"github.com/teris-io/shortid"
	"mime/multipart"
	"net/http"
	"path/filepath"
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
	dts := make([]model.Dataset, 0)
	dts = append(dts, model.Dataset{
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

	dt := model.Dataset{
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

	dt := model.Dataset{
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

	dt := model.Dataset{
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

// 数据上传
func DatasetUpload(ctx context.Context)  {
	user := ctx.Params().Get("username")
	dtid := ctx.Params().Get("dataset_id")

	uf := filepath.Join(config.Get().DataFolder, "uploads", user)
	utils.EnsurePathExist(uf)

	var filename string
	_, err :=ctx.UploadFormFiles(uf, func(i context.Context, header *multipart.FileHeader) {
		sid, _ := shortid.Generate()
		header.Filename = sid + "." + header.Filename
		filename = header.Filename
	})

	if err != nil{
		panic(err)
	}

	ctx.Application().Logger().Info("dataset uploaded:",user, dtid, filename)

	// 处理入库
	err = controls.DatasetParseAndStore(user, dtid, filepath.Join(uf,filename))
	if err != nil{
		panic(err)
	}

	ctx.JSON("ok")
}