package main

import (
	"log"
	"os"

	"github.com/00mark0/macva-press/api"
	"github.com/00mark0/macva-press/db/services"

	"context"

	"github.com/00mark0/macva-press/db/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	queries *db.Queries
	store   *db.Store
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	dbSource := os.Getenv("DB_URL")
	port := os.Getenv("SERVER_PORT")
	symmetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")

	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection pool created.")

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot connect to db!:", err)
	}
	log.Println("Connected to db.")

	// Initialize Redis
	redis.InitRedis()
	pong, err := redis.Client.Ping(redis.Ctx).Result()
	if err != nil {
		log.Fatalf("Cannot connect to Redis! %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)

	// Initialize the store and pass it into the server
	store := db.NewStore(conn)

	// Ensure an admin user exists before starting the server
	if err := api.BootstrapAdmin(store); err != nil {
		log.Fatal("Bootstrap failed:", err)
	}

	// Ensure global settings row exists before starting the server
	if err := api.BootstrapGlobalSettings(store); err != nil {
		log.Fatal("Bootstrap global settings failed:", err)
	}

	// Pass both store and Redis client into the server
	server, err := api.NewServer(store, symmetricKey, redis.Client)
	if err != nil {
		log.Fatal("Cannot create server!:", err)
	}

	err = server.Start(":" + port)
	if err != nil {
		log.Fatal("Cannot start server!:", err)
	}
}
