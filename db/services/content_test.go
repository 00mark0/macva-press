package db

import (
	"log"
	"testing"

	"context"

	//"github.com/00mark0/macva-press/utils"
	"github.com/00mark0/macva-press/utils"
	"github.com/go-loremipsum/loremipsum"
	"github.com/stretchr/testify/require"
)

var Loremipsumgen = loremipsum.New()

func createRandomContent(t *testing.T) Content {
	users, err := testQueries.GetAdminUsers(context.Background())
	require.NoError(t, err)

	category, err := testQueries.ListCategories(context.Background(), 10)

	arg := CreateContentParams{
		UserID:              users[0].UserID,
		CategoryID:          category[0].CategoryID,
		Title:               Loremipsumgen.Sentence(),
		ContentDescription:  Loremipsumgen.Paragraphs(10),
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := testQueries.CreateContent(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.UserID, content.UserID)
	require.Equal(t, arg.CategoryID, content.CategoryID)
	require.Equal(t, arg.Title, content.Title)
	require.Equal(t, arg.ContentDescription, content.ContentDescription)

	return content
}

func TestCreateContent(t *testing.T) {
	var contents []Content

	for i := 0; i < 10; i++ {
		contents = append(contents, createRandomContent(t))
	}

	for _, content := range contents {
		require.NotEmpty(t, content)
		require.NotEmpty(t, content.ContentID)
		require.NotEmpty(t, content.UserID)
		require.NotEmpty(t, content.CategoryID)
		require.NotEmpty(t, content.Title)
		require.NotEmpty(t, content.ContentDescription)
	}
}

func TestUpdateContent(t *testing.T) {
	content1 := createRandomContent(t)

	arg := UpdateContentParams{
		ContentID:           content1.ContentID,
		Title:               Loremipsumgen.Sentence(),
		ContentDescription:  Loremipsumgen.Paragraphs(8),
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content2, err := testQueries.UpdateContent(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.ContentID, content2.ContentID)
	require.Equal(t, arg.Title, content2.Title)
	require.Equal(t, arg.ContentDescription, content2.ContentDescription)
	require.Equal(t, arg.CommentsEnabled, content2.CommentsEnabled)
	require.Equal(t, arg.ViewCountEnabled, content2.ViewCountEnabled)
	require.Equal(t, arg.LikeCountEnabled, content2.LikeCountEnabled)
	require.Equal(t, arg.DislikeCountEnabled, content2.DislikeCountEnabled)
}

func TestPublishContent(t *testing.T) {
	for i := 0; i < 10; i++ {
		content1 := createRandomContent(t)

		content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
		require.NoError(t, err)

		require.Equal(t, content2.Status, "published")
		require.Equal(t, content2.PublishedAt.Valid, true)
		require.NotNil(t, content2.PublishedAt.Time)
		require.NotEqual(t, content2.Status, content1.Status)
	}
}

func TestSoftDeleteContent(t *testing.T) {
	content1 := createRandomContent(t)

	content2, err := testQueries.SoftDeleteContent(context.Background(), content1.ContentID)
	require.NoError(t, err)

	require.Equal(t, content1.IsDeleted.Bool, false)
	require.Equal(t, content2.IsDeleted.Bool, true)
}

func TestHardDeleteContent(t *testing.T) {
	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	count1, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count1)

	content3, err := testQueries.HardDeleteContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content3.ContentID)

	count2, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count1-1)
}

