package api

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ygoCardRestrictionServiceServer) GetEffectiveTimelineForFormat(ctx context.Context, req *ygo.Format) (*ygo.EffectiveTimeline, error) {
	format := req.Value
	logger, newCtx := util.NewLogger(ctx, "Format Timeline", slog.String("format", format))

	if !strings.EqualFold(format, "Genesys") {
		logger.Error("Format not supported")
		return nil, status.New(codes.InvalidArgument, "Format not supported").Err()
	}

	if effectiveDates, err := cardRestrictionRepo.GetDatesForFormat(newCtx, format); err != nil {
		return nil, err.Err()
	} else {
		today := time.Now().In(chicagoLocation).Truncate(24 * time.Hour)
		futureDates := []string{}
		var activeDate string

		for _, effectiveDateStr := range effectiveDates {
			effectiveDate, _ := time.Parse(effectiveDateStr, "2006-01-02")
			if effectiveDate.After(today) {
				futureDates = append(futureDates, effectiveDateStr)
			} else {
				activeDate = effectiveDateStr
				break
			}
		}

		return &ygo.EffectiveTimeline{
			AllDates:    effectiveDates,
			FutureDates: futureDates,
			ActiveDate:  activeDate}, nil
	}
}
