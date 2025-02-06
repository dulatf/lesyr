package postgres

import (
	"context"
	"errors"

	"github.com/dulatf/lesyr/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedRepository struct {
	db *pgxpool.Pool
}

func NewFeedRepository(db *pgxpool.Pool) *FeedRepository {
	return &FeedRepository{db: db}
}

func (r *FeedRepository) Create(ctx context.Context, feed *model.Feed) error {
	query := `
		INSERT INTO feeds (id, user_id, url, title, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		feed.ID,
		feed.UserID,
		feed.URL,
		feed.Title,
		feed.Description,
	).Scan(&feed.CreatedAt, &feed.UpdatedAt)
}

func (r *FeedRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Feed, error) {
	feed := &model.Feed{}

	query := `
		SELECT id, user_id, url, title, description, last_fetched_at, created_at, updated_at
		FROM feeds
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&feed.ID,
		&feed.UserID,
		&feed.URL,
		&feed.Title,
		&feed.Description,
		&feed.LastFetchedAt,
		&feed.CreatedAt,
		&feed.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return feed, nil
}

func (r *FeedRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Feed, error) {
	query := `
		SELECT id, user_id, url, title, description, last_fetched_at, created_at, updated_at
		FROM feeds
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*model.Feed
	for rows.Next() {
		feed := &model.Feed{}
		err := rows.Scan(
			&feed.ID,
			&feed.UserID,
			&feed.URL,
			&feed.Title,
			&feed.Description,
			&feed.LastFetchedAt,
			&feed.CreatedAt,
			&feed.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, rows.Err()
}

func (r *FeedRepository) Update(ctx context.Context, feed *model.Feed) error {
	query := `
		UPDATE feeds
		SET url = $1, title = $2, description = $3, last_fetched_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		feed.URL,
		feed.Title,
		feed.Description,
		feed.LastFetchedAt,
		feed.ID,
	).Scan(&feed.UpdatedAt)
}

func (r *FeedRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM feeds WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("feed not found")
	}

	return nil
}
