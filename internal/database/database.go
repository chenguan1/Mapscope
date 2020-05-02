package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

var db *gorm.DB

// 用于sqlite的连接池
var sqliteDbs map[string]*gorm.DB

func Initialize() error {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.database"))

	var err error
	db, err = gorm.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("init gorm db error, details: %s", err)
	}

	return nil
}

func Get() *gorm.DB {
	return db
}

// 打开sqlite数据库
func OpenSqlite(path string) (*gorm.DB, error) {
	if sqliteDbs == nil{
		sqliteDbs = make(map[string]*gorm.DB)
	}

	v,ok := sqliteDbs[path]
	if ok{
		return v, nil
	}

	d, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite3 err: %v", err)
	}

	sqliteDbs[path] = d
	return d, nil
}