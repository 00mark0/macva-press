package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

var (
	Loc  *time.Location
	Now  time.Time
	Date time.Time
)

func init() {
	var err error
	Loc, err = time.LoadLocation("Europe/Belgrade")
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}
	Now = time.Now().In(Loc)
	Date = time.Date(Now.Year(), Now.Month(), Now.Day(), 0, 0, 0, 0, Loc)
}

func (server *Server) homePage(ctx echo.Context) error {
	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in homePage:", err)
	}

	// Prepare meta information dynamically
	meta := components.Meta{
		Title:       "Mačva Press | Vaš izvor vesti", // More localized
		Description: "Najnovije vesti i dešavanja iz Mačve i Srbije – budite u toku sa svim bitnim informacijama.",
		Canonical:   BaseUrl, // Update with your actual domain
		OpenGraph: components.OpenGraphMeta{
			Title:       "Mačva Press | Vesti iz srca Mačve i Srbije",
			Description: "Pouzdane, tačne i pravovremene informacije o događajima u Mačvi i regionu.",
			URL:         BaseUrl, // Update with your actual domain
			Type:        "website",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Prepare an Open Graph image
		},
		Twitter: components.TwitterCardMeta{
			Card:        "summary_large_image",
			Title:       "Mačva Press | Prava vest u pravo vreme",
			Description: "Najnovije lokalne i regionalne vesti iz Mačve – obavešteni, povezani, korak ispred.",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Prepare a Twitter card image
			Creator:     "@MacvaNews",                                  // Optional: your Twitter handle
		},
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 11)
	if err != nil {
		log.Println("Error listing active ads in homePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in homePage:", err)
		return err
	}

	// Get the pre-rendered slider component
	prerenderedSlider, err := server.GenerateNewsSliderComponent(ctx.Request().Context())
	if err != nil {
		return err
	}

	// Render the Index template with the pre-rendered slider
	return Render(ctx, http.StatusOK, components.Index(userData, meta, activeAds, categories, prerenderedSlider))
}

// full page to be served
func (server *Server) adminDash(ctx echo.Context) error {
	user, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in adminDash:", err)
	}

	return Render(ctx, http.StatusOK, components.DashPage(user))
}

// htmx content insert
func (server *Server) adminDashContent(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.AdminDashboard())
}

func (server *Server) adminCats(ctx echo.Context) error {
	var req ListCatsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in adminCats:", err)
		return err
	}

	nextLimit := req.Limit + 10

	categories, err := server.store.ListCategories(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing categories in adminCats:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.AdminCategories(int(nextLimit), categories))
}

func (server *Server) createCategoryForm(ctx echo.Context) error {
	var createCategoryErr components.CreateCategoryErr

	return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCategoryErr))
}

func (server *Server) deleteCategoryModal(ctx echo.Context) error {
	categoryIDStr := ctx.Param("id")

	categoryID, err := utils.ParseUUID(categoryIDStr, "category ID")
	if err != nil {
		log.Println("Invalid category ID format in deleteCategoryModal:", err)
		return err
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error getting category in deleteCategoryModal:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.DeleteCategoryModal(category))
}

func (server *Server) updateCategoryForm(ctx echo.Context) error {
	categoryIDStr := ctx.Param("id")
	var updateCategoryErr components.UpdateCategoryErr

	categoryID, err := utils.ParseUUID(categoryIDStr, "category ID")
	if err != nil {
		log.Println("Invalid category ID format in updateCategoryForm:", err)
		return err
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error getting category in updateCategoryForm:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCategoryErr))
}

type ListPublishedLimitReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) adminArts(ctx echo.Context) error {
	var req ListPublishedLimitReq

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting content overview in adminArts:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in adminArts:", err)
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

	return Render(ctx, http.StatusOK, components.AdminArticles(overview, int(nextLimit), content, url))
}

func (server *Server) publishedContentList(ctx echo.Context) error {
	var req ListPublishedLimitReq

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in publishedContentList:", err)
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

	return Render(ctx, http.StatusOK, components.PublishedContentSort(int(nextLimit), content, url))
}

