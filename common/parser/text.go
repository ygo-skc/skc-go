package parser

import (
	"strings"
)

func TextContainsSubStr(text, substring string) bool {
	return OccurrencesOfQuotedSubStr(text, substring, true) == 1
}

func OccurrencesOfQuotedSubStr(text, substring string, exitOnFirstOccurrence bool) int {
	textLen := len(text)
	nameLen := len(substring)

	occurrences := 0

	for i := 0; i < textLen; i++ {
		if text[i] == '"' || text[i] == '\'' {
			start := i + 1
			end := start + nameLen

			if end >= textLen {
				break
			}

			if text[end] != text[i] {
				continue
			}

			if text[start:end] == substring {
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
