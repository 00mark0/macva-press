// /api/server.go

package api

import (
	"fmt"
	"os"

	"github.com/00mark0/macva-press/db/redis"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/token"
	"github.com/labstack/echo/v4"
	redisClient "github.com/redis/go-redis/v9"
)

var BaseUrl = os.Getenv("BASE_URL")

// Server serves HTTP requests
type Server struct {
	store           *db.Store
	tokenMaker      token.Maker
	cacheService    *redis.CacheService // Store the cache service here
	router          *echo.Echo
	uploadSemaphore chan struct{}
}

// NewServer creates an HTTP server and sets up routing.
func NewServer(store *db.Store, symmetricKey string, redisClient *redisClient.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// Create a CacheService instance from the redis client
	cacheService := redis.NewCacheService(redisClient)

	server := &Server{
		store:        store,
		tokenMaker:   tokenMaker,
		cacheService: cacheService, // Pass CacheService to server
	}

	server.setupRouter()

	return server, nil
}

// Start runs the HTTP server.
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}