func (server *Server) draftContentList(ctx echo.Context) error {
	var req ListPublishedLimitReq

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContent(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing draft content in draftContentList:", err)
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

	return Render(ctx, http.StatusOK, components.DraftContentSort(int(nextLimit), content, url))
}

func (server *Server) deletedContentList(ctx echo.Context) error {
	var req ListPublishedLimitReq

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContent(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing deleted content in deletedContentList:", err)
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

	return Render(ctx, http.StatusOK, components.DeletedContentSort(int(nextLimit), content, url))
}

type ListUsersLimitReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) adminUsers(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20
	activeCount, err := server.store.GetActiveUsersCount(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting active users count in adminUsers:", err)
		return err
	}

	bannedCount, err := server.store.GetBannedUsersCount(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting banned users count in adminUsers:", err)
		return err
	}

	delCount, err := server.store.GetDeletedUsersCount(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting deleted users count in adminUsers:", err)
		return err
	}

	overview := components.UsersOverview{
		ActiveUsersCount:  int(activeCount),
		BannedUsersCount:  int(bannedCount),
		DeletedUsersCount: int(delCount),
	}

	var users []components.UsersRes

	data, err := server.store.GetActiveUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in adminUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active?limit="

	return Render(ctx, http.StatusOK, components.AdminUsers(overview, int(nextLimit), users, url))
}

func (server *Server) activeUsersList(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in activeUsersList:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active?limit="

	return Render(ctx, http.StatusOK, components.ActiveUsersSort(int(nextLimit), users, url))
}

func (server *Server) bannedUsersList(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in bannedUsersList:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/banned?limit="

	return Render(ctx, http.StatusOK, components.BannedUsersSort(int(nextLimit), users, url))
}

func (server *Server) deletedUsersList(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in deletedUsersList:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/deleted?limit="

	return Render(ctx, http.StatusOK, components.DelUsersSort(int(nextLimit), users, url))
}

type ListAdsReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) adminAds(ctx echo.Context) error {
	var req ListAdsReq

	nextLimit := req.Limit + 20

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing active ads in adminAds:", err)
		return err
	}

	url := "/api/admin/ads/active?limit="

	return Render(ctx, http.StatusOK, components.AdminAds(int(nextLimit), activeAds, url))
}

func (server *Server) activeAdsList(ctx echo.Context) error {
	var req ListAdsReq

	nextLimit := req.Limit + 20

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing active ads in activeAdsList:", err)
		return err
	}

	url := "/api/admin/ads/active?limit="

	return Render(ctx, http.StatusOK, components.ActiveAdsSort(int(nextLimit), activeAds, url))
}

func (server *Server) inactiveAdsList(ctx echo.Context) error {
	var req ListAdsReq

	nextLimit := req.Limit + 20

	inactiveAds, err := server.store.ListInactiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing inactive ads in inactiveAdsList:", err)
		return err
	}

	url := "/api/admin/ads/inactive?limit="

	return Render(ctx, http.StatusOK, components.InactiveAdsSort(int(nextLimit), inactiveAds, url))
}

func (server *Server) scheduledAdsList(ctx echo.Context) error {
	var req ListAdsReq

	nextLimit := req.Limit + 20

	scheduledAds, err := server.store.ListScheduledAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing scheduled ads in scheduledAdsList:", err)
		return err
	}

	url := "/api/admin/ads/scheduled?limit="

	return Render(ctx, http.StatusOK, components.ScheduledAdsSort(int(nextLimit), scheduledAds, url))
}

func (server *Server) createAdModal(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.CreateAdModal(""))
}

