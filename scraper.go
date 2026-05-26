package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CheckAllCompetitors(config Config) {
	for _, comp := range config.Competitors {
		fmt.Println("👉 जाँच रहा हूँ:", comp.Name)
		articles, err := FetchNewArticles(comp)
		if err != nil {
			log.Printf("❌ %s के लिए त्रुटि: %v", comp.Name, err)
			continue
		}
		for _, art := range articles {
			AddToExcel(config.ExcelFile, art)
			fmt.Printf("  ✅ नया: %s - %s\n", art.Title, art.URL)

			if config.AutoGenerate {
				// AI जनरेशन और पोस्टिंग शुरू (goroutine में, ताकि तेज़ी से हो)
				go TriggerAIContent(art)
			}
		}
		if len(articles) == 0 {
			fmt.Println("  ℹ️ कोई नया आर्टिकल नहीं।")
		}
	}
	SaveSeenLinks()
}

func FetchNewArticles(comp Competitor) ([]Article, error) {
	resp, err := http.Get(comp.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var newArticles []Article
	if seenLinks[comp.Name] == nil {
		seenLinks[comp.Name] = make(map[string]bool)
	}

	doc.Find(comp.ArticleSelector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		fullURL := ToAbsoluteURL(comp.URL, href)
		title := strings.TrimSpace(s.Text())

		if seenLinks[comp.Name][fullURL] {
			return
		}
		seenLinks[comp.Name][fullURL] = true

		dateStr := ""
		if comp.DateSelector != "" {
			dateEl := s.Parent().Find(comp.DateSelector).First()
			if dateEl.Length() == 0 {
				dateEl = s.Closest("article").Find(comp.DateSelector).First()
			}
			if comp.DateAttr != "" {
				dateStr, _ = dateEl.Attr(comp.DateAttr)
			} else {
				dateStr = strings.TrimSpace(dateEl.Text())
			}
		}

		art := Article{
			Competitor: comp.Name,
			Title:      title,
			URL:        fullURL,
			Date:       dateStr,
		}
		newArticles = append(newArticles, art)
	})

	return newArticles, nil
}