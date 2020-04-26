package global

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB

func initDb() error {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.database"))

	var err error
	db,err = gorm.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("init gorm db error, details: %s", err)
	}

	return nil
}

func GetDb() *gorm.DB {
	return db
}