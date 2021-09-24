package models

import (
	"time"

	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
)

// Item represents a hacker news item
type Item struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	URL       string    `json:"url"`
	Score     int       `json:"score"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

func Itop(item Item) *pb.Item {
	return &pb.Item{
		Id:        int32(item.ID),
		Type:      item.Type,
		Content:   item.Content,
		Url:       item.URL,
		Score:     int32(item.Score),
		Title:     item.Title,
		CreatedAt: item.CreatedAt.Unix(),
		CreatedBy: item.CreatedBy,
	}
}

func Ptoi(item *pb.Item) Item {
	return Item{
		ID:        int(item.Id),
		Type:      item.Type,
		Content:   item.Content,
		URL:       item.Url,
		Score:     int(item.Score),
		Title:     item.Title,
		CreatedAt: time.Unix(item.CreatedAt, 0),
		CreatedBy: item.CreatedBy,
	}
}
