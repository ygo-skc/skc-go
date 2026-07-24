package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ygo-skc/skc-go/common/v3/model"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
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

func queryCard(ctx context.Context, logger *slog.Logger, query string, args []any) (*ygo.Card, *status.Status) {
	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)

	if err := skcDBConn.QueryRowContext(ctx, query, args...).Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
		return nil, handleQueryError(logger, err)
	}

	card := model.NewYGOCardProtoBuilder(id, name).
		WithColor(color).
		WithAttribute(attribute).
		WithEffect(effect).
		WithMonsterType(monsterType).
		WithAttack(atk).
		WithDefense(def).
		Build()
	return card, nil
}

func parseCardRows[T []*ygo.Card | map[string]*ygo.Card](ctx context.Context, rows *sql.Rows, dataStructure *T, collector func(*T, *ygo.Card)) *status.Status {
	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)
	for rows.Next() {
		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return handleRowParsingError(util.RetrieveLogger(ctx), err)
		}
		collector(
			dataStructure,
			model.NewYGOCardProtoBuilder(id, name).
				WithColor(color).
				WithAttribute(attribute).
				WithEffect(effect).
				WithMonsterType(monsterType).
				WithAttack(atk).
				WithDefense(def).
				Build())
	}

	if err := rows.Err(); err != nil {
		return handleQueryError(util.RetrieveLogger(ctx), err)
	}
	return nil
}

func collectWithList(cardList *[]*ygo.Card, card *ygo.Card) {
	*cardList = append(*cardList, card)
}

func collectWithMapUsingIDKey(cards *map[string]*ygo.Card, card *ygo.Card) {
	(*cards)[model.CardIDAsKey(card)] = card
}

func collectWithMapUsingNameKey(cards *map[string]*ygo.Card, card *ygo.Card) {
	(*cards)[model.CardNameAsKey(card)] = card
}

