package repository

import (
	"context"
	"errors"

	"github.com/dulatf/lesyr/internal/model"
	"github.com/google/uuid"
)

var (
	ErrDuplicateArticle = errors.New("article already exists")
)

type UserRepository interface {
	// CreateOrUpdateUser creates a new user or updates existing one from OAuth data
	CreateOrUpdateUser(ctx context.Context, email, provider, providerUserID, name, avatarURL string) (*model.User, error)

	// GetByID retrieves a user by their internal ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// GetByProviderID retrieves a user by their OAuth provider and provider-specific ID
	GetByProviderID(ctx context.Context, provider, providerUserID string) (*model.User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type FeedRepository interface {
	Create(ctx context.Context, feed *model.Feed) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Feed, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Feed, error)
	Update(ctx context.Context, feed *model.Feed) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]*model.Feed, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	GetByFeedID(ctx context.Context, feedID uuid.UUID) ([]*model.Article, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error)
	MarkAsRead(ctx context.Context, userID, articleID uuid.UUID) error
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Article, error)
}
