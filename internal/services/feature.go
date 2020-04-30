package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/thirdparty/paulmach/orb/encoding/wkt"
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"strings"

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

	ft,err := geojson.UnmarshalFeature(gjson)
	if err != nil{
		return fmt.Errorf("FeatureUpdate failed, err: %v", err)
	}

	// geom
	wt := wkt.MarshalString(ft.Geometry)
	fmt.Println(wt)

	// 属性
	setvalue := ""
	fds := strings.Split(dt.Fields, ",")
	ps := map[string]interface{}(ft.Properties)
	for _,key := range fds{
		v,ok := ps[key]
		if ok{
			setvalue = setvalue + fmt.Sprintf("%v='%v',",key,v)
		}
	}
	setvalue = setvalue + fmt.Sprintf("geom = st_geomfromtext('%v',%v)", wt, 4326)

	// 更新
	sql := fmt.Sprintf("update %s\nset %s \nwhere gid = %v", dt.TableName, setvalue, feature_id)
	fmt.Println(sql)

	dbrst := db.Raw(sql)
	if dbrst.Error != nil{
		return fmt.Errorf("FeatureUpdate failed, err: %v", dbrst.Error)
	}
	/*if db.Raw(sql).RowsAffected == 0{
		return fmt.Errorf("FeatureUpdate failed, err: %v", "rows affected == 0")
	}*/

	return nil
}