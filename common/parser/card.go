package parser

import "strings"

func CardTextContainsName(cardText, cardName string) bool {
	if occurrences := OccurrenceOfNameInText(cardText, cardName, true); occurrences == 1 {
		return true
	}
	return false
}

func OccurrenceOfNameInText(cardText, cardName string, exitOnFirstOccurrence bool) int {
	runes := []rune(cardText)
	nameRunes := []rune(cardName)
	textLen := len(runes)
	nameLen := len(nameRunes)

	occurrences := 0

	for i := 0; i < textLen; i++ {
		if runes[i] == '"' || runes[i] == '\'' {
			start := i + 1
			end := start + nameLen

			if end >= textLen {
				break
			}

			if runes[end] != runes[i] {
				continue
			}

			if string(runes[start:end]) == cardName {
				if exitOnFirstOccurrence {
					return 1
				}

				i = end
				occurrences++
			}
		}
	}
	return occurrences
}

type QuotedToken = string

// cleans up a quoted string found in card text so its easier to parse
func CleanupToken(t *QuotedToken) {
	*t = strings.TrimSpace(*t)
	*t = strings.ReplaceAll(*t, `".`, "")
	*t = strings.ReplaceAll(*t, `",`, "")
	*t = strings.ReplaceAll(*t, "'.", "")
	*t = strings.ReplaceAll(*t, "',", "")

	*t = strings.Trim(*t, "'")
	*t = strings.Trim(*t, `"`)
}
