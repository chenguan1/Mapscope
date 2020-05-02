package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// 根据字体名获取字体
func FontGet(name string) (*models.Font, error) {
	var ft models.Font
	err := database.Get().Where(models.Font{Id: name}).Find(&ft).Error
	if err != nil {
		return nil, fmt.Errorf("FontGet err: %v", err)
	}
	return &ft, nil
}

// 获取某用户的所有dataset
func FontList(user string) ([]models.Font, error) {
	var fts []models.Font
	err := database.Get().Where(models.Font{Owner: user}).Find(&fts).Error
	if err != nil {
		return nil, fmt.Errorf("FontList err: %v", err)
	}
	return fts, nil
}

func FontDelete(fontname string) error {
	return database.Get().Delete(models.Font{Id:fontname}).Error
}


// 载入字体
func FontLoad(path string) (*models.Font, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != models.PBFONTEXT{
		return nil, fmt.Errorf("FontLoad err: ext %s is unsurpported now.", ext)
	}

	fstat,err := os.Stat(path)
	if err != nil{
		return nil, fmt.Errorf("FontLoad err: %v", err)
	}

	ext = filepath.Ext(path)
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, ext)

	out := &models.Font{
		Id:          name,
		Name:        name,
		Owner:       "mapscope",
		Path:        path,
		Size:        fstat.Size(),
		Type:        models.PBFONTEXT,
		Public:      1,
		Compression: 0,
	}

	return out,nil
}

// 获取字体切片
func fontPbf(name, fontrange string) ([]byte, error) {
	ft,err := FontGet(name)
	if err != nil{
		return nil, err
		/*if ft,err = FontGet(models.DEFAULTFONT); err != nil{
			if err != nil {
				return nil, fmt.Errorf("font %s not exist, use default %s failed. err: %v", name, models.DEFAULTFONT, err)
			}
		}*/
	}

	db3,err := database.OpenSqlite(ft.Path)
	if err != nil{
		return nil, err
	}

	var data []byte
	err = db3.DB().QueryRow("select data from fonts where range = ?", fontrange).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// 获取Glyphs
func FontGlyphs(fonts []string, fontrange string) ([]byte, error) {
	var pbf []byte
	switch len(fonts) {
	case 0:
		return nil, fmt.Errorf("FontGlyphs, fontstack is nil.")
	case 1:
		data, err := fontPbf(fonts[0],fontrange)
		if err != nil {
			//return  nil, fmt.Errorf("FontGlyphs err: %v", err)
			fmt.Printf("FontGlyphs err: %v\n", err)
			if data, err = fontPbf(models.DEFAULTFONT,fontrange); err != nil{
				return  nil, fmt.Errorf("FontGlyphs font not found, use default err: %v", err)
			}
		}
		pbf = data
	default: //multi
		var fs []*models.Font
		hasdefault := false
		haslost := false
		for _, font := range fonts {
			if font == models.DEFAULTFONT {
				hasdefault = true
			}

			f,err := FontGet(font)
			if err != nil {
				haslost = true
				continue
			}
			fs = append(fs, f)
		}
		//没有默认字体且有丢失字体,则加载默认字体
		if !hasdefault && haslost {
			f,err := FontGet(models.DEFAULTFONT)
			if err == nil {
				fs = append(fs, f)
			}
		}

		contents := make([][]byte, len(fs))
		var wg sync.WaitGroup
		//need define func, can't use sugar ":="
		getFontPBF := func(f *models.Font, fontrange string, index int) {
			//fallbacks unchanging
			defer wg.Done()
			d, err := database.OpenSqlite(f.Path)
			if err != nil{
				fmt.Printf("FontGlyphs getFontPBF %v err: %v", f.Path ,err)
				return
			}
			err = d.DB().QueryRow("select data from fonts where range = ?", fontrange).Scan(&contents[index])
			if err != nil {
				fmt.Printf("FontGlyphs getFontPBF %v err: %v", f.Path ,err)
				if err == sql.ErrNoRows {
					return
				}
				return
			}
		}
		for i, f := range fs {
			wg.Add(1)
			go getFontPBF(f, fontrange, i)
		}
		wg.Wait()

		//if  getFontPBF can't get content,the buffer array is nil, remove the nils
		var buffers [][]byte
		var bufFonts []string
		for i, buf := range contents {
			if nil == buf {
				continue
			}
			buffers = append(buffers, buf)
			bufFonts = append(bufFonts, fonts[i])
		}
		if len(buffers) != len(bufFonts) {
			fmt.Printf("FontGlyphs getFontPBF len(buffers) != len(fonts)")
		}
		if 0 == len(buffers) {
			return nil, fmt.Errorf("FontGlyphs, empty pbf font")
		}
		if 1 == len(buffers) {
			pbf = buffers[0]
		} else {
			c, err := combine(buffers, bufFonts)
			if err != nil {
				fmt.Printf("FontGlyphs combine buffers error:", err)
				pbf = buffers[0]
			} else {
				pbf = c
			}
		}
	}

	return pbf, nil
}

//Combine 多字体请求合并
func combine(buffers [][]byte, fontstack []string) ([]byte, error) {
	coverage := make(map[uint32]bool)
	result := &models.Glyphs{}
	for i, buf := range buffers {
		pbf := &models.Glyphs{}
		err := proto.Unmarshal(buf, pbf)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		if stacks := pbf.GetStacks(); stacks != nil && len(stacks) > 0 {
			stack := stacks[0]
			if 0 == i {
				for _, gly := range stack.Glyphs {
					coverage[gly.GetId()] = true
				}
				result = pbf
			} else {
				for _, gly := range stack.Glyphs {
					if !coverage[gly.GetId()] {
						result.Stacks[0].Glyphs = append(result.Stacks[0].Glyphs, gly)
						coverage[gly.GetId()] = true
					}
				}
				result.Stacks[0].Name = proto.String(result.Stacks[0].GetName() + "," + stack.GetName())
			}
		}

		if fontstack != nil {
			result.Stacks[0].Name = proto.String(strings.Join(fontstack, ","))
		}
	}

	glys := result.Stacks[0].GetGlyphs()

	sort.Slice(glys, func(i, j int) bool {
		return glys[i].GetId() < glys[j].GetId()
	})

	return proto.Marshal(result)
}