type CardRepository interface {
	GetCardColorIDs(context.Context) (map[string]uint32, *status.Status)

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
func (imp YGOCardRepository) GetCardColorIDs(ctx context.Context) (map[string]uint32, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card colors")

	rows, err := skcDBConn.QueryContext(ctx, cardColorIDsQuery)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cardColorIDs := make(map[string]uint32, 18)
	for rows.Next() {
		var colorId uint32
		var cardColor string

		if err := rows.Scan(&colorId, &cardColor); err != nil {
			return nil, handleRowParsingError(logger, err)
		}

		cardColorIDs[cardColor] = colorId
	}

	if err := rows.Err(); err != nil {
		return nil, handleQueryError(logger, err)
	}

	logger.Info("Retrieved card colors", slog.Int("count", len(cardColorIDs)))
	return cardColorIDs, nil
}

func (imp YGOCardRepository) GetCardByID(ctx context.Context, cardID string) (*ygo.Card, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card data", slog.String("card_id", cardID))

	args := make([]any, 1)
	args[0] = cardID
	query := fmt.Sprintf(cardByCardIDQuery, cardAttributes)

	c, err := queryCard(ctx, logger, query, args)
	if err != nil && err.Code() == codes.NotFound {
		logger.Info("Card ID not found", slog.String("card_id", cardID))
	}
	return c, err
}

func (imp YGOCardRepository) GetCardsByIDs(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card data", slog.Any("card_ids", cardIDs))

	args, numCards := buildVariableQuerySubjects(cardIDs)
	query := fmt.Sprintf(cardsByCardIDsQuery, cardAttributes, variablePlaceholders(numCards))

	rows, err := skcDBConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cards := make(map[string]*ygo.Card, 0)
	if err := parseCardRows(ctx, rows, &cards, collectWithMapUsingIDKey); err != nil {
		return nil, err
	}
	return &ygo.Cards{
		CardInfo:         cards,
		UnknownResources: model.FindMissingKeys(cards, cardIDs),
	}, nil
}

// Uses card names to find instance of card
func (imp YGOCardRepository) GetCardsByNames(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card data", slog.Int("card_name_count", len(cardNames)))

	args, numCards := buildVariableQuerySubjects(cardNames)
	query := fmt.Sprintf(cardsByCardNamesQuery, cardAttributes, variablePlaceholders(numCards))

	rows, err := skcDBConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cards := make(map[string]*ygo.Card, 0)
	if err := parseCardRows(ctx, rows, &cards, collectWithMapUsingNameKey); err != nil {
		return nil, err
	}
	return &ygo.Cards{
		CardInfo:         cards,
		UnknownResources: model.FindMissingKeys(cards, cardNames),
	}, nil
}

func (imp YGOCardRepository) GetCardsReferencingNameInEffect(ctx context.Context, namesOfCards []string) (*ygo.CardList, *status.Status) {
	numCards := len(namesOfCards)
	logger := util.RetrieveLogger(ctx)
	if numCards == 0 {
		logger.Info("No card names provided, returning empty list of references")
		return &ygo.CardList{Cards: []*ygo.Card{}}, nil
	}
	logger.Info("Retrieving cards referencing card names in effect text", slog.Any("names", namesOfCards))

	fullTextNames := make([]string, numCards)
	for ind, name := range namesOfCards {
		fullTextNames[ind] = convertToFullText(name)
	}

	query := fmt.Sprintf(searchCardUsingEffectQuery, cardAttributes)
	rows, err := skcDBConn.QueryContext(ctx, query, strings.Join(fullTextNames, " "))
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cards := make([]*ygo.Card, 0)
	if err := parseCardRows(ctx, rows, &cards, collectWithList); err != nil {
		return nil, err
	}
	return &ygo.CardList{Cards: cards}, nil
}

func (imp YGOCardRepository) GetArchetypalCardsUsingCardName(ctx context.Context, archetypeName string) (*ygo.CardList, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving cards by archetype name", slog.String("archetype", archetypeName))
	searchTerm := `%` + archetypeName + `%`

	query := fmt.Sprintf(archetypalCardsUsingCardNameQuery, cardAttributes)
	rows, err := skcDBConn.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cards := make([]*ygo.Card, 0)
	if err := parseCardRows(ctx, rows, &cards, collectWithList); err != nil {
		return nil, err
	}
	return &ygo.CardList{Cards: cards}, nil
}

func (imp YGOCardRepository) GetExplicitArchetypalInclusions(ctx context.Context, archetypeName string) (*ygo.CardList, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving explicit archetype inclusions", slog.String("archetype", archetypeName))

	subQuery := fmt.Sprintf(archetypeInclusionSubQuery, cardAttributes, archetypeName)
	query := fmt.Sprintf(archetypalCardsUsingCardTextQuery, subQuery, archetypeName)
	rows, err := skcDBConn.QueryContext(ctx, query)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cards := make([]*ygo.Card, 0)
	if err := parseCardRows(ctx, rows, &cards, collectWithList); err != nil {
		return nil, err
	}
	return &ygo.CardList{Cards: cards}, nil
}
func (imp YGOCardRepository) GetExplicitArchetypalExclusions(ctx context.Context, archetypeName string) (*ygo.CardList, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving explicit archetype exclusions", slog.String("archetype", archetypeName))

	subQuery := fmt.Sprintf(archetypeExclusionSubQuery, cardAttributes, archetypeName)
	query := fmt.Sprintf(nonArchetypalCardsUsingCardTextQuery, subQuery, archetypeName)
	rows, err := skcDBConn.QueryContext(ctx, query)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	cards := make([]*ygo.Card, 0)
	if err := parseCardRows(ctx, rows, &cards, collectWithList); err != nil {
		return nil, err
	}
	return &ygo.CardList{Cards: cards}, nil
}

func (imp YGOCardRepository) GetRandomCard(ctx context.Context, blacklistedCards []string) (*ygo.Card, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving random card", slog.Int("blacklisted_count", len(blacklistedCards)))

	// pick correct query based on contents of blacklistedCards
	numBlackListed := len(blacklistedCards)
	var query string
	var args []any
	if numBlackListed == 0 {
		query = fmt.Sprintf(randomCardQuery, cardAttributes)
	} else {
		args, _ = buildVariableQuerySubjects(blacklistedCards)
		query = fmt.Sprintf(randomCardWithBlackListQuery, cardAttributes, variablePlaceholders(numBlackListed))
	}

	c, err := queryCard(ctx, logger, query, args)
	if err == nil {
		logger.Info("Random card selected", slog.String("card_id", c.Id), slog.String("card_name", c.Name))
	}
	return c, err
}