func (server *Server) updateAdModal(ctx echo.Context) error {
	adIDStr := ctx.Param("id")

	adID, err := utils.ParseUUID(adIDStr, "ad ID")
	if err != nil {
		log.Println("Invalid ad ID format in updateAdModal:", err)
		return err
	}

	ad, err := server.store.GetAd(ctx.Request().Context(), adID)
	if err != nil {
		log.Println("Error getting ad in updateAdModal:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UpdateAdModal("", ad))
}

func (server *Server) loginPage(ctx echo.Context) error {
	var loginErr components.LoginErr

	return Render(ctx, http.StatusOK, components.Login(loginErr))
}

func (server *Server) createArticlePage(ctx echo.Context) error {
	categories, err := server.store.ListCategories(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Failed to get create article page in createArticlePage:", err)
		return err
	}

	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Failed to get tags for create article page in createArticlePage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.CreateArticle(categories, tags))
}

func (server *Server) updateArticlePage(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")

	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in updateArticlePage:", err)
		return err
	}

	content, err := server.store.GetContentDetails(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Failed to get content for update article page in updateArticlePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Failed to get update article page in updateArticlePage:", err)
		return err
	}

	media, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Failed to get media for update article page in updateArticlePage:", err)
		return err
	}

	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Failed to get tags for update article page in updateArticlePage:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Failed to get tags for update article page in updateArticlePage:", err)
		return err
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "content_id",
		Value:    contentIDStr,
		Path:     "/",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return Render(ctx, http.StatusOK, components.UpdateArticle(content, categories, media, tags, contentTags))
}

func (server *Server) adminSettings(ctx echo.Context) error {
	// Get global settings or create if they don't exist
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil || len(globalSettings) == 0 {
		// If there's an error or no settings exist, create new settings
		newSettings, err := server.store.CreateGlobalSettings(ctx.Request().Context())
		if err != nil {
			log.Println("Error creating global settings in adminSettings:", err)
			return err
		}
		globalSettings = []db.GlobalSetting{newSettings}
	}

	user, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in homePage:", err)
	}

	// Create props for the AdminSettings component
	props := components.AdminSettingsProps{
		// User settings from auth payload
		UserID:   user.UserID.String(),
		Username: user.Username,
		Pfp:      user.Pfp,

		// Global settings from the first record
		DisableComments: globalSettings[0].DisableComments,
		DisableLikes:    globalSettings[0].DisableLikes,
		DisableDislikes: globalSettings[0].DisableDislikes,
		DisableViews:    globalSettings[0].DisableViews,
		DisableAds:      globalSettings[0].DisableAds,
	}

	// Render the AdminSettings component with the props
	return Render(ctx, http.StatusOK, components.AdminSettings(props))
}

func (server *Server) passwordResetPage(ctx echo.Context) error {
	token := ctx.Param("token")

	// Validate the token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return Render(ctx, http.StatusOK, components.PasswordReset("", "Link za resetovanje lozinke je nevažeći."))
	}

	// Verify that user_id exists in the claims
	if _, exists := claims["user_id"]; !exists {
		return Render(ctx, http.StatusOK, components.PasswordReset("", "Link za resetovanje lozinke je nevažeći."))
	}

	// Token is valid, show the password reset form
	return Render(ctx, http.StatusOK, components.PasswordReset(token, ""))
}

func (server *Server) registerPage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.RegisterPage(""))
}

