package db

import (
	"fmt"
	"strings"
)

const (
	cardAttributes = "card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense"

	dbVersionQuery    = "SELECT VERSION()"
	cardColorIDsQuery = "SELECT color_id, card_color from card_colors ORDER BY color_id"

	cardByCardIDQuery   = "SELECT %s FROM card_info WHERE card_number = ?"
	cardsByCardIDsQuery = "SELECT %s FROM card_info WHERE card_number IN (%s)"

	cardsByCardNamesQuery      = "SELECT %s FROM card_info WHERE card_name IN (%s)"
	searchCardUsingEffectQuery = "SELECT %s FROM card_info WHERE MATCH(card_effect) AGAINST(? IN BOOLEAN MODE) AND card_name NOT IN (%s) ORDER BY color_id, card_name"

	archetypeInclusionSubQuery = `SELECT %s FROM card_info WHERE MATCH (card_effect) AGAINST ('+"This card is always treated as" +"%s"' IN BOOLEAN MODE)`
	archetypeExclusionSubQuery = `SELECT %s FROM card_info WHERE MATCH (card_effect) AGAINST ('+"This card is not treated as" +"%s"'  IN BOOLEAN MODE)`

	archetypalCardsUsingCardNameQuery    = "SELECT %s FROM card_info WHERE card_name LIKE BINARY ? ORDER BY card_name"
	archetypalCardsUsingCardTextQuery    = `SELECT a.* FROM (%s) a WHERE a.card_effect REGEXP 'always treated as a.*"%s".* card' ORDER BY card_name`
	nonArchetypalCardsUsingCardTextQuery = `SELECT a.* FROM (%s) a WHERE a.card_effect REGEXP 'not treated as.*"%s".* card' ORDER BY card_name`

	randomCardQuery              = "SELECT %s FROM card_info WHERE card_color != 'Token' ORDER BY RAND() LIMIT 1"
	randomCardWithBlackListQuery = "SELECT %s FROM card_info WHERE card_number NOT IN (%s) AND card_color != 'Token' ORDER BY RAND() LIMIT 1"
)

const (
	productDetailsQuery   = "SELECT product_id, product_locale, product_name, product_type, product_sub_type, product_release_date FROM products where product_id = ?"
	cardsByProductIDQuery = "SELECT %s, product_position, card_rarity FROM product_contents WHERE product_id= ? ORDER BY product_position"
	productInfoByIDs      = "SELECT product_id, product_locale, product_name, product_type, product_sub_type, product_release_date, product_content_total FROM product_info WHERE product_id IN (%s)"
)

func convertToFullText(subject string) string {
	fullTextSubject := spaceRegex.ReplaceAllString(strings.ReplaceAll(subject, "-", " "), " +")
	return fmt.Sprintf(`"+%s"`, fullTextSubject) // match phrase, not all words in text will match only consecutive matches of words in phrase
}

func buildVariableQuerySubjects(subjects []string) ([]interface{}, int) {
	numSubjects := len(subjects)
	args := make([]interface{}, numSubjects)

	for index, subject := range subjects {
		args[index] = subject
	}

	return args, numSubjects
}

func variablePlaceholders(totalFields int) string {
	switch totalFields {
	case 0:
		return ""
	case 1:
		return "?"
	default:
		return fmt.Sprintf("?%s", strings.Repeat(", ?", totalFields-1))
	}
}
