package utils

import (
	"os"
	"path/filepath"
)

func EnsurePathExist(path string)  {
	path,_ = filepath.Abs(path)
	if _,err := os.Stat(path); os.IsNotExist(err){
		os.MkdirAll(path, 0666)
	}
}
