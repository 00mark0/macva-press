package api

import (
	"log"
	"os"

	"github.com/00mark0/macva-press/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (server *Server) setupRouter() {
	router := echo.New()

	router.Use(middleware.Gzip())

	// Create rate limiters for different types of operations
	// For authentication - very strict limits to prevent brute force attacks
	authLimiter, err := CreateRateLimiter("10-M") // 10 requests per minute
	if err != nil {
		log.Fatal("Failed to create auth rate limiter:", err)
	}

	// For comment operations - moderate limits to prevent comment spam
	commentLimiter, err := CreateRateLimiter("30-M") // 30 requests per minute
	if err != nil {
		log.Fatal("Failed to create comment rate limiter:", err)
	}

	// For content interactions - higher limits for normal user activity
	contentLimiter, err := CreateRateLimiter("60-M") // 60 requests per minute
	if err != nil {
		log.Fatal("Failed to create content rate limiter:", err)
	}

	// For user settings - moderate limits
	userSettingsLimiter, err := CreateRateLimiter("30-M") // 30 requests per minute
	if err != nil {
		log.Fatal("Failed to create user settings rate limiter:", err)
	}

	// For admin operations - higher limits for admin work
	adminLimiter, err := CreateRateLimiter("200-M") // 200 requests per minute
	if err != nil {
		log.Fatal("Failed to create admin rate limiter:", err)
	}

	// For search operations - prevent search abuse
	searchLimiter, err := CreateRateLimiter("60-M") // 60 requests per minute
	if err != nil {
		log.Fatal("Failed to create search rate limiter:", err)
	}

	// Initialize custom validator from validator.go
	router.Validator = NewCustomValidator()

	// Run cron job to create daily analytics
	go server.scheduleDailyAnalytics()

	// Run cron job to deactivate expired ads
	go server.deactivateAds()

	// Serve static files
	router.Static("/static", "static")

	if os.Getenv("DEV_MODE") == "true" {
		router.Use(utils.NoCacheMiddleware)
	}

	// Set authRoutes instead of router to any routes that require middleware
	authRoutes := router.Group("")
	authRoutes.Use(server.authMiddleware(server.tokenMaker))

	adminRoutes := router.Group("")
	adminRoutes.Use(server.adminMiddleware(server.tokenMaker))

	// ==== Page Routes (No rate limiting) ====

	// Admin Page Routes - no rate limiting for page views
	adminRoutes.GET("/admin", server.adminDash)
	adminRoutes.GET("/admin/hx-admin", server.adminDashContent)
	adminRoutes.GET("/admin/categories", server.adminCats)
	adminRoutes.GET("/admin/create-cat-form", server.createCategoryForm)
	adminRoutes.GET("/admin/delete-cat-modal/:id", server.deleteCategoryModal)
	adminRoutes.GET("/admin/update-cat-form/:id", server.updateCategoryForm)
	adminRoutes.GET("/admin/content", server.adminArts)
	adminRoutes.GET("/admin/content/create", server.createArticlePage)
	adminRoutes.GET("/admin/content/update/:id", server.updateArticlePage)
	adminRoutes.GET("/admin/pub-content", server.publishedContentList)
	adminRoutes.GET("/admin/draft-content", server.draftContentList)
	adminRoutes.GET("/admin/del-content", server.deletedContentList)
	adminRoutes.GET("/admin/users", server.adminUsers)
	adminRoutes.GET("/admin/active-users", server.activeUsersList)
	adminRoutes.GET("/admin/banned-users", server.bannedUsersList)
	adminRoutes.GET("/admin/deleted-users", server.deletedUsersList)
	adminRoutes.GET("/admin/ads", server.adminAds)
	adminRoutes.GET("/admin/active-ads", server.activeAdsList)
	adminRoutes.GET("/admin/inactive-ads", server.inactiveAdsList)
	adminRoutes.GET("/admin/scheduled-ads", server.scheduledAdsList)
	adminRoutes.GET("/admin/create-ad-modal", server.createAdModal)
	adminRoutes.GET("/admin/update-ad-modal/:id", server.updateAdModal)
	adminRoutes.GET("/admin/settings", server.adminSettings)

	// Auth Pages - no rate limiting for page views
	router.GET("/login", server.loginPage)
	router.GET("/register", server.registerPage)
	router.GET("/reset-lozinke/:token", server.passwordResetPage)
	router.GET("/zaboravljena-lozinka", server.requestPassResetPage)
	router.GET("/potvrdi-email/:token", server.emailVerifiedPage)

	// User Page Routes - no rate limiting for page views
	router.GET("/", server.homePage)
	router.GET("/pretraga", server.searchResultsPage)
	router.GET("/kategorije/:category/:id", server.categoriesPage)
	router.GET("/oznake/:tag/:id", server.tagPage)
	router.GET("/:article/:id", server.articlePage)
	authRoutes.GET("/podesavanja", server.userSettingsPage)

	// ==== API Routes with Rate Limiting ====

	// ---- Authentication API (Strict Limiting) ----
	authApiRoutes := router.Group("/api")
	authApiRoutes.Use(server.RateLimitMiddleware(authLimiter))

	authApiRoutes.POST("/login", server.login)
	authApiRoutes.POST("/register", server.register)
	authApiRoutes.POST("/reset-password", server.resetPassword)
	authApiRoutes.POST("/send-password-reset-form", server.requestPassResetFromForm)

	// Auth routes that require auth middleware
	authAuthRoutes := authRoutes.Group("/api")
	authAuthRoutes.Use(server.RateLimitMiddleware(authLimiter))

	authAuthRoutes.POST("/logout", server.logOut)
	authAuthRoutes.POST("/send-password-reset", server.requestPassReset)

	// ---- Comment API ----
	commentApiRoutes := router.Group("/api")

	// Comment listing routes - public (No rate limiting)
	commentApiRoutes.GET("/content/comments/:id", server.listContentComments)
	commentApiRoutes.GET("/content/comments/:id/score", server.listContentCommentsScore)
	commentApiRoutes.GET("/comments/:id/reply-info", server.listRepliesInfo)
	commentApiRoutes.GET("/comments/:id/more-replies", server.listCommentReplies)

	// Comment write routes - require auth
	commentAuthRoutes := authRoutes.Group("/api")
	commentAuthRoutes.Use(server.RateLimitMiddleware(commentLimiter))

	commentAuthRoutes.POST("/content/comments/:id", server.createComment)
	commentAuthRoutes.POST("/comments/:id/upvote", server.handleUpvoteComment)
	commentAuthRoutes.POST("/comments/:id/downvote", server.handleDownvoteComment)
	commentAuthRoutes.POST("/comments/:id/reply", server.createReply)
	commentAuthRoutes.PUT("/comments/:id/edit", server.updateComment)
	commentAuthRoutes.DELETE("/comments/:id", server.deleteComment)

	// ---- Content Interaction API (Moderate Limiting) ----
	contentApiRoutes := authRoutes.Group("/api")
	contentApiRoutes.Use(server.RateLimitMiddleware(contentLimiter))

	contentApiRoutes.POST("/content/like/:id", server.handleLikeContent)
	contentApiRoutes.POST("/content/dislike/:id", server.handleDislikeContent)
	router.POST("/api/increment-ads-clicks", server.incrementDailyAdsClicks) // Public route

	// ---- User Settings API (Moderate Limiting) ----
	userSettingsRoutes := authRoutes.Group("/api/admin/settings")
	userSettingsRoutes.Use(server.RateLimitMiddleware(userSettingsLimiter))

	userSettingsRoutes.PUT("/username/:id", server.updateUsername)
	userSettingsRoutes.PUT("/pfp/:id", server.updatePfp)

	// Cookie deletion
	authRoutes.DELETE("/api/cookie", server.deleteCookie)

	// ---- Search API (Moderate Limiting) ----
	searchApiRoutes := router.Group("/api")
	searchApiRoutes.Use(server.RateLimitMiddleware(searchLimiter))

	searchApiRoutes.GET("/search", server.loadMoreSearch)

	// ---- Public Read-only API (No Rate Limiting) ----
	// These routes are typically used for page loads and don't need rate limiting
	router.GET("/api/content/other", server.listOtherContent)
	router.GET("/api/news-slider", server.newsSlider)
	router.GET("/api/content/popular", server.listTrendingContentUser)
	router.GET("/api/content/categories", server.categoriesWithContent)
	router.GET("/api/category/content/recent/:id", server.listRecentCategoryContent)
	router.GET("/api/category/:id/tags/content", server.listContentByTagsUnderCategory)
	router.GET("/api/tag/content/recent/:id", server.listAllContentByTag)
	router.GET("/api/content/media/:id", server.listMediaForArticlePage)

	// ---- Admin API (Higher Limits) ----
	// Admin content management API
	adminApiRoutes := adminRoutes.Group("/api/admin")
	adminApiRoutes.Use(server.RateLimitMiddleware(adminLimiter))

	// Admin overview
	adminApiRoutes.GET("/trending", server.listTrendingContent)
	adminApiRoutes.GET("/analytics", server.getDailyAnalytics)

	// Admin categories
	adminApiRoutes.GET("/categories", server.listCats)
	adminApiRoutes.POST("/category", server.createCategory)
	adminApiRoutes.DELETE("/category/:id", server.deleteCategory)
	adminApiRoutes.PUT("/category/:id", server.updateCategory)

	// Admin articles
	adminApiRoutes.GET("/content/published", server.listPubContent)
	adminApiRoutes.GET("/content/published/oldest", server.listPubContentOldest)
	adminApiRoutes.GET("/content/published/title", server.listPubContentTitle)
	adminApiRoutes.GET("/content/draft", server.listDraftContent)
	adminApiRoutes.GET("/content/draft/oldest", server.listDraftContentOldest)
	adminApiRoutes.GET("/content/draft/title", server.listDraftContentTitle)
	adminApiRoutes.GET("/content/deleted", server.listDelContent)
	adminApiRoutes.GET("/content/deleted/oldest", server.listDelContentOldest)
	adminApiRoutes.GET("/content/deleted/title", server.listDelContentTitle)
	adminApiRoutes.GET("/content/published/search", server.listSearchPubContent)
	adminApiRoutes.GET("/content/draft/search", server.listSearchDraftContent)
	adminApiRoutes.GET("/content/deleted/search", server.listSearchDelContent)
	adminApiRoutes.PUT("/content/archive/:id", server.archivePubContent)
	adminApiRoutes.DELETE("/content/:id", server.deleteContent)
	adminApiRoutes.PUT("/content/publish/:id", server.publishDraftContent)
	adminApiRoutes.PUT("/content/unarchive/:id", server.unarchiveContent)
	adminApiRoutes.PUT("/content/:id", server.updateContent)
	adminApiRoutes.POST("/content/draft", server.createContent)
	adminApiRoutes.POST("/content/publish", server.createAndPublishContent)

	// Admin Media
	adminApiRoutes.GET("/media", server.listMediaForContent)
	adminApiRoutes.POST("/media/upload/new", server.addMediaToNewContent)
	adminApiRoutes.POST("/media/upload/:id", server.addMediaToUpdateContent)
	adminApiRoutes.DELETE("/media/remove/:id", server.deleteMedia)

	// Admin Tags
	adminApiRoutes.GET("/tags", server.listTags)
	adminApiRoutes.GET("/tags/search", server.listSearchTags)
	adminApiRoutes.GET("/tags/:id", server.listTagsByContent)
	adminApiRoutes.POST("/tags", server.createTag)
	adminApiRoutes.POST("/tags/add", server.addTagToContent)
	adminApiRoutes.POST("/tags/add/:id", server.addTagToContentUpdate)
	adminApiRoutes.DELETE("/tags/content/remove/:id", server.removeTagFromContent)
	adminApiRoutes.DELETE("/tags/content/remove/:content_id/:tag_id", server.removeTagFromContentUpdate)
	adminApiRoutes.DELETE("/tags/remove/:id", server.deleteTag)

	// Admin Users
	adminApiRoutes.GET("/users/active", server.listActiveUsers)
	adminApiRoutes.GET("/users/active/oldest", server.listActiveUsersOldest)
	adminApiRoutes.GET("/users/active/title", server.listActiveUsersTitle)
	adminApiRoutes.GET("/users/banned", server.listBannedUsers)
	adminApiRoutes.GET("/users/banned/oldest", server.listBannedUsersOldest)
	adminApiRoutes.GET("/users/banned/title", server.listBannedUsersTitle)
	adminApiRoutes.GET("/users/deleted", server.listDeletedUsers)
	adminApiRoutes.GET("/users/deleted/oldest", server.listDeletedUsersOldest)
	adminApiRoutes.GET("/users/deleted/title", server.listDeletedUsersTitle)
	adminApiRoutes.GET("/users/active/search", server.searchActiveUsers)
	adminApiRoutes.GET("/users/banned/search", server.searchBannedUsers)
	adminApiRoutes.GET("/users/deleted/search", server.searchArchivedUsers)
	adminApiRoutes.PUT("/users/ban/:id", server.banUser)
	adminApiRoutes.PUT("/users/unban/:id", server.unbanUser)
	adminApiRoutes.PUT("/users/archive/:id", server.deleteUser)

	// Admin settings
	adminApiRoutes.PUT("/global-settings", server.updateGlobalSettings)
	adminApiRoutes.PUT("/reset-global-settings", server.resetGlobalSettings)

	// Admin ads
	adminApiRoutes.GET("/ads/active", server.listActiveAds)
	adminApiRoutes.GET("/ads/inactive", server.listInactiveAds)
	adminApiRoutes.POST("/ads", server.createAd)
	adminApiRoutes.DELETE("/ads/:id", server.deleteAd)
	adminApiRoutes.PUT("/ads/:id", server.updateAd)
	adminApiRoutes.PUT("/ads/deactivate/:id", server.deactivateAd)

	server.router = router
}
