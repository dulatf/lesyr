package handler

import (
	"context"
	"fmt"

	"github.com/dulatf/lesyr/internal/model"
	"github.com/dulatf/lesyr/internal/repository"
	"github.com/dulatf/lesyr/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FeedHandler struct {
	feedRepo    repository.FeedRepository
	feedFetcher *service.FeedFetcher
}

func NewFeedHandler(feedRepo repository.FeedRepository, feedFetcher *service.FeedFetcher) *FeedHandler {
	return &FeedHandler{
		feedRepo:    feedRepo,
		feedFetcher: feedFetcher,
	}
}

// CreateFeed handles the creation of a new feed
func (h *FeedHandler) CreateFeed(c *fiber.Ctx) error {
	// TODO: Get real user ID from JWT token
	userID, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID")
	}

	feed := new(model.Feed)
	if err := c.BodyParser(feed); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate feed URL and get metadata
	parsedFeed, err := h.feedFetcher.ValidateFeed(feed.URL)
	if err != nil {
		if feedErr, ok := err.(*service.FeedError); ok {
			return fiber.NewError(feedErr.Code, feedErr.Message)
		}
		return fiber.NewError(fiber.StatusBadRequest, "Invalid feed URL")
	}

	// Set feed metadata
	feed.ID = uuid.New()
	feed.UserID = userID
	feed.Title = parsedFeed.Title
	feed.Description = parsedFeed.Description

	// First create the feed in the database
	if err := h.feedRepo.Create(c.Context(), feed); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create feed")
	}
	// Kick off the article fetching process in the background
	go func() {
		if err := h.feedFetcher.RefreshFeed(context.Background(), feed.ID); err != nil {
			fmt.Printf("Error refreshing feed: %v\n", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(feed)
}

// RefreshFeed allows manually triggering a feed refresh
func (h *FeedHandler) RefreshFeed(c *fiber.Ctx) error {
	feedID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid feed ID")
	}

	// Trigger refresh in background
	go func() {
		if err := h.feedFetcher.RefreshFeed(context.Background(), feedID); err != nil {
			fmt.Printf("Error refreshing feed: %v\n", err)
		}
	}()

	return c.SendStatus(fiber.StatusAccepted)
}

// ListFeeds returns all feeds for the authenticated user
func (h *FeedHandler) ListFeeds(c *fiber.Ctx) error {
	// TODO: Get real user ID from JWT token
	userID, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID")
	}

	feeds, err := h.feedRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch feeds")
	}

	return c.JSON(fiber.Map{
		"feeds": feeds,
	})
}

// GetFeed returns a specific feed by ID
func (h *FeedHandler) GetFeed(c *fiber.Ctx) error {
	feedID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid feed ID")
	}

	feed, err := h.feedRepo.GetByID(c.Context(), feedID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch feed")
	}

	if feed == nil {
		return fiber.NewError(fiber.StatusNotFound, "Feed not found")
	}

	// TODO: Verify that the feed belongs to the authenticated user

	return c.JSON(feed)
}

// DeleteFeed removes a feed by ID
func (h *FeedHandler) DeleteFeed(c *fiber.Ctx) error {
	feedID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid feed ID")
	}

	// TODO: Verify that the feed belongs to the authenticated user

	if err := h.feedRepo.Delete(c.Context(), feedID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete feed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
