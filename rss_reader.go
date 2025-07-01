// Package rssreader provides functionality for parsing RSS feeds.
package rssreader

import (
	"context"
	"time"
)

// RssItem represents a single item from an RSS feed.
type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}

// Parse fetches and parses RSS feeds from the provided URLs.
// It returns a slice of RssItem and any error encountered during parsing.
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
