package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/thirdparty/paulmach/orb/encoding/wkt"
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"time"
	//gogeo "github.com/paulmach/go.geojson"
)

// 获取features
func FeatureListByGeojson(dataset_id string, offset, count int) (interface{}, error) {
	var dt models.Dataset
	db := database.Get()
	err := db.Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		return nil, fmt.Errorf("FeatureGetByGeojson failed, err: %v", err)
	}

	type gjItem struct {
		Gj []byte `gorm:"column:gj"`
	}

	var gis []gjItem
	//err = db.Raw("select st_asgeojson(t.*) as gj from ? t where gid = ?", dt.TableName, feature_id).Scan(&gi).Error
	sql := fmt.Sprintf("select st_asgeojson(t.*) as gj from %v t offset %v limit %v", dt.TableName, offset, count)
	fmt.Println(sql)

	err = db.Raw(sql).Scan(&gis).Error
	if err != nil {
		return nil, fmt.Errorf("FeatureGetByGeojson query failed, err: %v", err)
	}

	var jsons []interface{}
	for _, g := range gis {
		var gj interface{}
		err = json.Unmarshal(g.Gj, &gj)
		if err != nil {
			return nil, fmt.Errorf("FeatureGetByGeojson unmarshal json failed, err: %v", err)
		}
		jsons = append(jsons, gj)
	}

	return jsons, nil
}

// 获取feature
func FeatureGetByGeojson(dataset_id string, feature_id string) (interface{}, error) {
	var dt models.Dataset
	db := database.Get()
	err := db.Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		return nil, fmt.Errorf("FeatureGetByGeojson failed, err: %v", err)
	}

	type gjItem struct {
		Gj []byte `gorm:"column:gj"`
	}

	var gi gjItem
	//err = db.Raw("select st_asgeojson(t.*) as gj from ? t where gid = ?", dt.TableName, feature_id).Scan(&gi).Error
	sql := fmt.Sprintf("select st_asgeojson(t.*) as gj from %v t where gid = %v", dt.TableName, feature_id)
	fmt.Println(sql)

	err = db.Raw(sql).Scan(&gi).Error
	if err != nil {
		return nil, fmt.Errorf("FeatureGetByGeojson query failed, err: %v", err)
	}

	var gj interface{}
	err = json.Unmarshal(gi.Gj, &gj)
	if err != nil {
		return nil, fmt.Errorf("FeatureGetByGeojson unmarshal json failed, err: %v", err)
	}

	return gj, nil
}

// 删除feature
func FeatureDelete(dataset_id string, feature_id string) error {
	var dt models.Dataset

	tx := database.Get().Begin()

	err := tx.Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("FeatureDelete failed, err: %v", err)
	}

	sqlfmt := `delete from %s where gid = %v`

	sql := fmt.Sprintf(sqlfmt, dt.TableName, feature_id)

	dbrt := tx.Exec(sql)
	err = dbrt.Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("FeatureDelete failed, err: %v", err)
	}

	if dbrt.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("FeatureDelete failed, err: %v", "rows affected == 0")
	}

	dt.Edited++
	dt.Modified = time.Now()
	tx.Save(&dt)

	tx.Commit()

	return nil
}

// 更新feature
func FeatureUpdate(dataset_id string, feature_id string, gjson []byte) error {
	var dt models.Dataset
	tx := database.Get().Begin()
	err := tx.Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		return fmt.Errorf("FeatureUpdate failed, err: %v", err)
	}

	ft, err := geojson.UnmarshalFeature(gjson)
	if err != nil {
		return fmt.Errorf("FeatureUpdate failed, err: %v", err)
	}

	// geom
	wt := wkt.MarshalString(ft.Geometry)
	fmt.Println(wt)

	// 属性
	setvalue := ""
	fds := dt.Fields.Keys()
	ps := map[string]interface{}(ft.Properties)
	for _, key := range fds {
		v, ok := ps[key]
		if ok {
			setvalue = setvalue + fmt.Sprintf("%v='%v',", key, v)
		}
	}
	setvalue = setvalue + fmt.Sprintf("geom = st_geomfromtext('%v',%v)", wt, 4326)

	// 更新
	//sql := fmt.Sprintf("update %s\nset %s \nwhere gid = %v", dt.TableName, setvalue, feature_id)
	//fmt.Println(sql)

	sqlfmt := `
update %s 
set %s
where gid = %v
`
	sql := fmt.Sprintf(sqlfmt, dt.TableName, setvalue, feature_id)

	dbrt := tx.Exec(sql)
	err = dbrt.Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("FeatureUpdate failed, err: %v", err)
	}

	if dbrt.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("FeatureUpdate failed, err: %v", "rows affected == 0")
	}

	dt.Edited++
	dt.Modified = time.Now()
	tx.Save(&dt)

	tx.Commit()

	return nil
}

// 插入feature
func FeatureInsert(dataset_id string, gjson []byte) error {
	var dt models.Dataset
	tx := database.Get().Begin()
	err := tx.Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		return fmt.Errorf("FeatureInsert failed, err: %v", err)
	}

	ft, err := geojson.UnmarshalFeature(gjson)
	if err != nil {
		return fmt.Errorf("FeatureInsert failed, err: %v", err)
	}

	// geom
	wt := wkt.MarshalString(ft.Geometry)
	fmt.Println(wt)

	// 属性
	setfiled := ""
	setvalue := ""
	fds := dt.Fields.Keys()
	ps := map[string]interface{}(ft.Properties)
	for _, key := range fds {
		v, ok := ps[key]
		if ok {
			setfiled = setfiled + key + ","
			setvalue = setvalue + fmt.Sprintf("'%v',", v)
		}
	}
	setfiled = setfiled + "geom"
	setvalue = setvalue + fmt.Sprintf("st_geomfromtext('%v',%v)", wt, 4326)

	sqlfmt := `insert into %s (%s) values (%s)`
	sql := fmt.Sprintf(sqlfmt, dt.TableName, setfiled, setvalue)

	fmt.Println(sql)

	dbrt := tx.Exec(sql)
	err = dbrt.Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("FeatureUpdate failed, err: %v", err)
	}

	if dbrt.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("FeatureUpdate failed, err: %v", "rows affected == 0")
	}

	dt.Edited++
	dt.Modified = time.Now()
	tx.Save(&dt)
	tx.Commit()

	return nil
}
