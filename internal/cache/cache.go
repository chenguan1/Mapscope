package cache

import (
	"Mapscope/internal/config"
	"Mapscope/internal/utils"
	"database/sql"
	"github.com/faabiosr/cachego"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"time"
)

var cache cachego.Cache
var db *sql.DB

const DEFAULT_LIFE_TIME = time.Hour * 24 * 365

func Initialize() error {
	var err error
	cachePath := config.PathCaches()
	utils.EnsurePathExist(cachePath)
	cacheFile := filepath.Join(cachePath, "mapscope.cache")
	db, err = sql.Open("sqlite3", cacheFile)
	if err != nil {
		return err
	}

	cache, err = cachego.NewSqlite3(db, "cache")
	return err
}

func Destroy() {
	db.Close()
}

func Save(key string, data interface{}) error {
	return nil
}
