package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	cUtil "github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cardAttributes = `
card_number,
card_color,
card_name,
card_attribute,
card_effect,
monster_type,
monster_attack,
monster_defense`

	dbVersionQuery = "SELECT VERSION()"

	cardColorIDsQuery = `
SELECT
	color_id,
	card_color
FROM
	card_colors
ORDER BY
	color_id`

	cardByCardIDQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	card_number = ?`
	cardsByCardIDsQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	card_number IN (%s)`

	cardsByCardNamesQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	card_name IN (%s)`
	searchCardUsingEffectQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	MATCH (card_effect) AGAINST (? IN BOOLEAN MODE)
ORDER BY
	color_id,
	card_name`

	archetypeInclusionSubQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	MATCH (card_effect) AGAINST ('+"This card is always treated as" +"%s"' IN BOOLEAN MODE)`
	archetypeExclusionSubQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	MATCH (card_effect) AGAINST ('+"This card is not treated as" +"%s"' IN BOOLEAN MODE)`

	archetypalCardsUsingCardNameQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	card_name LIKE BINARY ?
ORDER BY
	card_name`
	archetypalCardsUsingCardTextQuery = `
SELECT
	a.*
FROM
	(%s) a
WHERE
	a.card_effect REGEXP 'always treated as a.*"%s".* card'
ORDER BY
	card_name`
	nonArchetypalCardsUsingCardTextQuery = `
SELECT
	a.*
FROM
	(%s) a
WHERE
	a.card_effect REGEXP 'not treated as.*"%s".* card'
ORDER BY
	card_name`

	randomCardQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	card_color != 'Token'
ORDER BY
	RAND()
LIMIT
	1`
	randomCardWithBlackListQuery = `
SELECT
	%s
FROM
	card_info
WHERE
	card_number NOT IN (%s)
	AND card_color != 'Token'
ORDER BY
	RAND()
