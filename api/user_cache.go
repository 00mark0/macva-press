package api

import (
	"log"
	"time"

	"github.com/00mark0/macva-press/db/redis"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
	"github.com/labstack/echo/v4"
)

func (server *Server) getUserFromCacheOrDb(ctx echo.Context, cookieName string) (db.GetUserByIDRow, error) {
	var userData db.GetUserByIDRow

	cookie, err := ctx.Cookie(cookieName)
	if err != nil {
		log.Println("Error getting cookie in getUserFromCookieOrCache:", err)
		return userData, err
	}

	payload, err := server.tokenMaker.VerifyToken(cookie.Value)
	if err != nil {
		log.Println("Error verifying token in getUserFromCookieOrCache:", err)
		return userData, err
	}

	cacheKey := redis.GenerateKey("user", payload.UserID)

	cacheHit, err := server.cacheService.Get(ctx.Request().Context(), cacheKey, &userData)
	if err != nil {
		log.Printf("Error fetching user from cache: %v", err)
	}
	if cacheHit {
		return userData, nil
	}

	log.Printf("Cache miss for user: %s", cacheKey)
	userID, err := utils.ParseUUID(payload.UserID, "userID")
	if err != nil {
		log.Println("Error parsing user_id in getUserFromCookieOrCache:", err)
		return userData, err
	}

	userData, err = server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return userData, err
	}

	if err := server.cacheService.Set(ctx.Request().Context(), cacheKey, &userData, 10*time.Minute); err != nil {
		log.Printf("Error setting user in cache: %v", err)
	}

	return userData, nil
}
