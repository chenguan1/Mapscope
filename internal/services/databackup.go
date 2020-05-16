package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"fmt"
	"time"
)

// dataset backip
func DatasetBackup(dataset_id string, force bool) (*models.DataBackup, error) {
	dt, err := DatasetGet(dataset_id)
	if err != nil {
		return nil, fmt.Errorf("DatasetBackup err: %v", err)
	}

	// 当前状态
	if dt.IsEditing() {
		return nil, fmt.Errorf("DatasetBackup err: is editing.")
	}

	// 是否已经备份过
	db := database.Get()
	var bks []models.DataBackup
	err = db.Where(models.DataBackup{Version: dt.Version}).Find(&bks).Error
	if err != nil {
		return nil, fmt.Errorf("database cnn failed. %v", err)
	}

	bk := models.DataBackup{
		Dataset:     dt.Id,
		Source:      dt.Source,
		Version:     dt.Version,
		TableName:   fmt.Sprintf("%s_v%v", dt.TableName, dt.Version),
		Created:     time.Now(),
		Description: dt.Description,
	}

	if len(bks) > 0 {
		if !force {
			return nil, fmt.Errorf("database version %v has been backuped.", dt.Version)
		}
		bk = bks[0]
		bk.Created = time.Now()
	}

	tx := db.Begin()

	// 删除备份记录
	tx.DropTableIfExists(bk.TableName)
	// 备份
	sql := fmt.Sprintf(`CREATE TABLE %v as (SELECT * FROM %v)`, bk.TableName, dt.TableName)
	if err = tx.Exec(sql).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("dataset backup err: %v", err)
	}

	if err = tx.Save(&bk).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("dataset backup err.: %v", err)
	}

	dt.Modified = time.Now()
	tx.Save(&dt)

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("dataset backup err..: %v", err)
	}

	return &bk, nil
}

// 删除某dataset相关的所有记录
func DataBackupDeleteAll(dataset_id string) error {
	dt, err := DatasetGet(dataset_id)
	if err != nil {
		return fmt.Errorf("DatasetBackup delete err: %v", err)
	}

	db := database.Get()

	// all backups
	var bks []models.DataBackup
	err = db.Where(models.DataBackup{Dataset: dt.Id, Version: dt.Version}).Find(&bks).Error
	if err != nil {
		return fmt.Errorf("database cnn failed. %v", err)
	}

	if len(bks) == 0 {
		return nil
	}

	tx := db.Begin()

	for _, bk := range bks {
		if err = tx.DropTableIfExists(bk.TableName).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("DatasetBackup delete err.: %v", err)
		}
		if err = tx.Delete(bk).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("DatasetBackup delete err..: %v", err)
		}
	}

	tx.Commit()

	return nil
}

// 获取某dataset的所有备份
func DatasetBackupList(dataset_id string) ([]models.DataBackup, error) {
	dt, err := DatasetGet(dataset_id)
	if err != nil {
		return nil, fmt.Errorf("DatasetBackupList err: %v", err)
	}

	db := database.Get()

	// all backups
	var bks []models.DataBackup
	err = db.Where(models.DataBackup{Dataset: dt.Id}).Find(&bks).Error
	if err != nil {
		return nil, fmt.Errorf("DatasetBackupList err: %v", err)
	}

	return bks, nil
}
