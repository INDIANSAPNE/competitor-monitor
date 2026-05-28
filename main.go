package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	os.MkdirAll("data", 0755)

	LoadPromptFolder("schemas", schemaCache, "Schemas")
	LoadPromptFolder("templates", templateCache, "Templates")

	config, err := LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Config लोड करने में त्रुटि: %v", err)
	}

	LoadSeenLinks()
	LoadCoveredTopics()
	PrepareExcel(config.ExcelFile)

	fmt.Println("🔍 जाँच शुरू...")
	CheckAllCompetitors(config)
	fmt.Println("✅ एकल जाँच पूरी! AI जनरेशन पूरा होने की प्रतीक्षा...")

	// AI जनरेशन और पोस्टिंग पूरी होने के लिए पर्याप्त समय
	time.Sleep(90 * time.Second)
	fmt.Println("✅ सभी कार्य समाप्त।")
}