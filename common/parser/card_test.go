package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardTextContainsName(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		testName              string
		cardText              string
		cardName              string
		textShouldContainName bool
	}{
		{
			testName:              "Card name found in card text",
			cardText:              `"Neos" is a powerful monster from and his name is "Neos"`,
			cardName:              "Neos",
			textShouldContainName: true,
		},
		{
			testName:              "Card name NOT found in card text",
			cardText:              `"Neos2" is a powerful monster from and his name is "Neos3"`,
			cardName:              "Neos",
			textShouldContainName: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			assert.Equal(tt.textShouldContainName, CardTextContainsName(tt.cardText, tt.cardName))
		})
	}
}

func TestOccurrenceOfNameInText(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		testName            string
		cardText            string
		cardName            string
		expectedOccurrences int
	}{
		{
			testName:            "Two complete quotes",
			cardText:            `"Neos" is a powerful monster from and his name is "Neos"`,
			cardName:            "Neos",
			expectedOccurrences: 2,
		},
		{
			testName:            "Unfinished quote at end of card text",
			cardText:            `"Neos" is a powerful monster from and his name is "Neos`,
			cardName:            "Neos",
			expectedOccurrences: 1,
		},
		{
			testName:            "Quoted strings in text does not match card name",
			cardText:            `"Neos2" is a powerful monster from and his name is "Neos3"`,
			cardName:            "Neos",
			expectedOccurrences: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			assert.Equal(tt.expectedOccurrences, OccurrenceOfNameInText(tt.cardText, tt.cardName, false))
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
