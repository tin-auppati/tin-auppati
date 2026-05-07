package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// GitHubUser represents the structure of the GitHub user API response we care about.
type GitHubUser struct {
	PublicRepos       int `json:"public_repos"`
	TotalPrivateRepos int `json:"total_private_repos"`
	Followers         int `json:"followers"`
}

// FetchGitHubStats fetches repository and follower statistics.
// If GITHUB_TOKEN is present, it uses the authenticated "/user" endpoint to fetch private stats.
// Otherwise, it falls back to the public "/users/<username>" endpoint.
func FetchGitHubStats(username string) (GitHubUser, error) {
	var url string
	token := os.Getenv("GITHUB_TOKEN")

	if token != "" {
		url = "https://api.github.com/user"
	} else {
		log.Println("Note: GITHUB_TOKEN is missing. Falling back to public endpoint (private repos will show as 0).")
		url = fmt.Sprintf("https://api.github.com/users/%s", username)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("failed to create request: %w", err)
	}

	// GitHub API requires a User-Agent header to avoid being blocked
	req.Header.Set("User-Agent", "Go-Automator")

	// Apply GITHUB_TOKEN for authorization if provided
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

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
