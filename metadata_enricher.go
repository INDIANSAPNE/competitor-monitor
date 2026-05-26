package main

import (
	"encoding/json"
	"strings"
)

func EnrichMetadata(title, url, category string) (string, string, error) {
	promptTemplate := GetSchemaPromptForCategory(category)
	prompt := strings.ReplaceAll(promptTemplate, "{TITLE}", title)
	prompt = strings.ReplaceAll(prompt, "{URL}", url)

	responseJSON, err := callDeepSeek(prompt, true) // jsonMode = true
	if err != nil {
		return "", "", err
	}

	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(responseJSON), &meta); err != nil {
		meta = make(map[string]interface{})
	}

	primaryKW := title
	if kw, ok := meta["PRIMARY_KEYWORD"].(string); ok {
		primaryKW = kw
	} else if kw, ok := meta["PRIMARY_TOPIC"].(string); ok {
		primaryKW = kw
	}

	return primaryKW, responseJSON, nil
}