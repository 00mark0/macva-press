package db

import (
	"fmt"
	"testing"
	"time"

	"context"

	"github.com/00mark0/macva-press/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomAd(t *testing.T) Ad {
	arg := CreateAdParams{
		Title:       pgtype.Text{String: Loremipsumgen.Sentence(), Valid: true},
		Description: pgtype.Text{String: Loremipsumgen.Paragraph(), Valid: true},
		ImageUrl:    pgtype.Text{String: fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)), Valid: true},
		TargetUrl:   pgtype.Text{String: "https://google.com", Valid: true},
		Placement:   pgtype.Text{String: utils.RandomPlacement(), Valid: true},
		Status:      pgtype.Text{String: "inactive", Valid: true},
		StartDate: pgtype.Timestamptz{
			Valid: false,
		},
		EndDate: pgtype.Timestamptz{
			Valid: false,
		},
	}

	ad, err := testQueries.CreateAd(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ad)
	require.NotEmpty(t, ad.ID)

	require.Equal(t, arg.Title.String, ad.Title.String)
	require.Equal(t, arg.Description.String, ad.Description.String)
	require.Equal(t, arg.ImageUrl.String, ad.ImageUrl.String)
	require.Equal(t, arg.TargetUrl.String, ad.TargetUrl.String)
	require.Equal(t, arg.Placement.String, ad.Placement.String)
	require.Equal(t, arg.Status.String, ad.Status.String)

	return ad
}

func TestCreateAd(t *testing.T) {
	var ads []Ad

	for i := 0; i < 5; i++ {
		ad := createRandomAd(t)

		ads = append(ads, ad)
	}

	for _, ad := range ads {
		require.NotEmpty(t, ad)
		require.NotEmpty(t, ad.ID)
		require.NotEmpty(t, ad.Title)
		require.NotEmpty(t, ad.Description)
		require.NotEmpty(t, ad.ImageUrl)
		require.NotEmpty(t, ad.TargetUrl)
		require.NotEmpty(t, ad.Placement)
		require.NotEmpty(t, ad.Status)
	}
}

func TestUpdateAd(t *testing.T) {
	ad := createRandomAd(t)

	arg := UpdateAdParams{
		Title: pgtype.Text{String: utils.RandomString(10), Valid: true},
		Description: pgtype.Text{
			String: utils.RandomString(100),
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
			Valid:  true,
		},
		TargetUrl: pgtype.Text{String: "https://yourmom.zip", Valid: true},
		Placement: pgtype.Text{String: utils.RandomPlacement(), Valid: true},
		Status:    pgtype.Text{String: "active", Valid: true},
		StartDate: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
		ID: ad.ID,
	}

	updatedAd, err := testQueries.UpdateAd(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, arg.ID.Bytes, updatedAd.ID.Bytes)
	require.Equal(t, arg.Title.String, updatedAd.Title.String)
	require.Equal(t, arg.Description.String, updatedAd.Description.String)
	require.Equal(t, arg.ImageUrl.String, updatedAd.ImageUrl.String)
	require.Equal(t, arg.TargetUrl.String, updatedAd.TargetUrl.String)
	require.Equal(t, arg.Placement.String, updatedAd.Placement.String)
	require.Equal(t, arg.Status.String, updatedAd.Status.String)
	require.WithinDuration(t, arg.StartDate.Time, updatedAd.StartDate.Time, time.Second)
	require.WithinDuration(t, arg.EndDate.Time, updatedAd.EndDate.Time, time.Second)
}

func TestDeactivateAd(t *testing.T) {
	ad := createRandomAd(t)

	arg := UpdateAdParams{
		Title: pgtype.Text{String: utils.RandomString(10), Valid: true},
		Description: pgtype.Text{
			String: utils.RandomString(100),
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
			Valid:  true,
		},
		TargetUrl: pgtype.Text{String: "https://yourmom.zip", Valid: true},
		Placement: pgtype.Text{String: utils.RandomPlacement(), Valid: true},
		Status:    pgtype.Text{String: "active", Valid: true},
		StartDate: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
		ID: ad.ID,
	}

	updatedAd, err := testQueries.UpdateAd(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, arg.ID.Bytes, updatedAd.ID.Bytes)
	require.Equal(t, arg.Title.String, updatedAd.Title.String)
	require.Equal(t, arg.Description.String, updatedAd.Description.String)
	require.Equal(t, arg.ImageUrl.String, updatedAd.ImageUrl.String)
	require.Equal(t, arg.TargetUrl.String, updatedAd.TargetUrl.String)
	require.Equal(t, arg.Placement.String, updatedAd.Placement.String)
	require.Equal(t, arg.Status.String, updatedAd.Status.String)
	require.WithinDuration(t, arg.StartDate.Time, updatedAd.StartDate.Time, time.Second)
	require.WithinDuration(t, arg.EndDate.Time, updatedAd.EndDate.Time, time.Second)

	deactivatedAd, err := testQueries.DeactivateAd(context.Background(), updatedAd.ID)
	require.NoError(t, err)

	require.Equal(t, updatedAd.ID.Bytes, deactivatedAd.ID.Bytes)
	require.Equal(t, updatedAd.Title.String, deactivatedAd.Title.String)
	require.Equal(t, updatedAd.Description.String, deactivatedAd.Description.String)
	require.Equal(t, updatedAd.ImageUrl.String, deactivatedAd.ImageUrl.String)
	require.Equal(t, updatedAd.TargetUrl.String, deactivatedAd.TargetUrl.String)
	require.Equal(t, updatedAd.Placement.String, deactivatedAd.Placement.String)
	require.Equal(t, ad.Status.String, deactivatedAd.Status.String)
	require.WithinDuration(t, ad.StartDate.Time, deactivatedAd.StartDate.Time, time.Second)
	require.WithinDuration(t, ad.EndDate.Time, deactivatedAd.EndDate.Time, time.Second)
}

