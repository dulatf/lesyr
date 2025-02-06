package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dulatf/lesyr/internal/model"
	"github.com/dulatf/lesyr/internal/repository"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

type FeedError struct {
	Message string
	Code    int
}

func (e *FeedError) Error() string {
	return e.Message
}

type FeedFetcher struct {
	feedRepo    repository.FeedRepository
	articleRepo repository.ArticleRepository
	parser      *gofeed.Parser
}

func NewFeedFetcher(feedRepo repository.FeedRepository, articleRepo repository.ArticleRepository) *FeedFetcher {
	return &FeedFetcher{
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
		parser:      gofeed.NewParser(),
	}
}

// ValidateFeed validates a feed URL and returns its metadata
func (f *FeedFetcher) ValidateFeed(url string) (*gofeed.Feed, error) {
	parsedFeed, err := f.parser.ParseURL(url)
	if err != nil {
		return nil, &FeedError{
			Message: fmt.Sprintf("invalid feed URL: %v", err),
			Code:    400,
		}
	}
	return parsedFeed, nil
}

// RefreshFeed fetches the latest articles for a feed and stores them
// This is the unified method that handles both initial feed population
// and subsequent updates
func (f *FeedFetcher) RefreshFeed(ctx context.Context, feedID uuid.UUID) error {
	// Get feed from database
	feed, err := f.feedRepo.GetByID(ctx, feedID)
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}
	if feed == nil {
		return fmt.Errorf("feed not found: %v", feedID)
	}

	// Fetch and parse feed content
	parsedFeed, err := f.parser.ParseURL(feed.URL)
	if err != nil {
		return fmt.Errorf("error parsing feed: %v", err)
	}

	// Update feed metadata if changed
	metadataChanged := false
	if feed.Title != parsedFeed.Title || feed.Description != parsedFeed.Description {
		feed.Title = parsedFeed.Title
		feed.Description = parsedFeed.Description
		metadataChanged = true
	}

	// Update last fetched time
	now := time.Now()
	feed.LastFetchedAt = &now

	if metadataChanged {
		if err := f.feedRepo.Update(ctx, feed); err != nil {
			return fmt.Errorf("error updating feed metadata: %v", err)
		}
	}

	// Keep track of processed articles
	newArticles := 0
	duplicates := 0

	// Process articles in reverse chronological order
	for _, item := range parsedFeed.Items {
		// Skip articles without a published date or URL
		if item.PublishedParsed == nil || item.Link == "" {
			continue
		}

		article := &model.Article{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       item.Title,
			URL:         item.Link,
			Content:     item.Description,
			PublishedAt: *item.PublishedParsed,
		}

		err := f.articleRepo.Create(ctx, article)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateArticle) {
				duplicates++
				continue
			}
			// Log error but continue processing other articles
			fmt.Printf("error storing article: %v\n", err)
			continue
		}
		newArticles++
	}

	fmt.Printf("Feed %s processed: %d new articles, %d duplicates\n",
		feed.Title, newArticles, duplicates)

	return nil
}
