package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// WakaLanguage represents a programming language returned by WakaTime
type WakaLanguage struct {
	Name    string  `json:"name"`
	Text    string  `json:"text"`
	Percent float64 `json:"percent"`
}

// WakaStatsData is the data container for language stats
type WakaStatsData struct {
	Languages []WakaLanguage `json:"languages"`
}

// WakaStatsResponse is the top-level WakaTime stats JSON structure
type WakaStatsResponse struct {
	Data WakaStatsData `json:"data"`
}

// getLanguageEmoji maps popular programming languages to expressive emojis
func getLanguageEmoji(name string) string {
	switch name {
	case "Go", "Golang":
		return "🐹"
	case "TypeScript":
		return "🟦"
	case "JavaScript":
		return "🟨"
	case "Python":
		return "🐍"
	case "Rust":
		return "🦀"
	case "HTML":
		return "🌐"
	case "CSS":
		return "🎨"
	case "Markdown":
		return "📝"
	case "C++":
		return "🦕"
	case "C#":
		return "💜"
	case "Java":
		return "☕"
	case "Ruby":
		return "💎"
	case "PHP":
		return "🐘"
	case "Docker", "Dockerfile":
		return "🐳"
	case "SQL":
		return "🗄️"
	case "YAML":
		return "⚙️"
	case "JSON":
		return "📄"
	case "Shell", "Bash", "ShellScript":
		return "🐚"
	default:
		return "💻"
	}
}

// FetchWakaTimeStats connects to WakaTime API to retrieve coding language statistics.
func FetchWakaTimeStats() (string, error) {
	apiKey := os.Getenv("WAKATIME_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("WAKATIME_API_KEY is empty or missing")
	}

	url := "https://wakatime.com/api/v1/users/current/stats/last_7_days"

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// WakaTime requires API key base64-encoded as HTTP Basic authentication
	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedKey))
	req.Header.Set("User-Agent", "Go-Automator")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("WakaTime API returned status %d (%s)", resp.StatusCode, resp.Status)
	}

	var statsResp WakaStatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statsResp); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	languages := statsResp.Data.Languages
	if len(languages) == 0 {
		return "_No coding stats recorded for the last 7 days._", nil
	}

	// Format top 5 coding languages
	limit := 5
	if len(languages) < limit {
		limit = len(languages)
	}

	var markdown string
	for i := 0; i < limit; i++ {
		lang := languages[i]
		emoji := getLanguageEmoji(lang.Name)
		markdown += fmt.Sprintf("- %s **%s**: %s (%.1f%%)\n", emoji, lang.Name, lang.Text, lang.Percent)
	}

	// Trim trailing newline if present
	if len(markdown) > 0 && markdown[len(markdown)-1] == '\n' {
		markdown = markdown[:len(markdown)-1]
	}

	return markdown, nil
}
