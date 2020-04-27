package utils

import (
	"Mapscope/internal/models"
	"github.com/kataras/iris/v12/context"
	"github.com/teris-io/shortid"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 将上传到的文件保存到指定目录中，
// 如果文件中已经存在同名的文件，则加上一个前缀
func SaveFormFiles(ctx context.Context, folder string) []models.FileUped {
	files := make([]models.FileUped,0)

	ctx.UploadFormFiles(folder, func(i context.Context, header *multipart.FileHeader) {
		sid, _ := shortid.Generate()
		ext := filepath.Ext(header.Filename)
		name := strings.TrimSuffix(header.Filename, ext)
		ext = strings.ToLower(ext)
		ff := models.FileUped{
			Sid:sid,
			Name:name,
			Ext:ext,
			Path:filepath.Join(folder,header.Filename),
			Size:header.Size,
		}

		if _,err := os.Stat(ff.Path); os.IsExist(err){
			header.Filename = sid + header.Filename
			ff.Path = filepath.Join(folder,header.Filename)
		}

		files = append(files,ff)
	})

	return files
}
