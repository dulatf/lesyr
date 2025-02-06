package postgres

import (
	"context"
	"errors"

	"github.com/dulatf/lesyr/internal/model"
	"github.com/dulatf/lesyr/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArticleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(ctx context.Context, article *model.Article) error {
	query := `
		INSERT INTO articles (id, feed_id, title, url, content, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (feed_id, url) DO NOTHING
		RETURNING created_at`

	err := r.db.QueryRow(ctx, query,
		article.ID,
		article.FeedID,
		article.Title,
		article.URL,
		article.Content,
		article.PublishedAt,
	).Scan(&article.CreatedAt)

	if err != nil {
		// Check if it's a unique violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrDuplicateArticle
		}
		return err
	}

	return nil
}

func (r *ArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	article := &model.Article{}

	query := `
		SELECT id, feed_id, title, url, content, published_at, created_at
		FROM articles
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&article.ID,
		&article.FeedID,
		&article.Title,
		&article.URL,
		&article.Content,
		&article.PublishedAt,
		&article.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return article, nil
}

func (r *ArticleRepository) GetByFeedID(ctx context.Context, feedID uuid.UUID) ([]*model.Article, error) {
	query := `
		SELECT id, feed_id, title, url, content, published_at, created_at
		FROM articles
		WHERE feed_id = $1
		ORDER BY published_at DESC`

	rows, err := r.db.Query(ctx, query, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		article := &model.Article{}
		err := rows.Scan(
			&article.ID,
			&article.FeedID,
			&article.Title,
			&article.URL,
			&article.Content,
			&article.PublishedAt,
			&article.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, rows.Err()
}

func (r *ArticleRepository) MarkAsRead(ctx context.Context, userID, articleID uuid.UUID) error {
	query := `
		INSERT INTO read_status (user_id, article_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, article_id) DO NOTHING`

	_, err := r.db.Exec(ctx, query, userID, articleID)
	return err
}

func (r *ArticleRepository) GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Article, error) {
	query := `
		SELECT a.id, a.feed_id, a.title, a.url, a.content, a.published_at, a.created_at
		FROM articles a
		JOIN feeds f ON a.feed_id = f.id
		LEFT JOIN read_status rs ON rs.article_id = a.id AND rs.user_id = $1
		WHERE f.user_id = $1 AND rs.user_id IS NULL
		ORDER BY a.published_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		article := &model.Article{}
		err := rows.Scan(
			&article.ID,
			&article.FeedID,
			&article.Title,
			&article.URL,
			&article.Content,
			&article.PublishedAt,
			&article.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, rows.Err()
}
