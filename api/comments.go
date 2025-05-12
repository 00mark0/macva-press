package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
	"github.com/labstack/echo/v4"
)

func (server *Server) listContentComments(ctx echo.Context) error {
	var req ListAdsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listContentComments:", err)
		return err
	}

	contentIDStr := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in listContentComments:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in listContentComments:", err)
	}

	userReactions, err := server.getUserReactionsForContentWithCache(ctx.Request().Context(), contentID, userData.UserID)
	if err != nil {
		log.Println("Error getting user reactions in listContentComments:", err)
		return err
	}

	nextLimit := req.Limit + 10

	comments, err := server.getCommentsWithCache(ctx.Request().Context(), contentID, nextLimit)
	if err != nil {
		log.Println("Error listing comments in listContentComments:", err)
		return err
	}

	url := fmt.Sprintf("/api/content/comments/%s?limit=", contentIDStr)

	commentCount, err := server.getCommentCountWithCache(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting comment count in listContentComments:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleComments(
		contentIDStr,
		comments,
		userData,
		userReactions,
		int(nextLimit),
		url,
		int(commentCount),
	))
}

func (server *Server) listContentCommentsScore(ctx echo.Context) error {
	var req ListAdsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listContentCommentsScore:", err)
		return err
	}

	contentIDStr := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in listContentCommentsScore:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in listContentCommentsScore:", err)
	}

	userReactions, err := server.getUserReactionsForContentWithCache(ctx.Request().Context(), contentID, userData.UserID)
	if err != nil {
		log.Println("Error getting user reactions in listContentCommentsScore:", err)
		return err
	}

	nextLimit := req.Limit + 10

	comments, err := server.getCommentsByScoreWithCache(ctx.Request().Context(), contentID, nextLimit)
	if err != nil {
		log.Println("Error listing comments in listContentCommentsScore:", err)
		return err
	}

	url := fmt.Sprintf("/api/content/comments/%s/score?limit=", contentIDStr)

	commentCount, err := server.getCommentCountWithCache(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting comment count in listContentCommentsScore:", err)
		return err
	}

	var convertedComments []db.ListContentCommentsRow
	for _, comment := range comments {
		convertedComments = append(convertedComments, db.ListContentCommentsRow{
			CommentID:       comment.CommentID,
			ContentID:       comment.ContentID,
			UserID:          comment.UserID,
			CommentText:     comment.CommentText,
			Score:           comment.Score,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
			IsDeleted:       comment.IsDeleted,
			ParentCommentID: comment.ParentCommentID,
			Username:        comment.Username,
			Pfp:             comment.Pfp,
			Role:            comment.Role,
		})
	}

	return Render(ctx, http.StatusOK, components.ArticleComments(
		contentIDStr,
		convertedComments,
		userData,
		userReactions,
		int(nextLimit),
		url,
		int(commentCount),
	))
}

type CreateCommentReq struct {
	CommentText string `form:"comment_text" validate:"required,min=1,max=10000"`
}

func (server *Server) createComment(ctx echo.Context) error {
	var req CreateCommentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createComment:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in createComment:", err)
		return err
	}

	contentIDStr := ctx.Param("id")

	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in createComment:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in createComment:", err)
	}

	_, err = server.store.CreateComment(ctx.Request().Context(), db.CreateCommentParams{
		ContentID:   contentID,
		UserID:      userData.UserID,
		CommentText: req.CommentText,
	})
	if err != nil {
		log.Println("Error creating comment:", err)
		return err
	}

	err = server.incrementDailyComments(ctx)
	if err != nil {
		log.Println(err)
	}

	err = server.store.IncrementCommentCount(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error incrementing comment count in createComment:", err)
		return err
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Printf("Failed to invalidate comment-related cache: %v", err)
	}

	return server.listContentComments(ctx)
}

