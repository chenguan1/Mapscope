package cache

import (
	"Mapscope/internal/config"
	"Mapscope/internal/utils"
	"database/sql"
	"github.com/faabiosr/cachego"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

var cache cachego.Cache

func InitCache() {
	cachePath := config.PathCaches()
	utils.EnsurePathExist(cachePath)
	cacheFile := filepath.Join(cachePath,"mapscope.cache")
	db, _ := sql.Open("sqlite3", cacheFile)

	cache, _ = cachego.NewSqlite3(db, "cache")
}

func Save(key string, data interface{}) error {
	return nil
}