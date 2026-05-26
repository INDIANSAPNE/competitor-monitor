package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// डेटा फोल्डर बनाएँ अगर नहीं है
	os.MkdirAll("data", 0755)

	// स्कीमा और टेम्पलेट फोल्डर से प्रॉम्प्ट लोड करें
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
	time.Sleep(30 * time.Second)
	fmt.Println("✅ एकल जाँच पूरी!")
}
