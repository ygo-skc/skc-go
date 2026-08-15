package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardTextContainsName(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		testName              string
		text                  string
		substring             string
		textShouldContainName bool
	}{
		{
			testName:              "Substring found in text",
			text:                  `"Neos" is a powerful monster from and his name is "Neos"`,
			substring:             "Neos",
			textShouldContainName: true,
		},
		{
			testName:              "Substring NOT found in text",
			text:                  `"Neos2" is a powerful monster from and his name is "Neos3"`,
			substring:             "Neos",
			textShouldContainName: false,
		},
		{
			testName:              "Multi-byte substring found in text",
			text:                  `"日本語" is quoted once`,
			substring:             "日本語",
			textShouldContainName: true,
		},
		{
			testName:              "Accented substring NOT matched by unaccented text",
			text:                  `"Cafe" is missing its accent`,
			substring:             "Café",
			textShouldContainName: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			assert.Equal(tt.textShouldContainName, TextContainsSubStr(tt.text, tt.substring))
		})
	}
}

func TestOccurrenceOfNameInText(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		testName            string
		text                string
		substring           string
		expectedOccurrences int
	}{
		{
			testName:            "Two complete quotes",
			text:                `"Neos" is a powerful monster from and his name is "Neos"`,
			substring:           "Neos",
			expectedOccurrences: 2,
		},
		{
			testName:            "Unfinished quote at end of text",
			text:                `"Neos" is a powerful monster from and his name is "Neos`,
			substring:           "Neos",
			expectedOccurrences: 1,
		},
		{
			testName:            "Quoted strings in text does not match name",
			text:                `"Neos2" is a powerful monster from and his name is "Neos3"`,
			substring:           "Neos",
			expectedOccurrences: 0,
		},
		{
			testName:            "Multi-byte substring quoted twice",
			text:                `"龍" and "龍" appear twice`,
			substring:           "龍",
			expectedOccurrences: 2,
		},
		{
			testName:            "Accented substring matched exactly",
			text:                `"Café" is accented`,
			substring:           "Café",
			expectedOccurrences: 1,
		},
		{
			testName:            "Accented substring NOT matched by unaccented text",
			text:                `"Cafe" is missing its accent`,
			substring:           "Café",
			expectedOccurrences: 0,
		},
		{
			testName:            "Substring surrounded by multi-byte runes",
			text:                `龍"Neos"龍 padded by multi-byte runes`,
			substring:           "Neos",
			expectedOccurrences: 1,
		},
		{
			testName:            "Substring with inner apostrophe",
			text:                `"Magicians' Souls" has an inner apostrophe`,
			substring:           "Magicians' Souls",
			expectedOccurrences: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			assert.Equal(tt.expectedOccurrences, OccurrencesOfQuotedSubStr(tt.text, tt.substring, false))
		})
	}
}

func TestCleanupToken(t *testing.T) {
	// setup
	assert := assert.New(t)

	testData := []string{`HERO".`, `HERO",`, `"HERO`, ` HERO `, "HERO'.", "HERO',", "'HERO"}
	for _, data := range testData {
		CleanupToken(&data)
		assert.Equal("HERO", data, "Token not cleaned up correctly")
	}

	// edge case 1 - inner single quote should not be removed
	edge1 := "Magicians' Souls"
	CleanupToken(&edge1)
	assert.Equal("Magicians' Souls", edge1, "Edge case 1 (inner single quote should not be removed) - failed")
}
