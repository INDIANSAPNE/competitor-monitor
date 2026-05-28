package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

			// मेटाडेटा एनरिच करें (टैग्स सहित)
			err = EnrichMetadata(&art, category)
			if err != nil {
				log.Printf("  ❌ मेटाडेटा एनरिचमेंट त्रुटि: %v", err)
				continue
			}

			// फ़िल्टर चेक (केवल अनुमत कैटेगरी और वैकेंसी > 100)
			if !shouldGenerateDraft(art) {
				fmt.Printf("  ⏭️ फ़िल्टर आउट: %s [%s]\n", art.Title, art.Category)
				continue
			}

			// डीडुप्लीकेशन चेक
			shouldProceed, existingID := IsTopicNewOrUpdatable(art.PrimaryKeyword)
			if !shouldProceed {
				fmt.Printf("  ⏭️ पहले से कवर टॉपिक – छोड़ दिया\n")
				continue
			}

			// एक्सेल में जोड़ें
			AddToExcel(config.ExcelFile, art)
			fmt.Printf("  ✅ नया: %s [%s] (टैग्स: %v)\n", art.Title, art.Category, art.Tags)

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

// shouldGenerateDraft अनुमत कैटेगरी और वैकेंसी > 100 (नौकरियों के लिए) चेक करता है
func shouldGenerateDraft(art Article) bool {
	allowedCategories := map[string]bool{
		"सरकारी नौकरियाँ":          true,
		"प्राइवेट नौकरियाँ":        true,
		"इंटरनेशनल नौकरियाँ":      true,
		"परिणाम (Results)":         true,
		"प्रवेश पत्र (Admit Card)": true,
		"आंसर की (Answer Key)":     true,
		"कटऑफ (Cutoff)":            true,
		"सिलेबस":                   true,
		"प्रीवियस पेपर":            true,
		"स्कॉलरशिप":                true,
		"करियर गाइडेंस":            true,
		"स्किल एंड रोज़गार":        true,
		"महिला / दिव्यांग":         true,
		"विशेष जानकारी":            true,
		"सामान्य":                  true,
	}

	if !allowedCategories[art.Category] {
		return false
	}

	// नौकरी कैटेगरी के लिए वैकेंसी चेक
	jobCategories := map[string]bool{
		"सरकारी नौकरियाँ":     true,
		"प्राइवेट नौकरियाँ":   true,
		"इंटरनेशनल नौकरियाँ": true,
	}
	if jobCategories[art.Category] {
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(art.ExtraDataJSON), &meta); err != nil {
			return false
		}
		vac := getVacancy(meta)
		return vac > 100
	}
	return true
}

func getVacancy(meta map[string]interface{}) int {
	keys := []string{"TOTAL_VACANCIES", "TOTAL_VACANCY", "VACANCY", "NO_OF_OPENINGS"}
	for _, key := range keys {
		if val, ok := meta[key]; ok {
			switch v := val.(type) {
			case float64:
				return int(v)
			case string:
				if i, err := strconv.Atoi(v); err == nil {
					return i
				}
			}
		}
	}
	return 0
}