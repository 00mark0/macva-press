package api

import (
	"context"
	"log"
	"time"

	"github.com/00mark0/macva-press/db/redis"
	"github.com/00mark0/macva-press/db/services"
	"github.com/jackc/pgx/v5/pgtype"
	redisClient "github.com/redis/go-redis/v9"
)

func (server *Server) getCommentsWithCache(ctx context.Context, contentID pgtype.UUID, limit int32) ([]db.ListContentCommentsRow, error) {
	// Generate the cache key
	cacheKey := redis.GenerateKey("comments", contentID, limit)

	// Try to get comments from cache
	var comments []db.ListContentCommentsRow
	cacheHit, err := server.cacheService.Get(ctx, cacheKey, &comments)
	if err != nil {
		log.Printf("Error fetching comments from cache: %v", err)
	}

	if cacheHit {
		// Cache hit, return cached comments
		return comments, nil
	}

	// Cache miss, fetch from database
	log.Printf("Cache miss for comments: %s", cacheKey)
	arg := db.ListContentCommentsParams{
		ContentID: contentID,
		Limit:     limit,
	}
	comments, err = server.store.ListContentComments(ctx, arg)
	if err != nil {
		return nil, err
	}

	// Store in cache for future use
	err = server.cacheService.Set(ctx, cacheKey, comments, 10*time.Minute)
	if err != nil {
		log.Printf("Error caching comments: %v", err)
	}

	return comments, nil
}

func (server *Server) getCommentsByScoreWithCache(ctx context.Context, contentID pgtype.UUID, limit int32) ([]db.ListContentCommentsByScoreRow, error) {
	// Generate the cache key
	cacheKey := redis.GenerateKey("comments_by_score", contentID, limit)

	// Try to get comments from cache
	var comments []db.ListContentCommentsByScoreRow
	cacheHit, err := server.cacheService.Get(ctx, cacheKey, &comments)
	if err != nil {
		log.Printf("Error fetching comments from cache: %v", err)
	}

	if cacheHit {
		// Cache hit, return cached comments
		return comments, nil
	}

	// Cache miss, fetch from database
	log.Printf("Cache miss for comments: %s", cacheKey)
	arg := db.ListContentCommentsByScoreParams{
		ContentID: contentID,
		Limit:     limit,
	}
	comments, err = server.store.ListContentCommentsByScore(ctx, arg)
	if err != nil {
		return nil, err
	}

	// Store in cache for future use
	err = server.cacheService.Set(ctx, cacheKey, comments, 10*time.Minute)
	if err != nil {
		log.Printf("Error caching comments: %v", err)
	}

	return comments, nil
}

func (server *Server) getReplyCountAndAdminPfp(ctx context.Context, parentCommentID pgtype.UUID) (int64, string, error) {
	// Generate cache keys
	checkedCacheKey := redis.GenerateKey("comments_checked_admin_replies", parentCommentID)
	adminPfpCacheKey := redis.GenerateKey("comments_admin_pfp", parentCommentID)
	countCacheKey := redis.GenerateKey("comments_reply_count", parentCommentID)

	// Check if we have the reply count cached
	var replyCount int64
	countCacheHit, err := server.cacheService.Get(ctx, countCacheKey, &replyCount)
	if err != nil && err != redisClient.Nil {
		log.Printf("Error checking cached reply count for ParentCommentID: %s: %v", parentCommentID, err)
	}

	// If count not in cache, fetch from database
	if !countCacheHit {
		log.Printf("Cache miss for reply count for ParentCommentID: %s, fetching from DB", parentCommentID)
		replyCount, err = server.store.GetReplyCount(ctx, parentCommentID)
		if err != nil {
			log.Printf("Error fetching reply count for ParentCommentID: %s: %v", parentCommentID, err)
			return 0, "", err
		}

		// Cache the reply count
		err = server.cacheService.Set(ctx, countCacheKey, replyCount, 10*time.Minute)
		if err != nil {
			log.Printf("Error caching reply count for ParentCommentID: %s: %v", parentCommentID, err)
		}
	}

	// Check if we've already checked this parent comment for admin replies
	var hasAdminReply bool
	adminCheckCacheHit, err := server.cacheService.Get(ctx, checkedCacheKey, &hasAdminReply)
	if err != nil && err != redisClient.Nil {
		log.Printf("Error checking admin reply status for ParentCommentID: %s: %v", parentCommentID, err)
	}

	var adminPfp string

	if adminCheckCacheHit {
		// We already know if this comment has admin replies
		if hasAdminReply {
			// Only try to get the cached pfp if we know there's an admin reply
			pfpCacheHit, err := server.cacheService.Get(ctx, adminPfpCacheKey, &adminPfp)
			if err != nil && err != redisClient.Nil {
				log.Printf("Error fetching cached admin pfp for ParentCommentID: %s: %v", parentCommentID, err)
			}

			if pfpCacheHit && adminPfp != "" {
				// Cache hit: use cached admin pfp
			} else {
				// Cache miss: admin pfp needs to be rescanned
				log.Printf("Admin pfp cache miss for ParentCommentID: %s, rescanning replies", parentCommentID)
				adminPfp = scanForAdminPfp(ctx, server, parentCommentID)
			}
		}
	} else {
		// If no cache hit, need to scan replies for admin replies
		log.Printf("Cache miss for checked_admin_replies for ParentCommentID: %s, scanning for admin", parentCommentID)
		adminPfp = scanForAdminPfp(ctx, server, parentCommentID)

		// Cache the checked status
		hasAdminReply = adminPfp != ""
		err = server.cacheService.Set(ctx, checkedCacheKey, hasAdminReply, 10*time.Minute)
		if err != nil {
			log.Printf("Error caching admin reply status for ParentCommentID: %s: %v", parentCommentID, err)
		}
	}

	return replyCount, adminPfp, nil
}

