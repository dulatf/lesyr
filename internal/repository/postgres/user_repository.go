package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dulatf/lesyr/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Updated signature: now accepts name and avatarURL.
func (r *UserRepository) CreateOrUpdateUser(ctx context.Context, email, provider, providerUserID, name, avatarURL string) (*model.User, error) {
	fmt.Fprintf(os.Stderr, "Checking existing user with provider %s and providerUserID %s\n", provider, providerUserID)
	existingUser, err := r.GetByProviderID(ctx, provider, providerUserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error checking existing user: %v", err)
	}
	fmt.Fprintf(os.Stderr, "Existing user: %v\n", existingUser)
	if existingUser != nil {
		// Update existing user if any fields have changed
		if existingUser.Email != email || existingUser.Name != name || existingUser.AvatarURL != avatarURL {
			query := `
				UPDATE users 
				SET email = $1, name = $2, avatar_url = $3, updated_at = CURRENT_TIMESTAMP
				WHERE id = $4
				RETURNING id, email, provider, provider_user_id, name, avatar_url, created_at, updated_at`

			err := r.db.QueryRow(ctx, query, email, name, avatarURL, existingUser.ID).Scan(
				&existingUser.ID,
				&existingUser.Email,
				&existingUser.Provider,
				&existingUser.ProviderUserID,
				&existingUser.Name,
				&existingUser.AvatarURL,
				&existingUser.CreatedAt,
				&existingUser.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("error updating user: %v", err)
			}
		}
		return existingUser, nil
	}

	// Create new user
	user := &model.User{
		ID:             uuid.New(),
		Email:          email,
		Provider:       provider,
		ProviderUserID: providerUserID,
		Name:           name,
		AvatarURL:      avatarURL,
	}

	query := `
		INSERT INTO users (id, email, provider, provider_user_id, name, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	err = r.db.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.Provider,
		user.ProviderUserID,
		user.Name,
		user.AvatarURL,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) GetByProviderID(ctx context.Context, provider, providerUserID string) (*model.User, error) {
	user := &model.User{}
	fmt.Fprintf(os.Stderr, "Getting user with provider %s and providerUserID %s\n", provider, providerUserID)
	query := `
		SELECT id, email, provider, provider_user_id,
		       COALESCE(name, ''), COALESCE(avatar_url, ''),
		       created_at, updated_at
		FROM users
		WHERE provider = $1 AND provider_user_id = $2`

	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&user.ID,
		&user.Email,
		&user.Provider,
		&user.ProviderUserID,
		&user.Name,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user: %v\n", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, email, provider, provider_user_id,
		       COALESCE(name, ''), COALESCE(avatar_url, ''),
		       created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Provider,
		&user.ProviderUserID,
		&user.Name,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, email, provider, provider_user_id,
		       COALESCE(name, ''), COALESCE(avatar_url, ''),
		       created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Provider,
		&user.ProviderUserID,
		&user.Name,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
