package utils

import (
	"encoding/json"
	"io"
	"os"
)

func LoadJson(filePath string, model interface{}) *interface{} {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	r := io.Reader(f)
	if err = json.NewDecoder(r).Decode(&model); err != nil {
		panic(err)
	}

	f.Close()
	return &model
}

func DumpJson(filePath string, model interface{}) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
	defer fp.Close()
	if err != nil {
		return false, err
	}

	data, err := json.Marshal(model)
	if err != nil {
		return false, err
	}

	_, err = fp.Write(data)
	if err != nil {
		return false, err
	}

	return true, nil
}