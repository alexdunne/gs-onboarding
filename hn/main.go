package hn

import "time"

type ItemType string

const (
	StoryType ItemType = "story"
	JobType   ItemType = "job"
)

type Item struct {
	ID        int       `json:"id"`
	Type      ItemType  `json:"type"`
	Content   string    `json:"content"`
	URL       string    `json:"url"`
	Score     int       `json:"score"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

type Items []*Item
