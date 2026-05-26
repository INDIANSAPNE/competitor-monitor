package main

import (
	"fmt"
	"log"
)

func CheckAllCompetitors(config Config) {
	for _, comp := range config.Competitors {
		fmt.Println("👉 जाँच रहा हूँ:", comp.Name)
		articles, err := FetchArticlesFromRSS(comp)
		if err != nil {
			log.Printf("❌ %s के लिए त्रुटि: %v", comp.Name, err)
			continue
		}

		for _, art := range articles {
			// कैटेगरी पहचानें
			category, err := ClassifyCategory(art.Title, art.URL)
			if err != nil {
				log.Printf("  ⚠️ कैटेगरी नहीं पहचानी: %v, डिफॉल्ट 'सामान्य'", err)
				category = "सामान्य"
			}
			art.Category = category

			// स्कीमा भरें
			primaryKW, extraJSON, err := EnrichMetadata(art.Title, art.URL, category)
			if err != nil {
				log.Printf("  ❌ मेटाडेटा एनरिचमेंट त्रुटि: %v", err)
				continue
			}
			art.PrimaryKeyword = primaryKW
			art.ExtraDataJSON = extraJSON

			// डीडुप्लीकेशन चेक
			shouldProceed, existingID := IsTopicNewOrUpdatable(primaryKW)
			if !shouldProceed {
				fmt.Printf("  ⏭️ पहले से कवर टॉपिक – छोड़ दिया\n")
				continue
			}

			// एक्सेल में जोड़ें
			AddToExcel(config.ExcelFile, art)
			fmt.Printf("  ✅ नया: %s [%s]\n", art.Title, art.Category)

			// AI जनरेशन और पोस्टिंग
			if config.AutoGenerate {
				go TriggerAIContent(art, existingID)
			}
		}
		if len(articles) == 0 {
			fmt.Println("  ℹ️ कोई नया आर्टिकल नहीं।")
		}
	}
	SaveSeenLinks()
	SaveCoveredTopics()
}