package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GitHubUser represents the structure of the GitHub user API response we care about.
type GitHubUser struct {
	PublicRepos int `json:"public_repos"`
	Followers   int `json:"followers"`
}

// FetchGitHubStats fetches public repository and follower statistics for the given username.
func FetchGitHubStats(username string) (GitHubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("failed to create request: %w", err)
	}

	// GitHub API requires a User-Agent header to avoid being blocked
	req.Header.Set("User-Agent", "Go-Automator")

	resp, err := client.Do(req)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitHubUser{}, fmt.Errorf("API returned status %d (%s)", resp.StatusCode, resp.Status)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return GitHubUser{}, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return user, nil
}
