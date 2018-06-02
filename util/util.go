package util

import (
	"encoding/json"
	"io/ioutil"
)

func LoadJsonDataFromFile(filePath string, object interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, object); len(data) > 0 && err != nil {
		return err
	}
	return nil
}
