package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/token"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

	"github.com/00mark0/macva-press/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"os"

	"github.com/labstack/echo/v4"
)

type userResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Pfp      string `json:"pfp"`
	Role     string `json:"role"`
}

type loginUserReq struct {
	Email      string `json:"email" form:"email" validate:"required,email"`
	Password   string `json:"password" form:"password" validate:"required"`
	RememberMe bool   `json:"remember_me" form:"remember_me"`
}

type loginUserRes struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) login(ctx echo.Context) error {
	var req loginUserReq
	var loginErr components.LoginErr
	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in login:", err)
		return err
	}

	// Run validation
	if err := ctx.Validate(req); err != nil {
		var loginErr components.LoginErr

		// Loop through validation errors and handle each field separately
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Email":
				// Custom error message for email validation
				loginErr = "Email je obavezan i mora biti validan."
			case "Password":
				// Custom error message for password validation
				loginErr = "Lozinka je obavezna."
			default:
				loginErr = "Nevažecí podaci za prijavu."
			}
		}

		// Render the login form with the custom error message
		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	user, err := server.store.GetUserByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		loginErr = "Nevažeći podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	if user.Banned.Bool {
		loginErr = "Nevažecí podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	if !user.EmailVerified.Bool {
		loginErr = "Email nije verifikovan. Poslat je nov link za verifikaciju na vašu adresu."

		token, err := utils.GenerateToken(jwt.MapClaims{
			"user_id": user.UserID.String(),
		}, time.Hour*24)
		if err != nil {
			log.Println("Error generating token in requestPassReset:", err)
			return err
		}

		verifyEmailLink := fmt.Sprintf("%s/potvrdi-email/%s", BaseUrl, token)

		err = utils.SendEmailVerificationEmail(req.Email, verifyEmailLink)
		if err != nil {
			log.Println("Error sending email verification email in login:", err)
			return err
		}

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))

	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		loginErr = "Nevažecí podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	durationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Println("Error parsing duration in login:", err)
		return err
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.UserID.String(),
		user.Username,
		user.Email,
		user.Pfp,
		user.Role,
		user.EmailVerified.Bool,
		user.Banned.Bool,
		user.IsDeleted.Bool,
		duration,
	)
	if err != nil {
		log.Println("Error creating token in login:", err)
		return err
	}

	refreshTokenDurationStr := os.Getenv("REFRESH_TOKEN_DURATION")
	refreshTokenDuration, err := time.ParseDuration(refreshTokenDurationStr)
	if err != nil {
		log.Println("Error parsing duration in login:", err)
		return err
	}

	log.Printf("Remember me value: %v", req.RememberMe)
	if req.RememberMe {
		extendedDurationStr := os.Getenv("REMEMBER_ME_DURATION") // Fetch from .env
		extendedDuration, err := time.ParseDuration(extendedDurationStr)
		if err != nil {
			log.Println("Error parsing extended refresh token duration:", err)
			return err
		}
		refreshTokenDuration = extendedDuration
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.UserID.String(),
		user.Username,
		user.Email,
		user.Pfp,
		user.Role,
		user.EmailVerified.Bool,
		user.Banned.Bool,
		user.IsDeleted.Bool,
		refreshTokenDuration,
	)
	if err != nil {
		log.Println("Error creating token in login:", err)
		return err
	}

	session, err := server.store.CreateSession(ctx.Request().Context(), db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: refreshTokenPayload.ID, Valid: true},
		UserID:       user.UserID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request().UserAgent(),
		ClientIp:     ctx.RealIP(),
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshTokenPayload.ExpiredAt, Valid: true},
	})
	if err != nil {
		log.Println("Error creating session in login:", err)
		return err
	}

	// Set token as a secure, HTTP-only cookie
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(duration),
		Path:     "/",
		HttpOnly: true,
	})

	ctx.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenDuration),
		Path:     "/",
		HttpOnly: true,
	})

	ctx.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Expires:  time.Now().Add(refreshTokenDuration),
		Path:     "/",
		HttpOnly: true,
	})

	ctx.Response().Header().Set("HX-Redirect", "/")
	return ctx.NoContent(http.StatusOK)
}

