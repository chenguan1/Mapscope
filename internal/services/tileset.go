package services

import (
	"Mapscope/internal/models"
	"fmt"
	"github.com/consbio/mbtileserver/mbtiles"
	"os"
	"time"
)

func TilesetLoad(path string) (*models.Tileset, error) {
	mb, err := mbtiles.NewDB(path)
	if err != nil {
		return nil, fmt.Errorf("Valid tileset could not be opened: %q", err)
	}
	defer mb.Close()

	finfo, err := os.Stat(path)
	if err != nil{
		return nil, fmt.Errorf("tileset info fetch failed: %q", err)
	}

	meta,err := mb.ReadMetadata()
	if err != nil{
		return nil, fmt.Errorf("metadata of tileset fetch failed: %q", err)
	}

	tm := &models.TilesetMeta{}
	if err = tm.Parse(meta); err != nil{
		return nil, fmt.Errorf("parse metadata table failed: %q", err)
	}

	ts := &models.Tileset{
		Id:"",
		Name:tm.Name,
		Type:tm.Type,
		Filesize:finfo.Size(),
		Format: tm.Format,
		Center: tm.Center,
		Description:tm.Description,
		Visibility:"public",
		Status:"???",
		Created: finfo.ModTime(),
		Modified: time.Now(),
		Dataset: "",
	}

	return ts, nil
}
