package models

import (
	"Mapscope/internal/database"
	"time"
	//_ "github.com/mattn/go-sqlite3" // import sqlite3 driver
	// "github.com/paulmach/orb/encoding/wkb"
)

// Task 数据导入信息预览
type Task struct {
	ID        string        `json:"id" form:"id" binding:"required" gorm:"primary_key"`
	Base      string        `json:"base" form:"base" gorm:"index"`
	Name      string        `json:"name" form:"name"`
	Type      TaskType      `json:"type" form:"type"`
	Owner     string        `json:"owner" form:"owner"`
	Progress  int           `json:"progress" form:"progress"`
	Status    string        `json:"status"`
	Error     string        `json:"error" `
	Pipe      chan struct{} `json:"-" form:"-" gorm:"-"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}


func (task *Task) Save() error {
	db := database.Get()
	err := db.Create(task).Error
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) Update() error {
	db := database.Get()
	err := db.Model(&Task{}).Update(task).Error
	if err != nil {
		return err
	}
	return nil
}