func (server *Server) logOut(ctx echo.Context) error {
	sessionCookie, err := ctx.Cookie("session_id")
	if err != nil {
		log.Println("Error getting session cookie in logOut:", err)
		return err
	}

	sessionIDStr := sessionCookie.Value
	sessionID, err := utils.ParseUUID(sessionIDStr, "session ID")
	if err != nil {
		log.Println("Invalid session ID format in logOut:", err)
		return err
	}

	err = server.store.DeleteSession(ctx.Request().Context(), sessionID)
	if err != nil {
		log.Println("Error deleting session in logOut:", err)
		return err
	}

	// Clear all cookies
	clearCookie := func(name string) {
		ctx.SetCookie(&http.Cookie{
			Name:   name,
			Value:  "",
			Path:   "/",
			MaxAge: -1, // Expire immediately
		})
	}
	clearCookie("access_token")
	clearCookie("refresh_token")
	clearCookie("session_id")

	ctx.Response().Header().Set("HX-Redirect", "/")
	return ctx.NoContent(http.StatusOK)
}

type ListUsersReq struct {
	Limit int32 `query:"limit"`
}

type SearchUserReq struct {
	SearchTerm string `query:"search_term" validate:"required"`
	Limit      int32  `query:"limit"`
}

func (server *Server) listActiveUsers(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listActiveUsers:", err)
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

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listActiveUsersOldest(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveUsersOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsersOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listActiveUsersOldest:", err)
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

	url := "/api/admin/users/active/oldest?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listActiveUsersTitle(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveUsersTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsersTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listActiveUsersTitle:", err)
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

	url := "/api/admin/users/active/title?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listBannedUsers(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listBannedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listBannedUsers:", err)
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

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listBannedUsersOldest(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listBannedUsersOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsersOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listBannedUsersOldest:", err)
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

	url := "/api/admin/users/banned/oldest?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listBannedUsersTitle(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listBannedUsersTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsersTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listBannedUsersTitle:", err)
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

	url := "/api/admin/users/banned/title?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listDeletedUsers(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDeletedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listDeletedUsers:", err)
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

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listDeletedUsersOldest(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDeletedUsersOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsersOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listDeletedUsersOldest:", err)
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

	url := "/api/admin/users/deleted/oldest?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listDeletedUsersTitle(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDeletedUsersTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsersTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listDeletedUsersTitle:", err)
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

	url := "/api/admin/users/deleted/title?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) searchActiveUsers(ctx echo.Context) error {
	var req SearchUserReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchActiveUsers:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchActiveUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchActiveUsersParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	var users []components.UsersRes

	data, err := server.store.SearchActiveUsers(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching users in searchActiveUsers:", err)
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

	url := "/api/admin/users/active/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) searchBannedUsers(ctx echo.Context) error {
	var req SearchUserReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchBannedUsers:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchBannedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchBannedUsersParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	var users []components.UsersRes

	data, err := server.store.SearchBannedUsers(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching users in searchBannedUsers:", err)
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

	url := "/api/admin/users/banned/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) searchArchivedUsers(ctx echo.Context) error {
	var req SearchUserReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchArchivedUsers:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchArchivedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchDeletedUsersParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	var users []components.UsersRes

	data, err := server.store.SearchDeletedUsers(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching users in searchArchivedUsers:", err)
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

	url := "/api/admin/users/deleted/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) banUser(ctx echo.Context) error {
	userIDStr := ctx.Param("id")
	userID, err := utils.ParseUUID(userIDStr, "user_id")
	if err != nil {
		log.Println("Error parsing user id in banUser:", err)
		return err
	}

	err = server.store.BanUser(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error banning user in banUser:", err)
		return err
	}

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

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "user*")
	if err != nil {
		log.Println("Error deleting cache in banUser:", err)
		return err
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Println("Error deleting cache in banUser:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UsersNav(overview))
}

func (server *Server) unbanUser(ctx echo.Context) error {
	userIDStr := ctx.Param("id")
	userID, err := utils.ParseUUID(userIDStr, "user_id")
	if err != nil {
		log.Println("Error parsing user id in unbanUser:", err)
		return err
	}

	err = server.store.UnbanUser(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error unbanning user in unbanUser:", err)
		return err
	}

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

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "user*")
	if err != nil {
		log.Println("Error deleting cache in unbanUser:", err)
		return err
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Println("Error deleting cache in banUser:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UsersNav(overview))
}

func (server *Server) deleteUser(ctx echo.Context) error {
	userIDStr := ctx.Param("id")
	userID, err := utils.ParseUUID(userIDStr, "user_id")
	if err != nil {
		log.Println("Error parsing user id in deleteUser:", err)
		return err
	}

	err = server.store.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error deleting user in deleteUser:", err)
		return err
	}

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

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "user*")
	if err != nil {
		log.Println("Error deleting user from cache in deleteUser:", err)
		return err
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Println("Error deleting cache in banUser:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UsersNav(overview))
}

type UserInfoReq struct {
	Username string `form:"username" validate:"required,min=3,max=20"`
	Pfp      string `form:"pfp"`
}

func (server *Server) updateUsername(ctx echo.Context) error {
	var req UserInfoReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updateUsername:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Korisničko ime mora biti između 3 i 20 karaktera."

		return Render(ctx, http.StatusOK, components.UpdateError(message))
	}

	userIDStr := ctx.Param("id")
	userID, err := utils.ParseUUID(userIDStr, "user_id")
	if err != nil {
		log.Println("Error parsing user id in updateUsername:", err)
		return err
	}

	user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error getting user in updateUsername:", err)
		return err
	}

	arg := db.UpdateUserParams{
		UserID:   userID,
		Username: req.Username,
		Pfp:      user.Pfp,
	}

	err = server.store.UpdateUser(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating user in updateUsername:", err)
		return err
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "user*")
	if err != nil {
		log.Println("Error deleting cache in updateUsername:", err)
		return err
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Println("Error deleting cache in banUser:", err)
		return err
	}

	message := "Korisničko ime je uspešno promenjeno."
	return Render(ctx, http.StatusOK, components.UpdateSuccess(message))
}

type UpdatePfpReq struct {
	Pfp string `form:"pfp"`
}

func (server *Server) updatePfp(ctx echo.Context) error {
	var req UpdatePfpReq
	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updatePfp:", err)
		return err
	}

	userIDStr := ctx.Param("id")
	userID, err := utils.ParseUUID(userIDStr, "user_id")
	if err != nil {
		log.Println("Error parsing user id in updatePfp:", err)
		return err
	}

	user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error getting user in updatePfp:", err)
		return err
	}

	// Delete previous profile picture if it exists and is in static/pfp directory
	if user.Pfp != "" && strings.HasPrefix(user.Pfp, "static/pfp/") {
		// Check if file exists before attempting to delete
		if _, err := os.Stat(user.Pfp); err == nil {
			if err := os.Remove(user.Pfp); err != nil {
				// Log the error but continue with the update
				log.Println("Warning: couldn't delete previous profile picture:", err)
			} else {
				log.Println("Successfully deleted previous profile picture:", user.Pfp)
			}
		}
	}

	file, err := ctx.FormFile("pfp")
	if err != nil {
		log.Println("Error retrieving uploaded file in updatePfp:", err)
		return err
	}

	uploadsDir := "static/pfp"
	// Ensure directory exists
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Println("Error creating directory in updatePfp:", err)
			return err
		}
	}

	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	src, err := file.Open()
	if err != nil {
		log.Println("Error opening file in updatePfp:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in updatePfp:", err)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Println("Error copying file data in updatePfp:", err)
		return err
	}

	convertedPath, err := ConvertToWebPWithResize(filePath, 400, 400, 80)
	if err != nil {
		log.Println("Error converting file to WebP in createAd:", err)
	} else {
		filePath = convertedPath
	}

	// Ensure ImageUrl does not have a double leading slash
	imagePath := filePath
	if !strings.HasPrefix(filePath, "/") {
		imagePath = "/" + filePath
	}

	arg := db.UpdateUserParams{
		UserID:   userID,
		Username: user.Username,
		Pfp:      imagePath,
	}

	err = server.store.UpdateUser(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating user in updatePfp:", err)
		// If update fails, delete the newly uploaded file to avoid orphaned files
		os.Remove(filePath)
		return err
	}

	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "user*")
	if err != nil {
		log.Println("Error deleting cache in updatePfp:", err)
		// If cache deletion fails, delete the newly uploaded file to avoid orphaned files
		return err
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), "comments*")
	if err != nil {
		log.Println("Error deleting cache in banUser:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.AdminPfp(filePath))
}

func (server *Server) requestPassReset(ctx echo.Context) error {
	payload := ctx.Get(authorizationPayloadKey).(*token.Payload)

	token, err := utils.GenerateToken(jwt.MapClaims{
		"user_id": payload.UserID,
	}, time.Hour)
	if err != nil {
		log.Println("Error generating token in requestPassReset:", err)
		return err
	}

	resetLink := fmt.Sprintf("%s/reset-lozinke/%s", BaseUrl, token)

	err = utils.SendPasswordResetEmail(payload.Email, resetLink)
	if err != nil {
		log.Println("Error sending password reset email in requestPassReset:", err)
		message := "Dogodila se greška prilikom slanja linka za promenu lozinke."

		return Render(ctx, http.StatusOK, components.UpdateError(message))
	}

	message := "Link za promenu lozinke je poslat na vašu email adresu."
	return Render(ctx, http.StatusOK, components.UpdateSuccess(message))
}

type ReqPassResetFormReq struct {
	Email string `form:"email" validate:"required,email"`
}

func (server *Server) requestPassResetFromForm(ctx echo.Context) error {
	var req ReqPassResetFormReq
	var reqPassResetFormErr components.RequestPassResetErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in requestPassResetFromForm:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Email":
				reqPassResetFormErr = "Email mora biti validan."
			}
		}

		return Render(ctx, http.StatusOK, components.RequestPassResetForm(reqPassResetFormErr))
	}

	user, err := server.store.GetUserByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		log.Println("Error getting user in requestPassResetFromForm:", err)
		return err
	}

	token, err := utils.GenerateToken(jwt.MapClaims{
		"user_id": user.UserID.String(),
	}, time.Hour)
	if err != nil {
		log.Println("Error generating token in requestPassResetFromForm:", err)
		return err
	}

	resetLink := fmt.Sprintf("%s/reset-lozinke/%s", BaseUrl, token)

	err = utils.SendPasswordResetEmail(req.Email, resetLink)
	if err != nil {
		log.Println("Error sending password reset email in requestPassResetFromForm:", err)

		reqPassResetFormErr = "Dogodila se greška prilikom slanja linka za promenu lozinke."
		return Render(ctx, http.StatusOK, components.RequestPassResetForm(reqPassResetFormErr))
	}

	reqPassResetFormErr = "Link za promenu lozinke je poslat na vašu email adresu."
	return Render(ctx, http.StatusOK, components.RequestPassResetForm(reqPassResetFormErr))
}

