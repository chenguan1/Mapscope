package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/thirdparty/paulmach/orb/encoding/wkt"
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb/geojson"

	//gogeo "github.com/paulmach/go.geojson"
)

// 获取feature
func FeatureGetByGeojson(dataset_id string, feature_id string) (interface{}, error) {
	var dt models.Dataset
	db := database.Get()
	err := db.Where(models.Dataset{Id:dataset_id}).Find(&dt).Error
	if err != nil{
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
	if err != nil{
		return nil, fmt.Errorf("FeatureGetByGeojson query failed, err: %v", err)
	}

	var gj interface{}
	err = json.Unmarshal(gi.Gj,&gj)
	if err != nil{
		return nil, fmt.Errorf("FeatureGetByGeojson unmarshal json failed, err: %v", err)
	}

	return gj, nil
}


// 删除feature

// 更新feature
func FeatureUpdate(dataset_id string, feature_id string, gjson []byte) error {
	var dt models.Dataset
	db := database.Get()
	err := db.Where(models.Dataset{Id:dataset_id}).Find(&dt).Error
	if err != nil{
		return fmt.Errorf("FeatureUpdate failed, err: %v", err)
	}

	/*ft, err := geojson.UnmarshalFeature(gjson)
	if err != nil{
		return fmt.Errorf("FeatureUpdate unmarshal feature failed, err: %v", err)
	}*/

	ft,err := geojson.UnmarshalFeature(gjson)
	wt := wkt.MarshalString(ft.Geometry)
	fmt.Println(wt)



	return nil
}