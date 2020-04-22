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

func LikelyEncoding(file string) string {
	stat, err := os.Stat(file)
	if err != nil {
		return ""
	}
	bufsize := int64(0)
	if stat.Size() < bufsize {
		bufsize = stat.Size()
	}
	r, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer r.Close()
	buf := make([]byte, bufsize)
	rn, err := r.Read(buf)
	if err != nil {
		return ""
	}
	return Mostlike(buf[:rn])
}