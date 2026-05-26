package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

// seenLinks को ग्लोबल रखें ताकि scraper और main दोनों में एक्सेस हो
var seenLinks = map[string]map[string]bool{}

func LoadSeenLinks() {
	file, err := ioutil.ReadFile("seen_links.json")
	if err == nil {
		json.Unmarshal(file, &seenLinks)
	}
}

func SaveSeenLinks() {
	data, _ := json.MarshalIndent(seenLinks, "", "  ")
	ioutil.WriteFile("seen_links.json", data, 0644)
}

func PrepareExcel(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f := excelize.NewFile()
		sheet := "Articles"
		f.SetSheetName("Sheet1", sheet)
		headers := []string{"Competitor", "Title", "URL", "Date"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
		}
		f.SaveAs(filename)
		log.Println("📄 नई एक्सेल फ़ाइल बनाई:", filename)
	}
}

func AddToExcel(filename string, art Article) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Println("एक्सेल खोलने में त्रुटि:", err)
		return
	}
	sheet := "Articles"
	rows, _ := f.GetRows(sheet)
	rowNum := len(rows) + 1

	values := []interface{}{art.Competitor, art.Title, art.URL, art.Date}
	for col, val := range values {
		cell, _ := excelize.CoordinatesToCellName(col+1, rowNum)
		f.SetCellValue(sheet, cell, val)
	}
	f.SaveAs(filename)
}

func ToAbsoluteURL(baseURL, href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	if strings.HasPrefix(href, "/") {
		parts := strings.Split(baseURL, "/")
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2] + href
		}
	}
	return baseURL + "/" + href
}