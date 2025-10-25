package api

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ygoScoreServiceServer) GetScoresByFormatAndDate(ctx context.Context, req *ygo.RestrictedContentRequest) (*ygo.ScoresForFormatAndDate, error) {
	format := req.Format
	effectiveDate := req.EffectiveDate

	logger, newCtx := util.NewLogger(ctx, "Scores By Format & Date",
		slog.String("format", format),
		slog.String("effective_date", effectiveDate),
	)

	if entries, numEntries, err := scoreRepo.GetScoresByFormatAndDate(newCtx, format, effectiveDate); err != nil {
		return nil, err.Err()
	} else {
		if numEntries == 0 {
			logger.Error("Cannot find format and date combination")
			return nil, status.New(codes.NotFound, "Format and date combination DNE").Err()
		}

		return &ygo.ScoresForFormatAndDate{
			Format:        format,
			EffectiveDate: effectiveDate,
			Entries:       entries,
			TotalEntries:  numEntries,
		}, nil
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

func (s *ygoScoreServiceServer) GetCardScoresByIDs(ctx context.Context, req *ygo.ResourceIDs) (*ygo.CardScores, error) {
	_, newCtx := util.NewLogger(ctx, "Multi-card Score")

	today := time.Now().In(chicagoLocation)
	todaysDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, chicagoLocation)
	if scores, err := scoreRepo.GetCardScoresByIDs(newCtx, req.IDs, todaysDate, parser); err != nil {
		return nil, err.Err()
	} else {
		return &ygo.CardScores{
			CardInfo:         scores,
			UnknownResources: model.FindMissingKeys(scores, model.CardIDs(req.IDs)),
		}, nil
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
