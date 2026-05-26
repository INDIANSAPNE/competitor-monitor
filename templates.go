package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

var (
	templateCache = map[string]string{}
	schemaCache   = map[string]string{}
)

func LoadPromptFolder(folder string, cache map[string]string, folderName string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatalf("%s फोल्डर नहीं पढ़ सके: %v", folderName, err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		category := strings.TrimSuffix(file.Name(), ".txt")
		content, err := ioutil.ReadFile(filepath.Join(folder, file.Name()))
		if err != nil {
			log.Printf("⚠️ %s में %s नहीं पढ़ सके", folderName, file.Name())
			continue
		}
		cache[category] = string(content)
	}
	log.Printf("✅ %s से %d प्रॉम्प्ट लोड हुए", folderName, len(cache))
}

func GetSchemaPromptForCategory(category string) string {
	if prompt, ok := schemaCache[category]; ok {
		return prompt
	}
	if prompt, ok := schemaCache["सामान्य"]; ok {
		return prompt
	}
	return "टाइटल और URL से JSON बनाएं: {TITLE} {URL}"
}

func GetTemplateForCategory(category string) string {
	if prompt, ok := templateCache[category]; ok {
		return prompt
	}
	if prompt, ok := templateCache["सामान्य"]; ok {
		return prompt
	}
	return "एक SEO आर्टिकल लिखें: {TITLE}\nJSON: {EXTRA_DATA}"
}