// Handler for upvoting a comment
func (server *Server) handleUpvoteComment(ctx echo.Context) error {
	commentIDStr := ctx.Param("id")

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in handleUpvoteComment:", err)
	}

	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in handleUpvoteComment:", err)
		return err
	}

	// Check if the user already has a reaction
	userReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	// Handle reaction logic based on whether we found a reaction and what it was
	if err == nil {
		// User has an existing reaction
		if userReaction.Reaction == "like" {
			// If already liked, remove the reaction
			_, err = server.store.DeleteCommentReaction(ctx.Request().Context(), db.DeleteCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
			})
			if err != nil {
				log.Println("Error deleting comment reaction from like to remove like:", err)
				return err
			}
		} else if userReaction.Reaction == "dislike" {
			// If disliked, change to like
			_, err := server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
				Reaction:  "like",
			})
			if err != nil {
				log.Println("Error changing reaction from dislike to like:", err)
				return err
			}
		}
	} else {
		// No reaction yet, add a like
		_, err := server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
			CommentID: commentID,
			UserID:    userData.UserID,
			Reaction:  "like",
		})
		if err != nil {
			log.Println("Error adding new like reaction:", err)
			return err
		}
	}

	// Update the comment's score
	updatedComment, err := server.store.UpdateCommentScore(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error updating comment score:", err)
		return err
	}

	// Get the updated user reaction for the response
	reactionStatus := ""
	updatedUserReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	if err == nil {
		reactionStatus = updatedUserReaction.Reaction
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Printf("Failed to invalidate comment-related cache: %v", err)
	}

	// Render just the comment actions part
	return Render(ctx, http.StatusOK, components.CommentActions(updatedComment, userData, reactionStatus))
}

// Handler for downvoting a comment
func (server *Server) handleDownvoteComment(ctx echo.Context) error {
	commentIDStr := ctx.Param("id")

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in handleDownvoteComment:", err)
	}

	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in handleDownvoteComment:", err)
		return err
	}

	// Check if the user already has a reaction
	userReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	if err == nil {
		// User has an existing reaction
		if userReaction.Reaction == "dislike" {
			// If already disliked, remove the reaction
			_, err = server.store.DeleteCommentReaction(ctx.Request().Context(), db.DeleteCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
			})
			if err != nil {
				log.Println("Error deleting comment reaction from dislike to remove dislike:", err)
				return err
			}
		} else if userReaction.Reaction == "like" {
			// If liked, change to dislike
			_, err = server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
				Reaction:  "dislike",
			})
			if err != nil {
				log.Println("Error changing reaction from like to dislike:", err)
				return err
			}
		}
	} else {
		// No reaction yet, add a dislike
		_, err = server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
			CommentID: commentID,
			UserID:    userData.UserID,
			Reaction:  "dislike",
		})
		if err != nil {
			log.Println("Error adding new dislike reaction:", err)
			return err
		}
	}

	// Update the comment's score
	updatedComment, err := server.store.UpdateCommentScore(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error updating comment score:", err)
		return err
	}

	// Get the updated user reaction for the response
	reactionStatus := ""
	updatedUserReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	if err == nil {
		reactionStatus = updatedUserReaction.Reaction
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Printf("Failed to invalidate comment-related cache: %v", err)
	}

	// Render just the comment actions part
	return Render(ctx, http.StatusOK, components.CommentActions(updatedComment, userData, reactionStatus))
}

type CreateReplyReq struct {
	ReplyText string `form:"reply_text" validate:"required,min=1,max=10000"`
}

func (server *Server) createReply(ctx echo.Context) error {
	var req CreateReplyReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createReply:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in createReply:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in createReply:", err)
	}

	parentCommentIDStr := ctx.Param("id")
	parentCommentID, err := utils.ParseUUID(parentCommentIDStr, "parent comment ID")
	if err != nil {
		log.Println("Invalid parent comment ID format in createReply:", err)
		return err
	}

	parentComment, err := server.store.GetCommentByID(ctx.Request().Context(), parentCommentID)
	if err != nil {
		log.Println("Error getting parent comment in createReply:", err)
		return err
	}

	arg := db.CreateReplyParams{
		ParentCommentID: parentCommentID,
		UserID:          userData.UserID,
		ContentID:       parentComment.ContentID,
		CommentText:     req.ReplyText,
	}

	replyToReplyIDStr := ctx.FormValue("reply_to_reply")

	if replyToReplyIDStr != "" {
		replyToReplyID, err := utils.ParseUUID(replyToReplyIDStr, "reply to reply ID")
		if err != nil {
			log.Println("Invalid reply to reply ID format in createReply:", err)
			return err
		}

		replyToReplyComment, err := server.store.GetCommentByID(ctx.Request().Context(), replyToReplyID)
		if err != nil {
			log.Println("Error getting reply to reply comment in createReply:", err)
			return err
		}

		replyToReplyCommentUser, err := server.store.GetUserByID(ctx.Request().Context(), replyToReplyComment.UserID)
		if err != nil {
			log.Println("Error getting user in createReply:", err)
			return err
		}

		arg = db.CreateReplyParams{
			ParentCommentID: parentCommentID,
			UserID:          userData.UserID,
			ContentID:       parentComment.ContentID,
			CommentText:     fmt.Sprintf("@%s %s", replyToReplyCommentUser.Username, req.ReplyText),
		}
	}

	comment, err := server.store.CreateReply(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error creating reply:", err)
		return err
	}

	err = server.incrementDailyComments(ctx)
	if err != nil {
		log.Println(err)
	}

	err = server.store.IncrementCommentCount(ctx.Request().Context(), parentComment.ContentID)
	if err != nil {
		log.Println("Error incrementing comment count in createReply:", err)
		return err
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Printf("Failed to invalidate comment-related cache: %v", err)
	}

	convertedComment := db.ListContentCommentsRow{
		CommentID:       comment.CommentID,
		ContentID:       comment.ContentID,
		UserID:          comment.UserID,
		CommentText:     comment.CommentText,
		Score:           comment.Score,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		IsDeleted:       comment.IsDeleted,
		ParentCommentID: comment.ParentCommentID,
		Username:        userData.Username,
		Pfp:             userData.Pfp,
		Role:            userData.Role,
	}

	return Render(ctx, http.StatusOK, components.CommentReplyItem(convertedComment, userData, ""))
}

