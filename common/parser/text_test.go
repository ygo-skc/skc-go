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
