package utils

import (
	"github.com/kataras/iris/v12/context"
	"github.com/teris-io/shortid"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 保存上传的文件
type FormFile struct {
	Path       string
	OriginName string
	Ext        string // .zip .geojson
	Size       int64
}

// 将上传到的文件保存到指定目录中，
// 如果文件中已经存在同名的文件，则加上一个前缀
func SaveFormFiles(ctx context.Context, folder string) []FormFile {
	files := make([]FormFile,0)

	ctx.UploadFormFiles(folder, func(i context.Context, header *multipart.FileHeader) {
		ff := FormFile{}
		ff.OriginName = header.Filename
		ff.Path = filepath.Join(folder,header.Filename)
		if _,err := os.Stat(ff.Path); os.IsExist(err){
			sid, _ := shortid.Generate()
			header.Filename = sid + header.Filename
			ff.Path = filepath.Join(folder,header.Filename)
		}

		ff.Ext = strings.ToLower(filepath.Ext(ff.OriginName))
		ff.Size = header.Size

		files = append(files,ff)
	})

	return files
}
