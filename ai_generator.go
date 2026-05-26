package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GenerateSEOTopicContent(topic, competitorURL string) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY सेट नहीं है")
	}

	prompt := fmt.Sprintf(`आप एक पेशेवर SEO कंटेंट राइटर हैं। नीचे दिए गए प्रतिस्पर्धी आर्टिकल का टॉपिक और यूआरएल दिया गया है। इस टॉपिक पर एक बेहतरीन, अनोखा, और SEO-अनुकूल हिंदी आर्टिकल लिखें। आर्टिकल कम से कम 500 शब्दों का हो, जिसमें उचित हेडिंग्स (H2, H3) हों, कीवर्ड का स्वाभाविक उपयोग हो, और पाठक के लिए बेहद उपयोगी हो। कंटेंट ओरिजिनल हो, कॉपी न हो। प्रतिस्पर्धी लिंक: %s
टॉपिक: %s`, competitorURL, topic)

	reqBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 2000,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		choice := choices[0].(map[string]interface{})
		message := choice["message"].(map[string]interface{})
		content := message["content"].(string)
		return content, nil
	}
	return "", fmt.Errorf("AI से प्रतिक्रिया नहीं मिली: %s", string(body))
}

func TriggerAIContent(art Article) {
	fmt.Printf("  🤖 AI जनरेट कर रहा है: %s\n", art.Title)
	content, err := GenerateSEOTopicContent(art.Title, art.URL)
	if err != nil {
		log.Printf("  ❌ AI जनरेशन त्रुटि: %v\n", err)
		return
	}
	// अब वेबसाइट पर पोस्ट करें
	err = PostToWebsite(art.Title, content)
	if err != nil {
		log.Printf("  ❌ पोस्ट त्रुटि: %v\n", err)
	} else {
		fmt.Printf("  ✅ पोस्ट सफल: %s\n", art.Title)
	}
}