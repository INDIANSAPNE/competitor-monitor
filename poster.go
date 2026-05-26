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

func PostToWebsite(title, content string) error {
	wpURL := os.Getenv("WP_URL")
	wpUser := os.Getenv("WP_USER")
	wpPass := os.Getenv("WP_PASS")

	if wpURL == "" || wpUser == "" || wpPass == "" {
		return fmt.Errorf("WP_URL, WP_USER, WP_PASS एनवायरनमेंट वेरिएबल सेट नहीं हैं")
	}

	apiURL := strings.TrimRight(wpURL, "/") + "/wp-json/wp/v2/posts"
	postData := map[string]interface{}{
		"title":   title,
		"content": content,
		"status":  "draft", // "publish" से सीधे पब्लिश होगा
	}
	jsonData, _ := json.Marshal(postData)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(wpUser, wpPass)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("पोस्ट रिक्वेस्ट त्रुटि: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("पोस्ट नहीं बन सका (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}