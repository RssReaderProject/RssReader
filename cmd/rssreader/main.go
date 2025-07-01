// Main entry point for the RSS reader CLI tool.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	rssreader "github.com/RssReaderProject/RssReader"
)

func main() {
	// Define command line flags
	var (
		urls    = flag.String("urls", "", "Comma-separated list of RSS feed URLs")
		format  = flag.String("format", "json", "Output format: json, text")
		timeout = flag.Duration("timeout", 30*time.Second, "Timeout for fetching feeds")
		help    = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Println("RSS Reader - A simple RSS feed parser")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/rssreader/main.go -urls=\"https://example.com/feed.xml,https://another.com/rss\"")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return
	}

	// Check if URLs are provided
	if *urls == "" {
		log.Print("Error: -urls flag is required. Use -help for usage information.")
		os.Exit(1)
	}

	// Parse URLs from comma-separated string
	urlList := []string{}
	if *urls != "" {
		// Simple comma splitting - in a real app you might want more sophisticated parsing
		for _, url := range strings.Split(*urls, ",") {
			url = strings.TrimSpace(url)
			if url != "" {
				urlList = append(urlList, url)
			}
		}
	}

	if len(urlList) == 0 {
		log.Print("Error: No valid URLs provided")
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)

	// Parse RSS feeds
	items, err := rssreader.Parse(ctx, urlList)
	if err != nil {
		log.Printf("Error parsing RSS feeds: %v", err)
		os.Exit(1)
	}

	defer cancel()

	// Output results based on format
	switch *format {
	case "json":
		err := outputJSON(items)
		if err != nil {
			log.Printf("Error outputting JSON: %v", err)
		}
	case "text":
		outputText(items)
	default:
		log.Printf("Unknown format: %s. Supported formats: json, text", *format)
	}
}

func outputJSON(items []rssreader.RssItem) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}
	return nil
}

func outputText(items []rssreader.RssItem) {
	for i, item := range items {
		fmt.Printf("=== Item %d ===\n", i+1)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Source: %s\n", item.Source)
		fmt.Printf("Source URL: %s\n", item.SourceURL)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Publish Date: %s\n", item.PublishDate.Format(time.RFC3339))
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Println()
	}
}