type PasswordResetReq struct {
	Token           string `form:"token" validate:"required"`
	Password        string `form:"password" validate:"required,password"`
	ConfirmPassword string `form:"confirmPassword" validate:"required,eqfield=Password"`
}

func (server *Server) resetPassword(ctx echo.Context) error {
	var req PasswordResetReq
	var resetErr components.ResetErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in resetPassword:", err)
		return err
	}

	// Run validation
	if err := ctx.Validate(req); err != nil {
		// Loop through validation errors and handle each field separately
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Token":
				// Custom error message for email validation
				resetErr = "Link za resetovanje lozinke je nevažeći."
			case "Password":
				// Custom error message for password validation
				resetErr = "Lozinka mora imati najmanje 8 karaktera, uključujući jedno veliko slovo, jedno malo slovo, jedan broj i jedan specijalni karakter."
			case "ConfirmPassword":
				// Custom error message for password validation
				resetErr = "Potvrda lozinke mora biti ista kao i originalna lozinka."
			}
		}

		// Render the login form with the custom error message
		return Render(ctx, http.StatusOK, components.ResetForm(req.Token, resetErr))
	}

	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		resetErr = "Link za resetovanje lozinke je nevažeci."

		return Render(ctx, http.StatusOK, components.ResetForm(req.Token, resetErr))
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		log.Println("Error extracting user_id from claims in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.ResetForm(req.Token, "Link za resetovanje lozinke je nevažeći."))
	}

	// Parse user ID to UUID
	userID, err := utils.ParseUUID(userIDStr, "user_id")
	if err != nil {
		log.Println("Error parsing user_id to UUID in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.ResetForm(req.Token, "Link za resetovanje lozinke je nevažeci."))
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return Render(ctx, http.StatusOK, components.ResetForm(req.Token, "Greška pri obradi lozinke."))
	}

	arg := db.UpdateUserPasswordParams{
		UserID:   userID,
		Password: hashedPassword,
	}

	err = server.store.UpdateUserPassword(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating user password in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.ResetForm(req.Token, "Greška pri obradi lozinke."))
	}

	// Password updated successfully
	return Render(ctx, http.StatusOK, components.ResetSuccess())
}

