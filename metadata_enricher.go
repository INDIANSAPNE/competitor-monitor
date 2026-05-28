package main

import (
	"encoding/json"
	"strings"
)

func EnrichMetadata(art *Article, category string) error {
	promptTemplate := GetSchemaPromptForCategory(category)
	prompt := strings.ReplaceAll(promptTemplate, "{TITLE}", art.Title)
	prompt = strings.ReplaceAll(prompt, "{URL}", art.URL)

	responseJSON, err := callDeepSeek(prompt, true)
	if err != nil {
		return err
	}

	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(responseJSON), &meta); err != nil {
		return err
	}

	// प्राइमरी कीवर्ड
	primaryKW := art.Title
	if kw, ok := meta["PRIMARY_KEYWORD"].(string); ok {
		primaryKW = kw
	} else if kw, ok := meta["PRIMARY_TOPIC"].(string); ok {
		primaryKW = kw
	}
	art.PrimaryKeyword = primaryKW
	art.ExtraDataJSON = responseJSON

	// टैग्स निकालें
	art.Tags = []string{}
	if tagsRaw, ok := meta["TAGS"]; ok {
		switch v := tagsRaw.(type) {
		case []interface{}:
			for _, t := range v {
				if s, ok := t.(string); ok {
					art.Tags = append(art.Tags, s)
				}
			}
		}
	}
	// कम से कम 3 टैग्स सुनिश्चित करें
	if len(art.Tags) < 3 {
		if art.Category != "" {
			art.Tags = append(art.Tags, art.Category)
		}
		if art.PrimaryKeyword != "" {
			art.Tags = append(art.Tags, art.PrimaryKeyword)
		}
		if len(art.Tags) < 3 {
			art.Tags = append(art.Tags, "सामान्य")
		}
	}
	// 5 से अधिक हो तो काटें
	if len(art.Tags) > 5 {
		art.Tags = art.Tags[:5]
	}

	return nil
}