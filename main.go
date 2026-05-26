package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Config लोड करने में त्रुटि: %v", err)
	}

	LoadSeenLinks()
	PrepareExcel(config.ExcelFile)

	fmt.Println("🔍 जाँच शुरू...")
	CheckAllCompetitors(config)
	fmt.Println("✅ एकल जाँच पूरी! अगली जाँच GitHub Actions द्वारा शेड्यूल होगी।")
}