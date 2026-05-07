package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const readmePath = "../README.md"

// updateSection replaces everything between <!-- START_SECTION:<sectionName> -->
// and <!-- END_SECTION:<sectionName> --> with the provided newData.
// It handles missing tags gracefully by returning the original content.
func updateSection(content string, sectionName string, newData string) string {
	startTag := fmt.Sprintf("<!-- START_SECTION:%s -->", sectionName)
	endTag := fmt.Sprintf("<!-- END_SECTION:%s -->", sectionName)

	startIdx := strings.Index(content, startTag)
	if startIdx == -1 {
		log.Printf("Error: start tag %q not found in file content", startTag)
		return content
	}

	endIdx := strings.Index(content, endTag)
	if endIdx == -1 {
		log.Printf("Error: end tag %q not found in file content", endTag)
		return content
	}

	if startIdx >= endIdx {
		log.Printf("Error: start tag %q found after or at end tag %q", startTag, endTag)
		return content
	}

	// Extract the content up to the end of the start tag
	prefix := content[:startIdx+len(startTag)]
	// Extract the content starting from the end tag
	suffix := content[endIdx:]

	// Insert the new data with proper newline spacing
	return prefix + "\n" + newData + "\n" + suffix
}

func main() {
	// Load local environment variables from .env file if available
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: Error loading .env file, relying on system/environment variables.")
	}

	// Read the target markdown file
	contentBytes, err := os.ReadFile(readmePath)
	if err != nil {
		log.Fatalf("Failed to read README file at %s: %v", readmePath, err)
	}

	content := string(contentBytes)

	// Update the "test" section
	currentTimeStr := time.Now().Format("Jan 2, 2006, 3:04 PM")
	testData := "**Hello from Go Automator! Update time:** " + currentTimeStr
	updatedContent := updateSection(content, "test", testData)

	// Fetch GitHub stats for user dynamically from GITHUB_USERNAME env var
	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		username = "tin-auppati"
	}
	stats, err := FetchGitHubStats(username)
	if err != nil {
		log.Printf("Warning: failed to fetch GitHub stats: %v", err)
	} else {
		// Format the statistics into markdown including private repo count
		statsData := fmt.Sprintf("- 📦 **Public Repos:** %d\n- 🔒 **Private Repos:** %d\n- 👥 **Followers:** %d", stats.PublicRepos, stats.TotalPrivateRepos, stats.Followers)
		// Update the "github_stats" section
		updatedContent = updateSection(updatedContent, "github_stats", statsData)
	}

	// Fetch WakaTime stats
	wakaStats, err := FetchWakaTimeStats()
	if err != nil {
		log.Printf("Warning: failed to fetch WakaTime stats: %v", err)
	} else {
		// Update the "wakatime" section
		updatedContent = updateSection(updatedContent, "wakatime", wakaStats)
	}

	// If no tags were found or nothing changed, exit early
	if updatedContent == content {
		log.Println("Warning: No changes were made to the file content. Verify tag presence/names.")
		return
	}

	// Write the updated content back with 0644 permissions
	err = os.WriteFile(readmePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write updated content to %s: %v", readmePath, err)
	}

	fmt.Printf("Successfully updated %s sections 'test', 'github_stats', and 'wakatime' at %s\n", readmePath, currentTimeStr)
}
