package utils

import (
	"archive/zip"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func EnsurePathExist(path string) {
	path, _ = filepath.Abs(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
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

// 解压zip文件，到指定的目录中，返回解压后的所有文件
func UnzipFile(zipfile string, outdir string) ([]string, error) {
	files := make([]string, 0)
	if outdir == "" {
		outdir = strings.TrimSuffix(zipfile, filepath.Ext(zipfile))
	}
	EnsurePathExist(outdir)

	zr, err := zip.OpenReader(zipfile)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	for _, f := range zr.File {
		name := f.Name
		if f.NonUTF8 {
			ns, err := simplifiedchinese.GB18030.NewDecoder().String(f.Name)
			if err == nil {
				name = ns
			} else {
				logrus.Warn("unzip file error,", f.Name)
				return nil, err
			}
		}
		pn := filepath.Join(outdir, name)
		if f.FileInfo().IsDir() {
			EnsurePathExist(pn)
			continue
		}
		ext := filepath.Ext(name)
		pn = strings.TrimSuffix(pn, ext) + strings.ToLower(ext)
		w, err := os.Create(pn)
		if err != nil {
			logrus.Warnf("Cannot unzip %s: %v", zipfile, err)
			return nil, err
		}
		defer w.Close()
		r, err := f.Open()
		if err != nil {
			logrus.Warnf("Cannot unzip %s: %v", zipfile, err)
			return nil, err
		}
		defer r.Close()
		_, err = io.Copy(w, r)
		if err != nil {
			logrus.Warnf("Cannot unzip %s: %v", zipfile, err)
			return nil, err
		}
		files = append(files, pn)
	}
	return files, nil
}

// 根据后缀名过滤 exts .json.geojson.shp.zip
func FilterByExt(files []string, exts string) []string {
	fs := make([]string, 0)
	for _, f := range files {
		if strings.Contains(exts, strings.ToLower(filepath.Ext(f))) {
			fs = append(fs, f)
		}
	}
	return fs
}

// 拷贝文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func PathExist(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