// Helper function to scan for admin pfp to reduce code duplication
func scanForAdminPfp(ctx context.Context, server *Server, parentCommentID pgtype.UUID) string {
	// Fetch all replies (using a reasonable high limit)
	allReplies, err := server.store.ListCommentReplies(ctx, db.ListCommentRepliesParams{
		ParentCommentID: parentCommentID,
		Limit:           10000, // Reasonable high limit
	})
	if err != nil {
		log.Printf("Error fetching all replies for ParentCommentID: %s: %v", parentCommentID, err)
		return ""
	}

	// Look for admin replies
	adminPfpCacheKey := redis.GenerateKey("comments_admin_pfp", parentCommentID)
	for _, reply := range allReplies {
		if reply.Role == "admin" {
			adminPfp := reply.Pfp

			// Cache the admin pfp
			err := server.cacheService.Set(ctx, adminPfpCacheKey, adminPfp, 10*time.Minute)
			if err != nil {
				log.Printf("Error caching admin pfp for ParentCommentID: %s: %v", parentCommentID, err)
			}
			return adminPfp
		}
	}

	return ""
}

func (server *Server) getUserReactionsForContentWithCache(ctx context.Context, contentID, userID pgtype.UUID) (map[string]string, error) {
	cacheKey := redis.GenerateKey("comments_user_reactions", contentID, userID)

	var userReactions map[string]string
	cacheHit, err := server.cacheService.Get(ctx, cacheKey, &userReactions)
	if err != nil {
		log.Printf("Error fetching user reactions from cache: %v", err)
	}
	if cacheHit {
		return userReactions, nil
	}

	log.Printf("Cache miss for user reactions: %s", cacheKey)
	arg := db.GetUserReactionsForContentCommentsParams{
		ContentID: contentID,
		UserID:    userID,
	}
	reactions, err := server.store.GetUserReactionsForContentComments(ctx, arg)
	if err != nil {
		return nil, err
	}

	userReactions = make(map[string]string)
	for _, reaction := range reactions {
		userReactions[reaction.CommentID.String()] = reaction.Reaction
	}

	err = server.cacheService.Set(ctx, cacheKey, userReactions, 10*time.Minute)
	if err != nil {
		log.Printf("Error caching user reactions: %v", err)
	}

	return userReactions, nil
}

func (server *Server) getCommentCountWithCache(ctx context.Context, contentID pgtype.UUID) (int64, error) {
	cacheKey := redis.GenerateKey("comments_count", contentID)

	var commentCount int64
	cacheHit, err := server.cacheService.Get(ctx, cacheKey, &commentCount)
	if err != nil {
		log.Printf("Error fetching comment count from cache: %v", err)
	}
	if cacheHit {
		return commentCount, nil
	}

	log.Printf("Cache miss for comment count: %s", cacheKey)
	count, err := server.store.GetCommentCountForContent(ctx, contentID)
	if err != nil {
		return 0, err
	}

	err = server.cacheService.Set(ctx, cacheKey, count, 10*time.Minute)
	if err != nil {
		log.Printf("Error caching comment count: %v", err)
	}

	return count, nil
}

func (server *Server) listCommentRepliesWithCache(ctx context.Context, parentCommentID pgtype.UUID, limit int32) ([]db.ListCommentRepliesRow, error) {
	cacheKey := redis.GenerateKey("comments_replies", parentCommentID, limit)

	var replies []db.ListCommentRepliesRow
	cacheHit, err := server.cacheService.Get(ctx, cacheKey, &replies)
	if err != nil {
		log.Printf("Error fetching comment replies from cache: %v", err)
	}
	if cacheHit {
		return replies, nil
	}

	log.Printf("Cache miss for listing comment replies: %s", cacheKey)
	replies, err = server.store.ListCommentReplies(ctx, db.ListCommentRepliesParams{
		ParentCommentID: parentCommentID,
		Limit:           limit,
	})
	if err != nil {
		return nil, err
	}

	err = server.cacheService.Set(ctx, cacheKey, replies, 10*time.Minute)
	if err != nil {
		log.Printf("Error caching comment replies: %v", err)
	}

	return replies, nil
}
