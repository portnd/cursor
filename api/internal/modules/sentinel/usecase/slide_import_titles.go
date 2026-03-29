package usecase

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxSuggestedTaskTitleRunes = 240

// suggestedTaskTitleFromSlideText builds a one-line task title from visible slide body text.
// Intentionally does not use the slide's structural/placeholder title — backlog tasks should reflect on-slide content.
func suggestedTaskTitleFromSlideText(body string, slideIndex int) string {
	fields := strings.Fields(body)
	if len(fields) == 0 {
		return fmt.Sprintf("Slide %d", slideIndex)
	}
	out := strings.Join(fields, " ")
	if utf8.RuneCountInString(out) > maxSuggestedTaskTitleRunes {
		r := []rune(out)
		return string(r[:maxSuggestedTaskTitleRunes]) + "…"
	}
	return out
}