func (server *Server) listRepliesInfo(ctx echo.Context) error {
	commentIDStr := ctx.Param("id")

	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in listRepliesInfo:", err)
		return err
	}

	replyCount, adminPfp, err := server.getReplyCountAndAdminPfp(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error getting reply count and admin pfp:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.CommentReplyInfo(int(replyCount), adminPfp, commentIDStr))
}

func (server *Server) listCommentReplies(ctx echo.Context) error {
	var req ListAdsReq

	parentCommentIDStr := ctx.Param("id")
	parentCommentID, err := utils.ParseUUID(parentCommentIDStr, "parent comment ID")
	if err != nil {
		log.Println("Invalid parent comment ID format in listCommentReplies:", err)
		return err
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listCommentReplies:", err)
		return err
	}

	nextLimit := req.Limit + 10

	replies, err := server.listCommentRepliesWithCache(ctx.Request().Context(), parentCommentID, nextLimit)
	if err != nil {
		log.Println("Error listing comment replies:", err)
		return err
	}

	var convertedReplies []db.ListContentCommentsRow
	for _, reply := range replies {
		convertedReplies = append(convertedReplies, db.ListContentCommentsRow{
			CommentID:       reply.CommentID,
			ContentID:       reply.ContentID,
			UserID:          reply.UserID,
			CommentText:     reply.CommentText,
			Score:           reply.Score,
			CreatedAt:       reply.CreatedAt,
			UpdatedAt:       reply.UpdatedAt,
			IsDeleted:       reply.IsDeleted,
			ParentCommentID: reply.ParentCommentID,
			Username:        reply.Username,
			Pfp:             reply.Pfp,
			Role:            reply.Role,
		})
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in listCommentReplies:", err)
	}

	userReactions, err := server.getUserReactionsForContentWithCache(ctx.Request().Context(), replies[0].ContentID, userData.UserID)
	if err != nil {
		log.Println("Error fetching user reactions:", err)
	}

	url := fmt.Sprintf("/api/comments/%s/more-replies?limit=", parentCommentIDStr)

	return Render(ctx, http.StatusOK, components.CommentReplyList(convertedReplies, userData, userReactions, int(nextLimit), url))
}

type UpdateCommentReq struct {
	CommentText string `form:"edit_text" validate:"required,min=1,max=10000"`
}

func (server *Server) updateComment(ctx echo.Context) error {
	var req UpdateCommentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updateComment:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in updateComment:", err)
		return err
	}

	commentIDStr := ctx.Param("id")
	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in updateComment:", err)
		return err
	}

	arg := db.UpdateCommentParams{
		CommentID:   commentID,
		CommentText: req.CommentText,
	}

	commentData, err := server.store.UpdateComment(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating comment:", err)
		return err
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Printf("Failed to invalidate comment-related cache: %v", err)
	}

	return Render(ctx, http.StatusOK, components.EditCommentResponse(commentData))
}

func (server *Server) deleteComment(ctx echo.Context) error {
	commentIDStr := ctx.Param("id")
	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in deleteComment:", err)
		return err
	}

	_, err = server.store.DeleteComment(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error deleting comment:", err)
		return err
	}

	err = server.decrementDailyComments(ctx)
	if err != nil {
		log.Println(err)
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Printf("Failed to invalidate comment-related cache: %v", err)
	}

	return ctx.NoContent(http.StatusNoContent)
}
