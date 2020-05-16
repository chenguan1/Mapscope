package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// 属性字段名-类型
type Fields map[string]string

func NewFields() Fields {
	fds := make(map[string]string)
	return fds
}

func (f Fields) Set(key, value string) {
	f[key] = value
}

func (f Fields) Get(key string) (string, bool) {
	v, ok := f[key]
	if ok {
		return v, true
	} else {
		return "", false
	}
}

func (f Fields) Value() (value driver.Value, err error) {

	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}

func (f *Fields) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Invalid Scan Source")
	}

	return json.Unmarshal(s, f)
}

func (fs Fields) Keys() []string {
	//fm := map[string]string(fs)
	j := 0
	keys := make([]string, len(fs))
	for k := range fs {
		keys[j] = k
		j++
	}
	return keys
}

func (fs Fields) Types() []string {
	j := 0
	types := make([]string, len(fs))
	for _, v := range fs {
		types[j] = v
		j++
	}

	return types
}