func TestGetContentDetails(t *testing.T) {
	content1 := createRandomContent(t)

	tag := createRandomTag(t)

	err := testQueries.AddTagToContent(context.Background(), AddTagToContentParams{
		ContentID: content1.ContentID,
		TagID:     tag.TagID,
	})
	require.NoError(t, err)

	content2, err := testQueries.GetContentDetails(context.Background(), content1.ContentID)
	require.NoError(t, err)

	require.Equal(t, content1.ContentID, content2.ContentID)
	require.Equal(t, content1.UserID, content2.UserID)
	require.Equal(t, content1.CategoryID, content2.CategoryID)
	require.Equal(t, content1.Title, content2.Title)
	require.Equal(t, content1.ContentDescription, content2.ContentDescription)
	require.Equal(t, content1.CommentsEnabled, content2.CommentsEnabled)
	require.Equal(t, content1.ViewCountEnabled, content2.ViewCountEnabled)
	require.Equal(t, content1.LikeCountEnabled, content2.LikeCountEnabled)
	require.Equal(t, content1.DislikeCountEnabled, content2.DislikeCountEnabled)
	require.Equal(t, content1.Status, content2.Status)
	require.Equal(t, content1.ViewCount, content2.ViewCount)
	require.Equal(t, content1.LikeCount, content2.LikeCount)
	require.Equal(t, content1.DislikeCount, content2.DislikeCount)
	require.Equal(t, content1.CommentCount, content2.CommentCount)
	require.Equal(t, content1.CreatedAt, content2.CreatedAt)
	require.Equal(t, content1.UpdatedAt, content2.UpdatedAt)
	require.Equal(t, content1.PublishedAt, content2.PublishedAt)
	require.Equal(t, content1.IsDeleted, content2.IsDeleted)
	log.Println(content2.Username)
	log.Println(content2.CategoryName)
	log.Println(content2.Tags)
}

func TestGetPublishedContentCount(t *testing.T) {
	count1, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count1)

	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	count2, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count1+1)
}

