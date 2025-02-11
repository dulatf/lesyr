package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Provider       string    `json:"provider"`         // e.g., 'github', 'google'
	ProviderUserID string    `json:"provider_user_id"` // ID from the OAuth provider
	Name           string    `json:"name"`
	AvatarURL      string    `json:"avatar_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
	LastError     string     `json:"last_error,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Article struct {
	ID          uuid.UUID `json:"id"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
}
