package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"fmt"
	"strconv"
	"time"
)

// 创建数据表备份
/*
1. 备份数据表
2. 登记DataBackup表
*/
func DataBackupCreate(dt *models.Dataset, description string) (*models.DataBackup, error) {
	newTableName := dt.TableName + "_" + strconv.Itoa(dt.Version)
	db := database.Get()
	err := db.DropTableIfExists(newTableName).Error
	if err != nil {
		return nil, fmt.Errorf("DataBackupCreate, drop table failed, err: %v", err)
	}

	err = db.Exec("create table ? as (select * from ?)", newTableName, dt.TableName).Error
	if err != nil {
		return nil, fmt.Errorf("DataBackupCreate, backup table failed, err: %v", err)
	}

	dp := &models.DataBackup{
		Dataset:     dt.Id,
		Source:      dt.Source,
		Version:     dt.Version,
		TableName:   newTableName,
		Created:     time.Now(),
		Description: description,
	}

	err = db.Save(dp).Error
	if err != nil {
		return nil, fmt.Errorf("DataBackupCreate, save backup info failed, err: %v", err)
	}

	return dp, nil
}
