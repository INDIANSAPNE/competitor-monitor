package main

type Article struct {
	Competitor     string
	Title          string
	URL            string
	Date           string
	Category       string
	PrimaryKeyword string
	ExtraDataJSON  string
	Tags           []string
}

type CoveredTopic struct {
	PrimaryKeyword string `json:"primary_keyword"`
	WPPostID       int    `json:"wp_post_id"`
	Status         string `json:"status"`
}

type Competitor struct {
	Name   string `json:"name"`
	RSSURL string `json:"rss_url"`
}