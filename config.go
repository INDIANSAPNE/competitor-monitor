package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Competitors  []Competitor `json:"competitors"`
	ExcelFile    string       `json:"excel_file"`
	AutoGenerate bool         `json:"auto_generate"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}