# RSS Reader Package

A Go package for parsing RSS feeds asynchronously from multiple URLs.

## Overview

This package provides functionality to parse RSS feeds from multiple URLs concurrently and return structured data. It uses the latest stable version of Go and follows Go best practices.

## Features

- Asynchronous RSS feed parsing from multiple URLs
- Structured RSS item data with all essential fields
- Comprehensive test coverage

## API

### Types

```go
type RssItem struct {
    Title        string
    Source       string
    SourceURL    string
    Link         string
    PublishDate  time.Time
    Description  string
}
```

### Methods

```go
func Parse(ctx context.Context, urls []string) ([]RssItem, error)
```

- **Parameters**: `ctx` - context for cancellation and timeout; `urls` - Array of RSS feed URLs to parse
- **Returns**: Array of `RssItem` structs generated from all provided RSS posts
- **Behavior**: Parses feeds asynchronously for better performance

## Requirements

- Latest stable version of Go (see https://go.dev/dl/)

## Installation

```bash
go get github.com/yourusername/RssReader
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/yourusername/RssReader"
)

func main() {
    ctx := context.Background()
    urls := []string{
        "https://example.com/feed.xml",
        "https://another-blog.com/rss",
    }
    
    items, err := RssReader.Parse(ctx, urls)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, item := range items {
        fmt.Printf("Title: %s\n", item.Title)
        fmt.Printf("Source: %s\n", item.Source)
        fmt.Printf("Published: %s\n", item.PublishDate.Format(time.RFC3339))
        fmt.Printf("Link: %s\n", item.Link)
        fmt.Printf("---\n")
    }
}
```

## Development

### Prerequisites

- Go 1.24+

### Running Tests

```bash
go test ./...
```

### Running Linter

```bash
golangci-lint run
```

### Running Security Checks

```bash
gosec ./...
```

### All Checks

```bash
make all
```

## Project Structure

```
RssReader/
├── README.md
├── go.mod
├── go.sum
├── .github/
│   └── workflows/
│       └── ci.yml
├── .golangci.yml
└── rssreader.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

MIT

## Testing

All changes must be covered with tests. The project includes:

- Unit tests for all exported functions
- Integration tests for RSS parsing
- Benchmark tests for performance validation

Run the full test suite with:

```bash
go test -v -race -cover ./...
``` 