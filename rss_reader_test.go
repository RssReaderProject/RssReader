package rssreader

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testServer creates a mock HTTP server that returns the given RSS content
func testServer(content string, contentType string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(content))
	}))
}

func TestParse_EmptyURLs(t *testing.T) {
	ctx := context.Background()
	items, err := Parse(ctx, []string{})

	if err != nil {
		t.Errorf("Expected no error for empty URLs, got: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected empty slice for empty URLs, got %d items", len(items))
	}
}

func TestParse_SingleValidFeed(t *testing.T) {
	// Create a mock server that returns valid RSS
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
    <link>http://example.com</link>
    <description>A test RSS feed</description>
    <item>
      <title>Test Article 1</title>
      <link>http://example.com/article1</link>
      <description>This is a test article</description>
      <pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
    </item>
    <item>
      <title>Test Article 2</title>
      <link>http://example.com/article2</link>
      <description>This is another test article</description>
      <pubDate>Tue, 03 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// Check first item
	if items[0].Title != "Test Article 1" {
		t.Errorf("Expected title 'Test Article 1', got '%s'", items[0].Title)
	}
	if items[0].Source != "Test Feed" {
		t.Errorf("Expected source 'Test Feed', got '%s'", items[0].Source)
	}
	if items[0].SourceURL != server.URL {
		t.Errorf("Expected source URL '%s', got '%s'", server.URL, items[0].SourceURL)
	}
	if items[0].Link != "http://example.com/article1" {
		t.Errorf("Expected link 'http://example.com/article1', got '%s'", items[0].Link)
	}
	if items[0].Description != "This is a test article" {
		t.Errorf("Expected description 'This is a test article', got '%s'", items[0].Description)
	}

	// Check second item
	if items[1].Title != "Test Article 2" {
		t.Errorf("Expected title 'Test Article 2', got '%s'", items[1].Title)
	}
}

func TestParse_MultipleValidFeeds(t *testing.T) {
	// Create two mock servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Feed 1</title>
    <item>
      <title>Article from Feed 1</title>
      <link>http://example.com/feed1/article</link>
      <description>Article from first feed</description>
      <pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Feed 2</title>
    <item>
      <title>Article from Feed 2</title>
      <link>http://example.com/feed2/article</link>
      <description>Article from second feed</description>
      <pubDate>Tue, 03 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server2.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server1.URL, server2.URL})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items total, got %d", len(items))
	}

	// Check that we have items from both feeds
	feed1Found := false
	feed2Found := false
	for _, item := range items {
		if item.Source == "Feed 1" && item.Title == "Article from Feed 1" {
			feed1Found = true
		}
		if item.Source == "Feed 2" && item.Title == "Article from Feed 2" {
			feed2Found = true
		}
	}

	if !feed1Found {
		t.Error("Expected to find item from Feed 1")
	}
	if !feed2Found {
		t.Error("Expected to find item from Feed 2")
	}
}

func TestParse_InvalidURL(t *testing.T) {
	ctx := context.Background()
	items, err := Parse(ctx, []string{"http://invalid-url-that-does-not-exist.com/feed"})

	// Should get items (empty slice) but also an error
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}

	if len(items) != 0 {
		t.Errorf("Expected empty items slice for invalid URL, got %d items", len(items))
	}
}

func TestParse_MixedValidAndInvalidURLs(t *testing.T) {
	// Create a valid mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Valid Feed</title>
    <item>
      <title>Valid Article</title>
      <link>http://example.com/article</link>
      <description>Valid article</description>
      <pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL, "http://invalid-url-that-does-not-exist.com/feed"})

	// Should get items from valid feed but also an error
	if err == nil {
		t.Error("Expected error for mixed valid/invalid URLs, got nil")
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item from valid feed, got %d", len(items))
	}

	if items[0].Title != "Valid Article" {
		t.Errorf("Expected title 'Valid Article', got '%s'", items[0].Title)
	}
}

func TestParse_ServerError(t *testing.T) {
	// Create a server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "Internal Server Error")
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err == nil {
		t.Error("Expected error for server error, got nil")
	}

	if len(items) != 0 {
		t.Errorf("Expected empty items slice for server error, got %d items", len(items))
	}
}

