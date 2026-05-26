package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/xuri/excelize/v2"
)
var seenLinks = map[string]map[string]bool{}

func LoadSeenLinks() {
	data, err := ioutil.ReadFile("data/seen_links.json")
	if err == nil {
		json.Unmarshal(data, &seenLinks)
	}
}

func SaveSeenLinks() {
	data, _ := json.MarshalIndent(seenLinks, "", "  ")
	ioutil.WriteFile("data/seen_links.json", data, 0644)
}

func PrepareExcel(filename string) {
	path := "data/" + filename
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f := excelize.NewFile()
		sheet := "Articles"
		f.SetSheetName("Sheet1", sheet)
		headers := []string{"Competitor", "Title", "URL", "Date", "Category", "Primary Keyword", "Extra Data (JSON)"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
		}
		f.SaveAs(path)
		log.Println("📄 नई एक्सेल फ़ाइल बनाई:", path)
	}
}

func AddToExcel(filename string, art Article) {
	path := "data/" + filename
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Println("एक्सेल खोलने में त्रुटि:", err)
		return
	}
	sheet := "Articles"
	rows, _ := f.GetRows(sheet)
	rowNum := len(rows) + 1

	values := []interface{}{art.Competitor, art.Title, art.URL, art.Date, art.Category, art.PrimaryKeyword, art.ExtraDataJSON}
	for col, val := range values {
		cell, _ := excelize.CoordinatesToCellName(col+1, rowNum)
		f.SetCellValue(sheet, cell, val)
	}
	f.SaveAs(path)
}

func callDeepSeek(prompt string, jsonMode bool) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY सेट नहीं है")
	}

	reqBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 2000,
	}
	if jsonMode {
		reqBody["response_format"] = map[string]string{"type": "json_object"}
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("DeepSeek API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return "", fmt.Errorf("कोई उत्तर नहीं")
	}
	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	return message["content"].(string), nil
}