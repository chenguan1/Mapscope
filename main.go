package main

import (
	"Mapscope/global"
	"Mapscope/models"
	"Mapscope/server"
)

func main() {
	err := global.InitGlobal()
	if err != nil{
		panic(err)
	}
	defer global.Clear()

	db := global.GetDb()
	db.AutoMigrate(&models.Datasource{})
	db.AutoMigrate(&models.Dataset{})


	server.Run()
}
