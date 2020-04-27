package routes

import (
	"Mapscope/internal/models"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io"
	"net/http"
	"os"
)

func TilesetSourceCreate(ctx context.Context) {
	//user := ctx.Params().Get("username")
	tssid := ctx.Params().Get("id")

	file, info, err := ctx.FormFile("file")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}

	defer file.Close()
	fname := info.Filename

	// Create a file with the same name
	// assuming that you have a folder named 'uploads'
	out, err := os.OpenFile("./data/uploads/"+fname,
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	defer out.Close()
	io.Copy(out, file)

	tss := models.TilesetSource{
		Id:         tssid,
		FileSize:   info.Size,
		Files:      1,
		SourceSize: info.Size,
	}

	ctx.JSON(tss)
}

func TilesetSourceRetrieve(ctx context.Context) {
	//user := ctx.Params().Get("username")
	tssid := ctx.Params().Get("id")

	tss := models.TilesetSource{
		Id:       tssid,
		Files:    2,
		Size:     2048,
		SizeNice: "2.0KB",
	}

	ctx.JSON(tss)
}

func TilesetSourceList(ctx context.Context) {
	//user := ctx.Params().Get("username")
	tssid := ctx.Params().Get("id")

	tss := models.TilesetSource{
		Id:       tssid,
		Files:    2,
		Size:     2048,
		SizeNice: "2.0KB",
	}

	list := make([]models.TilesetSource, 0)
	list = append(list, tss)

	ctx.JSON(list)
}

func TilesetSourceDelete(ctx context.Context) {
	user := ctx.Params().Get("username")
	tssid := ctx.Params().Get("id")

	ctx.Application().Logger().Debug(user, tssid)

	ctx.StatusCode(http.StatusNoContent)
}