type RegisterReq struct {
	Username        string `form:"username" validate:"required,username"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,password"`
	ConfirmPassword string `form:"confirmPassword" validate:"required,eqfield=Password"`
}

func (server *Server) register(ctx echo.Context) error {
	var req RegisterReq
	var registerErr components.RegisterErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in register:", err)
		return err
	}

	// Run validation
	if err := ctx.Validate(req); err != nil {
		// Loop through validation errors and handle them
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Username":
				switch fieldErr.Tag() {
				case "required":
					registerErr = "Korisničko ime je obavezno."
				case "username":
					registerErr = "Korisničko ime mora biti između 3 i 20 karaktera, ne sme početi brojem, i može sadržavati samo slova, brojeve, donje crtice i crtice."
				}
			case "Email":
				registerErr = "Email mora biti validan."
			case "Password":
				registerErr = "Lozinka mora imati najmanje 8 karaktera, uključujući jedno veliko slovo, jedno malo slovo, jedan broj i jedan specijalni karakter."
			case "ConfirmPassword":
				registerErr = "Potvrda lozinke mora biti ista kao i originalna lozinka."
			}
		}

		// Render the form with the custom error message
		return Render(ctx, http.StatusOK, components.RegisterForm(registerErr))
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return Render(ctx, http.StatusOK, components.RegisterForm("Greška pri obradi lozinke."))
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error creating user in register:", err)
		return Render(ctx, http.StatusOK, components.RegisterForm("Greška pri pravljenju naloga."))
	}

	token, err := utils.GenerateToken(jwt.MapClaims{
		"user_id": user.UserID.String(),
	}, time.Hour*24)
	if err != nil {
		log.Println("Error generating token in requestPassReset:", err)
		return err
	}

	verifyEmailLink := fmt.Sprintf("%s/potvrdi-email/%s", BaseUrl, token)

	err = utils.SendEmailVerificationEmail(req.Email, verifyEmailLink)
	if err != nil {
		log.Println("Error sending email verification email in register:", err)
		return Render(ctx, http.StatusOK, components.RegisterForm("Greška pri slanju email-a za verifikaciju naloga."))
	}

	return Render(ctx, http.StatusOK, components.RegisterSuccess())
}
