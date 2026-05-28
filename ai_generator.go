package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func GenerateSEOTopicContent(art Article) (string, error) {
	promptTemplate := GetTemplateForCategory(art.Category)
	prompt := strings.ReplaceAll(promptTemplate, "{EXTRA_DATA}", art.ExtraDataJSON)
	prompt = strings.ReplaceAll(prompt, "{TITLE}", art.Title)
	prompt = strings.ReplaceAll(prompt, "{COMPETITOR_TITLE}", art.Title)
	prompt = strings.ReplaceAll(prompt, "{PRIMARY_KEYWORD}", art.PrimaryKeyword)

	return callDeepSeek(prompt, false)
}

func TriggerAIContent(art Article, existingPostID int) {
	fmt.Printf("  🤖 AI जनरेट कर रहा है: %s\n", art.Title)
	content, err := GenerateSEOTopicContent(art)
	if err != nil {
		log.Printf("  ❌ AI जनरेशन त्रुटि: %v\n", err)
		return
	}

	// SEO_TITLE को वर्डप्रेस पोस्ट टाइटल के रूप में उपयोग करें
	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(art.ExtraDataJSON), &meta); err == nil {
		if seoTitle, ok := meta["SEO_TITLE"].(string); ok && seoTitle != "" {
			art.Title = seoTitle
		}
	}

	newID, err := PostToWebsite(art, content, existingPostID)
	if err != nil {
		log.Printf("  ❌ पोस्ट त्रुटि: %v\n", err)
	} else {
		if newID > 0 {
			MarkTopicCovered(art.PrimaryKeyword, newID, "draft")
		} else {
			MarkTopicCovered(art.PrimaryKeyword, existingPostID, "draft")
		}
		fmt.Printf("  ✅ पोस्ट सफल (ड्राफ्ट): %s\n", art.Title)
	}
}