package api

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ygoScoreServiceServer) GetDatesForFormat(ctx context.Context, req *ygo.ResourceName) (*ygo.Dates, error) {
	_, newCtx := util.NewLogger(ctx, "Dates for Format")

	format := req.Value
	if !strings.EqualFold(format, "Genesys") {
		return nil, status.New(codes.InvalidArgument, "Format not supported").Err()
	}

	if dates, err := scoreRepo.GetDatesForFormat(newCtx, format); err != nil {
		return nil, err.Err()
	} else {
		return &ygo.Dates{Dates: dates}, nil
	}
}

func (s *ygoScoreServiceServer) GetCardScoreByID(ctx context.Context, req *ygo.ResourceID) (*ygo.CardScore, error) {
	_, newCtx := util.NewLogger(ctx, "Card Score")

	if scoreHistory, err := scoreRepo.GetCardScoreByID(newCtx, req.ID); err != nil {
		return nil, err.Err()
	} else {
		currentScoreByFormat, uniqueFormats, scheduledChanges := parseScoreHistory(scoreHistory)
		return &ygo.CardScore{
			CurrentScoreByFormat: currentScoreByFormat,
			UniqueFormats:        uniqueFormats,
			ScoreHistory:         scoreHistory,
			ScheduledChanges:     scheduledChanges,
		}, nil
	}
}

func (s *ygoScoreServiceServer) GetCardScoresByIDs(ctx context.Context, req *ygo.ResourceIDs) (*ygo.CardScores, error) {
	// _, newCtx := util.NewLogger(ctx, "Card Score")

	return nil, nil
}

func parseScoreHistory(scoresHistory []*ygo.ScoreInstance) (map[string]uint32, []string, []string) {
	today := time.Now().In(chicagoLocation)
	todaysDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, chicagoLocation)

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

	return scoreByFormat, uniqueFormats, scheduledChanges
}
