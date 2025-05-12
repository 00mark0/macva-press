package db

import (
	//"log"
	"testing"

	"github.com/00mark0/macva-press/utils"

	"context"

	"fmt"
	"github.com/stretchr/testify/require"
)

func createMedia(t *testing.T) []Medium {
	content := createRandomContent(t)
	var media []Medium

	medium1, err := testQueries.InsertMedia(context.Background(), InsertMediaParams{
		ContentID:    content.ContentID,
		MediaType:    "image",
		MediaUrl:     fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
		MediaCaption: "Random Image",
		MediaOrder:   1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, medium1)
	media = append(media, medium1)

	medium2, err := testQueries.InsertMedia(context.Background(), InsertMediaParams{
		ContentID:    content.ContentID,
		MediaType:    "image",
		MediaUrl:     fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
		MediaCaption: "Random Image",
		MediaOrder:   2,
	})
	require.NoError(t, err)
	require.NotEmpty(t, medium2)
	media = append(media, medium2)

	medium3, err := testQueries.InsertMedia(context.Background(), InsertMediaParams{
		ContentID:    content.ContentID,
		MediaType:    "video",
		MediaUrl:     "https://samplelib.com/lib/preview/mp4/sample-5s.mp4",
		MediaCaption: "Random Video",
		MediaOrder:   3,
	})
	require.NoError(t, err)
	require.NotEmpty(t, medium3)
	media = append(media, medium3)

	return media
}

func TestInsertMedia(t *testing.T) {
	media := createMedia(t)

	for _, medium := range media {
		require.NotEmpty(t, medium)
		require.NotEmpty(t, medium.MediaID)
		require.NotEmpty(t, medium.ContentID)
		require.NotEmpty(t, medium.MediaType)
		require.NotEmpty(t, medium.MediaUrl)
		require.NotEmpty(t, medium.MediaCaption)
		require.NotEmpty(t, medium.MediaOrder)
	}
}

func TestUpdateMedia(t *testing.T) {
	media := createMedia(t)

	updatedMedia, err := testQueries.UpdateMedia(context.Background(), UpdateMediaParams{
		MediaUrl:     media[0].MediaUrl,
		MediaCaption: "Specific Image",
		MediaOrder:   media[0].MediaOrder,
		MediaID:      media[0].MediaID,
	})
	require.NoError(t, err)
	require.Equal(t, "Random Image", media[0].MediaCaption)
	require.Equal(t, "Specific Image", updatedMedia.MediaCaption)
}

// this one tests both the ListMediaForContent and DeleteMedia
func TestDeleteMedia(t *testing.T) {
	content := createRandomContent(t)

	medium, err := testQueries.InsertMedia(context.Background(), InsertMediaParams{
		ContentID:    content.ContentID,
		MediaType:    "image",
		MediaUrl:     fmt.Sprintf("https://picsum.photos/600/400?random=%d", utils.RandomInt(1, 1000)),
		MediaCaption: "Random Image",
		MediaOrder:   1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, medium)

	mediaList, err := testQueries.ListMediaForContent(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.NotEmpty(t, mediaList)
	require.Len(t, mediaList, 1)
	require.Equal(t, medium.MediaID, mediaList[0].MediaID)

	err = testQueries.DeleteMedia(context.Background(), medium.MediaID)
	require.NoError(t, err)

	mediaListDel, err := testQueries.ListMediaForContent(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Empty(t, mediaListDel)
}

func TestBatchUpdateMediaOrder(t *testing.T) {
	media := createMedia(t)
	require.Equal(t, media[0].MediaOrder, int32(1))
	require.Equal(t, media[2].MediaOrder, int32(3))

	err := testQueries.BatchUpdateMediaOrder(context.Background(), BatchUpdateMediaOrderParams{
		Media1ID:    media[0].MediaID,
		Media1Order: 3,
		Media2ID:    media[2].MediaID,
		Media2Order: 1,
	})

	media2, err := testQueries.ListMediaForContent(context.Background(), media[0].ContentID)
	require.NoError(t, err)

	for _, medium := range media2 {
		if medium.MediaType == "video" {
			require.Equal(t, medium.MediaOrder, int32(1))
		}
	}
}
