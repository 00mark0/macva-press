package db

import (
	"log"
	"testing"

	"context"
	//"github.com/00mark0/macva-press/utils"
	"github.com/00mark0/macva-press/utils"
	"github.com/stretchr/testify/require"
)

func createTagInteractive(name string) Tag {
	tag, err := testQueries.CreateTag(context.Background(), name)
	if err != nil {
		log.Println(err)
	}

	return tag
}

func createRandomTag(t *testing.T) Tag {
	tag, err := testQueries.CreateTag(context.Background(), Loremipsumgen.Word())
	require.NoError(t, err)
	require.NotEmpty(t, tag)

	return tag
}

func TestCreateTag(t *testing.T) {
	var tags []Tag

	for i := 0; i < 10; i++ {
		tags = append(tags, createRandomTag(t))
	}

	for _, tag := range tags {
		require.NotEmpty(t, tag)
		require.NotEmpty(t, tag.TagID)
		require.NotEmpty(t, tag.TagName)
	}
}

func TestUpdateTag(t *testing.T) {
	tag1 := createRandomTag(t)

	arg := UpdateTagParams{
		TagID:   tag1.TagID,
		TagName: Loremipsumgen.Word(),
	}

	tag2, err := testQueries.UpdateTag(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, tag1.TagID, tag2.TagID)
	require.NotEqual(t, tag1.TagName, tag2.TagName)
}

func TestDeleteTag(t *testing.T) {
	tag1 := createRandomTag(t)
	err := testQueries.DeleteTag(context.Background(), tag1.TagID)
	require.NoError(t, err)
}

func TestGetTag(t *testing.T) {
	tag1 := createRandomTag(t)

	tag2, err := testQueries.GetTag(context.Background(), tag1.TagID)
	require.NoError(t, err)
	require.Equal(t, tag1.TagID, tag2.TagID)
	require.Equal(t, tag1.TagName, tag2.TagName)
}

func TestListTags(t *testing.T) {
	tags, err := testQueries.ListTags(context.Background(), 5)
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.LessOrEqual(t, len(tags), 5)

	for _, tag := range tags {
		require.NotEmpty(t, tag)
		require.NotEmpty(t, tag.TagID)
		require.NotEmpty(t, tag.TagName)
	}
}

func TestSearchTags(t *testing.T) {
	searchTerm := "hleb_" + utils.RandomString(5)
	createTagInteractive(searchTerm)

	arg := SearchTagsParams{
		Limit:  5,
		Search: searchTerm,
	}

	tags, err := testQueries.SearchTags(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.LessOrEqual(t, len(tags), 1)
	require.Equal(t, tags[0].TagName, arg.Search)
}

// this tests both AddTagToContent and GetTagsByContent
func TestAddTagToContent(t *testing.T) {
	content := createRandomContent(t)
	tag := createRandomTag(t)

	arg := AddTagToContentParams{
		ContentID: content.ContentID,
		TagID:     tag.TagID,
	}

	err := testQueries.AddTagToContent(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, content.ContentID, arg.ContentID)
	require.Equal(t, tag.TagID, arg.TagID)

	tags, err := testQueries.GetTagsByContent(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Equal(t, tags[0].TagID, tag.TagID)
	require.Equal(t, tags[0].TagName, tag.TagName)
}

func TestRemoveTagFromContent(t *testing.T) {
	content := createRandomContent(t)
	tag := createRandomTag(t)

	arg := AddTagToContentParams{
		ContentID: content.ContentID,
		TagID:     tag.TagID,
	}

	err := testQueries.AddTagToContent(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, content.ContentID, arg.ContentID)
	require.Equal(t, tag.TagID, arg.TagID)

	arg2 := RemoveTagFromContentParams{
		ContentID: content.ContentID,
		TagID:     tag.TagID,
	}

	err = testQueries.RemoveTagFromContent(context.Background(), arg2)
	require.NoError(t, err)

	tags, err := testQueries.GetTagsByContent(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Empty(t, tags)
}
