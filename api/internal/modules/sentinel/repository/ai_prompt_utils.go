package repository

import (
	"regexp"
	"strings"
)

var (
	// Match base64 data URIs (images, files) — can be hundreds of KB each
	reBase64DataURI = regexp.MustCompile(`src="data:[^"]*"`)
	// Match HTML tags
	reHTMLTag = regexp.MustCompile(`<[^>]+>`)
	// Match multiple whitespace/newlines
	reMultiSpace = regexp.MustCompile(`\s{2,}`)
)

// stripDescriptionForAI cleans a task description for AI prompt consumption:
// 1. Removes base64 data URIs (embedded images/files — can be MBs of text)
// 2. Strips HTML tags → plain text
// 3. Collapses whitespace
// 4. Truncates to maxLen if still too long
//
// This is essential because task descriptions can contain rich content with
// embedded base64 images that are 1MB+ of text — useless for AI estimation
// but would exceed any model's context window.
func stripDescriptionForAI(desc string, maxLen int) string {
	if desc == "" {
		return ""
	}

	// Step 1: Remove base64 data URIs (biggest win — single image can be 500KB+)
	desc = reBase64DataURI.ReplaceAllString(desc, `src="[image]"`)

	// Step 2: Strip HTML tags → plain text
	desc = reHTMLTag.ReplaceAllString(desc, " ")

	// Step 3: Decode common HTML entities
	desc = strings.ReplaceAll(desc, "&nbsp;", " ")
	desc = strings.ReplaceAll(desc, "&amp;", "&")
	desc = strings.ReplaceAll(desc, "&lt;", "<")
	desc = strings.ReplaceAll(desc, "&gt;", ">")
	desc = strings.ReplaceAll(desc, "&quot;", `"`)
	desc = strings.ReplaceAll(desc, "&#39;", "'")

	// Step 4: Collapse whitespace
	desc = reMultiSpace.ReplaceAllString(desc, " ")
	desc = strings.TrimSpace(desc)

	// Step 5: Truncate if still too long
	if maxLen > 0 && len(desc) > maxLen {
		desc = desc[:maxLen] + "\n[...truncated]"
	}

	return desc
}

// clamp restricts v to [min, max]. Used for AI factor scores (1-10).
func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
