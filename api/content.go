package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
)

var ThumbnailURL = "/static/assets/f4a004a8-8612-4152-b96e-2212646d7bdf-article-placeholder.webp"

func (server *Server) listPubContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listPubContent:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in listPubContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listPubContentOldest(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listPubContentOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimitOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in listPubContentOldest:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/oldest?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listPubContentTitle(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listPubContentTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimitTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in listPubContentTitle:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDraftContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDraftContent:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContent(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing draft content in listDraftContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDraftContentOldest(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDraftContentOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContentOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing draft content in listDraftContentOldest:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft/oldest?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDraftContentTitle(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDraftContentTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContentTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing draft content in listDraftContentTitle:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDelContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDelContent:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContent(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing deleted content in listDelContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDelContentOldest(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDelContentOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContentOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing deleted content in listDelContentOldest:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted/oldest?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDelContentTitle(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDelContentTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContentTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing deleted content in listDelContentTitle:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

type SearchContentReq struct {
	SearchTerm string `query:"search_term" validate:"required,min=3"`
	Limit      int32  `query:"limit"`
}

func (server *Server) listSearchPubContent(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listSearchPubContent:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in listSearchPubContent:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	data, err := server.store.SearchContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching content in listSearchPubContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listSearchDraftContent(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listSearchDraftContent:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in listSearchDraftContent:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchDraftContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	data, err := server.store.SearchDraftContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching content in listSearchDraftContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listSearchDelContent(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listSearchDelContent:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in listSearchDelContent:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchDelContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	data, err := server.store.SearchDelContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching content in listSearchDelContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) archivePubContent(ctx echo.Context) error {
	id := ctx.Param("id")
	pgUUID, err := utils.ParseUUID(id, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in archivePubContent:", err)
		return err
	}

	_, err = server.store.SoftDeleteContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error archiving content in archivePubContent:", err)
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting content overview in archivePubContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

func (server *Server) deleteContent(ctx echo.Context) error {
	id := ctx.Param("id")
	pgUUID, err := utils.ParseUUID(id, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in deleteContent:", err)
		return err
	}

	// delete media files associated with content
	media, err := server.store.ListMediaForContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error listing media while deleting content:", err)
		return err
	}

	if len(media) > 0 {
		for _, v := range media {
			// Remove the file from filesystem
			// The filepath is stored with leading slash, so trim it for filesystem operations
			filePath := strings.TrimPrefix(v.MediaUrl, "/")
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error removing file from filesystem at %s: %v", filePath, err)
				// Consider whether to return this error or continue
			}
		}
	}

	_, err = server.store.HardDeleteContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error deleting content in deleteContent:", err)
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting content overview in deleteContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

func (server *Server) publishDraftContent(ctx echo.Context) error {
	id := ctx.Param("id")
	pgUUID, err := utils.ParseUUID(id, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in publishDraftContent:", err)
		return err
	}

	_, err = server.store.PublishContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error publishing content in publishDraftContent:", err)
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting content overview in publishDraftContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

func (server *Server) unarchiveContent(ctx echo.Context) error {
	id := ctx.Param("id")
	pgUUID, err := utils.ParseUUID(id, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in unarchiveContent:", err)
		return err
	}

	_, err = server.store.UnarchiveContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error unarchiving content in unarchiveContent:", err)
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting content overview in unarchiveContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

type UpdateContentReq struct {
	ContentID           string  `query:"content_id"`
	Title               *string `form:"title" validate:"required"`
	ContentDescription  *string `form:"content_description" validate:"required"`
	CategoryID          *string `form:"category_id"`
	CommentsEnabled     *bool   `form:"comments_enabled"`
	ViewCountEnabled    *bool   `form:"view_count_enabled"`
	LikeCountEnabled    *bool   `form:"like_count_enabled"`
	DislikeCountEnabled *bool   `form:"dislike_count_enabled"`
}

func (server *Server) updateContent(ctx echo.Context) error {
	contentIDString := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDString, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in updateContent:", err)
		return err
	}

	var req UpdateContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updateContent:", err)
		return err
	}

	content, err := server.store.GetContentDetails(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting content details in updateContent:", err)
		return err
	}

	var parsedCategoryID uuid.UUID
	var categoryID pgtype.UUID
	if req.CategoryID != nil {
		parsedCategoryID, err = uuid.Parse(*req.CategoryID)
		if err != nil {
			log.Println("Invalid category ID format in updateContent:", err)
			return err
		}
		categoryID = pgtype.UUID{
			Bytes: parsedCategoryID,
			Valid: true,
		}
	} else {
		categoryID = content.CategoryID
	}

	arg := db.UpdateContentParams{
		ContentID: contentID,
		Title: func() string {
			if req.Title != nil {
				return *req.Title
			}

			return content.Title
		}(),
		ContentDescription: func() string {
			if req.ContentDescription != nil {
				return *req.ContentDescription
			}
			return content.ContentDescription
		}(),
		CategoryID: func() pgtype.UUID {
			if req.CategoryID != nil {
				return categoryID
			}
			return content.CategoryID
		}(),
		CommentsEnabled: func() bool {
			if req.CommentsEnabled != nil {
				return *req.CommentsEnabled
			}
			return content.CommentsEnabled
		}(),
		ViewCountEnabled: func() bool {
			if req.ViewCountEnabled != nil {
				return *req.ViewCountEnabled
			}
			return content.ViewCountEnabled
		}(),
		LikeCountEnabled: func() bool {
			if req.LikeCountEnabled != nil {
				return *req.LikeCountEnabled
			}
			return content.LikeCountEnabled
		}(),
		DislikeCountEnabled: func() bool {
			if req.DislikeCountEnabled != nil {
				return *req.DislikeCountEnabled
			}
			return content.DislikeCountEnabled
		}(),
	}

	if req.Title != nil && *req.Title == "" {
		message := "Naslov je obavezan."

		log.Println(message)
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.CategoryID != nil && *req.CategoryID == "" {
		message := "Kategorija je obavezna."

		log.Println(message)
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.ContentDescription != nil && *req.ContentDescription == "" {
		message := "Sadržaj je obavezan."

		log.Println(message)
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	_, err = server.store.UpdateContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating content in updateContent:", err)
		return err
	}

	message := "Sadržaj uspešno ažuriran."
	return Render(ctx, http.StatusOK, components.ArticleSuccess(message))
}

type CreateContentReq struct {
	CategoryID         string `form:"category_id"`
	Title              string `form:"title"`
	ContentDescription string `form:"content_description"`
}

func (server *Server) createContent(ctx echo.Context) error {
	var req CreateContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createContent:", err)
		return err
	}

	if req.Title == "" {
		message := "Naslov je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.CategoryID == "" {
		message := "Kategorija je obavezna."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.ContentDescription == "" {
		message := "Sadržaj je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "access_token")
	if err != nil {
		log.Println("Error getting user data in createContent:", err)
	}

	categoryID, err := utils.ParseUUID(req.CategoryID, "category ID")
	if err != nil {
		log.Println("Invalid category ID format in createContent:", err)
		return err
	}

	arg := db.CreateContentParams{
		UserID:              userData.UserID,
		CategoryID:          categoryID,
		Title:               req.Title,
		ContentDescription:  req.ContentDescription,
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := server.store.CreateContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Failed to create content"

		return Render(ctx, http.StatusInternalServerError, components.ArticleError(message))
	}

	ctx.SetCookie(&http.Cookie{
		Name:    "content_id",
		Value:   content.ContentID.String(),
		MaxAge:  60 * 60 * 24 * 365,
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})

	message := "Uspešno ste sačuvali novi sadržaj."

	return Render(ctx, http.StatusOK, components.ArticleSuccess(message))
}

func (server *Server) createAndPublishContent(ctx echo.Context) error {
	var req CreateContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createAndPublishContent:", err)
		return err
	}

	if req.Title == "" {
		message := "Naslov je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.CategoryID == "" {
		message := "Kategorija je obavezna."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.ContentDescription == "" {
		message := "Sadržaj je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "access_token")
	if err != nil {
		log.Println("Error getting user data in createAndPublishContent:", err)
	}

	categoryID, err := utils.ParseUUID(req.CategoryID, "category ID")
	if err != nil {
		log.Println("Invalid category ID format in createAndPublishContent:", err)
		return err
	}

	arg := db.CreateContentParams{
		UserID:              userData.UserID,
		CategoryID:          categoryID,
		Title:               req.Title,
		ContentDescription:  req.ContentDescription,
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := server.store.CreateContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom cuvanja sadržaja."

		return Render(ctx, http.StatusInternalServerError, components.ArticleError(message))
	}

	_, err = server.store.PublishContent(ctx.Request().Context(), content.ContentID)
	if err != nil {
		message := "Greška prilikom objavljivanja sadržaja."

		return Render(ctx, http.StatusInternalServerError, components.ArticleError(message))
	}

	ctx.SetCookie(&http.Cookie{
		Name:    "content_id",
		Value:   content.ContentID.String(),
		MaxAge:  60 * 60 * 24 * 365,
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})

	message := "Uspešno ste sačuvali i objavili novi sadržaj."

	return Render(ctx, http.StatusOK, components.ArticleSuccess(message))
}

func (server *Server) loadMoreSearch(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in loadMoreSearch:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in loadMoreSearch:", err)
		return ctx.NoContent(http.StatusNoContent)
	}

	nextLimit := req.Limit + 20

	arg := db.SearchContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	searchCount, err := server.store.GetSearchContentCount(ctx.Request().Context(), req.SearchTerm)
	if err != nil {
		log.Println("Error getting search count in loadMoreSearch:", err)
		return err
	}

	searchResults, err := server.store.SearchContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching content in loadMoreSearch:", err)
		return err
	}

	for i := range searchResults {
		if searchResults[i].Thumbnail.String == "" {
			searchResults[i].Thumbnail = pgtype.Text{String: ThumbnailURL, Valid: true}
		}
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in loadMoreSearch:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.SearchResults(searchResults, searchCount, req.SearchTerm, int(nextLimit), globalSettings[0]))
}

func (server *Server) listOtherContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listPubContent:", err)
		return err
	}

	nextLimit := req.Limit + 5

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in listPubContent:", err)
		return err
	}

	var content []components.ListPublishedContentRes
	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:  v.ContentID.String(),
			UserID:     v.UserID.String(),
			CategoryID: v.CategoryID.String(),
			Title:      v.Title,
			Thumbnail: func() string {
				if v.Thumbnail.Valid && v.Thumbnail.String != "" {
					return v.Thumbnail.String
				}

				return ThumbnailURL
			}(),
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         utils.TimeAgo(v.PublishedAt.Time.In(Loc)),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in listPubContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.OtherContent(content, int(nextLimit), globalSettings[0]))
}

type CategoryWithContent struct {
	Category db.Category
	Content  []db.ListContentByCategoryRow
}

// This function just generates the component and can be reused
func (server *Server) GenerateNewsSliderComponent(ctx context.Context) (templ.Component, error) {
	// Fetch categories (limit to 20 for example)
	categories, err := server.store.ListCategories(ctx, 20)
	if err != nil {
		log.Print("Error fetching categories in newsSlider:", err)
		return nil, err
	}

	// Create a map to store content by category ID
	contentByCategory := make(map[pgtype.UUID][]db.ListContentByCategoryRow)
	// Filtered categories that actually have content
	var filteredCategories []db.Category

	// Fetch content for each category
	for _, category := range categories {
		contentParams := db.ListContentByCategoryParams{
			CategoryID: category.CategoryID,
			Limit:      1, // Limit number of articles per category
			Offset:     0,
		}
		content, err := server.store.ListContentByCategory(ctx, contentParams)
		if err != nil {
			log.Printf("Error fetching content for category %v: %v", category.CategoryName, err)
			continue // Skip this category if content fetch fails
		}
		if len(content) == 0 {
			continue // Skip this category if no content is found
		}
		// Store content in the map using category ID as the key
		contentByCategory[category.CategoryID] = content
		// Add category to filtered list since it has content
		filteredCategories = append(filteredCategories, category)
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx)
	if err != nil {
		log.Println("Error getting global settings in newsSlider:", err)
		return nil, err
	}

	// Return the component
	return components.NewsSlider(filteredCategories, contentByCategory, globalSettings[0]), nil
}

func (server *Server) newsSlider(ctx echo.Context) error {
	component, err := server.GenerateNewsSliderComponent(ctx.Request().Context())
	if err != nil {
		return err
	}

	return Render(ctx, http.StatusOK, component)
}

func (server *Server) categoriesWithContent(ctx echo.Context) error {
	categories, err := server.store.ListCategories(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Error listing categories in categoriesWithContent:", err)
		return err
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(categories), func(i, j int) {
		categories[i], categories[j] = categories[j], categories[i]
	})

	// Pick up to 4 categories (or fewer if there aren't enough)
	limit := 4
	if len(categories) < limit {
		limit = len(categories) // Take whatever is available
	}
	randomCategories := categories[:limit]

	// Prepare the final content structure
	var allCategoryContent []components.ContentDataSlice

	// For each random category, fetch its content
	for _, category := range randomCategories {
		// Fetch content for this category (limit to 20 items per category)
		contentItems, err := server.store.ListContentByCategoryLimit(ctx.Request().Context(), db.ListContentByCategoryLimitParams{
			CategoryID: category.CategoryID,
			Limit:      20,
		})
		if err != nil {
			log.Printf("Error fetching content for category %s: %v", category.CategoryName, err)
			continue // Skip this category if there's an error
		}

		// Skip empty results
		if len(contentItems) == 0 {
			continue
		}

		// Convert DB content items to ContentData
		var categoryContent []components.ContentData
		for _, item := range contentItems {
			categoryContent = append(categoryContent, components.ContentData{
				ContentID:    item.ContentID,
				UserID:       item.UserID,
				CategoryID:   item.CategoryID,
				CategoryName: category.CategoryName, // Use the category name from the category object
				Title:        item.Title,
				Thumbnail: func() pgtype.Text {
					if item.Thumbnail.Valid && item.Thumbnail.String != "" {
						return item.Thumbnail
					}

					return pgtype.Text{String: ThumbnailURL, Valid: true}
				}(),
				ContentDescription:  item.ContentDescription,
				CommentsEnabled:     item.CommentsEnabled,
				ViewCountEnabled:    item.ViewCountEnabled,
				LikeCountEnabled:    item.LikeCountEnabled,
				DislikeCountEnabled: item.DislikeCountEnabled,
				Status:              item.Status,
				ViewCount:           item.ViewCount,
				LikeCount:           item.LikeCount,
				DislikeCount:        item.DislikeCount,
				CommentCount:        item.CommentCount,
				CreatedAt:           item.CreatedAt,
				UpdatedAt:           item.UpdatedAt,
				PublishedAt:         item.PublishedAt,
				IsDeleted:           item.IsDeleted,
			})
		}

		// Add this category's content to the final result
		if len(categoryContent) > 0 {
			allCategoryContent = append(allCategoryContent, components.ContentDataSlice{
				Content: categoryContent,
			})
		}
	}

	// If we don't have any content, return an empty page
	if len(allCategoryContent) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	// Get global settings
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in categoriesWithContent:", err)
		return err
	}

	// Render the template with all category content
	return Render(ctx, http.StatusOK, components.CategoriesWithContent(allCategoryContent, globalSettings[0]))
}

func (server *Server) handleLikeContent(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in handleLikeContent:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user data in handleLikeContent:", err)
	}

	// Check if the user already has a reaction using the efficient query
	userReaction, err := server.store.FetchUserContentReaction(ctx.Request().Context(), db.FetchUserContentReactionParams{
		ContentID: contentID,
		UserID:    userData.UserID,
	})

	// Handle reaction logic based on whether we found a reaction and what it was
	if err == nil {
		// User has an existing reaction
		if userReaction.Reaction == "like" {
			// If already liked, remove the reaction
			_, err = server.store.DeleteContentReaction(ctx.Request().Context(), db.DeleteContentReactionParams{
				ContentID: contentID,
				UserID:    userData.UserID,
			})
			if err != nil {
				log.Println("Error deleting content reaction from like to remove like:", err)
				return err
			}
			// Decrement daily likes since user removed their like
			err = server.decrementDailyLikes(ctx)
			if err != nil {
				log.Println("Error decrementing daily likes:", err)
				return err
			}
		} else if userReaction.Reaction == "dislike" {
			// If disliked, change to like
			_, err := server.store.InsertOrUpdateContentReaction(ctx.Request().Context(), db.InsertOrUpdateContentReactionParams{
				ContentID: contentID,
				UserID:    userData.UserID,
				Reaction:  "like",
			})
			if err != nil {
				log.Println("Error changing reaction from dislike to like:", err)
				return err
			}
			// Increment daily likes since user changed from dislike to like
			err = server.incrementDailyLikes(ctx)
			if err != nil {
				log.Println("Error incrementing daily likes:", err)
				return err
			}
		}
	} else {
		// No reaction yet, add a like
		_, err := server.store.InsertOrUpdateContentReaction(ctx.Request().Context(), db.InsertOrUpdateContentReactionParams{
			ContentID: contentID,
			UserID:    userData.UserID,
			Reaction:  "like",
		})
		if err != nil {
			log.Println("Error adding new like reaction:", err)
			return err
		}
		// Increment daily likes since user is adding a like for the first time
		err = server.incrementDailyLikes(ctx)
		if err != nil {
			log.Println("Error incrementing daily likes:", err)
			return err
		}
	}

	// Update the content's like/dislike counts
	_, err = server.store.UpdateContentLikeDislikeCount(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error updating content like and dislike count:", err)
		return err
	}

	// Get the updated user reaction for the response
	updatedUserReaction, err := server.store.FetchUserContentReaction(ctx.Request().Context(), db.FetchUserContentReactionParams{
		ContentID: contentID,
		UserID:    userData.UserID,
	})

	reactionStatus := ""
	if err == nil {
		reactionStatus = updatedUserReaction.Reaction
	}

	// Get global settings for rendering
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings:", err)
		return err
	}

	// Get updated content details for rendering
	content, err := server.store.GetContentDetails(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting content details:", err)
		return err
	}

	// Render the updated article stats component
	return Render(ctx, http.StatusOK, components.ArticleStats(content, globalSettings[0], reactionStatus))
}

func (server *Server) handleDislikeContent(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in handleDislikeContent:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user data in handleDislikeContent:", err)
	}

	// Check if the user already has a reaction using the efficient query
	userReaction, err := server.store.FetchUserContentReaction(ctx.Request().Context(), db.FetchUserContentReactionParams{
		ContentID: contentID,
		UserID:    userData.UserID,
	})

	// Handle reaction logic based on whether we found a reaction and what it was
	if err == nil {
		// User has an existing reaction
		if userReaction.Reaction == "dislike" {
			// If already disliked, remove the reaction
			_, err = server.store.DeleteContentReaction(ctx.Request().Context(), db.DeleteContentReactionParams{
				ContentID: contentID,
				UserID:    userData.UserID,
			})
			if err != nil {
				log.Println("Error deleting content reaction from dislike to remove dislike:", err)
				return err
			}
			// Decrement daily dislikes since user removed their dislike
			err = server.decrementDailyDislikes(ctx)
			if err != nil {
				log.Println("Error decrementing daily dislikes:", err)
				return err
			}
		} else if userReaction.Reaction == "like" {
			// If liked, change to dislike
			_, err := server.store.InsertOrUpdateContentReaction(ctx.Request().Context(), db.InsertOrUpdateContentReactionParams{
				ContentID: contentID,
				UserID:    userData.UserID,
				Reaction:  "dislike",
			})
			if err != nil {
				log.Println("Error changing reaction from like to dislike:", err)
				return err
			}
			// Increment daily dislikes since user changed from like to dislike
			err = server.incrementDailyDislikes(ctx)
			if err != nil {
				log.Println("Error incrementing daily dislikes:", err)
				return err
			}
		}
	} else {
		// No reaction yet, add a dislike
		_, err := server.store.InsertOrUpdateContentReaction(ctx.Request().Context(), db.InsertOrUpdateContentReactionParams{
			ContentID: contentID,
			UserID:    userData.UserID,
			Reaction:  "dislike",
		})
		if err != nil {
			log.Println("Error adding new dislike reaction:", err)
			return err
		}
		// Increment daily dislikes since user is adding a dislike for the first time
		err = server.incrementDailyDislikes(ctx)
		if err != nil {
			log.Println("Error incrementing daily dislikes:", err)
			return err
		}
	}

	// Update the content's like/dislike counts
	_, err = server.store.UpdateContentLikeDislikeCount(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error updating content like and dislike count:", err)
		return err
	}

	// Get the updated user reaction for the response
	updatedUserReaction, err := server.store.FetchUserContentReaction(ctx.Request().Context(), db.FetchUserContentReactionParams{
		ContentID: contentID,
		UserID:    userData.UserID,
	})

	reactionStatus := ""
	if err == nil {
		reactionStatus = updatedUserReaction.Reaction
	}

	// Get global settings for rendering
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings:", err)
		return err
	}

	// Get updated content details for rendering
	content, err := server.store.GetContentDetails(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting content details:", err)
		return err
	}

	// Render the updated article stats component
	return Render(ctx, http.StatusOK, components.ArticleStats(content, globalSettings[0], reactionStatus))
}
