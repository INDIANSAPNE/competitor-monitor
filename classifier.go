package main

import (
	"fmt"
	"strings"
)

func ClassifyCategory(title, url string) (string, error) {
	categories := make([]string, 0, len(schemaCache))
	for cat := range schemaCache {
		categories = append(categories, cat)
	}
	if len(categories) == 0 {
		return "", fmt.Errorf("कोई कैटेगरी स्कीमा नहीं है")
	}

	prompt := fmt.Sprintf(`निम्नलिखित आर्टिकल टाइटल और URL को देखकर बताएं कि यह इनमें से किस कैटेगरी में आता है: %s
सिर्फ कैटेगरी का ठीक वही नाम लिखें जो लिस्ट में है, और कुछ न लिखें।
टाइटल: %s
URL: %s`, strings.Join(categories, ", "), title, url)

	response, err := callDeepSeek(prompt, false)
	if err != nil {
		return "", err
	}
	cat := strings.TrimSpace(response)
	if _, ok := schemaCache[cat]; ok {
		return cat, nil
	}
	// फ़ॉलबैक
	if _, ok := schemaCache["सामान्य"]; ok {
		return "सामान्य", nil
	}
	return "", fmt.Errorf("पहचानी गई कैटेगरी '%s' स्कीमा में नहीं है", cat)
}