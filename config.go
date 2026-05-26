package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Competitors []Competitor `json:"competitors"`
	ExcelFile   string       `json:"excel_file"`
	AutoGenerate bool        `json:"auto_generate"`
}

type Competitor struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	ArticleSelector string `json:"article_selector"`
	DateSelector    string `json:"date_selector"`
	DateAttr        string `json:"date_attr"`
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