func (server *Server) emailVerifiedPage(ctx echo.Context) error {
	token := ctx.Param("token")

	// Validate the token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		log.Println("Error validating token in emailVerifiedPage:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	// Verify that user_id exists in the claims
	if _, exists := claims["user_id"]; !exists {
		log.Println("Error extracting user_id from claims in emailVerifiedPage:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		log.Println("Error extracting user_id from claims in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	// Parse user ID to UUID
	var userID pgtype.UUID
	err = userID.Scan(userIDStr)
	if err != nil {
		log.Println("Error parsing user_id in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	err = server.store.SetEmailVerified(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error setting email_verified in emailVerifiedPage:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	return Render(ctx, http.StatusOK, components.VerificationSuccess())
}

func (server *Server) requestPassResetPage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.RequestPassReset())
}

func (server *Server) searchResultsPage(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchResultsPage:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchResultsPage:", err)
		return ctx.NoContent(http.StatusNoContent)
	}

	nextLimit := req.Limit + 20
	searchTerm := req.SearchTerm

	arg := db.SearchContentParams{
		Limit:      nextLimit,
		SearchTerm: searchTerm,
	}

	searchResults, err := server.store.SearchContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching content in searchResultsPage:", err)
		return err
	}

	for i := range searchResults {
		if searchResults[i].Thumbnail.String == "" {
			searchResults[i].Thumbnail = pgtype.Text{String: ThumbnailURL, Valid: true}
		}
	}

	searchResultsCount, err := server.store.GetSearchContentCount(ctx.Request().Context(), req.SearchTerm)
	if err != nil {
		log.Println("Error counting search results in searchResultsPage:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in homePage:", err)
	}

	// Prepare meta information dynamically for the search page
	meta := components.Meta{
		Title:       "Mačva Press | Pretraga", // Updated for the search page
		Description: "Pretražite najnovije vesti i dešavanja iz Mačve i Srbije – brzo i jednostavno.",
		Canonical:   BaseUrl + "/pretraga", // Updated for the search page URL
		OpenGraph: components.OpenGraphMeta{
			Title:       "Mačva Press | Pretraga vesti iz Mačve i Srbije",
			Description: "Pronađite relevantne vesti iz Mačve i Srbije pomoću naše pretrage.",
			URL:         BaseUrl + "/pretraga", // Updated for the search page URL
			Type:        "website",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Use the same image
		},
		Twitter: components.TwitterCardMeta{
			Card:        "summary_large_image",
			Title:       "Mačva Press | Pretraga vesti u Mačvi i Srbiji",
			Description: "Pretražujte najnovije vesti i informacije iz Mačve sa jednostavnim pretraživačem.",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Use the same image
			Creator:     "@MacvaNews",                                  // Optional: your Twitter handle
		},
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 11)
	if err != nil {
		log.Println("Error listing active ads in homePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in homePage:", err)
		return err
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in homePage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.SearchPage(userData, meta, activeAds, categories, searchResults, searchResultsCount, searchTerm, int(nextLimit), globalSettings[0]))
}

