package db

import (
	"testing"
	"time"

	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateDailyAnalytics(t *testing.T) {
	now := time.Now().UTC() // consistent timezone
	todayMidnight := now.Truncate(24 * time.Hour)

	analytics, err := testQueries.CreateDailyAnalytics(context.Background(), pgtype.Date{Time: todayMidnight, Valid: true})
	require.NoError(t, err)

	require.Equal(t, int32(0), analytics.TotalViews)
	require.Equal(t, int32(0), analytics.TotalLikes)
	require.Equal(t, int32(0), analytics.TotalDislikes)
	require.Equal(t, int32(0), analytics.TotalComments)
	require.Equal(t, int32(0), analytics.TotalAdsClicks)
}

func TestGetDailyAnalytics(t *testing.T) {
	now := time.Now().UTC() // consistent timezone
	todayMidnight := now.Truncate(24 * time.Hour)
	tomorrowMidnight := todayMidnight.Add(24 * time.Hour)

	arg := GetDailyAnalyticsParams{
		AnalyticsDate: pgtype.Date{
			Time:  todayMidnight,
			Valid: true,
		},
		AnalyticsDate_2: pgtype.Date{
			Time:  tomorrowMidnight,
			Valid: true,
		},
		Limit: 10,
	}

	analitycs, err := testQueries.GetDailyAnalytics(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, analitycs)
	require.Len(t, analitycs, 1)
}

func TestUpdateDailyAnalytics(t *testing.T) {
	tomorrow := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)
	dayAfterTomorrow := tomorrow.Add(24 * time.Hour)

	tomorrowAnalytic, err := testQueries.CreateDailyAnalytics(context.Background(), pgtype.Date{Time: tomorrow, Valid: true})
	require.NoError(t, err)

	getTAnalytic, err := testQueries.GetDailyAnalytics(context.Background(), GetDailyAnalyticsParams{
		AnalyticsDate: pgtype.Date{
			Time:  tomorrow,
			Valid: true,
		},
		AnalyticsDate_2: pgtype.Date{
			Time:  dayAfterTomorrow,
			Valid: true,
		},
		Limit: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, getTAnalytic)
	require.Len(t, getTAnalytic, 1)
	require.Equal(t, tomorrowAnalytic.AnalyticsDate, getTAnalytic[0].AnalyticsDate)

	updateTAnalyticArg := UpdateDailyAnalyticsParams{
		AnalyticsDate:  getTAnalytic[0].AnalyticsDate,
		TotalViews:     10,
		TotalLikes:     10,
		TotalDislikes:  10,
		TotalComments:  10,
		TotalAdsClicks: 10,
	}

	updatedTAnalytic, err := testQueries.UpdateDailyAnalytics(context.Background(), updateTAnalyticArg)
	require.NoError(t, err)

	require.Equal(t, getTAnalytic[0].AnalyticsDate, updatedTAnalytic.AnalyticsDate)
	require.Equal(t, int32(updateTAnalyticArg.TotalViews), updatedTAnalytic.TotalViews)
	require.Equal(t, int32(updateTAnalyticArg.TotalLikes), updatedTAnalytic.TotalLikes)
	require.Equal(t, int32(updateTAnalyticArg.TotalDislikes), updatedTAnalytic.TotalDislikes)
	require.Equal(t, int32(updateTAnalyticArg.TotalComments), updatedTAnalytic.TotalComments)
	require.Equal(t, int32(updateTAnalyticArg.TotalAdsClicks), updatedTAnalytic.TotalAdsClicks)
}

func TestAggregateAnalytics(t *testing.T) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)
	dayAfterTomorrow := tomorrow.Add(24 * time.Hour)

	analytics, err := testQueries.GetDailyAnalytics(context.Background(), GetDailyAnalyticsParams{
		AnalyticsDate: pgtype.Date{
			Time:  today,
			Valid: true,
		},
		AnalyticsDate_2: pgtype.Date{
			Time:  dayAfterTomorrow,
			Valid: true,
		},
		Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, analytics, 2)

	arg := AggregateAnalyticsParams{
		AnalyticsDate: pgtype.Date{
			Time:  today,
			Valid: true,
		},
		AnalyticsDate_2: pgtype.Date{
			Time:  dayAfterTomorrow,
			Valid: true,
		},
	}

	aggregate, err := testQueries.AggregateAnalytics(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, int64(10), aggregate.TotalViews)
	require.Equal(t, int64(10), aggregate.TotalLikes)
	require.Equal(t, int64(10), aggregate.TotalDislikes)
	require.Equal(t, int64(10), aggregate.TotalComments)
	require.Equal(t, int64(10), aggregate.TotalAdsClicks)
}
