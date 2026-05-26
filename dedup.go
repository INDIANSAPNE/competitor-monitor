package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	coveredTopics = map[string]CoveredTopic{}
	dedupMutex    sync.Mutex
)

const coveredTopicsFile = "data/covered_topics.json"

func LoadCoveredTopics() {
	data, err := ioutil.ReadFile(coveredTopicsFile)
	if err == nil {
		json.Unmarshal(data, &coveredTopics)
		log.Printf("📚 %d पहले से कवर किए गए टॉपिक लोड हुए", len(coveredTopics))
	}
}

func SaveCoveredTopics() {
	dedupMutex.Lock()
	defer dedupMutex.Unlock()
	data, _ := json.MarshalIndent(coveredTopics, "", "  ")
	ioutil.WriteFile(coveredTopicsFile, data, 0644)
}

func IsTopicNewOrUpdatable(primaryKeyword string) (bool, int) {
	dedupMutex.Lock()
	defer dedupMutex.Unlock()

	if existing, ok := coveredTopics[primaryKeyword]; ok {
		if existing.Status == "draft" && existing.WPPostID > 0 {
			return true, existing.WPPostID
		}
		return false, 0
	}
	return true, 0
}

func MarkTopicCovered(primaryKeyword string, postID int, status string) {
	dedupMutex.Lock()
	defer dedupMutex.Unlock()
	coveredTopics[primaryKeyword] = CoveredTopic{
		PrimaryKeyword: primaryKeyword,
		WPPostID:       postID,
		Status:         status,
	}
}