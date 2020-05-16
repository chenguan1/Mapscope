package utils

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
)

func FontName(path string) (string, error) {
	ttf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	font, err := truetype.Parse(ttf)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", font.Name(truetype.NameIDFontFamily), font.Name(truetype.NameIDFontSubfamily)), nil
}
