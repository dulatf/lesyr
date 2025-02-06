package config

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type OAuthConfig struct {
	GitHubClientID     string
	GitHubClientSecret string
	AuthCallbackURL    string
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func LoadOAuthConfig() (*OAuthConfig, error) {
	viper.AutomaticEnv()

	// Set defaults for development
	viper.SetDefault("GITHUB_CLIENT_ID", "")
	viper.SetDefault("GITHUB_CLIENT_SECRET", "")
	viper.SetDefault("AUTH_CALLBACK_URL", "http://localhost:8080/api/v1/auth/callback")

	return &OAuthConfig{
		GitHubClientID:     viper.GetString("GITHUB_CLIENT_ID"),
		GitHubClientSecret: viper.GetString("GITHUB_CLIENT_SECRET"),
		AuthCallbackURL:    viper.GetString("AUTH_CALLBACK_URL"),
	}, nil
}

func (c *OAuthConfig) GetGitHubUser(accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &user, nil
}
