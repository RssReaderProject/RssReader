package rssreader

import (
	"context"
	"time"
)

type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}

func Parse(ctx context.Context, urls []string) ([]RssItem, error) {
	items := []RssItem{
		{
			Title:       "Test",
			Source:      "Test",
			SourceURL:   "Test",
			Link:        "Test",
			PublishDate: time.Now(),
			Description: "Test",
		},
	}
	return items, nil
}