LIMIT
	1`
)

func queryCard(logger *slog.Logger, query string, args []interface{}) (*ygo.Card, *status.Status) {
	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)

	if err := skcDBConn.QueryRow(query, args...).Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
		return nil, handleQueryError(logger, err)
	}

	card := model.NewYGOCardProtoBuilder(id, name).WithColor(color).
		WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build()
	return card, nil
}

func parseRowsForCards(ctx context.Context, rows *sql.Rows, keyFn func(*ygo.Card) string) (map[string]*ygo.Card, *status.Status) {
	cards := make(map[string]*ygo.Card)

	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)
	for rows.Next() {
		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(util.RetrieveLogger(ctx), err)
		} else {
			card := model.NewYGOCardProtoBuilder(id, name).WithColor(color).
				WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build()
			cards[keyFn(card)] = card
		}
	}

	return cards, nil
}

func parseRowsForCardList(ctx context.Context, rows *sql.Rows) ([]*ygo.Card, *status.Status) {
	cardList := make([]*ygo.Card, 0)

	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)
	for rows.Next() {
		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(util.RetrieveLogger(ctx), err)
		} else {
			cardList = append(cardList,
				model.NewYGOCardProtoBuilder(id, name).WithColor(color).
					WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build())
		}
	}

	return cardList, nil
}

type CardRepository interface {
	GetCardColorIDs(context.Context) (*ygo.CardColors, *status.Status)

	GetCardByID(context.Context, string) (*ygo.Card, *status.Status)
	GetCardsByIDs(context.Context, model.CardIDs) (*ygo.Cards, *status.Status)

	GetCardsByNames(context.Context, model.CardNames) (*ygo.Cards, *status.Status)
	GetCardsReferencingNameInEffect(context.Context, []string) (*ygo.CardList, *status.Status)

	GetArchetypalCardsUsingCardName(context.Context, string) (*ygo.CardList, *status.Status)
	GetExplicitArchetypalInclusions(context.Context, string) (*ygo.CardList, *status.Status)
	GetExplicitArchetypalExclusions(context.Context, string) (*ygo.CardList, *status.Status)

	GetRandomCard(context.Context, []string) (*ygo.Card, *status.Status)
}
type YGOCardRepository struct{}

// Get IDs for all card colors currently in database.
func (imp YGOCardRepository) GetCardColorIDs(ctx context.Context) (*ygo.CardColors, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info("Retrieving card colors")

	if rows, err := skcDBConn.Query(cardColorIDsQuery); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		cardColorIDs := make(map[string]uint32, 18)
		for rows.Next() {
			var colorId uint32
			var cardColor string

			if err := rows.Scan(&colorId, &cardColor); err != nil {
				return nil, handleRowParsingError(logger, err)
			}

			cardColorIDs[cardColor] = colorId
		}

		logger.Info(fmt.Sprintf("Retrieved %d card colors", len(cardColorIDs)))
		return &ygo.CardColors{Values: cardColorIDs}, nil
	}
}

func (imp YGOCardRepository) GetCardByID(ctx context.Context, cardID string) (*ygo.Card, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card data using ID %v", cardID))

	args := make([]interface{}, 1)
	args[0] = cardID
	query := fmt.Sprintf(cardByCardIDQuery, cardAttributes)

	c, err := queryCard(logger, query, args)
	if err != nil && err.Code() == codes.NotFound {
		logger.Info("Card ID is not valid")
	}
	return c, err
}

func (imp YGOCardRepository) GetCardsByIDs(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card data using ID's: %v", cardIDs))

	args, numCards := buildVariableQuerySubjects(cardIDs)
	query := fmt.Sprintf(cardsByCardIDsQuery, cardAttributes, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCards(ctx, rows, model.CardIDAsKey); err != nil {
			return nil, err
		} else {
			return &ygo.Cards{
				CardInfo:         cards,
				UnknownResources: model.FindMissingKeys(cards, cardIDs),
			}, nil
		}
	}
}

// Uses card names to find instance of card
func (imp YGOCardRepository) GetCardsByNames(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card data using %d different name(s)", len(cardNames)))

	args, numCards := buildVariableQuerySubjects(cardNames)
	query := fmt.Sprintf(cardsByCardNamesQuery, cardAttributes, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCards(ctx, rows, model.CardNameAsKey); err != nil {
			return nil, err
		} else {
			return &ygo.Cards{
				CardInfo:         cards,
				UnknownResources: model.FindMissingKeys(cards, cardNames),
			}, nil
		}
	}
}

func (imp YGOCardRepository) GetCardsReferencingNameInEffect(ctx context.Context, namesOfCards []string) (*ygo.CardList, *status.Status) {
	numCards := len(namesOfCards)
	logger := cUtil.RetrieveLogger(ctx)
	if numCards == 0 {
		logger.Info("User did not provide any card names, responding w/ empty list of references")
		return &ygo.CardList{Cards: []*ygo.Card{}}, nil
	} else {
		logger.Info(fmt.Sprintf("Retrieving cards that reference one or more of the following cards by name in their text: %v", namesOfCards))
	}

	fullTextNames := make([]string, numCards)
	for ind, name := range namesOfCards {
		fullTextNames[ind] = convertToFullText(name)
	}

	query := fmt.Sprintf(searchCardUsingEffectQuery, cardAttributes)
	if rows, err := skcDBConn.Query(query, strings.Join(fullTextNames, " ")); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}

func (imp YGOCardRepository) GetArchetypalCardsUsingCardName(ctx context.Context, archetypeName string) (*ygo.CardList, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card data from DB for all cards that reference archetype %s in their name", archetypeName))
	searchTerm := `%` + archetypeName + `%`

	query := fmt.Sprintf(archetypalCardsUsingCardNameQuery, cardAttributes)
	if rows, err := skcDBConn.Query(query, searchTerm); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}

func (imp YGOCardRepository) GetExplicitArchetypalInclusions(ctx context.Context, archetypeName string) (*ygo.CardList, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving cards that are explicitly considered part of archetype %s", archetypeName))

	subQuery := fmt.Sprintf(archetypeInclusionSubQuery, cardAttributes, archetypeName)
	query := fmt.Sprintf(archetypalCardsUsingCardTextQuery, subQuery, archetypeName)
	if rows, err := skcDBConn.Query(query); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}
func (imp YGOCardRepository) GetExplicitArchetypalExclusions(ctx context.Context, archetypeName string) (*ygo.CardList, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving cards that are explicitly NOT considered part of archetype %s", archetypeName))

	subQuery := fmt.Sprintf(archetypeExclusionSubQuery, cardAttributes, archetypeName)
	query := fmt.Sprintf(nonArchetypalCardsUsingCardTextQuery, subQuery, archetypeName)
	if rows, err := skcDBConn.Query(query); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}

func (imp YGOCardRepository) GetRandomCard(ctx context.Context, blacklistedCards []string) (*ygo.Card, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving random card from DB. Client has provided %d blacklisted IDs", len(blacklistedCards)))

	// pick correct query based on contents of blacklistedCards
	numBlackListed := len(blacklistedCards)
	var query string
	var args []interface{}
	if numBlackListed == 0 {
		query = fmt.Sprintf(randomCardQuery, cardAttributes)
	} else {
		args, _ = buildVariableQuerySubjects(blacklistedCards)
		query = fmt.Sprintf(randomCardWithBlackListQuery, cardAttributes, variablePlaceholders(numBlackListed))
	}

	c, err := queryCard(logger, query, args)
	logger.Info(fmt.Sprintf("Random card determined to be; ID: %s, Name: %s", c.ID, c.Name))
	return c, err
}
