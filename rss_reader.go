// Package rssreader provides functionality for parsing RSS feeds.
package rssreader

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

// RssItem represents a single item from an RSS feed.
type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
	RssURL      string
}

// Parse fetches and parses RSS feeds from the provided URLs asynchronously.
// It returns a slice of RssItem and any error encountered during parsing.
func Parse(ctx context.Context, urls []string) ([]RssItem, error) {
	if len(urls) == 0 {
		return []RssItem{}, nil
	}

	// Create channels to collect results from goroutines
	resultChan := make(chan []RssItem, len(urls))
	errorChan := make(chan error, len(urls))

	// Create a wait group to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Parse each URL in a separate goroutine
	for _, url := range urls {
		wg.Add(1)
		go func(feedURL string) {
			defer wg.Done()

			items, err := parseSingleFeed(ctx, feedURL)
			if err != nil {
				errorChan <- fmt.Errorf("failed to parse feed %s: %w", feedURL, err)
				return
			}

			resultChan <- items
		}(url)
	}

	// Wait for all goroutines to complete on the main thread
	wg.Wait()

	// Close channels after all goroutines are done
	close(resultChan)
	close(errorChan)

	// Collect results
	var allItems []RssItem
	var errors []error

	// Collect items from result channel
	for items := range resultChan {
		allItems = append(allItems, items...)
	}

	// Collect errors from error channel
	for err := range errorChan {
		errors = append(errors, err)
	}

	// Return error if any occurred
	if len(errors) > 0 {
		return allItems, fmt.Errorf("encountered %d errors: %v", len(errors), errors)
	}

	// Sort all items by PublishDate across all feeds
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].PublishDate.Before(allItems[j].PublishDate) || allItems[i].PublishDate.Equal(allItems[j].PublishDate)
	})

	return allItems, nil
}

// parseSingleFeed parses a single RSS feed from the given URL
func parseSingleFeed(ctx context.Context, url string) ([]RssItem, error) {
	fp := gofeed.NewParser()

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	var items []RssItem
	for _, item := range feed.Items {
		rssItem := RssItem{
			Title:       item.Title,
			Source:      feed.Title,
			SourceURL:   url,
			Link:        item.Link,
			Description: item.Description,
			RssURL:      url,
		}

		// Handle publish date
		switch {
		case item.PublishedParsed != nil:
			rssItem.PublishDate = *item.PublishedParsed
		case item.UpdatedParsed != nil:
			rssItem.PublishDate = *item.UpdatedParsed
		default:
			rssItem.PublishDate = time.Time{}
		}

		items = append(items, rssItem)
	}

	return items, nil
}
