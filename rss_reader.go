package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RSS struct {
	Channel struct {
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func FetchArticlesFromRSS(comp Competitor) ([]Article, error) {
	if comp.RSSURL == "" {
		return nil, fmt.Errorf("RSS URL खाली है")
	}

	client := http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(comp.RSSURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	var articles []Article
	if seenLinks[comp.Name] == nil {
		seenLinks[comp.Name] = make(map[string]bool)
	}

	for _, item := range rss.Channel.Items {
		url := strings.TrimSpace(item.Link)
		title := strings.TrimSpace(item.Title)
		dateStr := strings.TrimSpace(item.PubDate)

		if seenLinks[comp.Name][url] {
			continue
		}
		seenLinks[comp.Name][url] = true

		art := Article{
			Competitor: comp.Name,
			Title:      title,
			URL:        url,
			Date:       dateStr,
		}
		articles = append(articles, art)
	}
	return articles, nil
}