func (server *Server) categoriesPage(ctx echo.Context) error {
	categorySlug := ctx.Param("slug")

	category, err := server.store.GetCategoryBySlug(ctx.Request().Context(), categorySlug)
	if err != nil {
		log.Println("Error getting category in categoriesPage:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in homePage:", err)
	}

	// Prepare meta information dynamically for the search page
	meta := components.Meta{
		Title:       "Mačva Press | " + category.CategoryName, // Već promenjeno za stranicu kategorija
		Description: "Istražite vesti po kategorijama i saznajte najnovija dešavanja iz Mačve i Srbije.",
		Canonical:   BaseUrl + "/kategorije/" + utils.Slugify(category.CategoryName), // Ažurirano za URL stranice kategorija
		OpenGraph: components.OpenGraphMeta{
			Title:       "Mačva Press | " + category.CategoryName,
			Description: "Pregledajte vesti iz različitih kategorija i pratite najvažnije teme iz Mačve i Srbije.",
			URL:         BaseUrl + "/kategorije/" + utils.Slugify(category.CategoryName), // Ažurirano za URL stranice kategorija
			Type:        "website",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Koristi istu sliku
		},
		Twitter: components.TwitterCardMeta{
			Card:        "summary_large_image",
			Title:       "Mačva Press | " + category.CategoryName,
			Description: "Pronađite najnovije vesti razvrstane po kategorijama i budite u toku sa aktuelnim dešavanjima.",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Koristi istu sliku
			Creator:     "@MacvaNews",                                  // Opcionalno: vaš Twitter nalog
		},
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 11)
	if err != nil {
		log.Println("Error listing active ads in homePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in homePage:", err)
		return err
	}

	recentCatComponent, err := server.GenerateRecentCatContentComponent(ctx)
	if err != nil {
		log.Println("Error generating recent category component in categoriesPage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.CategoriesPage(userData, meta, activeAds, categories, category, recentCatComponent))
}

func (server *Server) tagPage(ctx echo.Context) error {
	tagSlug := ctx.Param("slug")

	tag, err := server.store.GetTagBySlug(ctx.Request().Context(), tagSlug)
	if err != nil {
		log.Println("Error getting category in categoriesPage:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in homePage:", err)
	}

	// Prepare meta information dynamically for the search page
	meta := components.Meta{
		Title:       "Mačva Press | " + tag.TagName, // Već promenjeno za stranicu kategorija
		Description: "Istražite vesti po oznakama i saznajte najnovija dešavanja iz Mačve i Srbije.",
		Canonical:   BaseUrl + "/oznake/" + tag.Slug, // Ažurirano za URL stranice kategorija
		OpenGraph: components.OpenGraphMeta{
			Title:       "Mačva Press | " + tag.TagName,
			Description: "Pregledajte vesti iz različitih oznaka i pratite najvažnije teme iz Mačve i Srbije.",
			URL:         BaseUrl + "/oznake/" + tag.Slug, // Ažurirano za URL stranice kategorija
			Type:        "website",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Koristi istu sliku
		},
		Twitter: components.TwitterCardMeta{
			Card:        "summary_large_image",
			Title:       "Mačva Press | " + tag.TagName,
			Description: "Pronađite najnovije vesti razvrstane po oznakama i budite u toku sa aktuelnim dešavanjima.",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Koristi istu sliku
			Creator:     "@MacvaNews",                                  // Opcionalno: vaš Twitter nalog
		},
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 11)
	if err != nil {
		log.Println("Error listing active ads in homePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in homePage:", err)
		return err
	}

	recentTagsComponent, err := server.GenerateRecentTagContentComponent(ctx)
	if err != nil {
		log.Println("Error generating recent tags component in tagsPage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsPage(userData, meta, activeAds, categories, tag, recentTagsComponent))
}

func getOrCreateAnonID(c echo.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie("anon_token")
	if err == nil {
		anonID, err := uuid.Parse(cookie.Value)
		if err == nil {
			return anonID, nil
		}
	}
	// Create new anon_token
	newAnonID := uuid.New()
	newCookie := &http.Cookie{
		Name:     "anon_token",
		Value:    newAnonID.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(20 * 365 * 24 * time.Hour), // 20 years
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	c.SetCookie(newCookie)
	return newAnonID, nil
}

func (server *Server) handleViews(ctx echo.Context, contentIDStr, userIDStr string) {
	contentID, err := utils.ParseUUID(contentIDStr, "contentID")
	if err != nil {
		log.Println("Invalid contentID:", err)
		return
	}
	var viewerID pgtype.UUID
	if userIDStr != "" {
		viewerID, err = utils.ParseUUID(userIDStr, "userID")
		if err != nil {
			log.Println("Invalid userID:", err)
			return
		}
	} else {
		anonID, err := getOrCreateAnonID(ctx)
		if err != nil {
			log.Println("Error getting/creating anontoken:", err)
			return
		}
		viewerID, err = utils.ParseUUID(anonID.String(), "anonID")
		if err != nil {
			log.Println("Invalid anonID:", err)
			return
		}
	}
	params := db.GetViewParams{
		ContentID: contentID,
		UserID:    viewerID,
	}

	_, err = server.store.GetView(ctx.Request().Context(), params)
	if errors.Is(err, pgx.ErrNoRows) {
		if err := server.store.AddView(ctx.Request().Context(), db.AddViewParams{
			ContentID: params.ContentID,
			UserID:    params.UserID,
		}); err != nil {
			log.Println("Error adding view:", err)
			return // or skip daily increment
		}
		_, err = server.store.IncrementViewCount(ctx.Request().Context(), contentID)
		if err != nil {
			log.Println("Error incrementing view count:", err)
		}
		err = server.incrementDailyViews(ctx)
		if err != nil {
			log.Println(err)
		}

	} else if err != nil {
		log.Println("Error checking view:", err)
	}
}

func (server *Server) articlePage(ctx echo.Context) error {
	year := ctx.Param("year")
	month := ctx.Param("month")
	slug := ctx.Param("slug")

	if len(year) != 4 || len(month) != 2 {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid URL format")
	}

	if _, err := strconv.Atoi(year); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid year format")
	}

	m, err := strconv.Atoi(month)
	if err != nil || m < 1 || m > 12 {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid month format")
	}

	// Now slug can be any string, proceed
	article, err := server.store.GetContentBySlug(ctx.Request().Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Article not found")
		}
		return err
	}

	// Validate published_at year/month matches URL (optional sanity check)
	if article.PublishedAt.Time.Format("2006") != year || article.PublishedAt.Time.Format("01") != month {
		return echo.NewHTTPError(http.StatusNotFound, "Article date mismatch")
	}
	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in homePage:", err)
	}
	server.handleViews(ctx, article.ContentID.String(), userData.UserID.String())

	// Prepare meta information dynamically for the search page
	meta := components.Meta{
		Title:       utils.GenerateTitleTag(article.Title), // Već promenjeno za stranicu kategorija
		Description: utils.GenerateMetaDescription(article.ContentDescription),
		Canonical:   BaseUrl + "/" + utils.PrettyURL(article.Slug, article.PublishedAt.Time), // Ažurirano za URL stranice kategorija
		OpenGraph: components.OpenGraphMeta{
			Title:       utils.GenerateTitleTag(article.Title),
			Description: utils.GenerateMetaDescription(article.ContentDescription),
			URL:         BaseUrl + "/" + utils.PrettyURL(article.Slug, article.PublishedAt.Time), // Ažurirano za URL stranice kategorija
			Type:        "website",
			Image:       BaseUrl + article.Thumbnail.String, // Koristi istu sliku
		},
		Twitter: components.TwitterCardMeta{
			Card:        "summary_large_image",
			Title:       utils.GenerateTitleTag(article.Title),
			Description: utils.GenerateMetaDescription(article.ContentDescription),
			Image:       BaseUrl + article.Thumbnail.String, // Koristi istu sliku
			Creator:     "@MacvaNews",                       // Opcionalno: vaš Twitter nalog
		},
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 11)
	if err != nil {
		log.Println("Error listing active ads in homePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in homePage:", err)
		return err
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in homePage:", err)
		return err
	}

	userReaction := ""
	if userData.UserID.Valid {
		reaction, err := server.store.FetchUserContentReaction(ctx.Request().Context(), db.FetchUserContentReactionParams{
			ContentID: article.ContentID,
			UserID:    userData.UserID,
		})
		if err != nil {
			log.Println("No reaction for user in articlePage:", err)
		} else {
			userReaction = reaction.Reaction
		}
	}

	prerenderedArticleMediaSliderComponent, err := server.GenerateArticleMediaSliderComponent(ctx)
	if err != nil {
		log.Println("Error generating prerendered article media slider component in articlePage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticlePage(userData, meta, activeAds, categories, article, globalSettings[0], userReaction, activeAds, meta.Canonical, prerenderedArticleMediaSliderComponent))
}

func (server *Server) userSettingsPage(ctx echo.Context) error {
	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in userSettingsPage:", err)
	}

	userProps := components.UserSettingsProps{
		UserID:   userData.UserID.String(),
		Username: userData.Username,
		Pfp:      userData.Pfp,
	}

	// Prepare meta information dynamically for the search page
	meta := components.Meta{
		Title:       "Mačva Press | Korisnička Podešavanja", // Već promenjeno za stranicu kategorija
		Description: "Korisnička Podešavanja",
		Canonical:   BaseUrl + "/" + "podesavanja", // Ažurirano za URL stranice kategorija
		OpenGraph: components.OpenGraphMeta{
			Title:       "Mačva Press | Korisnička Podešavanja",
			Description: "Korisnička Podešavanja",
			URL:         BaseUrl + "/" + "podesavanja", // Ažurirano za URL stranice kategorija
			Type:        "website",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Koristi istu sliku
		},
		Twitter: components.TwitterCardMeta{
			Card:        "summary_large_image",
			Title:       "Mačva Press | Korisnička Podešavanja",
			Description: "Korisnička Podešavanja",
			Image:       BaseUrl + "/static/assets/macva-1-300x71.png", // Koristi istu sliku
			Creator:     "@MacvaNews",                                  // Opcionalno: vaš Twitter nalog
		},
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 11)
	if err != nil {
		log.Println("Error listing active ads in homePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in homePage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UserSettingsPage(userData, meta, activeAds, categories, userProps))
}
