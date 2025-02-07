package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dulatf/lesyr/internal/config"
	"github.com/dulatf/lesyr/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userRepo    repository.UserRepository
	oauthConfig *config.OAuthConfig
	jwtSecret   string
}

func NewAuthHandler(userRepo repository.UserRepository, oauthConfig *config.OAuthConfig, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
		jwtSecret:   jwtSecret,
	}
}

// InitiateGitHubAuth redirects to GitHub OAuth page
func (h *AuthHandler) InitiateGitHubAuth(c *fiber.Ctx) error {
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = h.oauthConfig.AuthCallbackURL
	}

	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		h.oauthConfig.GitHubClientID,
		redirectURI,
	)
	return c.Redirect(url)
}

// HandleGitHubCallback processes the OAuth callback from GitHub
func (h *AuthHandler) HandleGitHubCallback(c *fiber.Ctx) error {
	// Get the code from the callback
	code := c.Query("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "No code provided")
	}

	// Exchange code for access token
	accessToken, err := h.exchangeCodeForToken(code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to exchange code for token")
	}

	// Get user info from GitHub
	githubUser, err := h.oauthConfig.GetGitHubUser(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user info")
	}

	// Create or update user in our database
	user, err := h.userRepo.CreateOrUpdateUser(
		c.Context(),
		githubUser.Email,
		"github",
		fmt.Sprintf("%d", githubUser.ID),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create/update user")
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	frontendURL := "http://localhost:3000/auth/callback"
	redirectURL := fmt.Sprintf("%s?token=%s", frontendURL, token)
	return c.Redirect(redirectURL)
}

// AuthMiddleware validates JWT tokens
func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No authorization header")
	}

	// Check if development bypass is enabled
	if strings.HasPrefix(authHeader, "Development ") && c.App().Config().Prefork {
		// Extract user ID from development header
		userID := strings.TrimPrefix(authHeader, "Development ")
		c.Locals("userID", userID)
		return c.Next()
	}

	// Parse JWT token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Store user ID in context
		c.Locals("userID", claims["sub"])
		return c.Next()
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
}

// Helper function to exchange code for access token
func (h *AuthHandler) exchangeCodeForToken(code string) (string, error) {
	// Create request
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("client_id", h.oauthConfig.GitHubClientID)
	q.Add("client_secret", h.oauthConfig.GitHubClientSecret)
	q.Add("code", code)
	req.URL.RawQuery = q.Encode()

	// Set headers to get JSON response
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to exchange code: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Check for OAuth error
	if result.Error != "" {
		return "", fmt.Errorf("oauth error: %s", result.Error)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	return result.AccessToken, nil
}

// Helper function to generate JWT token
func (h *AuthHandler) generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(h.jwtSecret))
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

func (h *AuthHandler) GetAuthenticatedUser(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return err
	}
	user, err := h.userRepo.GetByID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}
	return c.JSON(UserResponse{
		ID:        user.ID.String(),
		Username:  user.Email,
		Email:     user.Email,
		AvatarURL: "",
	})
}
