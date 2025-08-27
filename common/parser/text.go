package parser

import "strings"

func TextContainsSubStr(text, substring string) bool {
	if occurrences := OccurrencesOfQuotedSubStr(text, substring, true); occurrences == 1 {
		return true
	}
	return false
}

func OccurrencesOfQuotedSubStr(text, substring string, exitOnFirstOccurrence bool) int {
	runes := []rune(text)
	nameRunes := []rune(substring)
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

			if string(runes[start:end]) == substring {
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

// cleans up a quoted strings
func CleanupToken(t *QuotedToken) {
	*t = strings.TrimSpace(*t)
	*t = strings.ReplaceAll(*t, `".`, "")
	*t = strings.ReplaceAll(*t, `",`, "")
	*t = strings.ReplaceAll(*t, "'.", "")
	*t = strings.ReplaceAll(*t, "',", "")

	*t = strings.Trim(*t, "'")
	*t = strings.Trim(*t, `"`)
}
