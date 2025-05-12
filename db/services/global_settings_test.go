package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGlobalSettings(t *testing.T) {
	globalSettings, err := testQueries.CreateGlobalSettings(context.Background())
	require.NoError(t, err)

	require.Equal(t, false, globalSettings.DisableComments)
	require.Equal(t, false, globalSettings.DisableLikes)
	require.Equal(t, true, globalSettings.DisableDislikes)
	require.Equal(t, false, globalSettings.DisableViews)
	require.Equal(t, false, globalSettings.DisableAds)
}

func TestGetGlobalSettings(t *testing.T) {
	globalSettings, err := testQueries.GetGlobalSettings(context.Background())
	require.NoError(t, err)

	require.Equal(t, false, globalSettings[0].DisableComments)
	require.Equal(t, false, globalSettings[0].DisableLikes)
	require.Equal(t, true, globalSettings[0].DisableDislikes)
	require.Equal(t, false, globalSettings[0].DisableViews)
	require.Equal(t, false, globalSettings[0].DisableAds)
}

func TestUpdateGlobalSettings(t *testing.T) {
	arg := UpdateGlobalSettingsParams{
		DisableComments: true,
		DisableLikes:    true,
		DisableDislikes: false,
		DisableViews:    true,
		DisableAds:      true,
	}

	err := testQueries.UpdateGlobalSettings(context.Background(), arg)
	require.NoError(t, err)

	globalSettings, err := testQueries.GetGlobalSettings(context.Background())
	require.NoError(t, err)

	require.Equal(t, arg.DisableComments, globalSettings[0].DisableComments)
	require.Equal(t, arg.DisableLikes, globalSettings[0].DisableLikes)
	require.Equal(t, arg.DisableDislikes, globalSettings[0].DisableDislikes)
	require.Equal(t, arg.DisableViews, globalSettings[0].DisableViews)
	require.Equal(t, arg.DisableAds, globalSettings[0].DisableAds)
}

func TestResetGlobalSettings(t *testing.T) {
	err := testQueries.ResetGlobalSettings(context.Background())
	require.NoError(t, err)

	globalSettings, err := testQueries.GetGlobalSettings(context.Background())
	require.NoError(t, err)

	require.Equal(t, false, globalSettings[0].DisableComments)
	require.Equal(t, false, globalSettings[0].DisableLikes)
	require.Equal(t, true, globalSettings[0].DisableDislikes)
	require.Equal(t, false, globalSettings[0].DisableViews)
	require.Equal(t, false, globalSettings[0].DisableAds)
}
