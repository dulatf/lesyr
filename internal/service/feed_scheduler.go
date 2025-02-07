// internal/service/scheduler.go
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dulatf/lesyr/internal/repository"
	"github.com/google/uuid"
)

type FeedScheduler struct {
	feedRepo    repository.FeedRepository
	feedFetcher *FeedFetcher
	workers     int
	interval    time.Duration
	jobs        chan uuid.UUID
	quit        chan struct{}
	wg          sync.WaitGroup
}

func NewFeedScheduler(
	feedRepo repository.FeedRepository,
	feedFetcher *FeedFetcher,
	workers int,
	interval time.Duration,
) *FeedScheduler {
	if workers <= 0 {
		workers = 5 // default number of workers
	}
	if interval <= 0 {
		interval = 15 * time.Minute // default refresh interval
	}

	return &FeedScheduler{
		feedRepo:    feedRepo,
		feedFetcher: feedFetcher,
		workers:     workers,
		interval:    interval,
		jobs:        make(chan uuid.UUID, 100), // buffer size of 100 jobs
		quit:        make(chan struct{}),
	}
}

// Start begins the scheduler and worker goroutines
func (s *FeedScheduler) Start(ctx context.Context) error {
	// Start worker goroutines
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}

	// Start scheduler goroutine
	s.wg.Add(1)
	go s.scheduler(ctx)

	return nil
}

// Stop gracefully shuts down the scheduler and workers
func (s *FeedScheduler) Stop() {
	close(s.quit)
	s.wg.Wait()
}

// scheduler periodically checks for feeds that need refreshing
func (s *FeedScheduler) scheduler(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.quit:
			return
		case <-ticker.C:
			if err := s.scheduleFeeds(ctx); err != nil {
				fmt.Printf("Error scheduling feeds: %v\n", err)
			}
		}
	}
}

// scheduleFeeds finds feeds that need refreshing and adds them to the job queue
func (s *FeedScheduler) scheduleFeeds(ctx context.Context) error {
	feeds, err := s.feedRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	now := time.Now()
	for _, feed := range feeds {
		// Skip if feed was fetched recently
		if feed.LastFetchedAt != nil && now.Sub(*feed.LastFetchedAt) < s.interval {
			continue
		}

		// Try to add feed to job queue, skip if queue is full
		select {
		case s.jobs <- feed.ID:
			fmt.Printf("Job scheduled: %s\n", feed.ID)
			// Successfully queued
		default:
			fmt.Printf("Job queue full, skipping feed: %s\n", feed.ID)
		}
	}

	return nil
}

// worker processes feed refresh jobs
func (s *FeedScheduler) worker(ctx context.Context, id int) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.quit:
			return
		case feedID := <-s.jobs:
			if err := s.refreshFeed(ctx, feedID); err != nil {
				fmt.Printf("Worker %d - Error refreshing feed %s: %v\n", id, feedID, err)
			}
		}
	}
}

// refreshFeed handles the feed refresh operation with retry logic
func (s *FeedScheduler) refreshFeed(ctx context.Context, feedID uuid.UUID) error {
	// Get the feed
	feed, err := s.feedRepo.GetByID(ctx, feedID)
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}
	if feed == nil {
		return fmt.Errorf("feed not found: %v", feedID)
	}

	// Attempt to refresh the feed
	err = s.feedFetcher.RefreshFeed(ctx, feedID)
	if err != nil {
		// Update feed with error status
		feed.LastError = err.Error()
		if err := s.feedRepo.Update(ctx, feed); err != nil {
			fmt.Printf("Error updating feed status: %v\n", err)
		}
		return fmt.Errorf("error refreshing feed: %v", err)
	}

	// Clear any previous error status
	if feed.LastError != "" {
		feed.LastError = ""
		if err := s.feedRepo.Update(ctx, feed); err != nil {
			fmt.Printf("Error updating feed status: %v\n", err)
		}
	}

	return nil
}