// this one tests both ListAds and DeleteAd
func TestDeleteAd(t *testing.T) {
	ad := createRandomAd(t)

	ads1, err := testQueries.ListAds(context.Background(), 100)
	require.NoError(t, err)
	require.NotEmpty(t, ads1)

	err = testQueries.DeleteAd(context.Background(), ad.ID)
	require.NoError(t, err)

	ads2, err := testQueries.ListAds(context.Background(), 100)
	require.NoError(t, err)
	require.NotEmpty(t, ads2)
	require.Equal(t, len(ads1)-1, len(ads2))
}

func TestGetAd(t *testing.T) {
	ad1 := createRandomAd(t)

	ad2, err := testQueries.GetAd(context.Background(), ad1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ad2)

	require.Equal(t, ad1.ID, ad2.ID)
	require.Equal(t, ad1.Title.String, ad2.Title.String)
	require.Equal(t, ad1.Description.String, ad2.Description.String)
	require.Equal(t, ad1.ImageUrl.String, ad2.ImageUrl.String)
	require.Equal(t, ad1.TargetUrl.String, ad2.TargetUrl.String)
	require.Equal(t, ad1.Placement.String, ad2.Placement.String)
	require.Equal(t, ad1.Status.String, ad2.Status.String)
	require.WithinDuration(t, ad1.StartDate.Time, ad2.StartDate.Time, time.Second)
	require.WithinDuration(t, ad1.EndDate.Time, ad2.EndDate.Time, time.Second)
}

func TestListInactiveAds(t *testing.T) {
	inactiveAds, err := testQueries.ListInactiveAds(context.Background(), 100)
	require.NoError(t, err)
	require.NotEmpty(t, inactiveAds)

	_ = createRandomAd(t)

	inactiveAds2, err := testQueries.ListInactiveAds(context.Background(), 100)
	require.NoError(t, err)
	require.NotEmpty(t, inactiveAds2)
	require.Equal(t, len(inactiveAds)+1, len(inactiveAds2))
}

func TestListActiveAds(t *testing.T) {
	activeAds, err := testQueries.ListActiveAds(context.Background(), 100)
	require.NoError(t, err)
	require.NotEmpty(t, activeAds)

	ad := createRandomAd(t)

	arg := UpdateAdParams{
		Title: pgtype.Text{String: utils.RandomString(10), Valid: true},
		Description: pgtype.Text{
			String: utils.RandomString(100),
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
			Valid:  true,
		},
		TargetUrl: pgtype.Text{String: "https://yourmom.zip", Valid: true},
		Placement: pgtype.Text{String: utils.RandomPlacement(), Valid: true},
		Status:    pgtype.Text{String: "active", Valid: true},
		StartDate: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
		ID: ad.ID,
	}

	updatedAd, err := testQueries.UpdateAd(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, ad.ID, updatedAd.ID)
	require.Equal(t, activeAds[0].Status.String, updatedAd.Status.String)

	activeAds2, err := testQueries.ListActiveAds(context.Background(), 100)
	require.NoError(t, err)
	require.NotEmpty(t, activeAds2)
	require.Equal(t, len(activeAds)+1, len(activeAds2))
}

func TestListAdsByPlacement(t *testing.T) {
	adsByPlacement, err := testQueries.ListAdsByPlacement(context.Background(), ListAdsByPlacementParams{
		Placement: pgtype.Text{String: "footer", Valid: true},
		Limit:     100,
	})
	require.NoError(t, err)
	require.NotEmpty(t, adsByPlacement)

	ad := createRandomAd(t)

	arg := UpdateAdParams{
		Title: pgtype.Text{String: utils.RandomString(10), Valid: true},
		Description: pgtype.Text{
			String: utils.RandomString(100),
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
			Valid:  true,
		},
		TargetUrl: pgtype.Text{String: "https://yourmom.zip", Valid: true},
		Placement: pgtype.Text{String: "footer", Valid: true},
		Status:    pgtype.Text{String: "active", Valid: true},
		StartDate: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
		ID: ad.ID,
	}

	updatedAd, err := testQueries.UpdateAd(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, ad.ID, updatedAd.ID)
	require.Equal(t, adsByPlacement[0].Placement.String, updatedAd.Placement.String)

	adsByPlacement2, err := testQueries.ListAdsByPlacement(context.Background(), ListAdsByPlacementParams{
		Placement: pgtype.Text{String: "footer", Valid: true},
		Limit:     100,
	})
	require.NoError(t, err)
	require.NotEmpty(t, adsByPlacement2)
	require.Equal(t, len(adsByPlacement)+1, len(adsByPlacement2))
}

func TestIncrementAdClicks(t *testing.T) {
	ad := createRandomAd(t)
	require.Zero(t, ad.Clicks.Int32)

	clicks, err := testQueries.IncrementAdClicks(context.Background(), ad.ID)
	require.NoError(t, err)
	require.NotZero(t, clicks.Int32)
	require.Equal(t, ad.Clicks.Int32+1, clicks.Int32)
}
