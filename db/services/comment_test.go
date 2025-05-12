package db

import (
	"log"
	"testing"

	"github.com/00mark0/macva-press/utils"

	"context"

	"github.com/stretchr/testify/require"
)

func createRandomComment(t *testing.T) Comment {
	categories, err := testQueries.ListCategories(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	var category Category

	cont := createRandomContent(t)

	for _, cat := range categories {
		if cont.CategoryID == cat.CategoryID {
			category = cat
		}
	}

	content, err := testQueries.ListContentByCategory(context.Background(), ListContentByCategoryParams{CategoryID: category.CategoryID, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, content)

	users, err := testQueries.GetActiveUsers(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	randomUser := utils.RandomInt(0, int64(len(users)-1))
	randomContent := utils.RandomInt(0, int64(len(content)-1))
	arg := CreateCommentParams{
		UserID:      users[randomUser].UserID,
		ContentID:   content[randomContent].ContentID,
		CommentText: Loremipsumgen.Sentence(),
	}

	comment, err := testQueries.CreateComment(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.UserID, comment.UserID)
	require.Equal(t, arg.ContentID, comment.ContentID)
	require.Equal(t, arg.CommentText, comment.CommentText)

	return comment
}

func createRandomComments(t *testing.T) []Comment {
	categories, err := testQueries.ListCategories(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	var category Category

	cont := createRandomContent(t)

	for _, cat := range categories {
		if cont.CategoryID == cat.CategoryID {
			category = cat
		}
	}

	content, err := testQueries.ListContentByCategory(context.Background(), ListContentByCategoryParams{CategoryID: category.CategoryID, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, content)

	users, err := testQueries.GetActiveUsers(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	randomContent := utils.RandomInt(0, int64(len(content)-1))
	var comments []Comment

	for i := 0; i < 10; i++ {
		arg := CreateCommentParams{
			UserID:      users[utils.RandomInt(0, int64(len(users)-1))].UserID,
			ContentID:   content[randomContent].ContentID,
			CommentText: Loremipsumgen.Sentence(),
		}

		comment, err := testQueries.CreateComment(context.Background(), arg)
		require.NoError(t, err)
		require.Equal(t, arg.UserID, comment.UserID)
		require.Equal(t, arg.ContentID, comment.ContentID)
		require.Equal(t, arg.CommentText, comment.CommentText)

		comments = append(comments, comment)
	}

	return comments
}

func TestCreateComment(t *testing.T) {
	comment := createRandomComment(t)
	require.NotEmpty(t, comment)
	require.NotEmpty(t, comment.CommentID)
	require.NotEmpty(t, comment.UserID)
	require.NotEmpty(t, comment.ContentID)
	require.NotEmpty(t, comment.CommentText)
}

func TestCreateComments(t *testing.T) {
	comments := createRandomComments(t)
	require.NotEmpty(t, comments)
	require.Equal(t, 10, len(comments))

	for _, comment := range comments {
		require.NotEmpty(t, comment)
		require.NotEmpty(t, comment.CommentID)
		require.NotEmpty(t, comment.UserID)
		require.NotEmpty(t, comment.ContentID)
		require.NotEmpty(t, comment.CommentText)
	}
}

func TestUpdateComment(t *testing.T) {
	comment := createRandomComment(t)

	arg := UpdateCommentParams{
		CommentText: utils.RandomString(10),
		CommentID:   comment.CommentID,
	}

	updatedComment, err := testQueries.UpdateComment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, comment.CommentText, updatedComment.CommentText)
	require.NotEqual(t, comment.UpdatedAt, updatedComment.UpdatedAt)
	require.Equal(t, comment.CommentID, updatedComment.CommentID)
	require.Equal(t, arg.CommentText, updatedComment.CommentText)
}

// this one tests both ListContentComments and SoftDeleteComment
func TestSoftDeleteComment(t *testing.T) {
	comment := createRandomComment(t)

	comments1, err := testQueries.ListContentComments(context.Background(), ListContentCommentsParams{ContentID: comment.ContentID, Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, comments1)

	deleted, err := testQueries.SoftDeleteComment(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.Equal(t, true, deleted.IsDeleted.Bool)
	require.Equal(t, comment.CommentID, deleted.CommentID)

	comments2, err := testQueries.ListContentComments(context.Background(), ListContentCommentsParams{ContentID: comment.ContentID, Limit: 10})
	require.NoError(t, err)
	require.Equal(t, len(comments2), len(comments1)-1)
}

// this one tests InsertOrUpdateCommentReaction, FetchCommentReactions, UpdateCommentScore and DeleteCommentReaction
func TestInsertOrUpdateCommentReaction(t *testing.T) {
	comment := createRandomComment(t)

	arg := InsertOrUpdateCommentReactionParams{
		CommentID: comment.CommentID,
		UserID:    comment.UserID,
		Reaction:  "like",
	}

	comment_id, err := testQueries.InsertOrUpdateCommentReaction(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, comment_id, comment.CommentID)

	reactions, err := testQueries.FetchCommentReactions(context.Background(), FetchCommentReactionsParams{
		CommentID: comment_id,
		Limit:     10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, reactions)
	require.Equal(t, 1, len(reactions))
	require.Equal(t, comment_id, reactions[0].CommentID)
	require.Equal(t, reactions[0].Reaction, "like")
	require.Equal(t, comment.Score, int32(0))

	updatedComment, err := testQueries.UpdateCommentScore(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment)
	require.Equal(t, comment_id, updatedComment.CommentID)
	require.Equal(t, updatedComment.Score, int32(1))

	comment_id2, err := testQueries.DeleteCommentReaction(context.Background(), DeleteCommentReactionParams{
		CommentID: comment.CommentID,
		UserID:    comment.UserID,
	})
	require.NoError(t, err)
	require.Equal(t, comment_id2, comment_id)

	reactions2, err := testQueries.FetchCommentReactions(context.Background(), FetchCommentReactionsParams{
		CommentID: comment_id2,
		Limit:     10,
	})
	require.NoError(t, err)
	require.Empty(t, reactions2)

	updatedComment2, err := testQueries.UpdateCommentScore(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment2)
	require.Equal(t, comment_id2, updatedComment2.CommentID)
	require.Equal(t, updatedComment2.Score, int32(0))
}

func createCommentInteractive(content Content, user User) Comment {
	comment, err := testQueries.CreateComment(context.Background(), CreateCommentParams{
		UserID:      user.UserID,
		ContentID:   content.ContentID,
		CommentText: Loremipsumgen.Sentence(),
	})
	if err != nil {
		log.Println(err)
	}
	return comment
}

func TestListContentCommentsByScore(t *testing.T) {
	content := createRandomContent(t)

	contentPub, err := testQueries.PublishContent(context.Background(), content.ContentID)
	require.NoError(t, err)
	require.Equal(t, contentPub.Status, "published")

	users, err := testQueries.GetActiveUsers(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	comment1 := createCommentInteractive(contentPub, users[0])
	comment2 := createCommentInteractive(contentPub, users[1])
	comment3 := createCommentInteractive(contentPub, users[2])

	_, err = testQueries.InsertOrUpdateCommentReaction(context.Background(), InsertOrUpdateCommentReactionParams{
		CommentID: comment1.CommentID,
		UserID:    users[0].UserID,
		Reaction:  "like",
	})

	_, err = testQueries.InsertOrUpdateCommentReaction(context.Background(), InsertOrUpdateCommentReactionParams{
		CommentID: comment1.CommentID,
		UserID:    users[1].UserID,
		Reaction:  "like",
	})

	_, err = testQueries.InsertOrUpdateCommentReaction(context.Background(), InsertOrUpdateCommentReactionParams{
		CommentID: comment1.CommentID,
		UserID:    users[2].UserID,
		Reaction:  "like",
	})

	_, err = testQueries.InsertOrUpdateCommentReaction(context.Background(), InsertOrUpdateCommentReactionParams{
		CommentID: comment2.CommentID,
		UserID:    users[0].UserID,
		Reaction:  "like",
	})

	_, err = testQueries.InsertOrUpdateCommentReaction(context.Background(), InsertOrUpdateCommentReactionParams{
		CommentID: comment2.CommentID,
		UserID:    users[1].UserID,
		Reaction:  "like",
	})

	_, err = testQueries.InsertOrUpdateCommentReaction(context.Background(), InsertOrUpdateCommentReactionParams{
		CommentID: comment3.CommentID,
		UserID:    users[0].UserID,
		Reaction:  "like",
	})

	updatedComment1, err := testQueries.UpdateCommentScore(context.Background(), comment1.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment1)
	require.Equal(t, comment1.CommentID, updatedComment1.CommentID)
	require.Equal(t, int32(3), updatedComment1.Score)

	updatedComment2, err := testQueries.UpdateCommentScore(context.Background(), comment2.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment2)
	require.Equal(t, comment2.CommentID, updatedComment2.CommentID)
	require.Equal(t, int32(2), updatedComment2.Score)

	updatedComment3, err := testQueries.UpdateCommentScore(context.Background(), comment3.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment3)
	require.Equal(t, comment3.CommentID, updatedComment3.CommentID)
	require.Equal(t, int32(1), updatedComment3.Score)

	comments, err := testQueries.ListContentCommentsByScore(context.Background(), ListContentCommentsByScoreParams{
		ContentID: contentPub.ContentID,
		Limit:     10,
	})

	require.Equal(t, comments[0].CommentID, comment1.CommentID)
	require.Equal(t, comments[1].CommentID, comment2.CommentID)
	require.Equal(t, comments[2].CommentID, comment3.CommentID)

	require.Equal(t, int32(3), comments[0].Score)
	require.Equal(t, int32(2), comments[1].Score)
	require.Equal(t, int32(1), comments[2].Score)
}
