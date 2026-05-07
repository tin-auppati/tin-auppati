package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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
	// Read the target markdown file
	contentBytes, err := os.ReadFile(readmePath)
	if err != nil {
		log.Fatalf("Failed to read README file at %s: %v", readmePath, err)
	}

	content := string(contentBytes)

	// Format current time nicely (e.g., "May 7, 2026, 10:52 AM")
	currentTimeStr := time.Now().Format("Jan 2, 2006, 3:04 PM")
	newData := "**Hello from Go Automator! Update time:** " + currentTimeStr

	// Update the "test" section
	updatedContent := updateSection(content, "test", newData)

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

	fmt.Printf("Successfully updated %s 'test' section at %s\n", readmePath, currentTimeStr)
}
