package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Bound [4]float64

func (b Bound) Value() (value driver.Value, err error) {

	data, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}

func (b *Bound) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Invalid Scan Source")
	}

	return json.Unmarshal(s, b)
}
