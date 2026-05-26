package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func PostToWebsite(art Article, content string, existingPostID int) (int, error) {
	wpURL := os.Getenv("WP_URL")
	wpUser := os.Getenv("WP_USER")
	wpPass := os.Getenv("WP_PASS")

	if wpURL == "" || wpUser == "" || wpPass == "" {
		return 0, fmt.Errorf("WP_URL, WP_USER, WP_PASS सेट नहीं हैं")
	}

	categoryID := getCategoryIDByName(art.Category, wpURL, wpUser, wpPass)

	postData := map[string]interface{}{
		"title":   art.Title,
		"content": content,
		"status":  "draft",
	}
	if categoryID > 0 {
		postData["categories"] = []int{categoryID}
	}

	jsonData, _ := json.Marshal(postData)

	var req *http.Request
	var err error
	client := &http.Client{Timeout: 30 * time.Second}

	if existingPostID > 0 {
		apiURL := fmt.Sprintf("%s/wp-json/wp/v2/posts/%d", strings.TrimRight(wpURL, "/"), existingPostID)
		req, err = http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	} else {
		apiURL := strings.TrimRight(wpURL, "/") + "/wp-json/wp/v2/posts"
		req, err = http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	}
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(wpUser, wpPass)

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("पोस्ट रिक्वेस्ट त्रुटि: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("पोस्ट नहीं बन सका (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if id, ok := result["id"].(float64); ok {
		return int(id), nil
	}
	return existingPostID, nil
}

func getCategoryIDByName(name, wpURL, wpUser, wpPass string) int {
	apiURL := fmt.Sprintf("%s/wp-json/wp/v2/categories?search=%s", strings.TrimRight(wpURL, "/"), name)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.SetBasicAuth(wpUser, wpPass)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		var categories []map[string]interface{}
		json.Unmarshal(body, &categories)
		if len(categories) > 0 {
			if id, ok := categories[0]["id"].(float64); ok {
				return int(id)
			}
		}
	}
	return 0
}