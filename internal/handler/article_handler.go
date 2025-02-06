package handler

import (
	"github.com/dulatf/lesyr/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ArticleHandler struct {
	articleRepo repository.ArticleRepository
	feedRepo    repository.FeedRepository
}

func NewArticleHandler(articleRepo repository.ArticleRepository, feedRepo repository.FeedRepository) *ArticleHandler {
	return &ArticleHandler{
		articleRepo: articleRepo,
		feedRepo:    feedRepo,
	}
}

// GetFeedArticles returns all articles for a specific feed
func (h *ArticleHandler) GetFeedArticles(c *fiber.Ctx) error {
	feedID, err := uuid.Parse(c.Params("feedId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid feed ID")
	}

	// TODO: Verify that the feed belongs to the authenticated user
	feed, err := h.feedRepo.GetByID(c.Context(), feedID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch feed")
	}
	if feed == nil {
		return fiber.NewError(fiber.StatusNotFound, "Feed not found")
	}

	articles, err := h.articleRepo.GetByFeedID(c.Context(), feedID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch articles")
	}

	return c.JSON(fiber.Map{
		"articles": articles,
	})
}

// GetArticle returns a specific article by ID
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
	articleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid article ID")
	}

	article, err := h.articleRepo.GetByID(c.Context(), articleID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch article")
	}
	if article == nil {
		return fiber.NewError(fiber.StatusNotFound, "Article not found")
	}

	// TODO: Verify that the article's feed belongs to the authenticated user

	return c.JSON(article)
}

// MarkArticleRead marks an article as read for the current user
func (h *ArticleHandler) MarkArticleRead(c *fiber.Ctx) error {
	// TODO: Get real user ID from JWT token
	userID, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID")
	}

	articleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid article ID")
	}

	if err := h.articleRepo.MarkAsRead(c.Context(), userID, articleID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to mark article as read")
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetUnreadArticles returns all unread articles for the current user
func (h *ArticleHandler) GetUnreadArticles(c *fiber.Ctx) error {
	// TODO: Get real user ID from JWT token
	userID, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID")
	}

	articles, err := h.articleRepo.GetUnreadByUserID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch unread articles")
	}

	return c.JSON(fiber.Map{
		"articles": articles,
	})
}
