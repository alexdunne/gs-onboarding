package hn

import "time"

// Item represents the API response structure of an item
type Item struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	URL       string    `json:"url"`
	Score     int       `json:"score"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	Dead      bool      `json:"dead"`
	Deleted   bool      `json:"deleted"`
}