func TestListPublishedContent(t *testing.T) {
	content, err := testQueries.ListPublishedContent(context.Background(), ListPublishedContentParams{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, content)
	require.LessOrEqual(t, len(content), 10)

	for _, cont := range content {
		require.NotEmpty(t, cont)
		require.NotEmpty(t, cont.ContentID)
		require.NotEmpty(t, cont.UserID)
		require.NotEmpty(t, cont.CategoryID)
		require.NotEmpty(t, cont.Title)
		require.NotEmpty(t, cont.ContentDescription)
		log.Println(cont.Username)
	}
}

func TestGetContentByCategoryCount(t *testing.T) {
	categories, err := testQueries.ListCategories(context.Background(), 10)
	content := createRandomContent(t)
	var category Category

	for _, cat := range categories {
		if content.CategoryID == cat.CategoryID {
			category = cat
		}
	}

	count1, err := testQueries.GetContentByCategoryCount(context.Background(), category.CategoryID)
	require.NoError(t, err)
	require.NotZero(t, count1)

	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	count2, err := testQueries.GetContentByCategoryCount(context.Background(), category.CategoryID)
	require.NoError(t, err)
	require.Equal(t, count2, count1+1)
}

func TestListContentByCategory(t *testing.T) {
	categories, err := testQueries.ListCategories(context.Background(), 10)
	var category Category

	content := createRandomContent(t)

	for _, cat := range categories {
		if content.CategoryID == cat.CategoryID {
			category = cat
		}
	}

	contents, err := testQueries.ListContentByCategory(context.Background(), ListContentByCategoryParams{CategoryID: category.CategoryID, Limit: 10, Offset: 0})

	require.NoError(t, err)
	require.NotEmpty(t, contents)
	require.LessOrEqual(t, len(contents), 10)

	for _, cont := range contents {
		require.NotEmpty(t, cont)
		require.NotEmpty(t, cont.ContentID)
		require.NotEmpty(t, cont.UserID)
		require.NotEmpty(t, cont.CategoryID)
		require.NotEmpty(t, cont.Title)
		require.NotEmpty(t, cont.ContentDescription)
		log.Println(cont.Username)
	}
}

func TestGetContentByTagCount(t *testing.T) {
	tag := createRandomTag(t)
	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	arg := AddTagToContentParams{
		ContentID: content2.ContentID,
		TagID:     tag.TagID,
	}

	err = testQueries.AddTagToContent(context.Background(), arg)
	require.NoError(t, err)

	count1, err := testQueries.GetContentByTagCount(context.Background(), tag.TagName)
	require.NoError(t, err)
	require.NotZero(t, count1)

	content3 := createRandomContent(t)

	content4, err := testQueries.PublishContent(context.Background(), content3.ContentID)
	require.NoError(t, err)
	require.Equal(t, content3.ContentID, content4.ContentID)

	arg2 := AddTagToContentParams{
		ContentID: content4.ContentID,
		TagID:     tag.TagID,
	}

	err = testQueries.AddTagToContent(context.Background(), arg2)
	require.NoError(t, err)

	count2, err := testQueries.GetContentByTagCount(context.Background(), tag.TagName)
	require.NoError(t, err)
	require.Equal(t, count2, count1+1)
}

func TestListContentByTag(t *testing.T) {
	tag := createRandomTag(t)
	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	arg := AddTagToContentParams{
		ContentID: content2.ContentID,
		TagID:     tag.TagID,
	}

	err = testQueries.AddTagToContent(context.Background(), arg)
	require.NoError(t, err)

	contents, err := testQueries.ListContentByTag(context.Background(), ListContentByTagParams{TagName: tag.TagName, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, contents)
	require.Equal(t, len(contents), 1)
	log.Println(contents[0].Username)
	log.Println(contents[0].CategoryName)
}

func createContentInteractive(title, description string) Content {
	users, err := testQueries.GetAdminUsers(context.Background())

	category, err := testQueries.ListCategories(context.Background(), 10)

	arg := CreateContentParams{
		UserID:              users[0].UserID,
		CategoryID:          category[0].CategoryID,
		Title:               title,
		ContentDescription:  description,
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := testQueries.CreateContent(context.Background(), arg)
	if err != nil {
		log.Println(err)
	}

	return content
}

func TestSearchContent(t *testing.T) {
	titleSearchTerm := "test_" + utils.RandomString(5)

	content := createContentInteractive(titleSearchTerm, Loremipsumgen.Paragraphs(10))

	content2, err := testQueries.PublishContent(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Equal(t, content.ContentID, content2.ContentID)
	require.Equal(t, content.Title, content2.Title)
	require.Equal(t, content.ContentDescription, content2.ContentDescription)

	arg := SearchContentParams{
		Limit:      10,
		SearchTerm: titleSearchTerm,
	}

	contents, err := testQueries.SearchContent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contents)
	require.Equal(t, len(contents), 1)
	require.Equal(t, contents[0].ContentID, content.ContentID)
	require.Equal(t, contents[0].Title, content.Title)
	require.Equal(t, contents[0].ContentDescription, content.ContentDescription)
	log.Println(contents[0].Username)
	log.Println(contents[0].CategoryName)
}

func TestIncrementViewCount(t *testing.T) {
	content := createRandomContent(t)

	viewCount, err := testQueries.IncrementViewCount(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Equal(t, viewCount, content.ViewCount+1)
}

// this one test both TestInsertOrUpdateContentReaction and FetchContentReactions
func TestInsertOrUpdateContentReaction(t *testing.T) {
	content := createRandomContent(t)
	users, err := testQueries.GetActiveUsers(context.Background(), 10)
	require.NoError(t, err)

	contentID, err := testQueries.InsertOrUpdateContentReaction(context.Background(), InsertOrUpdateContentReactionParams{
		ContentID: content.ContentID,
		UserID:    users[0].UserID,
		Reaction:  "like",
	})
	require.NoError(t, err)
	require.Equal(t, contentID, content.ContentID)

	contentID2, err := testQueries.InsertOrUpdateContentReaction(context.Background(), InsertOrUpdateContentReactionParams{
		ContentID: content.ContentID,
		UserID:    users[1].UserID,
		Reaction:  "dislike",
	})
	require.NoError(t, err)
	require.Equal(t, contentID2, content.ContentID)

	contentID3, err := testQueries.InsertOrUpdateContentReaction(context.Background(), InsertOrUpdateContentReactionParams{
		ContentID: content.ContentID,
		UserID:    users[1].UserID,
		Reaction:  "like",
	})
	require.NoError(t, err)
	require.Equal(t, contentID3, content.ContentID)

	reactions, err := testQueries.FetchContentReactions(context.Background(), FetchContentReactionsParams{
		ContentID: content.ContentID,
		Limit:     10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, reactions)
	require.Equal(t, len(reactions), 2)

	for _, reaction := range reactions {
		require.Equal(t, reaction.ContentID, content.ContentID)
		require.Equal(t, reaction.Reaction, "like")
		log.Println(reaction.Username)
	}
}

func TestDeleteContentReaction(t *testing.T) {
	content := createRandomContent(t)
	users, err := testQueries.GetActiveUsers(context.Background(), 10)
	require.NoError(t, err)

	contentID, err := testQueries.InsertOrUpdateContentReaction(context.Background(), InsertOrUpdateContentReactionParams{
		ContentID: content.ContentID,
		UserID:    users[0].UserID,
		Reaction:  "like",
	})
	require.NoError(t, err)
	require.Equal(t, contentID, content.ContentID)

	contentID2, err := testQueries.DeleteContentReaction(context.Background(), DeleteContentReactionParams{
		ContentID: content.ContentID,
		UserID:    users[0].UserID,
	})
	require.NoError(t, err)
	require.Equal(t, contentID2, content.ContentID)

	reactions, err := testQueries.FetchContentReactions(context.Background(), FetchContentReactionsParams{ContentID: content.ContentID, Limit: 10})
	require.NoError(t, err)
	require.Empty(t, reactions)
}

func TestUpdateContentLikeDislikeCount(t *testing.T) {
	content := createRandomContent(t)
	users, err := testQueries.GetActiveUsers(context.Background(), 10)

	contentID, err := testQueries.InsertOrUpdateContentReaction(context.Background(), InsertOrUpdateContentReactionParams{
		ContentID: content.ContentID,
		UserID:    users[0].UserID,
		Reaction:  "like",
	})
	require.NoError(t, err)
	require.Equal(t, contentID, content.ContentID)

	content2, err := testQueries.UpdateContentLikeDislikeCount(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Equal(t, content2.ContentID, content.ContentID)
	require.Equal(t, content2.LikeCount, content.LikeCount+1)
	require.Equal(t, content2.DislikeCount, content.DislikeCount)
}

func TestListTrendingContent(t *testing.T) {
	content1 := createRandomContent(t)
	content2 := createRandomContent(t)

	content3, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	content4, err := testQueries.PublishContent(context.Background(), content2.ContentID)
	require.NoError(t, err)

	require.Equal(t, content3.ContentID, content1.ContentID)
	require.Equal(t, content4.ContentID, content2.ContentID)

	viewCount, err := testQueries.IncrementViewCount(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, viewCount, content1.ViewCount+1)

	viewCount2, err := testQueries.IncrementViewCount(context.Background(), content2.ContentID)
	require.NoError(t, err)
	require.Equal(t, viewCount2, content2.ViewCount+1)

	trendingContent, err := testQueries.ListTrendingContent(context.Background(), ListTrendingContentParams{PublishedAt: content3.PublishedAt, Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, trendingContent)
	require.Equal(t, len(trendingContent), 2)

	for _, cont := range trendingContent {
		require.Equal(t, int32(1), cont.TotalInteractions)
	}
}

func TestGetContentOverview(t *testing.T) {
	count1, err := testQueries.GetContentOverview(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, count1)

	_ = createRandomContent(t)
	content2 := createRandomContent(t)
	content3 := createRandomContent(t)

	_, err = testQueries.PublishContent(context.Background(), content2.ContentID)

	_, err = testQueries.SoftDeleteContent(context.Background(), content3.ContentID)

	count2, err := testQueries.GetContentOverview(context.Background())
	require.NoError(t, err)

	require.Equal(t, count2.DraftCount, count1.DraftCount+1)
	require.Equal(t, count2.PublishedCount, count1.PublishedCount+1)
	require.Equal(t, count2.DeletedCount, count1.DeletedCount+1)
}

func TestListRelatedContent(t *testing.T) {
	content1 := createRandomContent(t)
	content2 := createRandomContent(t)

	content1Pub, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1Pub.ContentID, content1.ContentID)
	content2Pub, err := testQueries.PublishContent(context.Background(), content2.ContentID)
	require.NoError(t, err)
	require.Equal(t, content2Pub.ContentID, content2.ContentID)

	tag := createTagInteractive("test2")

	err = testQueries.AddTagToContent(context.Background(), AddTagToContentParams{
		ContentID: content1Pub.ContentID,
		TagID:     tag.TagID,
	})
	require.NoError(t, err)

	err = testQueries.AddTagToContent(context.Background(), AddTagToContentParams{
		ContentID: content2Pub.ContentID,
		TagID:     tag.TagID,
	})
	require.NoError(t, err)

	relatedContent, err := testQueries.ListRelatedContent(context.Background(), ListRelatedContentParams{
		ContentID: content1Pub.ContentID,
		Limit:     10,
	})

	for _, content := range relatedContent {
		log.Println(content.Title)
	}

	require.NoError(t, err)
	require.NotEmpty(t, relatedContent)
	require.Equal(t, 1, len(relatedContent))
	require.Equal(t, relatedContent[0].ContentID, content2Pub.ContentID)
}
