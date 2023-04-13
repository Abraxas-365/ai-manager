package tokenizer

import "strings"

func countTokens(text string) int {
	// Assuming a rough 4 characters per token for English
	tokenLength := 4
	// Remove extra spaces
	text = strings.TrimSpace(text)
	// Count the number of characters (ignoring spaces)
	charCount := 0
	for _, char := range text {
		if char != ' ' {
			charCount++
		}
	}
	// Calculate the number of tokens
	tokenCount := charCount / tokenLength
	return tokenCount
}