func TestParse_InvalidXML(t *testing.T) {
	// Create a server that returns invalid XML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "This is not valid XML")
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}

	if len(items) != 0 {
		t.Errorf("Expected empty items slice for invalid XML, got %d items", len(items))
	}
}

func TestParse_Timeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second) // Longer than the 30-second timeout
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel></channel></rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	items, err := Parse(ctx, []string{server.URL})

	if err == nil {
		t.Error("Expected error for timeout, got nil")
	}

	if len(items) != 0 {
		t.Errorf("Expected empty items slice for timeout, got %d items", len(items))
	}
}

func TestParse_ItemWithoutPublishedDate(t *testing.T) {
	// Create a server that returns RSS without pubDate
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
    <item>
      <title>Article without date</title>
      <link>http://example.com/article</link>
      <description>Article without published date</description>
    </item>
  </channel>
</rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	// Check that PublishDate is set to zero time when no date is available
	if !items[0].PublishDate.IsZero() {
		t.Errorf("Expected PublishDate to be zero time when no date available, got %v", items[0].PublishDate)
	}
}

func TestParse_ItemWithUpdatedDate(t *testing.T) {
	// Create a server that returns Atom format which has proper updated field
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Test Feed</title>
  <entry>
    <title>Article with updated date</title>
    <link href="http://example.com/article"/>
    <summary>Article with updated date</summary>
    <updated>2006-01-02T15:04:05Z</updated>
  </entry>
</feed>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	expectedDate := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	if !items[0].PublishDate.Equal(expectedDate) {
		t.Errorf("Expected PublishDate %v, got %v", expectedDate, items[0].PublishDate)
	}
}

func TestParse_ConcurrentRequests(t *testing.T) {
	// Create multiple servers that return different feeds
	servers := make([]*httptest.Server, 5)
	urls := make([]string, 5)

	for i := 0; i < 5; i++ {
		i := i // Capture loop variable
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/rss+xml")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Feed %d</title>
    <item>
      <title>Article from Feed %d</title>
      <link>http://example.com/feed%d/article</link>
      <description>Article from feed %d</description>
      <pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`, i+1, i+1, i+1, i+1)
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		}))
		urls[i] = servers[i].URL
		defer servers[i].Close()
	}

	ctx := context.Background()
	items, err := Parse(ctx, urls)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 5 {
		t.Fatalf("Expected 5 items total, got %d", len(items))
	}

	// Check that we have items from all feeds
	feedCounts := make(map[string]int)
	for _, item := range items {
		feedCounts[item.Source]++
	}

	for i := 1; i <= 5; i++ {
		feedName := fmt.Sprintf("Feed %d", i)
		if feedCounts[feedName] != 1 {
			t.Errorf("Expected 1 item from %s, got %d", feedName, feedCounts[feedName])
		}
	}
}

func TestParse_EmptyFeed(t *testing.T) {
	// Create a server that returns valid RSS but with no items
	server := testServer(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Empty Feed</title>
    <link>http://example.com</link>
    <description>An empty RSS feed</description>
  </channel>
</rss>`, "application/rss+xml")
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items from empty feed, got %d", len(items))
	}
}

func TestParse_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel></channel></rss>`)
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context immediately
	cancel()

	items, err := Parse(ctx, []string{server.URL})

	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}

	if len(items) != 0 {
		t.Errorf("Expected empty items slice for cancelled context, got %d items", len(items))
	}
}

func TestParse_ItemWithMissingFields(t *testing.T) {
	// Create a server that returns RSS with items missing some fields
	server := testServer(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
    <item>
      <title>Article with missing fields</title>
      <!-- Missing link and description -->
      <pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`, "application/rss+xml")
	defer server.Close()

	ctx := context.Background()
	items, err := Parse(ctx, []string{server.URL})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	// Check that missing fields are handled gracefully
	if items[0].Title != "Article with missing fields" {
		t.Errorf("Expected title 'Article with missing fields', got '%s'", items[0].Title)
	}
	if items[0].Link != "" {
		t.Errorf("Expected empty link for missing field, got '%s'", items[0].Link)
	}
	if items[0].Description != "" {
		t.Errorf("Expected empty description for missing field, got '%s'", items[0].Description)
	}
}
