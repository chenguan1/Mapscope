package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"fmt"
	"github.com/consbio/mbtileserver/mbtiles"
	"os"
	"time"
)

// 根据id获取tileset
func TilesetGet(tileset_id string) (*models.Tileset, error) {
	var ts models.Tileset
	err := database.Get().Where(models.Tileset{Id: tileset_id}).Find(&ts).Error
	if err != nil {
		return nil, fmt.Errorf("TilesetGet failed, err: %v", err)
	}
	return &ts, nil
}

func TilesetLoad(path string) (*models.Tileset, error) {
	mb, err := mbtiles.NewDB(path)
	if err != nil {
		return nil, fmt.Errorf("Valid tileset could not be opened: %q", err)
	}
	defer mb.Close()

	finfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("tileset info fetch failed: %q", err)
	}

	meta, err := mb.ReadMetadata()
	if err != nil {
		return nil, fmt.Errorf("metadata of tileset fetch failed: %q", err)
	}

	tm := &models.TilesetMeta{}
	if err = tm.Parse(meta); err != nil {
		return nil, fmt.Errorf("parse metadata table failed: %q", err)
	}

	ts := &models.Tileset{
		Id:          "",
		Name:        tm.Name,
		Type:        tm.Type,
		Filesize:    finfo.Size(),
		Format:      tm.Format,
		Center:      tm.Center,
		Description: tm.Description,
		Visibility:  "public",
		Public:      1,
		Status:      "success",
		Created:     finfo.ModTime(),
		Modified:    time.Now(),
		Metadata:    *tm,
		Dataset:     "",
		Path:path,
	}

	return ts, nil
}

// 获取某用户的所有tileset
func TilesetList(user string) ([]models.Tileset, error) {
	var tss []models.Tileset
	err := database.Get().Where(models.Tileset{Owner: user}).Find(&tss).Error
	if err != nil {
		return nil, fmt.Errorf("TilesetList err: %v", err)
	}
	return tss, nil
}

// tileset delete,
func TilesetDelete(tileset_id string) error {
	var err error
	var ts models.Tileset

	db := database.Get()
	err = db.Where(models.Tileset{Id: tileset_id}).Find(&ts).Error
	if err != nil {
		return fmt.Errorf("TilesetDelete err: %v", err)
	}

	tspath := ts.Path

	// 事务删除
	tx := db.Begin()

	// 1.remove tileset record
	if err = tx.Delete(ts).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Database delete tileset record err: %v", err)
	}

	// 2.删除文件
	err = os.Remove(tspath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("tileset delete failed: %v", err)
	}

	// 提交
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Database commit err: %v", err)
	}

	return nil
}