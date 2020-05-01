package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
)

func CsvReader(r io.Reader, encoding string) (*csv.Reader, error) {
	switch encoding {
	case "gbk", "big5", "gb18030":
		decoder := mahonia.NewDecoder(encoding)
		if decoder == nil {
			return csv.NewReader(r), fmt.Errorf(`create %s encoder error`, encoding)
		}
		dreader := decoder.NewReader(r)
		return csv.NewReader(dreader), nil
	default:
		return csv.NewReader(r), nil
	}
}
