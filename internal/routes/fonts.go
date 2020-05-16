package routes

import (
	"Mapscope/internal/config"
	"Mapscope/internal/models"
	"Mapscope/internal/services"
	"Mapscope/internal/utils"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"strings"
	"time"
)

/*
https://docs.mapbox.com/api/maps/#fonts
Fonts API errors
Response body message	HTTP status code	Description
Invalid Range	400	Check the font glyph range. This error can also occur with empty fonts.
Maximum of 10 font faces permitted	400	Too many font faces in the query.
Not Authorized - No Token	401	No token was used in the query.
Not Authorized - Invalid Token	401	Check the access token you used in the query.
Not Found	404	Check the font name or names you used in the query.
*/

func FontGlypRange(ctx context.Context) {
	//user := ctx.Params().Get("username")
	font := ctx.Params().Get("fonts")
	pbf := ctx.Params().Get("rangepbf") // 0-255.pbf

	//fmt.Printf("%v,%v,%v\n", user, font, pbf)

	res := utils.NewRes(ctx)

	fonts := strings.Split(font, ",")
	data, err := services.FontGlyphs(fonts, pbf)
	if err != nil {
		res.FailMsg("get font range failed.")
		return
	}

	lastModified := time.Now().UTC().Format("2006-01-02 03:04:05 PM")
	ctx.Header("Content-Type", "application/x-protobuf")
	ctx.Header("Last-Modified", lastModified)
	ctx.Binary(data)
}

func FontList(ctx context.Context) {
	//user := ctx.Params().Get("username")
	user := "mapscope"

	res := utils.NewRes(ctx)

	fts, err := services.FontList(user)
	if err != nil {
		res.FailMsg("Get font list failed.")
		return
	}

	fonts := make([]string, 0, len(fts))
	for _, f := range fts {
		fonts = append(fonts, f.Name)
	}

	res.DoneData(fonts)
}

// 上传字体
func FontAdd(ctx context.Context) {
	//user := ctx.Params().Get("username")
	user := "mapscope"

	uploadPath := config.PathUploads(user)
	utils.EnsurePathExist(uploadPath)
	files := utils.SaveFormFiles(ctx, uploadPath, true)

	ctx.Application().Logger().Debug(files)

	res := utils.NewRes(ctx)

	fonts := make([]models.Font, 0, len(files))
	for _, f := range files {
		switch strings.ToLower(f.Ext) {
		case ".pbfonts":
			ft, err := services.FontLoad(f.Path)
			if err != nil {
				ctx.Application().Logger().Errorf("FontAdd err: %v", err)
				res.FailMsg("load font file failed.")
				return
			}
			ft.Owner = user
			fonts = append(fonts, *ft)
		default:
			continue
		}
	}

	for _, ft := range fonts {
		if err := ft.Save(); err != nil {
			ctx.Application().Logger().Errorf("FontAdd err: %v", err)
			res.FailMsg("save font info failed.")
			return
		}
	}

	res.DoneData(&fonts)
}

// 上传字体
func FontDelete(ctx context.Context) {
	//user := ctx.Params().Get("username")
	font := ctx.Params().Get("fontname")

	res := utils.NewRes(ctx)

	err := services.FontDelete(font)
	if err != nil {
		ctx.Application().Logger().Errorf("Font delete err: %v", err)
		res.FailMsg(fmt.Sprintf("font %s delete failed.", font))
	}

	res.DoneMsg(fmt.Sprintf("font %s has been deleted.", font))
}
