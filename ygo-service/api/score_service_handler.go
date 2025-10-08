package api

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ygoScoreServiceServer) GetDatesForFormat(ctx context.Context, req *ygo.Format) (*ygo.Dates, error) {
	format := req.Value
	logger, newCtx := util.NewLogger(ctx, "Dates for Format", slog.String("format", format))

	if !strings.EqualFold(format, "Genesys") {
		logger.Error("Format not supported")
		return nil, status.New(codes.InvalidArgument, "Format not supported").Err()
	}

	if dates, err := scoreRepo.GetDatesForFormat(newCtx, format); err != nil {
		return nil, err.Err()
	} else {
		return &ygo.Dates{Dates: dates}, nil
	}
}

func (s *ygoScoreServiceServer) GetCardScoreByID(ctx context.Context, req *ygo.ResourceID) (*ygo.CardScore, error) {
	logger, newCtx := util.NewLogger(ctx, "Card Score", slog.String("card_id", req.ID))

	today := time.Now().In(chicagoLocation)
	todaysDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, chicagoLocation)
	if score, err := scoreRepo.GetCardScoreByID(newCtx, req.ID, todaysDate, parser); err != nil {
		return nil, err.Err()
	} else {
		if len(score.ScoreHistory) == 0 {
			logger.Error("Scores not retrieved since card ID DNE")
			return nil, status.New(codes.NotFound, "Resource not found").Err()
		}
		return score, nil
	}
}

func parser(score *ygo.CardScore, entry *ygo.ScoreEntry, todaysDate time.Time) {
	effectiveDate, _ := time.Parse("2006-01-02", entry.EffectiveDate)

	if _, exists := score.CurrentScoreByFormat[entry.Format]; !exists && effectiveDate.Before(todaysDate) {
		score.CurrentScoreByFormat[entry.Format] = entry.Score
	}

	if !slices.Contains(score.UniqueFormats, entry.Format) {
		score.UniqueFormats = append(score.UniqueFormats, entry.Format)
	}

	if effectiveDate.After(todaysDate) && !slices.Contains(score.ScheduledChanges, entry.Format) {
		score.ScheduledChanges = append(score.ScheduledChanges, fmt.Sprintf("%s|%s", entry.Format, entry.EffectiveDate))
	}

	score.ScoreHistory = append(score.ScoreHistory, entry)
}

func (s *ygoScoreServiceServer) GetCardScoresByIDs(ctx context.Context, req *ygo.ResourceIDs) (*ygo.CardScores, error) {
	_, newCtx := util.NewLogger(ctx, "Multi-card Score")

	if scoreHistory, err := scoreRepo.GetCardScoresByIDs(newCtx, req.IDs); err != nil {
		return nil, err.Err()
	} else {
		today := time.Now().In(chicagoLocation)
		todaysDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, chicagoLocation)
		scores := make(map[string]*ygo.CardScore, len(scoreHistory))

		for cardID, s := range scoreHistory {
			score := parseScoreHistory(s, todaysDate)
			score.ScoreHistory = s
			scores[cardID] = &score
		}

		return &ygo.CardScores{
			CardInfo:         scores,
			UnknownResources: model.FindMissingKeys(scoreHistory, model.CardIDs(req.IDs)),
		}, nil
	}
}

func parseScoreHistory(scoresHistory []*ygo.ScoreEntry, todaysDate time.Time) ygo.CardScore {
	scoreByFormat := make(map[string]uint32, 3)
	uniqueFormats := make([]string, 0, 3)
	scheduledChanges := make([]string, 0, 3)

	for _, instance := range scoresHistory {
		effectiveDate, _ := time.Parse("2006-01-02", instance.EffectiveDate)

		if _, exists := scoreByFormat[instance.Format]; !exists && effectiveDate.Before(todaysDate) {
			scoreByFormat[instance.Format] = instance.Score
		}

		if !slices.Contains(uniqueFormats, instance.Format) {
			uniqueFormats = append(uniqueFormats, instance.Format)
		}

		if effectiveDate.After(todaysDate) && !slices.Contains(scheduledChanges, instance.Format) {
			scheduledChanges = append(scheduledChanges, fmt.Sprintf("%s|%s", instance.Format, instance.EffectiveDate))
		}
	}

	return ygo.CardScore{
		CurrentScoreByFormat: scoreByFormat,
		UniqueFormats:        uniqueFormats,
		ScheduledChanges:     scheduledChanges,
	}
}
