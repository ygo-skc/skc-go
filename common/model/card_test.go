package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestIsCardNameFoundInTokens(t *testing.T) {
	// setup
	assert := assert.New(t)

	tokens := []QuotedToken{"Elemental HERO Air Neos", "HERO"}

	// exact match
	scenario := YGOCardREST{Name: "Elemental HERO Air Neos"}
	assert.Equal(true, IsCardNameInTokens(scenario, tokens), fmt.Sprintf("%s expected to be in tokens %v", scenario.Name, tokens))

	// this method ignores case - so should be in token array
	scenario = YGOCardREST{Name: "Elemental Hero Air Neos"}
	assert.Equal(true, IsCardNameInTokens(scenario, tokens), "Expected to find return true as method ignores case")

	// not found in token array
	scenario = YGOCardREST{Name: "HEROs"}
	assert.Equal(false, IsCardNameInTokens(scenario, tokens), fmt.Sprintf("%s not expected to be in tokens %v", scenario.Name, tokens))
}
